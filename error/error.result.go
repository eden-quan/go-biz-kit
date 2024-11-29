package errorutil

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	contextutil "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/context"
	"maps"
	http2 "net/http"
	"reflect"
	"sync"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	errorpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/error"
	"go.opentelemetry.io/otel/trace"

	common "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/common/def"
)

type TruncateToEmptyErrorInterface interface {
	IsTruncateToEmpty() error
}

type TruncateToEmptyError struct {
	error
}

func IsEmptyError(err error) bool {
	return errorpkg.IsEmptyError(err)
}

func (t *TruncateToEmptyError) IsTruncateToEmpty() error {
	return t.error
}

func IsTruncateToEmptyError(err error) (*TruncateToEmptyError, bool) {
	var t *TruncateToEmptyError
	ok := errors.As(err, &t)
	return t, ok
}

func ErrorResultMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			_, ok := transport.FromServerContext(ctx)

			if !ok {
				return handler(ctx, req)
			}

			reply, err = handler(ctx, req)
			//if reply == nil {
			//	return
			//}

			if errorpkg.IsEmptyError(err) {
				err = nil
			}

			// make sure we are in grpc server side
			span := trace.SpanFromContext(ctx)
			data := &common.Result{TraceId: span.SpanContext().TraceID().String()}

			// check if it has the flatten field of result
			processed, result := matchAndUpdate(reply, err, data)
			if processed {
				// 根据错误配置统一处理 Http 状态码
				if info, e := errorpkg.NewErrorMetaInfo(err); e == nil {
					leaf := info.Leaf()
					if httpCtx, ok := ExtractHttpContext(ctx); ok && leaf.Code != 0 {
						err = httpCtx.Result(int(leaf.Code), result)
						if err == nil {
							httpCtx.Reset(&EmptyWriter{origWriter: httpCtx.Response()}, httpCtx.Request())
						}
					}
				}

				err = &TruncateToEmptyError{error: err}
			}

			reply = result
			return
		}
	}
}

type EmptyWriter struct {
	origWriter http.ResponseWriter
}

func (e *EmptyWriter) Header() http2.Header {
	return e.origWriter.Header()
}

func (e *EmptyWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (e *EmptyWriter) WriteHeader(_ int) {}

type resultTypeChecker struct {
	//resultType *common.Result
	resultType interface{}
	fieldCount int
	fieldList  []string
	fieldMap   map[string]*reflect.StructField
}

var checkers []*resultTypeChecker = []*resultTypeChecker{}
var checkerOnce sync.Once

func initChecker() {

	initWithType := func(resultType interface{}) {
		checker := &resultTypeChecker{}
		checker.resultType = resultType
		checker.fieldMap = make(map[string]*reflect.StructField)

		value := reflect.ValueOf(checker.resultType)
		valueType := value.Elem().Type()
		checker.fieldCount = valueType.NumField()
		for i := 0; i < valueType.NumField(); i++ {
			f := valueType.Field(i)

			if !f.IsExported() {
				continue
			}

			checker.fieldList = append(checker.fieldList, f.Name)
			checker.fieldMap[f.Name] = &f
		}

		checkers = append(checkers, checker)
	}

	initWithType(&common.SimpleResult{})
	initWithType(&common.Result{})

}

func matchAndUpdateForType(checker *resultTypeChecker, result interface{}, value reflect.Value, valueType reflect.Type, err error,
	data *common.Result) bool {

	if valueType.NumField() < checker.fieldCount {
		return false
	}

	for k, v := range checker.fieldMap {
		f, ok := valueType.FieldByName(k)
		if !ok || f.Type.Name() != v.Type.Name() {
			return false
		}
	}

	processed := false
	if errInfo, e := errorpkg.NewErrorMetaInfo(err); e == nil && errInfo.Leaf() != nil {
		meta := errInfo.Leaf()
		data.Code = meta.BizCode
		data.Reason = meta.Reason
		data.Message = meta.Message
		if data.Message == "" {
			data.Message = errInfo.Error()
		}
		if data.Message == "" {
			data.Message = meta.Reason
		}
		data.ErrorChain = errInfo.ErrorStack()
		data.MetaData = maps.Clone(meta.CleanError().Metadata)
		processed = true
	}

	dataValue := reflect.ValueOf(data)

	// check success, try update the error fields
	elem := value.Elem()

	for k := range checker.fieldMap {
		field := elem.FieldByName(k)
		dataField := dataValue.Elem().FieldByName(k)

		if !field.IsZero() {
			continue
		}

		if field.CanSet() {
			field.Set(dataField)
		}
	}

	return processed

}

// matchAndUpdate 处理错误信息后返回该错误是否已被处理
func matchAndUpdate(reply interface{}, err error, data *common.Result) (processed bool, result interface{}) {
	processed = false
	checkerOnce.Do(initChecker)

	isNil := false

	if reply == nil {
		reply = &common.SimpleResult{}
		isNil = true
	}

	value := reflect.ValueOf(reply)
	if value.IsNil() {
		// 尝试构建一个默认值，确认是否符合基础字段
		t := reflect.TypeOf(reply)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		nV := reflect.New(t).Elem()
		reply = nV.Addr().Interface()
		value = reflect.ValueOf(reply)
	}

	valueType := value.Elem().Type()

	for _, checker := range checkers {
		if matchAndUpdateForType(checker, reply, value, valueType, err, data) {
			processed = true
			break
		}
	}

	result = reply
	if isNil {
		processed = false
	}
	return
}

func ExtractHttpContext(ctx context.Context) (http.Context, bool) {
	hCtx, ok := ctx.Value(contextutil.HttpContextKey).(http.Context)
	return hCtx, ok
}

package errorutil

import (
	"context"
	"errors"
	"maps"
	"reflect"
	"sync"

	errorpkg "github.com/eden-quan/go-kratos-pkg/error"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/trace"

	common "github.com/eden-quan/go-biz-kit/common/def"
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
			if reply == nil {
				return
			}

			if errorpkg.IsEmptyError(err) {
				err = nil
			}

			// make sure we are in grpc server side
			span := trace.SpanFromContext(ctx)
			data := &common.Result{TraceId: span.SpanContext().TraceID().String()}

			// check if it has the flatten field of result
			processed := matchAndUpdate(reply, err, data)
			if processed {
				err = &TruncateToEmptyError{error: err}
			}

			return
		}
	}
}

type resultTypeChecker struct {
	resultType *common.Result
	fieldCount int
	fieldList  []string
	fieldMap   map[string]*reflect.StructField
}

var checker resultTypeChecker
var checkerOnce sync.Once

func initChecker() {
	checker.resultType = &common.Result{}
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
}

// matchAndUpdate 处理错误信息后返回该错误是否已被处理
func matchAndUpdate(v interface{}, err error, data *common.Result) (processed bool) {
	processed = false
	checkerOnce.Do(initChecker)

	value := reflect.ValueOf(v)
	if value.IsNil() {
		return
	}

	valueType := value.Elem().Type()

	if valueType.NumField() < checker.fieldCount {
		return
	}

	for k, v := range checker.fieldMap {
		f, ok := valueType.FieldByName(k)
		if !ok || f.Type.Name() != v.Type.Name() {
			return
		}
	}

	if errInfo, e := errorpkg.NewErrorMetaInfo(err); e == nil && errInfo.Leaf() != nil {
		meta := errInfo.Leaf()
		data.Code = meta.BizCode
		data.Reason = meta.Reason
		data.Message = meta.Message
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

	return
}

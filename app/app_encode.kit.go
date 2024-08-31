package apputil

import (
	"maps"
	stdhttp "net/http"

	kjson "github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	apppkg "github.com/eden/go-kratos-pkg/app"
	errorpkg "github.com/eden/go-kratos-pkg/error"
	headerpkg "github.com/eden/go-kratos-pkg/header"

	common "github.com/eden/go-biz-kit/common/def"
	"github.com/eden/go-biz-kit/encoding/json"
	"github.com/eden/go-biz-kit/utils/booleans"
)

// RewriteJSONEncoding 覆盖 重写 json 响应
// hub-pkg + kratos的encoding
func RewriteJSONEncoding() {
	kjson.MarshalOptions = json.MarshalOptions
	kjson.UnmarshalOptions = json.UnmarshalOptions
}

func init() {
	// 确保全局使用统一的 Encoding 机制
	// 如果特定的模块需要自定义，则应该通过 encoding/json 的 Marshaller 接口进行自定义
	RewriteJSONEncoding()
}

var _ = http.DefaultErrorEncoder

// ErrorEncoder http.DefaultErrorEncoder
// TODO: move it to server.http.error.handler
func ErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	// 在websocket时日志干扰：http: superfluous response.WriteHeader call from xxx(file:line)
	// 在websocket时日志干扰：http: response.Write on hijacked connection from
	// is websocket
	if headerpkg.GetIsWebsocket(r.Header) {
		return
	}

	traceId := ""
	if tr, ok := transport.FromServerContext(r.Context()); ok {
		traceId = tr.ReplyHeader().Get(headerpkg.TraceID)
	}

	// 响应错误
	data := &common.Result{
		TraceId: traceId,
	}

	// default code is 500
	httpCode := 200

	if info, e := errorpkg.NewErrorMetaInfo(err); e == nil {
		// server side error
		meta := info.Leaf()
		httpCode = int(meta.GetCode())
		data.Code = meta.BizCode
		data.Reason = meta.Reason
		data.Message = meta.Message
		data.ErrorChain = info.ErrorStack()
		data.MetaData = maps.Clone(meta.CleanError().Metadata)
	} else if meta := errorpkg.MetaFromError(err); meta != nil {
		// client side error
		httpCode = int(meta.GetCode())
		data.Code = meta.BizCode
		data.Reason = meta.Reason
		data.Message = meta.DefaultMessage
		data.ErrorChain = meta.Message
		data.MetaData = maps.Clone(meta.CleanError().Metadata)
	} else {
		// 兼容没按规范使用错误的代码
		se := errorpkg.FromError(err)
		code := int(se.GetCode())
		data.Code = se.GetCode()
		data.Reason = se.GetReason()
		data.Message = se.GetMessage()
		data.MetaData = se.GetMetadata()

		if booleans.Any(
			code == stdhttp.StatusUnauthorized,
			code == stdhttp.StatusForbidden,
			code == stdhttp.StatusTooManyRequests,
			booleans.All(code >= 500, code < 600),
		) {
			httpCode = code
		}
	}

	if _, exists := data.MetaData[errorpkg.BizCodeKey]; exists {
		delete(data.MetaData, errorpkg.BizCodeKey)
	}
	if defaultMsg, exists := data.MetaData[errorpkg.DefaultMessageKey]; exists && data.Message == "" {
		data.Message = defaultMsg
	}

	codec, _ := http.CodecForRequest(r, "Accept")
	apppkg.SetResponseContentType(w, codec)

	// // return
	body, err := codec.Marshal(data)
	if err != nil {
		w.WriteHeader(stdhttp.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(httpCode)
	_, _ = w.Write(body)

	return
}

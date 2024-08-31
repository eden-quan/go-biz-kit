package json

import (
	stdjson "encoding/json"
	"errors"
	"reflect"

	"github.com/go-kratos/kratos/v2/encoding"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
		RecursionLimit: 100,
	}
)

func init() {
	encoding.RegisterCodec(defaultMarshaller)
}

// defaultMarshaller 为默认的序列化器，他使用全局的序列化/反序列化配置
// 如果全局配置出现变更，他的行为也会同时发生变更
var defaultMarshaller Marshaller = &marshalImpl{
	encode: &MarshalOptions,
	decode: &UnmarshalOptions,
}

// Marshaller 为 JSON 序列化/反序列化接口，encoding/json 提供该接口的全局及实例化两种级别接口
// 当需要自定义的序列化机制时，可获取 Marshaller 的实例
// 该接口默认配置为：
//   - 序列化时将枚举类型转换为数值 / 保留默认值/空值字段
//   - 反序列化时忽略未知的字段 / 如果遇到递归的字段则通过 RecursionLimit 进行限制 (默认最多递归 100 层)
type Marshaller interface {
	Marshal(obj interface{}) ([]byte, error)
	Unmarshal(data []byte, obj interface{}) error
	MarshalToObjectField(fieldName string, obj interface{}) ([]byte, error)
	Name() string
}

type marshalImpl struct {
	encode *protojson.MarshalOptions
	decode *protojson.UnmarshalOptions
}

func (m *marshalImpl) Name() string {
	return "json"
}

func (m *marshalImpl) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case stdjson.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return stdjson.Marshal(m)
	}
}

func (m *marshalImpl) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case stdjson.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}

		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return stdjson.Unmarshal(data, m)
	}
}

func (m *marshalImpl) MarshalToObjectField(fieldName string, obj interface{}) ([]byte, error) {
	if len(fieldName) == 0 {
		// TODO: 统一错误类型
		return nil, errors.New("marshal object to json with empty field name")
	}

	d := map[string]interface{}{fieldName: obj}
	return MarshalJSON(d)
}

func NewJsonEncoder() Marshaller {
	return &marshalImpl{
		encode: &MarshalOptions,
		decode: &UnmarshalOptions,
	}
}

// NewJsonEncoderProvider 提供 JSON 编码器的依赖注入, 需要自主获取 Marshaller 的用户也可通过
// 该函数获取序列化器的实例
func NewJsonEncoderProvider(encode *protojson.MarshalOptions, decode *protojson.UnmarshalOptions) Marshaller {
	return &marshalImpl{
		encode: encode,
		decode: decode,
	}
}

// MarshalJSON 编码 json
func MarshalJSON(v interface{}) ([]byte, error) {
	return defaultMarshaller.Marshal(v)
}

// UnmarshalJSON 解码json
func UnmarshalJSON(data []byte, v interface{}) error {
	return defaultMarshaller.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return MarshalJSON(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return UnmarshalJSON(data, v)
}

// MarshalToObjectField 将 v 作为匿名对象的 fieldName 字段的值, 序列化为匿名对象的 Json
func MarshalToObjectField(fieldName string, v interface{}) ([]byte, error) {
	return defaultMarshaller.MarshalToObjectField(fieldName, v)
}

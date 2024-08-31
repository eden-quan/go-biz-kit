package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/eden/go-biz-kit/encoding/json"
)

const (
	TagNameService = "conf_service"
	TagNamePath    = "conf_path"
	TagNameValue   = "conf_value"
	TagNameJson    = "json"
)

var AllTagName = []string{TagNameValue, TagNamePath, TagNameService, TagNameJson}

// traceInfo 为本次解析的结果，包括了本次配置中的服务信息以及需要跟踪的所有配置信息
type traceInfo struct {
	service    *traceObject   // service 为当前配置的服务名，后续用于合并配置, service 的 objectType 为 TagNameService
	objects    []*traceObject // objects 为当前配置对象的各个字段, 保存到 objects 的都是 path 类型的字段
	objectsMap map[string][]*traceObject
}

// traceObject 从需要建立跟踪的对象中提取配置信息
type traceObject struct {
	field       *reflect.StructField // field 为当前字段在对象中的字段信息
	value       *reflect.Value       // value 为该字段对应的指针，后续用于热更新使用, 非 path 类型的 object value 为 nil
	path        string               // path 为当前字段的配置路径，如果是 service，则为 service 名
	pathArray   []string             // pathArray 为 path 以 / 切割后的数组，用于后续分析监听路径使用
	objectType  string               // objectType 使用 TagName*** 来定义，支持服务，路径及值三种类型
	valueFields []*traceObject       // valueFields 只有在当前字段为 path 类型时存在，他保存了当前 path 下所需的字段信息，用于快速热更新对象
}

// newTraceObject 创建一个跟踪对象，每个跟踪对象表示一个需要监听的配置
func newTraceObject(field *reflect.StructField) *traceObject {
	obj := &traceObject{
		field:       field,
		value:       nil,
		path:        "",
		objectType:  "",
		valueFields: nil,
	}

	if field != nil {
		for _, tagName := range AllTagName {
			path, exists := field.Tag.Lookup(tagName)
			if exists {
				obj.path = path
				obj.pathArray = strings.Split(
					strings.Trim(path, "/"), "/",
				)
				obj.objectType = tagName
				break
			}
		}
	}

	if obj.objectType == "" {
		// the field has no config tag
		// TODO: write log
	}

	return obj
}

func (t *traceObject) addValueField(value *traceObject) {
	t.valueFields = append(t.valueFields, value)
}

func (t *traceObject) isPath() bool {
	return t.objectType == TagNamePath
}

func (t *traceObject) isValue() bool {
	return t.objectType == TagNameValue || t.objectType == TagNameJson
}

func (t *traceObject) setReflectValue(value *reflect.Value) {
	t.value = value
}

func (t *traceInfo) addPathField(objectField *traceObject) {
	t.objects = append(t.objects, objectField)
	if t.objectsMap == nil {
		t.objectsMap = make(map[string][]*traceObject)
	}
	t.objectsMap[objectField.path] = append(t.objectsMap[objectField.path], objectField)
}

// setService 设置当前配置的服务名，如果已设置过服务名，则返回重复设置服务名异常 ConfigDuplicateServiceError
func (t *traceInfo) setService(service *traceObject) error {
	if t.service != nil {
		return fmt.Errorf("init config with duplicate service %s and %s", t.service.path, service.path)
	}

	t.service = service
	return nil
}

// analyseObj 解析 object 中的字段，将 value 类型的配置加入自身的 valueFields, 将 path 类型的配置加入 traceInfo 的 objects 供后续分析
func (t *traceInfo) analyseObj(object reflect.Value, _ *traceObject) error {
	typeOfConfig := object.Type()
	if typeOfConfig.Kind() == reflect.Pointer {
		typeOfConfig = typeOfConfig.Elem()
	}

	if typeOfConfig.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < typeOfConfig.NumField(); i++ {
		f := typeOfConfig.Field(i)

		if !unicode.IsUpper([]rune(f.Name)[0]) {
			continue
		}

		//traceField := &info.service
		traceField := newTraceObject(&f)
		// 匿名类型说明使用者通过基础 Config 自定义其他配置，需要保存匿名类型的 Tag 用于后续合并配置
		if f.Anonymous {
			err := t.setService(traceField)
			if err != nil {
				return err
			}
		}

		// 过滤未定义元数据的字段
		if traceField.isPath() {
			t.addPathField(traceField)
		}

		if object.Kind() == reflect.Pointer {
			object = object.Elem()
		}

		fieldValue := object.Field(i)
		if fieldValue.Kind() == reflect.Pointer {
			fieldValue = fieldValue.Elem()
		}
		traceField.setReflectValue(&fieldValue)
		if traceField.isValue() {
			// value 意味着已经是叶子节点，无需再递归进行检查
			continue
		}

		err := t.analyseObj(fieldValue, traceField)
		if err != nil {
			return fmt.Errorf("get value fields from %s with error %s", traceField.path, err)
		}
	}

	return nil
}

// setValue 更新 path 对应对象的值, 如果 path 不在前置监听的路径内，则不做任何操作
func (t *traceInfo) setValue(path string, value []byte) error {
	objs, exists := t.objectsMap[path]
	if !exists {
		return nil
	}

	var err []error
	for _, obj := range objs {
		// TODO: 1. 增加错误处理，先使用临时数据 Unmarshal, 避免出错时影响原始数据
		// TODO: 2. 增加局部更新对象的能力
		addr := obj.value.Addr()
		p := addr.Interface()

		marshalErr := json.UnmarshalJSON(value, p)
		if marshalErr != nil {
			marshalErr = fmt.Errorf(
				"parse config from path %s to object %s failed with error %s, please check the registration",
				path, addr.Type().String(), marshalErr)
		}

		err = append(err, marshalErr)
	}

	return errors.Join(err...)
}

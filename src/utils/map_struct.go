package utils

import (
	"gameserver/utils/log"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Map2Struct(m map[string]interface{}, i interface{}) error {
	d, err := jsoniter.Marshal(m)
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(d, i)
	if err != nil {
		return err
	}
	return nil
}

// 定义一个方法,将结构体A转化为B
// target 接受方必须是指针类型
// 返回被拷贝方不存在的字段
func StructAtoB(target interface{}, src interface{}) {
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Invalid {
		log.Error("src Invalid")
		return
	}

	if reflect.ValueOf(target).Kind() == reflect.Invalid {
		log.Error("target Invalid")
		return
	}

	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		log.Error("target pointer required")
		return
	}

	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetType = reflect.TypeOf(target).Elem()


	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		if v, ok := srcValue.Type().FieldByName(field.Name); ok {
			if v.Type.Name() == field.Type.Name() {
				targetValue.Field(i).Set(srcValue.FieldByName(field.Name))
				continue
			}
		}

	}
}

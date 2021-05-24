package utils

import (
	"reflect"
	"github.com/astaxie/beego"
	"strings"
)

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetFieldTagName(structName interface{}, tagName string) map[string]string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	fieldNum := t.NumField()
	result := make(map[string]string)
	for i := 0; i < fieldNum; i++ {
		//result = append(result, t.Field(i).Name)
		result[t.Field(i).Name] = t.Field(i).Tag.Get(tagName)
	}
	return result
}

func LoadModelFromConfig(ptr interface{}) {
	//t := reflect.TypeOf(ptr)
	//if t.Kind() == reflect.Ptr {
	//	t = t.Elem()
	//	for i:=0;i<t.NumField();i++ {
	//		reflect.ValueOf(ptr).Elem().Field(i).Set(reflect.ValueOf(beego.AppConfig.String(t.Field(i).Name)))
	//	}
	//}
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).Set(reflect.ValueOf(beego.AppConfig.String(v.Type().Field(i).Name)))
		case reflect.Int:
			intVal, _ := beego.AppConfig.Int(v.Type().Field(i).Name)
			v.Field(i).Set(reflect.ValueOf(intVal))
		case reflect.Slice:
			arrVal := beego.AppConfig.Strings(v.Type().Field(i).Name)
			//fmt.Println(v.Type().Field(i).Name,arrVal)
			v.Field(i).Set(reflect.ValueOf(arrVal))
		}
	}
}

func GetSnakeStrMapOfStruct(value interface{}) map[string]interface{} {
	var resMap = make(map[string]interface{})
	var sType = reflect.TypeOf(value)
	var sValue = reflect.ValueOf(value)
	var numsOfFields = sType.NumField()
	for i := 0; i < numsOfFields; i ++ {
		snakeName := SnakeString(sType.Field(i).Name)
		tagVal := sType.Field(i).Tag.Get("field")
		if strings.Contains(tagVal , "drop") {
			continue
		}
		resMap[snakeName] = sValue.Field(i).Interface()
		//switch sType.Field(i).Type.Name() {
		//case "int":
		//	resMap[snakeName] = strconv.FormatInt(sValue.Field(i).Int(), 10)
		//case "string":
		//	resMap[snakeName] = sValue.Field(i).String()
		//case "QueryStatus":
		//	resMap[snakeName] = strconv.FormatInt(sValue.Field(i).Int(), 10)
		//case "QueryEngine":
		//	resMap[snakeName] = strconv.FormatInt(sValue.Field(i).Int(), 10)
		//default:
		//	resMap[snakeName] = fmt.Sprintf("%v", sValue.Field(i))
		//}
	}
	return resMap
}

func GetSnakeStrMapOfStructForMuti(value interface{}, idx int) map[string]interface{} {
	var resMap = make(map[string]interface{})
	var sType = reflect.TypeOf(value)
	var sValue = reflect.ValueOf(value)
	var numsOfFields = sType.NumField()
	for i := 0; i < numsOfFields; i ++ {
		snakeName := SnakeString(sType.Field(i).Name)
		tagVal := strings.Split(sType.Field(i).Tag.Get("field"), ";")[idx]
		if strings.Contains(tagVal , "drop") {
			continue
		}
		resMap[snakeName] = sValue.Field(i).Interface()
		//switch sType.Field(i).Type.Name() {
		//case "int":
		//	resMap[snakeName] = strconv.FormatInt(sValue.Field(i).Int(), 10)
		//case "string":
		//	resMap[snakeName] = sValue.Field(i).String()
		//case "QueryStatus":
		//	resMap[snakeName] = strconv.FormatInt(sValue.Field(i).Int(), 10)
		//case "QueryEngine":
		//	resMap[snakeName] = strconv.FormatInt(sValue.Field(i).Int(), 10)
		//default:
		//	resMap[snakeName] = fmt.Sprintf("%v", sValue.Field(i))
		//}
	}
	return resMap
}



func Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

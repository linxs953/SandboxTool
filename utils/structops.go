package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
)

func GetTagName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Print("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		}
		result = append(result, tagName)
	}
	return result
}

func GetTagValue(structName interface{}) []string {
	v := reflect.ValueOf(structName)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		log.Print("Check type error not Struct")
		return nil
	}
	fieldNum := v.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		value := v.Field(i).Interface()
		result = append(result, fmt.Sprintf("%v", value))
	}
	return result
}

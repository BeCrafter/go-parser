package goparser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// 类型定义
const (
	TypeString = "string"
	TypeInt64  = "int64"
	TypeBool   = "bool"
	TypeFloat  = "float"
	TypeObject = "object"
)

// castType 基础类型转换，支持 string int64 bool float object 几种类型
func castType(data interface{}, typeName string) (interface{}, error) {
	if typeName == "" {
		return data, nil
	}
	switch strings.ToLower(typeName) {
	case TypeString:
		return castToString(data)
	case TypeInt64:
		return castToInt64(data)
	case TypeBool:
		return castToBoolean(data)
	case TypeFloat:
		return castToFloat(data)
	case TypeObject:
		return data, nil
	default:
		return nil, fmt.Errorf("type cast failure, unexpected type: %s", typeName)
	}
}

func castToBoolean(data interface{}) (bool, error) {
	if data == nil {
		return false, nil
	}
	switch t := data.(type) {
	case bool:
		return t, nil
	case string:
		if t == "" {
			return false, nil
		}
		return strconv.ParseBool(t)
	default:
		return false, fmt.Errorf("type cast failure, unexpected boolean value: %v", data)
	}
}

func castToString(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	return fmt.Sprint(data), nil
}

func castToInt64(data interface{}) (interface{}, error) {
	if data == nil {
		return 0, nil
	}

	switch t := data.(type) {
	case int:
		return int64(t), nil
	case float32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case json.Number:
		return t.Int64()
	case string:
		if t == "" {
			return 0, nil
		}
		return strconv.ParseInt(t, 10, 64)
	}
	return strconv.ParseInt(fmt.Sprint(data), 10, 64)
}

func castToFloat(data interface{}) (interface{}, error) {
	if data == nil {
		return 0, nil
	}

	switch t := data.(type) {
	case int:
		return float64(t), nil
	case int16:
		return float32(t), nil
	case int32:
		return float32(t), nil
	case int64:
		return float64(t), nil
	case uint:
		return float64(t), nil
	case uint16:
		return float32(t), nil
	case uint32:
		return float32(t), nil
	case uint64:
		return float64(t), nil
	case float32:
		return t, nil
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	case string:
		if t == "" {
			return 0., nil
		}
		return strconv.ParseFloat(t, 64)
	}
	return strconv.ParseFloat(fmt.Sprint(data), 64)
}

// StringBuilder 高效字符串拼接
func StringBuilder(p ...interface{}) string {
	var b strings.Builder
	l := len(p)
	for i := 0; i < l; i++ {
		switch v := p[i].(type) {
		case string:
			b.WriteString(v)
		case int:
			b.WriteString(strconv.FormatInt(int64(v), 10))
		case int8:
			b.WriteString(strconv.FormatInt(int64(v), 10))
		case int16:
			b.WriteString(strconv.FormatInt(int64(v), 10))
		case int32:
			b.WriteString(strconv.FormatInt(int64(v), 10))
		case int64:
			b.WriteString(strconv.FormatInt(int64(v), 10))
		case uint:
			b.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint8:
			b.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint16:
			b.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint32:
			b.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint64:
			b.WriteString(strconv.FormatUint(uint64(v), 10))
		case float32:
			b.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
		case float64:
			b.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 64))
		case json.Number:
			b.WriteString(string(v))
		case map[string]interface{}:
			b.WriteString(fmt.Sprintf("%v", v))
		case bool:
			b.WriteString(strconv.FormatBool(v))
		}
	}

	return b.String()
}

// InArray 判断数组中是否存在某元素
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

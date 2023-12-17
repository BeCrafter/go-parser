package goparser

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

// Jsoniter 别名
var Jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary

// ChildrenItem ...
type ChildrenItem struct {
	Op    string      `json:"op"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

// 逻辑运算符定义
const (
	LogicAnd = "&&" // 逻辑且
	LogicOr  = "||" // 逻辑或
	LogicNot = "!"  // 逻辑非

	LogicAndName = "AND" // 逻辑且标识
	LogicOrName  = "OR"  // 逻辑或标识
	LogicNotName = "NOT" // 逻辑非标识
)

// LogicMaps 逻辑运算符标识转换字典
var LogicMaps = map[string]string{
	LogicAndName: LogicAnd,
	LogicOrName:  LogicOr,
	LogicNotName: LogicNot,
}

// Expression 基于JSON生成表达式字符串
func Expression(exp map[string]interface{}) (string, error) {
	if len(exp) < 2 {
		return "", errors.New("invalid params")
	}

	child, ok := exp["children"]
	if !ok {
		return "", errors.New("invalid params filed[child]")
	}

	childs, ok2 := child.([]interface{})
	if !ok2 {
		return "", errors.New("childs convert array fail")
	}

	var con string = ""
	if v, has := LogicMaps[exp["connector"].(string)]; has {
		con = v
	}

	res := ""
	for _, item := range childs {
		val, ok := item.(map[string]interface{})
		if !ok {
			return "", errors.New("childs convert array fail")
		}

		ibyte, err := Jsoniter.Marshal(item)
		if err != nil {
			return "", err
		}

		if _, ok := val["children"]; ok {
			ss, err := Expression(val)
			if err != nil {
				return "", err
			}

			if con == LogicNot {
				res = LogicNot + ss
			} else {
				if len(res) == 0 {
					res = ss
				} else {
					res = StringBuilder("(", res, " ", con, " ", ss, ")")
				}
			}
		} else {
			var child ChildrenItem
			if err := json.Unmarshal(ibyte, &child); err != nil {
				return "", err
			}

			t := convert(child)
			if len(res) == 0 {
				res = t
			} else {
				res = StringBuilder("(", res, " ", con, " ", t, ")")
			}
		}
	}

	return res, nil
}

// ExportFields 导出表达中参数名
func ExportFields(exp map[string]interface{}) ([]string, error) {
	res := make([]string, 0)
	if len(exp) < 2 {
		return res, errors.New("invalid params")
	}

	child, ok := exp["children"]
	if !ok {
		return res, errors.New("invalid params filed[child]")
	}

	childs, ok2 := child.([]interface{})
	if !ok2 {
		return res, errors.New("childs convert array fail")
	}

	for _, item := range childs {
		val, ok := item.(map[string]interface{})
		if !ok {
			return res, errors.New("childs convert array fail")
		}

		ibyte, err := Jsoniter.Marshal(item)
		if err != nil {
			return res, err
		}

		if _, ok := val["children"]; ok {
			ss, err := ExportFields(val)
			if err != nil {
				return res, err
			}

			for _, v := range ss {
				if !InArray(v, res) {
					res = append(res, v)
				}
			}
		} else {
			var child ChildrenItem
			if err := json.Unmarshal(ibyte, &child); err != nil {
				return res, err
			}
			if !InArray(child.Field, res) {
				res = append(res, child.Field)
			}
		}
	}

	return res, nil
}

// convert 结构体转换为表达式
func convert(item ChildrenItem) string {
	opMap := map[string]string{
		"NE":  "!=",
		"EQ":  "==",
		"GT":  ">",
		"LT":  "<",
		"GE":  ">=",
		"LE":  "<=",
		"ADD": "+",
		"SUB": "-",
		"MUL": "*",
		"QUO": "/",
	}

	var val interface{}
	switch item.Value.(type) {
	case string:
		val = StringBuilder("\"", item.Value, "\"")
	default:
		val = item.Value
	}

	return StringBuilder(item.Field, " ", opMap[item.Op], " ", val)
}

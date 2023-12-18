package goparser

import (
	"errors"
	"go/ast"
)

// Func 生命自定义函数类型
type Func func(args []ast.Expr, data map[string]interface{}) interface{}

// 注册可执行函数
var funcNameMap = make(map[string]Func, 10)

// init 自定义函数初始化
func init() {
	// 注册内置函数
	RegisterFunc("in_array", inArray)
}

// RegisterFunc 注册自定义函数
func RegisterFunc(name string, f Func) {
	funcNameMap[name] = f
}

// inArray 判断变量是否存在在数组中
func inArray(args []ast.Expr, data map[string]interface{}) interface{} {
	// 规则表达式中的变量
	param := Eval(args[0], data)
	vRange, ok := args[1].(*ast.CompositeLit)
	if !ok {
		return errors.New("func in_array 2ed params is not a composite lit")
	}

	// 规则表达式中数组里的元素
	eltNodes := make([]interface{}, 0, len(vRange.Elts))
	for _, p := range vRange.Elts {
		elt := Eval(p, data)
		eltNodes = append(eltNodes, elt)
	}

	for _, node := range eltNodes {
		switch node.(type) {
		case int64:
			param, err := castType(param, TypeInt64)
			if err != nil {
				return false
			}
			paramInt64, paramOk := param.(int64)
			nodeInt64, nodeOk := node.(int64)
			if !paramOk || !nodeOk {
				return false
			}
			if nodeInt64 == paramInt64 {
				return true
			}
		case string:
			param, err := castType(param, TypeString)
			if err != nil {
				return false
			}
			nodeString, nodeOk := node.(string)
			paramString, paramOk := param.(string)
			if !paramOk || !nodeOk {
				return false
			}
			if nodeString == paramString {
				return true
			}
		}
	}
	return false
}

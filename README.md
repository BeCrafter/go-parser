# goparser

基于 golang 原生语法解析器（[parser](https://pkg.go.dev/go/parser)）实现的轻量级规则引擎。支持操作：

- 规则匹配：`goparser.Match(ruleStr, params)`


## 如何使用

### 安装

```bash
go get -u github.com/BeCrafter/go-parser@master
```

### 使用

#### 表达式解析
```go
import "github.com/BeCrafter/go-parser"

ruleStr := "!(a == 1 && b == 2 && c == "test" && d == false)"

// 匹配变量
params := map[string]interface{}{
    "a": 1,
    "b": 2,
    "c": "test",
    "d": true,
}

result, err := goparser.Match(ruleStr, params)
fmt.Println(result)
```

#### 表达式生成
```go
import "github.com/BeCrafter/go-parser"

str = `{
    "connector": "NOT",
    "children": [
        {
            "connector": "AND",
            "children": [
                {
                    "op": "GE",
                    "value": 43,
                    "field": "age"
                },
                {
                    "op": "EQ",
                    "value": "haha",
                    "field": "name"
                }
            ]
        }
    ]
}`

var tmp map[string]interface{}
json.Unmarshal([]byte(str), &tmp)
if val, err := Expression(tmp); err != nil {
    fmt.Errorf("goParser expression failed, err=%v", err)
}else{
    fmt.Println("Expression:", val)
}
```

#### 注册自定义函数

```go

import "github.com/BeCrafter/go-parser"

// MaxDemo 自定义Max函数
func MaxDemo(args []ast.Expr, data map[string]interface{}) interface{} {
    num1 := goparser.Eval(args[0], data)
    num2 := goparser.Eval(args[1], data)
    if num1 > num2 {
        return num1
    }

    return num2
}

// 在库中注册 `max` 函数
goparser.RegisterFunc("max", MaxDemo)

```

### 其他说明

#### 支持类型

- int
- int64
- string
- bool

#### 支持操作

- `!表达式`：支持一元表达式
- `&&`：支持多个表达式逻辑与
- `||`：支持多个表达式逻辑或
- `()`：支持表达式括号包裹
- `==`：int、int64、string、bool支持
- `!=`：int、int64、string、bool支持
- `>`：int、int64支持
- `<`：int、int64支持
- `>=`：int、int64支持
- `<=`：int、int64支持
- `+`：int、int64支持
- `-`：int、int64支持
- `*`：int、int64支持
- `/`：int、int64支持
- `%`：int、int64支持

#### 性能对比

```bash
Benchmarkgoparser_Match-8            127189             8912   ns/op     // goparser
BenchmarkGval_Match-8                 63584            18358   ns/op     // gval
BenchmarkGovaluateParser_Match-8      13628            86955   ns/op     // govaluate
BenchmarkYqlParser_Match-8            10364           112481   ns/op     // yql
```
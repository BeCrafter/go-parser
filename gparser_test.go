package goparser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGoParser_Match(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[string]interface{}
		want bool
	}{
		{
			name: "test_case1",
			expr: "a == 1 && b == 2",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			want: true,
		},
		{
			name: "test_case2",
			expr: "a == 1 && b == 2",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": 3,
			},
			want: false,
		},
		{
			name: "test_case3",
			expr: "a == 1 && b == 2 || c == \"test\"",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": "test",
			},
			want: true,
		},
		{
			name: "test_case4",
			expr: "a == 1 && b == 2 && c == \"test\"",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": "test",
			},
			want: false,
		},
		{
			name: "test_case5",
			expr: "a == 1 && b == 2 && c == \"test\" && d == true",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: true,
		},
		{
			name: "test_case6",
			expr: "a == 1 && b == 2 && c == \"test\" && d == false",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: false,
		},
		{
			name: "test_case7",
			expr: "!(a == 1 && b == 2 && c == \"test\" && d == false)",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: true,
		},
		{
			name: "test_case8",
			expr: "!(a == 1 && b == 2) || (c == \"test\" && d == false)",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": false,
			},
			want: true,
		},
		{
			name: "test_case9",
			expr: "a == 1 && b == 2",
			data: nil,
			want: false,
		},
		{
			name: "test_case10",
			expr: "",
			data: nil,
			want: true,
		},
		{
			name: "test_case11",
			expr: "1 == 1",
			data: nil,
			want: false,
		},
		{
			name: "test_case12",
			expr: "",
			data: make(map[string]interface{}, 0),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Match(tt.expr, tt.data); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("goParser match failed, want=%v, got=%v, err=%v", tt.want, got, err)
			} else {
				fmt.Printf("Ouput: %+v\n", got)
			}
		})
	}
}

func BenchmarkGoParser_Match(b *testing.B) {
	// 规则表达式
	expr := `(a == 1 && b == "b" && in_array(c, []int{100,99,98,97})) || (d == false)`
	// 映射数据
	data := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": 100,
		"d": true,
	}
	for i := 0; i < b.N; i++ {
		if _, err := Match(expr, data); err != nil {
			fmt.Printf("goParser BenchmarkGoParser Match failed, err=%v", err)
		}
	}
}

func TestGoParser_Expression(t *testing.T) {
	tests := []string{
		`{
			"connector": "",
			"children": [
				{
					"op": "GE",
					"value": 43,
					"field": "age"
				}
			]
		}`,
		`{
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
		}`,
		`{
			"connector": "NOT",
			"children": [
				{
					"connector": "AND",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"op": "GT",
									"value": 43,
									"field": "age"
								},
								{
									"op": "EQ",
									"value": "haha",
									"field": "name"
								}
							]
						},
						{
							"connector": "OR",
							"children": [
								{
									"op": "GT",
									"value": 40,
									"field": "age"
								},
								{
									"op": "EQ",
									"value": "haha1",
									"field": "name"
								}
							]
						}
					]
				}
			]
		}`,
		`{
			"connector": "NOT",
			"children": [
				{
					"connector": "OR",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"op": "GT",
									"value": 43,
									"field": "age"
								},
								{
									"op": "EQ",
									"value": "haha",
									"field": "name"
								}
							]
						},
						{
							"connector": "",
							"children": [
								{
									"op": "GT",
									"value": 40,
									"field": "age"
								}
							]
						}
					]
				}
			]
		}`,
		`{
			"connector": "AND",
			"children": [
				{
					"connector": "",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						}
					]
				},
				{
					"connector": "AND",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						},
						{
							"op": "EQ",
							"value": "haha",
							"field": "name"
						}
					]
				},
				{
					"connector": "OR",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						},
						{
							"op": "EQ",
							"value": "haha",
							"field": "name"
						}
					]
				},
				{
					"connector": "NOT",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						}
					]
				},
				{
					"connector": "NOT",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"op": "GT",
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
				},
				{
					"connector": "NOT",
					"children": [
						{
							"connector": "OR",
							"children": [
								{
									"connector": "AND",
									"children": [
										{
											"op": "GT",
											"value": 43,
											"field": "age"
										},
										{
											"op": "EQ",
											"value": "haha",
											"field": "name"
										}
									]
								},
								{
									"connector": "",
									"children": [
										{
											"op": "GT",
											"value": 40,
											"field": "age"
										}
									]
								}
							]
						}
					]
				},
				{
					"connector": "NOT",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"connector": "AND",
									"children": [
										{
											"op": "GT",
											"value": 43,
											"field": "age"
										},
										{
											"op": "EQ",
											"value": "haha",
											"field": "name"
										}
									]
								},
								{
									"connector": "OR",
									"children": [
										{
											"op": "GT",
											"value": 40,
											"field": "age"
										},
										{
											"op": "EQ",
											"value": "haha1",
											"field": "name"
										}
									]
								}
							]
						}
					]
				}
			]
		}`,
	}

	for index, tt := range tests {
		tindex, _ := castToString(index)
		name := "test_case_" + tindex
		t.Run(name, func(t *testing.T) {
			var tmp map[string]interface{}
			if err := json.Unmarshal([]byte(tt), &tmp); err != nil {
				t.Errorf("goParser expression json unmarshal failed, err=%v", err)
			}
			if _, err := Expression(tmp); err != nil {
				t.Errorf("goParser expression failed, err=%v", err)
			}
		})
	}
}

func BenchmarkGoParser_Expression(b *testing.B) {
	// 字符串
	str := `{
		"connector": "NOT",
		"children": [
			{
				"connector": "AND",
				"children": [
					{
						"connector": "AND",
						"children": [
							{
								"op": "GT",
								"value": 43,
								"field": "age"
							},
							{
								"op": "EQ",
								"value": "haha",
								"field": "name"
							}
						]
					},
					{
						"connector": "OR",
						"children": [
							{
								"op": "GT",
								"value": 40,
								"field": "age"
							},
							{
								"op": "EQ",
								"value": "haha1",
								"field": "name"
							}
						]
					}
				]
			}
		]
	}`

	var tmp map[string]interface{}
	if err := json.Unmarshal([]byte(str), &tmp); err != nil {
		fmt.Printf("json unmarshal fail,err:%+v", err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := Expression(tmp); err != nil {
			fmt.Printf("BenchmarkGoParser_Expression fail,err:%+v", err)
		}
	}
}

func TestGoParser_ExportFields(t *testing.T) {
	tests := []string{
		`{
			"connector": "",
			"children": [
				{
					"op": "GE",
					"value": 43,
					"field": "age"
				}
			]
		}`,
		`{
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
		}`,
		`{
			"connector": "NOT",
			"children": [
				{
					"connector": "AND",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"op": "GT",
									"value": 43,
									"field": "age"
								},
								{
									"op": "EQ",
									"value": "haha",
									"field": "name"
								}
							]
						},
						{
							"connector": "OR",
							"children": [
								{
									"op": "GT",
									"value": 40,
									"field": "age"
								},
								{
									"op": "EQ",
									"value": "haha1",
									"field": "name"
								}
							]
						}
					]
				}
			]
		}`,
		`{
			"connector": "NOT",
			"children": [
				{
					"connector": "OR",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"op": "GT",
									"value": 43,
									"field": "age"
								},
								{
									"op": "EQ",
									"value": "haha",
									"field": "name"
								}
							]
						},
						{
							"connector": "",
							"children": [
								{
									"op": "GT",
									"value": 40,
									"field": "age"
								}
							]
						}
					]
				}
			]
		}`,
		`{
			"connector": "AND",
			"children": [
				{
					"connector": "",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						}
					]
				},
				{
					"connector": "AND",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						},
						{
							"op": "EQ",
							"value": "haha",
							"field": "name"
						}
					]
				},
				{
					"connector": "OR",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						},
						{
							"op": "EQ",
							"value": "haha",
							"field": "name"
						}
					]
				},
				{
					"connector": "NOT",
					"children": [
						{
							"op": "GT",
							"value": 43,
							"field": "age"
						}
					]
				},
				{
					"connector": "NOT",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"op": "GT",
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
				},
				{
					"connector": "NOT",
					"children": [
						{
							"connector": "OR",
							"children": [
								{
									"connector": "AND",
									"children": [
										{
											"op": "GT",
											"value": 43,
											"field": "age"
										},
										{
											"op": "EQ",
											"value": "haha",
											"field": "name"
										}
									]
								},
								{
									"connector": "",
									"children": [
										{
											"op": "GT",
											"value": 40,
											"field": "age"
										}
									]
								}
							]
						}
					]
				},
				{
					"connector": "NOT",
					"children": [
						{
							"connector": "AND",
							"children": [
								{
									"connector": "AND",
									"children": [
										{
											"op": "GT",
											"value": 43,
											"field": "age"
										},
										{
											"op": "EQ",
											"value": "haha",
											"field": "name"
										}
									]
								},
								{
									"connector": "OR",
									"children": [
										{
											"op": "GT",
											"value": 40,
											"field": "age"
										},
										{
											"op": "EQ",
											"value": "haha1",
											"field": "name"
										}
									]
								}
							]
						}
					]
				}
			]
		}`,
	}

	for index, tt := range tests {
		tindex, _ := castToString(index)
		name := "test_export_case_" + tindex
		t.Run(name, func(t *testing.T) {
			var tmp map[string]interface{}
			if err := json.Unmarshal([]byte(tt), &tmp); err != nil {
				t.Errorf("goParser expression json unmarshal failed, err=%v", err)
			}
			if v, err := ExportFields(tmp); err != nil {
				t.Errorf("goParser expression failed, err=%v", err)
			} else {
				fmt.Printf("Output: %+v\n", v)
			}
		})
	}
}

func BenchmarkGoParser_ExportFields(b *testing.B) {
	// 字符串
	str := `{
		"connector": "NOT",
		"children": [
			{
				"connector": "AND",
				"children": [
					{
						"connector": "AND",
						"children": [
							{
								"op": "GT",
								"value": 43,
								"field": "age"
							},
							{
								"op": "EQ",
								"value": "haha",
								"field": "name"
							}
						]
					},
					{
						"connector": "OR",
						"children": [
							{
								"op": "GT",
								"value": 40,
								"field": "age"
							},
							{
								"op": "EQ",
								"value": "haha1",
								"field": "name"
							}
						]
					}
				]
			}
		]
	}`

	var tmp map[string]interface{}
	if err := json.Unmarshal([]byte(str), &tmp); err != nil {
		fmt.Printf("json unmarshal fail,err:%+v", err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := ExportFields(tmp); err != nil {
			fmt.Printf("BenchmarkGoParser_ExportFields fail,err:%+v", err)
		}
	}
}

func TestGoParser_RunAll(t *testing.T) {
	// 表达式：!((age > 43 && name == "haha") && (age > 40 || name == "haha1"))
	str := `{
		"connector": "NOT",
		"children": [
			{
				"connector": "AND",
				"children": [
					{
						"connector": "AND",
						"children": [
							{
								"op": "GT",
								"value": 43,
								"field": "age"
							},
							{
								"op": "EQ",
								"value": "haha",
								"field": "name"
							}
						]
					},
					{
						"connector": "OR",
						"children": [
							{
								"op": "GT",
								"value": 40,
								"field": "age"
							},
							{
								"op": "EQ",
								"value": "haha1",
								"field": "name"
							}
						]
					}
				]
			}
		]
	}`

	var tmp map[string]interface{}
	if err := json.Unmarshal([]byte(str), &tmp); err != nil {
		t.Errorf("json unmarshal fail,err:%+v", err)
	}

	exp, err := Expression(tmp)
	if err != nil {
		t.Errorf("Expression fail,err:%+v", err)
	}

	fmt.Printf("\nOutput Expression:\n\n\t%s\n\n", exp)

	// 表达式：!((age > 43 && name == "haha") && (age > 40 || name == "haha1"))
	tests := []struct {
		name string
		expr string
		data map[string]interface{}
		want bool
	}{
		// 预期返回假
		{
			name: "test_case1",
			expr: exp,
			data: map[string]interface{}{
				"name": "haha",
				"age":  44,
			},
			want: false,
		},
		// 预期返回真
		{
			name: "test_case2",
			expr: exp,
			data: map[string]interface{}{
				"name": "haha2",
				"age":  44,
			},
			want: true,
		},
		// 参数值为nil,但不能传空map
		{
			name: "test_case3",
			expr: exp,
			data: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Match(tt.expr, tt.data); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("goParser match failed, want=%v, got=%v, err=%v", tt.want, got, err)
			} else {
				fmt.Printf("Ouput: %+v\n", got)
			}
		})
	}
}

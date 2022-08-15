// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package expressions

import (
	//"base/docs"
	//"datavalues"
	"github.com/CC11001100/vectorsql/src/base/docs"
	"github.com/CC11001100/vectorsql/src/datavalues"
)

// ConstantExpression 常量表达式，返回给定的值
type ConstantExpression struct {

	// 返回这个固定值，这个值需要是数据库中兼容的数据库类型
	value datavalues.IDataValue
}

// CONST 创建一个常量表达式
func CONST(v interface{}) IExpression {
	return NewConstantExpression(datavalues.ToValue(v))
}

func NewConstantExpression(v datavalues.IDataValue) *ConstantExpression {
	return &ConstantExpression{
		value: v,
	}
}

// Eval 常量不能被执行
func (e *ConstantExpression) Eval() error {
	return nil
}

// Update 常量没有参数
func (e *ConstantExpression) Update(params IParams) (datavalues.IDataValue, error) {
	return e.value, nil
}

// Merge 妈的常量值竟然能被覆盖掉？
func (e *ConstantExpression) Merge(arg IExpression) (datavalues.IDataValue, error) {
	other := arg.(*ConstantExpression)
	return other.value, nil
}

// Result 返回最终的常量
func (e *ConstantExpression) Result() datavalues.IDataValue {
	return e.value
}

// Walk 遍历表达式
func (e *ConstantExpression) Walk(visit Visit) error {
	return nil
}

func (e *ConstantExpression) String() string {
	return string(e.value.String())
}

func (e *ConstantExpression) Document() docs.Documentation {
	return nil
}

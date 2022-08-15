// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package planners

import (
	"encoding/json"
)

// 用来组合两个Plan 
type BinaryExpressionPlan struct {
	Name     string

	// 这里认为是函数才会有二元操作，所以需要一个函数的名字 
	FuncName string

	// 第一个操作单元 
	Left     IPlan

	// 
	Right    IPlan
}

func NewBinaryExpressionPlan(funcName string, left IPlan, right IPlan) *BinaryExpressionPlan {
	return &BinaryExpressionPlan{
		Name:     "BinaryExpressionPlan",
		FuncName: funcName,
		Left:     left,
		Right:    right,
	}
}

func (plan *BinaryExpressionPlan) Build() error {
	return nil
}

func (plan *BinaryExpressionPlan) Walk(visit Visit) error {
	return Walk(visit, plan.Left, plan.Right)
}

func (plan *BinaryExpressionPlan) String() string {
	out, err := json.MarshalIndent(plan, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(out)
}

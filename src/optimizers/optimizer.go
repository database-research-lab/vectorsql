// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package optimizers

import (
	"github.com/CC11001100/vectorsql/src/planners"
)

// 优化器的定义 
type Optimizer struct {

	// 优化器的名称 
	Name        string

	// 优化器的作用描述，但是似乎也没用到？ 
	Description string

	// 优化器的优化方法，对传入的执行计划做优化 
	Reassembler func(planners.IPlan)
}

// 默认使用的优化器  
var DefaultOptimizers = []Optimizer{

	ProjectPushDownOptimizer,

	PredicatePushDownOptimizer,

}

// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"base/errors"
	"planners"
)

// 查询执行器 
type SelectExecutor struct {

	// 执行器上下文 
	ctx  *ExecutorContext

	// 查询所对应的执行计划 
	plan *planners.SelectPlan

	// 执行树？是啥东西？ 
	tree *ExecutorTree
}

// 创建一个查询执行器  
func NewSelectExecutor(ctx *ExecutorContext, plan planners.IPlan) IExecutor {
	return &SelectExecutor{
		ctx:  ctx,
		// 从上下文创建一棵树？这是一颗什么样子的树呢？  
		tree: NewExecutorTree(ctx),
		plan: plan.(*planners.SelectPlan),
	}
}

// 真正开始执行select查询了 
func (executor *SelectExecutor) Execute() (*Result, error) {
	ectx := executor.ctx
	tree := executor.tree

	// select查询计划涉及到的子查询计划
	children := executor.plan.SubPlan.SubPlans

	// 把子查询计划都转换为执行器挂到执行器树上  
	for _, plan := range children {
		switch plan := plan.(type) {
		case *planners.TableValuedFunctionPlan:
			executor := NewTableValuedFunctionExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.ScanPlan:
			executor := NewScanExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.FilterPlan:
			executor := NewFilterExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.SelectionPlan:
			executor := NewSelectionExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.OrderByPlan:
			executor := NewOrderByExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.LimitPlan:
			executor := NewLimitExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.ProjectionPlan:
			executor := NewProjectionExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.SinkPlan:
			executor := NewSinkExecutor(ectx, plan)
			tree.Add(executor)
		default:
			return nil, errors.Errorf("Unsupported plan:%T", plan)
		}
	}
	pipeline, err := tree.BuildPipeline()
	if err != nil {
		return nil, err
	}
	pipeline.Run()

	result := NewResult()
	result.SetInput(pipeline.Last())
	return result, nil
}

func (executor *SelectExecutor) String() string {
	res := ""
	for _, t := range executor.tree.subExecutors {
		res += t.String()
		res += " -> "
	}
	return res
}

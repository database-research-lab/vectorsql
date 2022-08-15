// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"github.com/CC11001100/vectorsql/src/processors"
)

// 是为了把查询抽象成树方便执行吧？
type ExecutorTree struct {

	// 执行上下文 
	ctx          *ExecutorContext

	// 子执行器 
	subExecutors []IExecutor
}

// 从执行上下文创建执行树 
func NewExecutorTree(ctx *ExecutorContext) *ExecutorTree {
	return &ExecutorTree{
		ctx: ctx,
	}
}

// 为当前查询增加一个子查询 
func (tree *ExecutorTree) Add(executor IExecutor) {
	tree.subExecutors = append(tree.subExecutors, executor)
}

// 要把父查询和的结果和子查询的输入连接起来 
func (tree *ExecutorTree) BuildPipeline() (*processors.Pipeline, error) {
	ectx := tree.ctx

	// TODO 这里没太看明白具体是怎么连接起来的？ 
	pipeline := processors.NewPipeline(ectx.ctx)
	for _, executor := range tree.subExecutors {
		transform, err := executor.Execute()
		if err != nil {
			return nil, err
		}
		pipeline.Add(transform.In)
	}
	return pipeline, nil
}

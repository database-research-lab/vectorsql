// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"github.com/CC11001100/vectorsql/src/databases"
	"github.com/CC11001100/vectorsql/src/planners"
)

// 创建数据库的执行器 
type CreateDatabaseExecutor struct {
	ctx  *ExecutorContext
	plan *planners.CreateDatabasePlan
}

func NewCreateDatabaseExecutor(ctx *ExecutorContext, plan planners.IPlan) IExecutor {
	return &CreateDatabaseExecutor{
		ctx:  ctx,
		plan: plan.(*planners.CreateDatabasePlan),
	}
}

// 在执行器中根据传递的参数创建数据库 
func (executor *CreateDatabaseExecutor) Execute() (*Result, error) {
	ectx := executor.ctx
	ast := executor.plan.Ast

	databaseCtx := databases.NewDatabaseContext(ectx.log, ectx.conf)
	database, err := databases.DatabaseFactory(databaseCtx, ast)
	if err != nil {
		return nil, err
	}
	if err := database.Executor().CreateDatabase(); err != nil {
		return nil, err
	}

	result := NewResult()
	return result, nil
}

func (executor *CreateDatabaseExecutor) String() string {
	return ""
}

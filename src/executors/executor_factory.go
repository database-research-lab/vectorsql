// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"reflect"

	"planners"

	"base/errors"
)

// 执行器的入口  

// 为plan创建执行器的方法  
type executorCreator func(ctx *ExecutorContext, plan planners.IPlan) IExecutor

// 根据不同的执行计划有不同的执行器  
var table = map[reflect.Type]executorCreator{

	// use xxx; 切换数据库 
	reflect.TypeOf(&planners.UsePlan{}):            NewUseExecutor,

	// select 查询数据 
	reflect.TypeOf(&planners.SelectPlan{}):         NewSelectExecutor,

	// create database; 创建数据库 
	reflect.TypeOf(&planners.CreateDatabasePlan{}): NewCreateDatabaseExecutor,

	// 删除数据库 
	reflect.TypeOf(&planners.DropDatabasePlan{}):   NewDropDatabaseExecutor,

	// 创建表 
	reflect.TypeOf(&planners.CreateTablePlan{}):    NewCreateTableExecutor,

	// 删除表 
	reflect.TypeOf(&planners.DropTablePlan{}):      NewDropTableExecutor,

	// 查看数据库 
	reflect.TypeOf(&planners.ShowDatabasesPlan{}):  NewShowDatabasesExecutor,

	// 查看表 
	reflect.TypeOf(&planners.ShowTablesPlan{}):     NewShowTablesExecutor,

	// 插入数据 
	reflect.TypeOf(&planners.InsertPlan{}):         NewInsertExecutor,
}

func ExecutorFactory(ctx *ExecutorContext, plan planners.IPlan) (IExecutor, error) {
	creator, ok := table[reflect.TypeOf(plan)]
	if !ok {
		return nil, errors.Errorf("Couldn't get the executor:%T", plan)
	}
	return creator(ctx, plan), nil
}

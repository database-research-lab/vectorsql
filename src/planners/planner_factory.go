// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package planners

import (
	"github.com/CC11001100/vectorsql/src/base/errors"
	"github.com/CC11001100/vectorsql/src/parsers"
	"github.com/CC11001100/vectorsql/src/parsers/sqlparser"
)

// planner的入口，用于为SQL查询生成执行计划  

// 创建执行计划的方法 
type planCreator func(ast sqlparser.Statement) IPlan

// 不同的方法会走不同的生成Plan的方法 
var table = map[string]planCreator{

	// use xxx; 切换数据库的时候 
	sqlparser.NodeNameUse:            NewUsePlan,

	// select查询语句 
	sqlparser.NodeNameSelect:         NewSelectPlan,

	// create database 创建数据库 
	sqlparser.NodeNameDatabaseCreate: NewCreateDatabasePlan,

	// drop database 删除数据库 
	sqlparser.NodeNameDatabaseDrop:   NewDropDatabasePlan,

	// 创建表 
	sqlparser.NodeNameTableCreate:    NewCreateTablePlan,

	// 删除表 
	sqlparser.NodeNameTableDrop:      NewDropTablePlan,

	// 查看所有数据库 
	sqlparser.NodeNameShowDatabases:  NewShowDatabasesPlan,

	// 查看所有表 
	sqlparser.NodeNameShowTables:     NewShowTablesPlan,

	// 插入数据  
	sqlparser.NodeNameInsert:         NewInsertPlan,
}

// 
func PlanFactory(query string) (IPlan, error) {
	// 先把SQL语句解析为AST树 
	statement, err := parsers.Parse(query)
	if err != nil {
		return nil, err
	}

	// 根据不同的语句类型走不同的查询计划生成  
	creator, ok := table[statement.Name()]
	if !ok {
		return nil, errors.Errorf("Couldn't get the planner:%T", statement)
	}
	plan := creator(statement)
	return plan, plan.Build()
}

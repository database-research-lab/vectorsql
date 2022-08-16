// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datablocks

import (
	"github.com/CC11001100/vectorsql/src/columns"
	"github.com/CC11001100/vectorsql/src/datatypes"
	"github.com/CC11001100/vectorsql/src/planners"
)

// ProjectionByPlan 根据投影执行计划裁剪列
func (block *DataBlock) ProjectionByPlan(plan *planners.MapPlan) (*DataBlock, error) {
	projects := plan

	// 把投影中涉及到的列取出来
	// Build the project exprs.
	exprs, err := planners.BuildExpressions(projects)
	if err != nil {
		return nil, err
	}

	rows := block.NumRows()
	if rows == 0 {
		// If empty, returns header only.
		cols := make([]*columns.Column, len(exprs))
		for i, expr := range exprs {
			cols[i] = columns.NewColumn(expr.String(), datatypes.NewStringDataType())
		}
		return NewDataBlock(cols), nil
	} else {
		columnValues := make([]*DataBlockValue, len(exprs))
		for i, expr := range exprs {
			columnValue, err := block.DataBlockValue(expr.String())
			if err != nil {
				return nil, err
			}
			columnValues[i] = columnValue
		}
		return newDataBlock(block.seqs, columnValues), nil
	}
}

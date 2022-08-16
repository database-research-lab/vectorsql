// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datablocks

import (
	"github.com/CC11001100/vectorsql/src/columns"
	"github.com/CC11001100/vectorsql/src/datavalues"
)

type ColumnIndex struct {

	// 列名
	Name  string

	// 列顺序下表
	Index int
}

// DataBlockValue 用于表示数据块中的一个值
type DataBlockValue struct {

	// 这个值所对应的列信息
	column *columns.Column

	// 值对应的值...我草好绕
	values []datavalues.IDataValue
}

func NewDataBlockValue(col *columns.Column) *DataBlockValue {
	return &DataBlockValue{
		column: col,
		values: make([]datavalues.IDataValue, 0),
	}
}

func newDataBlockValueWithValues(col *columns.Column, values []datavalues.IDataValue) *DataBlockValue {
	return &DataBlockValue{
		column: col,
		values: values,
	}
}

func (v *DataBlockValue) ColumnName() string {
	return v.column.Name
}

func (v *DataBlockValue) DeepClone() *DataBlockValue {
	clone := &DataBlockValue{
		column: v.column,
		values: make([]datavalues.IDataValue, len(v.values)),
	}
	copy(clone.values, v.values)
	return clone
}

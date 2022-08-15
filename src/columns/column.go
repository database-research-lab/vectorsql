// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package columns

import (
	//"datatypes"
	"github.com/CC11001100/vectorsql/src/datatypes"
)

// Column 表示数据库中的一个列
type Column struct {

	// 列的名字
	Name string

	// 列的数据类型
	DataType datatypes.IDataType
}

func NewColumn(name string, datatype datatypes.IDataType) *Column {
	return &Column{
		Name:     name,
		DataType: datatype,
	}
}

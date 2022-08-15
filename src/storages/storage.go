// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package storages

import (
	"github.com/CC11001100/vectorsql/src/columns"
	"github.com/CC11001100/vectorsql/src/datastreams"
	"github.com/CC11001100/vectorsql/src/sessions"
)

// IStorage 存储层抽象定义
type IStorage interface {

	// Name 使用的存储介质的名字
	Name() string

	// Columns 列都有哪些
	Columns() []*columns.Column

	// GetInputStream 输入流
	GetInputStream(*sessions.Session) (datastreams.IDataBlockInputStream, error)

	// GetOutputStream 输出流
	GetOutputStream(*sessions.Session) (datastreams.IDataBlockOutputStream, error)

	// Close 关闭存储介质
	Close()
}

// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datastreams

import (
	//"datablocks"
	"github.com/CC11001100/vectorsql/src/datablocks"
)

// IDataBlockInputStream 块输入流，用于读取记录序列化之后的数据
type IDataBlockInputStream interface {
	Name() string

	// Read next block.
	// If there are no more blocks, return nil.
	Read() (*datablocks.DataBlock, error)

	Close()
}

// IDataBlockOutputStream 块输出流，用于保存记录序列化之后的数据
type IDataBlockOutputStream interface {
	Name() string
	Write(*datablocks.DataBlock) error
	Finalize() error
	Close()
	SampleBlock() *datablocks.DataBlock
}

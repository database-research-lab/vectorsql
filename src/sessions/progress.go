// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package sessions

import (
	//"base/sync2"
	"github.com/CC11001100/vectorsql/src/base/sync2"
)

// ProgressValues 一些统计信息，为了成本优化的时候使用
type ProgressValues struct {

	// 花费了多长时间  
	Cost            sync2.AtomicDuration

	// 读取了多少行 
	ReadRows        sync2.AtomicInt64

	// 读取了多少个字节 
	ReadBytes       sync2.AtomicInt64

	// 总共读取了多少行 
	TotalRowsToRead sync2.AtomicInt64

	// 写了多少行 
	WrittenRows     sync2.AtomicInt64

	// 写了多少个字节 
	WrittenBytes    sync2.AtomicInt64
}

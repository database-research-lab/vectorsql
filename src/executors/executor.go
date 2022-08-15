// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"github.com/CC11001100/vectorsql/src/datastreams"
	"github.com/CC11001100/vectorsql/src/processors"
)

// IExecutor 定义了执行器方法
type IExecutor interface {

	// String 执行器会有一个名字
	String() string

	// Execute 封装真正执行的逻辑
	Execute() (*Result, error)
}

// Result 表示执行器的执行结果
type Result struct {
	In  processors.IProcessor
	Out datastreams.IDataBlockOutputStream
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) SetInput(in processors.IProcessor) {
	r.In = in
}

func (r *Result) SetOutput(out datastreams.IDataBlockOutputStream) {
	r.Out = out
}

func (r *Result) Read() <-chan interface{} {
	return r.In.In().Recv()
}

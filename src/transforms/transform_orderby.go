// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package transforms

import (
	"github.com/CC11001100/vectorsql/src/datablocks"
	"github.com/CC11001100/vectorsql/src/planners"
	"github.com/CC11001100/vectorsql/src/processors"
	"github.com/CC11001100/vectorsql/src/sessions"
	"time"
)

// OrderByTransform 排序处理器
type OrderByTransform struct {
	ctx *TransformContext

	// order by 对应的执行计划
	plan *planners.OrderByPlan

	progressValues sessions.ProgressValues
	processors.BaseProcessor
}

func NewOrderByTransform(ctx *TransformContext, plan *planners.OrderByPlan) processors.IProcessor {
	return &OrderByTransform{
		ctx:           ctx,
		plan:          plan,
		BaseProcessor: processors.NewBaseProcessor("transform_orderby"),
	}
}

func (t *OrderByTransform) Execute() {
	var block *datablocks.DataBlock

	plan := t.plan
	out := t.Out()
	defer out.Close()

	// 获取order by的字段，如果出错了就中断执行万事大吉
	// Get all base fields by the expression.
	fields, err := planners.BuildVariableValues(plan)
	if err != nil {
		out.Send(err)
		return
	}

	// 先持有着每一条输出结果
	onNext := func(x interface{}) {
		switch y := x.(type) {
		case *datablocks.DataBlock:
			if block == nil {
				block = y
			} else {
				if err := block.Append(y); err != nil {
					out.Send(err)
				}
			}
		case error:
			out.Send(y)
		}
	}

	// 在所有的结果读取完毕之后，进行排序，截取需要的结果
	onDone := func() {
		if block != nil {
			start := time.Now()
			if err := block.OrderByPlan(fields, t.plan); err != nil {
				out.Send(err)
			} else {
				cost := time.Since(start)
				t.progressValues.Cost.Add(cost)
				t.progressValues.ReadBytes.Add(int64(block.TotalBytes()))
				t.progressValues.ReadRows.Add(int64(block.NumRows()))
				t.progressValues.TotalRowsToRead.Add(int64(block.NumRows()))
				out.Send(block)
			}
		}
	}
	t.Subscribe(onNext, onDone)
}

func (t *OrderByTransform) Stats() sessions.ProgressValues {
	return t.progressValues
}

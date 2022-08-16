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

// ProjectionTransform 选择列的处理器，用于收缩结果的规模
type ProjectionTransform struct {
	ctx *TransformContext

	// 投影所对应的执行计划
	plan *planners.ProjectionPlan

	progressValues sessions.ProgressValues
	processors.BaseProcessor
}

func NewProjectionTransform(ctx *TransformContext, plan *planners.ProjectionPlan) processors.IProcessor {
	return &ProjectionTransform{
		ctx:           ctx,
		plan:          plan,
		BaseProcessor: processors.NewBaseProcessor("transform_projection"),
	}
}

func (t *ProjectionTransform) Execute() {

	out := t.Out()
	defer out.Close()

	// 对接收到的每条结果收缩列规模
	onNext := func(x interface{}) {
		switch y := x.(type) {
		case *datablocks.DataBlock:
			start := time.Now()
			// 啊草咋都是在Block里实现的啊，总感觉这个职责划分不太对...
			if block, err := y.ProjectionByPlan(t.plan.Projections); err != nil {
				out.Send(err)
			} else {
				cost := time.Since(start)
				t.progressValues.Cost.Add(cost)
				t.progressValues.ReadBytes.Add(int64(block.TotalBytes()))
				t.progressValues.ReadRows.Add(int64(block.NumRows()))
				t.progressValues.TotalRowsToRead.Add(int64(block.NumRows()))
				out.Send(block)
			}
		default:
			// 啊哈哈哈如果是自己理解不了的，就直接透传给后面的处理器
			out.Send(x)
		}
	}
	t.Subscribe(onNext)
}

func (t *ProjectionTransform) Stats() sessions.ProgressValues {
	return t.progressValues
}

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

// LimitTransform 名字修改为正常的了，不要黑客文化了...
// 限制返回的结果条数
type LimitTransform struct {
	ctx *TransformContext

	// limit的执行计划
	plan *planners.LimitPlan

	// 统计相关的一些信息，可以忽略掉...
	progressValues sessions.ProgressValues
	processors.BaseProcessor
}

func NewLimitTransform(ctx *TransformContext, plan *planners.LimitPlan) processors.IProcessor {
	return &LimitTransform{
		ctx:           ctx,
		plan:          plan,
		BaseProcessor: processors.NewBaseProcessor("transform_limit"),
	}
}

func (t *LimitTransform) Execute() {
	var (
		limit  int
		offset int
	)

	//Todo support eval(variable)
	offset = t.plan.OffsetPlan.(*planners.ConstantPlan).Value.(int)
	limit = t.plan.RowcountPlan.(*planners.ConstantPlan).Value.(int)

	out := t.Out()
	defer out.Close()

	onNext := func(x interface{}) {
		switch y := x.(type) {
		case *datablocks.DataBlock:
			if x != nil {
				start := time.Now()
				cutOffset, cutLimit := y.Limit(offset, limit)
				offset -= cutOffset
				limit -= cutLimit
				x = y

				cost := time.Since(start)
				t.progressValues.Cost.Add(cost)
				t.progressValues.ReadBytes.Add(int64(y.TotalBytes()))
				t.progressValues.ReadRows.Add(int64(y.NumRows()))
				t.progressValues.TotalRowsToRead.Add(int64(y.NumRows()))
			}
		}
		out.Send(x)

		if offset < 0 || limit <= 0 {
			out.Close()
			return
		}
	}
	t.Subscribe(onNext)
}

func (t *LimitTransform) Stats() sessions.ProgressValues {
	return t.progressValues
}

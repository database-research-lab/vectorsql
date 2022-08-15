package optimizers

import (
	"github.com/CC11001100/vectorsql/src/planners"
)

// 优化器的入口

func Optimize(plan planners.IPlan, optimizers []Optimizer) planners.IPlan {
	for _, opt := range optimizers {
		opt.Reassembler(plan)
	}
	return plan
}

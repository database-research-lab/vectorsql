// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package optimizers

import (
	"github.com/CC11001100/vectorsql/src/planners"
)

// ProjectPushDownOptimizer 映射下推，不读取多余的列
var ProjectPushDownOptimizer = Optimizer{
	Name:        "ProjectPushDownOptimizer",
	Description: "Push projects to scan plan",
	Reassembler: func(plan planners.IPlan) {
		var scan *planners.ScanPlan
		var project *planners.ProjectionPlan

		visit := func(plan planners.IPlan) (kontinue bool, err error) {
			switch plan := plan.(type) {
			case *planners.ScanPlan:
				scan = plan
			case *planners.ProjectionPlan:
				project = plan
			}
			return true, nil
		}
		if err := planners.Walk(visit, plan); err != nil {
			return
		}

		// 歪日这么实现真的正确吗... 
		if scan != nil && project != nil {
			scan.Project = project
		}
	},
}

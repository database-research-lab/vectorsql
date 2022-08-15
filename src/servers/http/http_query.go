// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package http

import (
	"context"
	"io"

	"datablocks"
	"datastreams"
	"executors"
	"optimizers"
	"planners"
	"processors"
	"sessions"
)

func (s *HTTPHandler) processQuery(query string, rw io.Writer) (err error) {
	log := s.log
	conf := s.conf
	session := sessions.NewSession()
	defer session.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Debug("HTTPHandler-Query->Enter:%+v", query)

	// 为查询生成逻辑执行计划 
	// Logical plans.
	plan, err := planners.PlanFactory(query)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	// 优化生成的执行计划  
	plan = optimizers.Optimize(plan, optimizers.DefaultOptimizers)

	// 真正开始执行  
	// Executors.
	ectx := executors.NewExecutorContext(ctx, log, conf, session)
	executor, err := executors.ExecutorFactory(ectx, plan)
	if err != nil {
		log.Error("%+v", err)
		return err
	}

	result, err := executor.Execute()
	if err != nil {
		log.Error("%+v", err)
		return err
	}

	if err = s.processOrdinaryQuery(rw, session, result.In); err != nil {
		return
	}
	log.Debug("%v", executor.String())
	return nil
}

func (s *HTTPHandler) processOrdinaryQuery(rw io.Writer, session *sessions.Session, sink processors.IProcessor) error {
	log := s.log

	log.Debug("HTTPHandler->OrdinaryQuery->Enter")
	if sink != nil {
		for x := range sink.In().Recv() {
			switch x := x.(type) {
			case error:
				log.Error("%+v", x)
				return x
			case *datablocks.DataBlock:
				log.Debug("HTTPHandler->OrdinaryQuery->DataBlock: rows:%+v", x.NumRows())
				if err := s.sendData(rw, x); err != nil {
					return err
				}
			}
		}
	}
	log.Debug("HTTPHandler->OrdinaryQuery->Return")
	return nil
}

func (s *HTTPHandler) sendData(writer io.Writer, block *datablocks.DataBlock) error {
	// TODO Get the format from AST
	output := datastreams.NewCustomFormatBlockOutputStream(block, writer, "TSV")

	if err := output.Write(block); err != nil {
		return err
	}
	if err := output.Finalize(); err != nil {
		return err
	}
	return nil
}

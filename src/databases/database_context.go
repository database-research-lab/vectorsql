// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package databases

import (
	"github.com/CC11001100/vectorsql/src/base/xlog"
	"github.com/CC11001100/vectorsql/src/config"
)

type DatabaseContext struct {
	log  *xlog.Log
	conf *config.Config
}

func NewDatabaseContext(log *xlog.Log, conf *config.Config) *DatabaseContext {
	return &DatabaseContext{
		log:  log,
		conf: conf,
	}
}

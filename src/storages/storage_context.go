// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package storages

import (
	"github.com/CC11001100/vectorsql/src/base/xlog"
	"github.com/CC11001100/vectorsql/src/config"
	"github.com/CC11001100/vectorsql/src/storages/system"
)

type StorageContext struct {
	log               *xlog.Log
	conf              *config.Config
	tablesFillFunc    system.TablesFillFunc
	databasesFillFunc system.DatabasesFillFunc
}

func NewStorageContext(log *xlog.Log, conf *config.Config) *StorageContext {
	return &StorageContext{
		log:  log,
		conf: conf,
	}
}

func (ctx *StorageContext) SetTablesFillFunc(fn system.TablesFillFunc) {
	ctx.tablesFillFunc = fn
}

func (ctx *StorageContext) SetDatabasesFillFunc(fn system.DatabasesFillFunc) {
	ctx.databasesFillFunc = fn
}

// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package storages

import (
	"github.com/CC11001100/vectorsql/src/columns"
	mem "github.com/CC11001100/vectorsql/src/storages/memory"
)

const (
	MemoryStorageEngineName = "MEMORY"
)

func NewMemoryStorage(ctx *StorageContext, cols []*columns.Column) IStorage {
	mctx := mem.NewMemoryStorageContext(ctx.log, ctx.conf)
	return mem.NewMemoryStorage(mctx, cols)
}

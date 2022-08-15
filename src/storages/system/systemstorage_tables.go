// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package system

import (
	"github.com/CC11001100/vectorsql/src/base/errors"
	"github.com/CC11001100/vectorsql/src/columns"
	"github.com/CC11001100/vectorsql/src/datablocks"
	"github.com/CC11001100/vectorsql/src/datastreams"
	"github.com/CC11001100/vectorsql/src/datatypes"
	"github.com/CC11001100/vectorsql/src/sessions"
)

type SystemTablesStorage struct {
	ctx *SystemStorageContext
}

func NewSystemTablesStorage(ctx *SystemStorageContext) *SystemTablesStorage {
	return &SystemTablesStorage{
		ctx: ctx,
	}
}

func (storage *SystemTablesStorage) Name() string {
	return ""
}

func (storage *SystemTablesStorage) Columns() []*columns.Column {
	return []*columns.Column{
		{Name: "name", DataType: datatypes.NewStringDataType()},
		{Name: "database", DataType: datatypes.NewStringDataType()},
		{Name: "engine", DataType: datatypes.NewStringDataType()},
	}
}

func (storage *SystemTablesStorage) GetOutputStream(session *sessions.Session) (datastreams.IDataBlockOutputStream, error) {
	return nil, errors.New("Couldn't find outputstream")
}

func (storage *SystemTablesStorage) GetInputStream(session *sessions.Session) (datastreams.IDataBlockInputStream, error) {
	ctx := storage.ctx

	// Block.
	block := datablocks.NewDataBlock(storage.Columns())
	if err := ctx.tablesFillFunc(block); err != nil {
		return nil, err
	}

	// Stream.
	return datastreams.NewOneBlockInputStream(block), nil
}

func (storage *SystemTablesStorage) Close() {
}

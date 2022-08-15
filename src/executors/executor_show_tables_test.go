// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"github.com/CC11001100/vectorsql/src/columns"
	"github.com/CC11001100/vectorsql/src/datablocks"
	"github.com/CC11001100/vectorsql/src/datatypes"
	"github.com/CC11001100/vectorsql/src/mocks"
	"github.com/CC11001100/vectorsql/src/planners"
	"testing"



	"github.com/stretchr/testify/assert"
)

func TestShowTablessExecutor(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		err    string
		expect *datablocks.DataBlock
	}{

		{
			name:  "show tables",
			query: "show tables where `engine` like '%SYSTEM_%' and name like '%tab%' limit 2",
			expect: mocks.NewBlockFromSlice(
				[]*columns.Column{
					{Name: "name", DataType: datatypes.NewStringDataType()},
				},
				[]interface{}{"databases"},
				[]interface{}{"tables"},
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := mocks.NewMock()
			defer cleanup()

			plan, err := planners.PlanFactory(test.query)
			assert.Nil(t, err)

			ctx := NewExecutorContext(mock.Ctx, mock.Log, mock.Conf, mock.Session)
			executor, err := ExecutorFactory(ctx, plan)
			assert.Nil(t, err)

			result, err := executor.Execute()
			if test.err != "" {
				assert.Equal(t, test.err, err.Error())
			} else {
				assert.Nil(t, err)
				for x := range result.In.In().Recv() {
					expect := test.expect
					actual := x.(*datablocks.DataBlock)
					assert.True(t, mocks.DataBlockEqual(expect, actual))
				}
			}
		})
	}
}

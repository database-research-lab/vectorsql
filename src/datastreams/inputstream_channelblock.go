// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datastreams

import (
	"github.com/CC11001100/vectorsql/src/datablocks"
)

// ChannelBlockInputStream 输入流的channel实现，数据是从channel中读取的
type ChannelBlockInputStream struct {
	queue chan interface{}
}

func NewChannelBlockInputStream(queue chan interface{}) IDataBlockInputStream {
	return &ChannelBlockInputStream{
		queue: queue,
	}
}

func (stream *ChannelBlockInputStream) Name() string {
	return "ChannelBlockInputStream"
}

func (stream *ChannelBlockInputStream) Read() (*datablocks.DataBlock, error) {
	// 每次从channel中读取一条数据
	val, ok := <-stream.queue
	if ok {
		// channel中的数据类型支持error和DataBlock两种，其它的任何类型会被认为是nil丢弃掉
		switch t := val.(type) {
		// 可以把error写到channel中，在这个地方会兼容读出来
		case error:
			return nil, t
		case *datablocks.DataBlock:
			// 可以把数据块写入到channel中，这个地方会读出来
			return t, nil
		}
	}
	return nil, nil
}

func (stream *ChannelBlockInputStream) Close() {}

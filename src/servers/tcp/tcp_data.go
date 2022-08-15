// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package tcp

import (
	"github.com/CC11001100/vectorsql/src/base/humanize"
	"github.com/CC11001100/vectorsql/src/datastreams"
)

func (s *TCPHandler) processData(session *TCPSession) error {
	log := s.log

	stream := datastreams.NewNativeBlockInputStream(session.reader)
	defer stream.Close()

	block, err := stream.Read()
	if err != nil {
		return err
	}
	if block != nil {
		log.Debug("Receive client data block: rows:%v, columns:%v, size:%v", block.NumRows(), block.NumColumns(), humanize.Bytes(block.TotalBytes()))
		if !s.state.Empty() {
			return s.state.result.Out.Write(block)
		}
	} else {
		log.Debug("Receive nil client data block")
	}
	return nil
}

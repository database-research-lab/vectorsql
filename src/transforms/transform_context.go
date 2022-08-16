// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package transforms

import (
	"context"
	"github.com/CC11001100/vectorsql/src/base/xlog"
	"github.com/CC11001100/vectorsql/src/config"
	"github.com/CC11001100/vectorsql/src/sessions"
)

// TransformContext 转换器的上下文，似乎也没啥需要特别关注的东西...
type TransformContext struct {
	ctx              context.Context
	log              *xlog.Log
	conf             *config.Config
	progressCallback func(values *sessions.ProgressValues)
}

func NewTransformContext(ctx context.Context, log *xlog.Log, conf *config.Config) *TransformContext {
	return &TransformContext{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}
}

func (ctx *TransformContext) SetProgressCallback(fn func(pv *sessions.ProgressValues)) {
	ctx.progressCallback = fn
}

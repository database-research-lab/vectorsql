// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package processors

import (
	"context"
	"time"
)

// IProcessor 处理器？干嘛的处理啊？？？我好懵...
type IProcessor interface {

	// Name 处理器的名称
	Name() string

	// Pause TODO 暂停？啥意思？先不处理了吗？
	Pause()

	// Resume 恢复？听起来像是跟上面那个方法对应的...
	Resume()

	// Execute 当前处理器的执行逻辑，不同的处理器的主要区别就是在这里
	Execute()

	// In 返回当前处理器的输入流
	In() *InPort

	// Out 返回当前处理器的输出流
	Out() *OutPort

	// Duration 当前处理器持续执行了多长时间了
	Duration() time.Duration

	To(...IProcessor)

	From(...IProcessor)

	SetContext(context.Context)
}

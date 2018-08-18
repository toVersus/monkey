package vm

import (
	"github.com/toversus/monkey/code"
	"github.com/toversus/monkey/object"
)

// Frame holds execution relevant information.
type Frame struct {
	fn *object.CompiledFunction // compiled function referenced by the frame
	ip int                      // instruction pointer in this frame for this function.
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}

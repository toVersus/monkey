package vm

import (
	"github.com/toversus/monkey/code"
	"github.com/toversus/monkey/object"
)

// Frame holds execution relevant information.
type Frame struct {
	fn          *object.CompiledFunction // compiled function referenced by the frame.
	ip          int                      // instruction pointer in this frame for this function.
	basePointer int                      // pointer that points to the bottom of the stack of the current call frame.
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{
		fn:          fn,
		ip:          -1,
		basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}

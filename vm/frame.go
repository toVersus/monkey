package vm

import (
	"github.com/toversus/monkey/code"
	"github.com/toversus/monkey/object"
)

// Frame holds execution relevant information.
type Frame struct {
	cl          *object.Closure // cl encloses the Fn field
	ip          int             // instruction pointer in this frame for this function.
	basePointer int             // pointer that points to the bottom of the stack of the current call frame.
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}

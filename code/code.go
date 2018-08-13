package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Definition represents to share definition of opcode between compiler and VM.
// Name helps to make an Opcode human-readable and OperandWidths contains
// the number of bytes each operand takes up.
type Definition struct {
	Name          string
	OperandWidths []int
}

// Instructions represents first bytecode instruction, which consists of
// an opcode and an optional number of operands.
type Instructions []byte

// String is used to show nicely-formatted multi-line output that tells us
// the opcodes in human-readable form.
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Opcode has an arbitary but unique value and is the first byte in the instruction.
type Opcode byte

const (
	OpConstant Opcode = iota // the constant using the operand as an index and push it on to the stack.

	OpAdd // '+'

	OpPop // pop the topmost element off the stack.

	OpSub // '-'
	OpMul // '*'
	OpDiv // '/'

	OpTrue  // 'true'
	OpFalse // 'false'

	OpEqual       // '='
	OpNotEqual    // '!='
	OpGreaterThan // '>', no OpLessThan. OpLessThan is generated reordering of code.

	OpMinus // prefix expression '-'
	OpBang  // prefix expression '!'

	OpJumpNotTruthy // jump if conditional is not false nor null.
	OpJump          // jump to the index of instruction by using its operand

	OpNull // representation of 'nothing'

	OpGetGlobal // get binding for global variables
	OpSetGlobal // set binding for global variables
)

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // one operand, the number previously assigned to the constant.

	OpAdd: {"OpAdd", []int{}},

	OpPop: {"OpPop", []int{}},

	OpSub: {"OpSub", []int{}},
	OpMul: {"OpMul", []int{}},
	OpDiv: {"OpDiv", []int{}},

	OpTrue:  {"OpTrue", []int{}},
	OpFalse: {"OpFalse", []int{}},

	OpEqual:       {"OpEqual", []int{}},
	OpNotEqual:    {"OpNotEqual", []int{}},
	OpGreaterThan: {"OpGreaterThan", []int{}},

	OpMinus: {"OpMinus", []int{}},
	OpBang:  {"OpBang", []int{}},

	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},

	OpNull: {"OpNull", []int{}},

	OpGetGlobal: {"OpGetGlobal", []int{2}},
	OpSetGlobal: {"OpSetGlobal", []int{2}},
}

// Lookup gets to the definition of opcode.
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Make can create single bytecode instruction that's made up of an Opcode
// and an optional number of operands.
func Make(op Opcode, operands ...int) []byte {
	// Lookup function isn't used here due to avoiding error check
	// whenever building up bytecode instructions.
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

// ReadOperands decode operands of a bytecode instruction.
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

// ReadUint16 skips the definition lookup required by ReadOperands
// and is used directly by the VM.
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

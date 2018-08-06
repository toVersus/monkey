package code

import (
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

// Opcode has an arbitary but unique value and is the first byte in the instruction.
type Opcode byte

const (
	// OpConstant has one operand, the number previously assigned to the constant.
	// VM will retrieve the constant using the operand as an index
	// and push it on to the stack.
	OpConstant Opcode = iota
)

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
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

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
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
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

	OpArray // takes N elements off the stack

	OpHash // takes N keys and N values off the stack

	OpIndex // takes the objects to be indexed and serving as the index off the stack.

	OpCall        // function call expression
	OpReturnValue // implicit and explicit return statement
	OpReturn      // return statement which has no explicit return value

	OpGetLocal // get binding for local variables
	OpSetLocal // set binding for local variables

	OpGetBuiltin // get binding for builtin functions

	OpClosure // send a message to wrap the specified compiled function in an closure

	OpGetFree // get binding for free variables
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

	OpArray: {"OpArray", []int{2}},

	OpHash: {"OpHash", []int{2}},

	OpIndex: {"OpIndex", []int{}},

	OpCall:        {"OpCall", []int{1}},
	OpReturnValue: {"OpReturnValue", []int{}},
	OpReturn:      {"OpReturn", []int{}},

	OpGetLocal: {"OpGetLocal", []int{1}},
	OpSetLocal: {"OpSetLocal", []int{1}},

	OpGetBuiltin: {"OpGetBuiltin", []int{1}},

	OpClosure: {"OpClosure", []int{2, 1}}, // the constant index and the count of free variables

	OpGetFree: {"OpGetFree", []int{1}},
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

		case 1:
			instruction[offset] = byte(o)
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

		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

// ReadUint8 reads one byte and turns it into an uint8.
func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }

// ReadUint16 skips the definition lookup required by ReadOperands
// and is used directly by the VM.
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

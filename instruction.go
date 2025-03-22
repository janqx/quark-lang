package quark

import "fmt"

type Opcode uint8

type Operand uint32

/* 0~23位是operand, 24~32位是opcode */
type Instruction uint32

const (
	OpNop Opcode = iota

	OpLoadNull
	OpLoadTrue
	OpLoadFalse
	OpLoadConst
	OpLoadLocal
	OpLoadOuter
	OpLoadGlobal
	OpLoadIndex
	OpLoadAttribute

	OpStoreLocal
	OpStoreOuter
	OpStoreGlobal
	OpStoreIndex
	OpStoreAttribute

	OpUnaryBitNot
	OpUnaryNot
	OpUnaryPlus
	OpUnaryMinus

	OpBinaryAdd
	OpBinarySub
	OpBinaryMul
	OpBinaryDiv
	OpBinaryMod
	OpBinaryLT
	OpBinaryLTE
	OpBinaryGT
	OpBinaryGTE
	OpBinaryEQ
	OpBinaryNEQ
	OpBinaryBitAnd
	OpBinaryBitOr
	OpBinaryBitXor
	OpBinaryBitLhs
	OpBinaryBitRhs

	OpJump
	OpJumpIfFalse
	OpJumpIfFalseOrPop
	OpJumpIfTrueOrPop

	OpClosure
	OpCall
	OpReturn
	OpRemoveTop

	OpBuildList
	OpBuildDict

	OpImport
	OpExport

	OpDebugger
)

var OpcodeToString = [...]string{

	OpNop: "OpNop",

	OpLoadNull:      "OpLoadNull",
	OpLoadTrue:      "OpLoadTrue",
	OpLoadFalse:     "OpLoadFalse",
	OpLoadConst:     "OpLoadConst",
	OpLoadLocal:     "OpLoadLocal",
	OpLoadOuter:     "OpLoadOuter",
	OpLoadGlobal:    "OpLoadGlobal",
	OpLoadIndex:     "OpLoadIndex",
	OpLoadAttribute: "OpLoadAttribute",

	OpStoreLocal:     "OpStoreLocal",
	OpStoreOuter:     "OpStoreOuter",
	OpStoreGlobal:    "OpStoreGlobal",
	OpStoreIndex:     "OpStoreIndex",
	OpStoreAttribute: "OpStoreAttribute",

	OpUnaryBitNot: "OpUnaryBitNot",
	OpUnaryNot:    "OpUnaryNot",
	OpUnaryPlus:   "OpUnaryPlus",
	OpUnaryMinus:  "OpUnaryMinus",

	OpBinaryAdd:    "OpBinaryAdd",
	OpBinarySub:    "OpBinarySub",
	OpBinaryMul:    "OpBinaryMul",
	OpBinaryDiv:    "OpBinaryDiv",
	OpBinaryMod:    "OpBinaryMod",
	OpBinaryLT:     "OpBinaryLT",
	OpBinaryLTE:    "OpBinaryLTE",
	OpBinaryGT:     "OpBinaryGT",
	OpBinaryGTE:    "OpBinaryGTE",
	OpBinaryEQ:     "OpBinaryEQ",
	OpBinaryNEQ:    "OpBinaryNEQ",
	OpBinaryBitAnd: "OpBinaryBitAnd",
	OpBinaryBitOr:  "OpBinaryBitOr",
	OpBinaryBitXor: "OpBinaryBitXor",
	OpBinaryBitLhs: "OpBinaryBitLhs",
	OpBinaryBitRhs: "OpBinaryBitRhs",

	OpJump:             "OpJump",
	OpJumpIfFalse:      "OpJumpIfFalse",
	OpJumpIfFalseOrPop: "OpJumpIfFalseOrPop",
	OpJumpIfTrueOrPop:  "OpJumpIfTrueOrPop",

	OpClosure:   "OpClosure",
	OpCall:      "OpCall",
	OpReturn:    "OpReturn",
	OpRemoveTop: "OpRemoveTop",

	OpBuildList: "OpBuildList",
	OpBuildDict: "OpBuildDict",

	OpImport: "OpImport",
	OpExport: "OpExport",

	OpDebugger: "OpDebugger",
}

func (op Opcode) String() string {
	return OpcodeToString[op]
}

var InvalidOperand Operand = 0x00ffffff

func (operand Operand) isValid() bool {
	return operand < InvalidOperand
}

func (i Instruction) Opcode() Opcode {
	return Opcode(i >> 24)
}

func (i Instruction) Operand() Operand {
	return Operand(i & 0x00ffffff)
}

func (i Instruction) String() string {
	result := i.Opcode().String()
	if i.Operand().isValid() {
		result += fmt.Sprintf(" %d", i.Operand())
	}
	return result
}

func NewInstruction(opcode Opcode, operand Operand) Instruction {
	return Instruction(uint32(opcode)<<24) | Instruction(operand&0x00ffffff)
}

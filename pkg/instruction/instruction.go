package instruction

import "fmt"

//go:generate stringer -type=OpcodeClass,OpcodeType,OpcodeArithmetic,OpcodeJump,OpcodeMode,AtomicOperation,AtomicModifier,ImmSource,OpcodeSize,Register -linecomment -output=stringer.go

type Opcode uint8

func (ins Instruction) Opcode() Opcode {
	return Opcode(ins.Basic & 0xFF00_0000_0000_0000 >> 56)
}

type OpcodeCode interface {
	String() string
}

func (o Opcode) Code() OpcodeCode {
	typedOpcode := o.ToTyped()
	return typedOpcode.CodeOrMode()
}

type OpcodeClass uint8

const (
	// non-standard load operations
	BPF_LD OpcodeClass = 0x00
	// load into register operations
	BPF_LDX OpcodeClass = 0x01
	// store from immediate operations
	BPF_ST OpcodeClass = 0x02
	// store from register operations
	BPF_STX OpcodeClass = 0x03
	// 32-bit arithmetic operations
	BPF_ALU OpcodeClass = 0x04
	// 64-bit jump operations
	BPF_JMP OpcodeClass = 0x05
	// 32-bit jump operations
	BPF_JMP32 OpcodeClass = 0x06
	// 64-bit arithmetic operations
	BPF_ALU64 OpcodeClass = 0x07
)

func (o Opcode) Class() OpcodeClass {
	// the three LSB bits of the opcode field store the instruction class
	return OpcodeClass(o & 0b0111)
}

type OpcodeType uint8

const (
	LOAD_AND_STORE      OpcodeType = 0x0
	ARITHMETIC_AND_JUMP OpcodeType = 0x1
)

// Type is whether an opcode is a load and store instruction or an arithmetic
// and jump instruction
func (c OpcodeClass) Type() OpcodeType {
	switch c {
	case BPF_LD, BPF_LDX, BPF_ST, BPF_STX:
		return LOAD_AND_STORE
	case BPF_ALU, BPF_JMP, BPF_JMP32, BPF_ALU64:
		return ARITHMETIC_AND_JUMP
	default:
		panic("invalid opcode type")
	}
}

type TypedOpcode interface {
	CodeOrMode() OpcodeCode
}

func (o Opcode) ToTyped() TypedOpcode {
	switch o.Class() {
	case BPF_ALU, BPF_ALU64:
		return ArithmeticOpcode(o)
	case BPF_JMP, BPF_JMP32:
		return JumpOpcode(o)
	case BPF_LD, BPF_LDX, BPF_ST, BPF_STX:
		return LoadAndStoreOpcode(o)
	default:
		panic("invalid opcode type")
	}
}

type ArithmeticAndJumpOpcode interface {
	Code() OpcodeCode
	Source() OpcodeSource
}

type ArithmeticOpcode uint8

func (o ArithmeticOpcode) CodeOrMode() OpcodeCode {
	return o.Code()
}

type JumpOpcode uint8

func (o JumpOpcode) CodeOrMode() OpcodeCode {
	return o.Code()
}

type OpcodeSource uint8

const (
	// use 32-bit 'imm' value as source operand
	BPF_K OpcodeSource = 0x00
	// use 'src_reg' register value as source operand
	BPF_X OpcodeSource = 0x08
)

func (o JumpOpcode) Source() OpcodeSource {
	return OpcodeSource(o & 0b1000)
}

func (o ArithmeticOpcode) Source() OpcodeSource {
	return OpcodeSource(o & 0b1000)
}

type OpcodeArithmetic uint8

const (
	// dst += src
	BPF_ADD OpcodeArithmetic = 0x00
	// dst -= src
	BPF_SUB OpcodeArithmetic = 0x10
	// dst *= src
	BPF_MUL OpcodeArithmetic = 0x20
	// dst = (src != 0) ? (dst / src) : 0
	BPF_DIV OpcodeArithmetic = 0x30
	// dst |= src
	BPF_OR OpcodeArithmetic = 0x40
	// dst &= src
	BPF_AND OpcodeArithmetic = 0x50
	// dst <<= (src & mask)
	BPF_LSH OpcodeArithmetic = 0x60
	// dst >>= (src & mask)
	BPF_RSH OpcodeArithmetic = 0x70
	// dst = ~src
	BPF_NEG OpcodeArithmetic = 0x80
	// dst = (src != 0) ? (dst % src) : dst
	BPF_MOD OpcodeArithmetic = 0x90
	// dst ^= src
	BPF_XOR OpcodeArithmetic = 0xa0
	// dst = src
	BPF_MOV OpcodeArithmetic = 0xb0
	// sign extending dst >>= (src & mask)
	BPF_ARSH OpcodeArithmetic = 0xc0
	// byte swap operations
	BPF_END OpcodeArithmetic = 0xd0
)

func (o ArithmeticOpcode) Code() OpcodeArithmetic {
	return OpcodeArithmetic(o & 0b11110000)
}

type OpcodeJump uint8

const (
	// PC += offset
	// src = 0x0
	// BPF_JMP only
	BPF_JA OpcodeJump = 0x0
	// PC += offset if dst == src
	// any
	BPF_JEQ OpcodeJump = 0x1
	// PC += offset if dst > src
	// any
	// unsigned
	BPF_JGT OpcodeJump = 0x2
	// PC += offset if dst >= src
	// any
	// unsigned
	BPF_JGE OpcodeJump = 0x3
	// PC += offset if dst & src
	// any
	BPF_JSET OpcodeJump = 0x4
	// PC += offset if dst != src
	// any
	BPF_JNE OpcodeJump = 0x5
	// PC += offset if dst > src
	// any
	// signed
	BPF_JSGT OpcodeJump = 0x6
	// PC += offset if dst >= src
	// any
	// signed
	BPF_JSGE OpcodeJump = 0x7
	// call helper function by address
	// 0x0
	// see Helper functions
	BPF_CALL OpcodeJump = 0x8
	// call PC += offset
	// 0x1
	// see Program-local functions
	// BPF_CALL OpcodeJump = 0x8
	// call helper function by BTF ID
	// 0x2
	// see Helper functions
	// BPF_CALL OpcodeJump = 0x8
	// return
	// 0x0
	// BPF_JMP only
	BPF_EXIT OpcodeJump = 0x9
	// PC += offset if dst < src
	// any
	// unsigned
	BPF_JLT OpcodeJump = 0xa
	// PC += offset if dst <= src
	// any
	// unsigned
	BPF_JLE OpcodeJump = 0xb
	// PC += offset if dst < src
	// any
	// signed
	BPF_JSLT OpcodeJump = 0xc
	// PC += offset if dst <= src
	// any
	// signed
	BPF_JSLE OpcodeJump = 0xd
)

func (o JumpOpcode) Code() OpcodeJump {
	return OpcodeJump(o >> 4)
}

type LoadAndStoreOpcode uint8

func (o LoadAndStoreOpcode) CodeOrMode() OpcodeCode {
	return o.Mode()
}

func (ins Instruction) ImmSrc() ImmSource {
	return ImmSource(ins.Regs().SrcReg() & 0xF0 >> 4)
}

type OpcodeSize uint8

const (
	// word (4 bytes)
	BPF_W OpcodeSize = 0x00 // u32
	// half word (2 bytes)
	BPF_H OpcodeSize = 0x08 // u16
	// byte
	BPF_B OpcodeSize = 0x10 // u8
	// double word (8 bytes)
	BPF_DW OpcodeSize = 0x18 // u64
)

func (o LoadAndStoreOpcode) Size() OpcodeSize {
	return OpcodeSize(o & 0b11000)
}

type OpcodeMode uint8

const (
	// 64-bit immediate instructions
	BPF_IMM OpcodeMode = 0x00
	// legacy BPF packet access (absolute)
	BPF_ABS OpcodeMode = 0x20
	// legacy BPF packet access (indirect)
	BPF_IND OpcodeMode = 0x40
	// regular load and store operations
	BPF_MEM OpcodeMode = 0x60
	// sign-extension load operations
	BPF_MEMSX OpcodeMode = 0x80
	// atomic operations
	BPF_ATOMIC OpcodeMode = 0xc0
)

type AtomicModifier uint8

const (
	// modifier: return old value
	// The BPF_FETCH modifier is optional for simple atomic operations, and
	// always set for the complex atomic operations. If the BPF_FETCH flag is
	// set, then the operation also overwrites src with the value that was in
	// memory before it was modified.
	BPF_FETCH AtomicModifier = 0x01
)

type AtomicOperation OpcodeArithmetic

const (
	// atomic exchange
	// The BPF_XCHG operation atomically exchanges src with the value addressed
	// by dst + offset.
	BPF_XCHG AtomicOperation = 0xe0 | AtomicOperation(BPF_FETCH)
	// atomic compare and exchange
	// The BPF_CMPXCHG operation atomically compares the value addressed by dst
	// + offset with R0. If they match, the value addressed by dst + offset is
	// replaced with src. In either case, the value that was at dst + offset
	// before the operation is zero-extended and loaded back to R0.
	BPF_CMPXCHG AtomicOperation = 0xf0 | AtomicOperation(BPF_FETCH)
)

func (ins Instruction) AtomicOperationImm() AtomicOperation {
	return AtomicOperation(ins.Imm())
}

func (o LoadAndStoreOpcode) Mode() OpcodeMode {
	return OpcodeMode(o & 0b11100000)
}

type Regs uint8

// Regs is composed of the source and the destination register numbers
func (ins Instruction) Regs() Regs {
	return Regs(ins.Basic & 0x00FF_0000_0000_0000 >> 48)
}

type Register uint8

const (
	// return value from function calls, and exit value for eBPF programs
	BPF_R0 Register = iota // r0
	// argument 1 for function calls
	BPF_R1 // r1
	// argument 2 for function calls
	BPF_R2 // r2
	// argument 3 for function calls
	BPF_R3 // r3
	// argument 4 for function calls
	BPF_R4 // r4
	// argument 5 for function calls
	BPF_R5 // r5
	// callee saved register 1 that function calls will preserve
	BPF_R6 // r6
	// callee saved register 2 that function calls will preserve
	BPF_R7 // r7
	// callee saved register 3 that function calls will preserve
	BPF_R8 // r8
	// callee saved register 4 that function calls will preserve
	BPF_R9 // r9
	// read-only frame pointer to access stack
	BPF_R10 // r10
)

// SrcReg is the source register number (0-10), except where otherwise specified.
func (r Regs) SrcReg() Register {
	return Register(r & 0xF0 >> 4)
}

// DstReg is the destination register number (0-10)
func (r Regs) DstReg() Register {
	return Register(r & 0x0F)
}

type ImmSource uint8

const (
	// dst = imm64
	// imm type: integer
	// dst type: integer
	BPF_IMM0 ImmSource = 0x0
	// dst = map_by_fd(imm)
	// imm type: map fd
	// dst type: map
	BPF_IMM1 ImmSource = 0x1
	// dst = map_val(map_by_fd(imm)) + next_imm
	// imm type: map fd
	// dst type: data pointer
	BPF_IMM2 ImmSource = 0x2
	// dst = var_addr(imm)
	// imm type: variable id
	// dst type: data pointer
	BPF_IMM3 ImmSource = 0x3
	// dst = code_addr(imm)
	// imm type: integer
	// dst type: code pointer
	BPF_IMM4 ImmSource = 0x4
	// dst = map_by_idx(imm)
	// imm type: map index
	// dst type: map
	BPF_IMM5 ImmSource = 0x5
	// dst = map_val(map_by_idx(imm)) + next_imm
	// imm type: map index
	// dst type: data pointer
	BPF_IMM6 ImmSource = 0x6
)

type Offset int16

// Offset is the signed integer offset used with pointer arithmetic
func (ins Instruction) Offset() Offset {
	// it's little endian
	msb := ins.Basic & 0x0000_FF00_0000_0000 >> 40
	lsb := ins.Basic & 0x0000_00FF_0000_0000 >> 24
	return Offset(msb | lsb)
}

type Imm int32

// Imm is the signed integer immediate value
func (ins Instruction) Imm() Imm {
	// it's little endian
	b1 := ins.Basic & 0x0000_0000_FF00_0000 >> 24
	b2 := ins.Basic & 0x0000_0000_00FF_0000 >> 8
	b3 := ins.Basic & 0x0000_0000_0000_FF00 << 8
	b4 := ins.Basic & 0x0000_0000_0000_00FF << 24
	return Imm(b1 | b2 | b3 | b4)
}

func (ins Instruction) NextImm() Imm {
	// it's little endian
	b1 := ins.Pseudo & 0x0000_0000_FF00_0000 >> 24
	b2 := ins.Pseudo & 0x0000_0000_00FF_0000 >> 8
	b3 := ins.Pseudo & 0x0000_0000_0000_FF00 << 8
	b4 := ins.Pseudo & 0x0000_0000_0000_00FF << 24
	return Imm(b1 | b2 | b3 | b4)
}

type Imm64 int64

func (ins Instruction) Imm64() Imm64 {
	return Imm64((int64(ins.NextImm()) << 32) | int64(ins.Imm()))
}

type Instruction struct {
	Basic      uint64
	Pseudo     uint64
	Extended64 bool
}

func (ins Instruction) String() string {
	if ins.Extended64 {
		return fmt.Sprintf("%016x %016x", ins.Basic, ins.Pseudo)
	}
	return fmt.Sprintf("%016x", ins.Basic)
}

func NewInstruction(ins uint64) Instruction {
	return Instruction{
		Basic: ins,
	}
}

func (ins Instruction) NeedPseudoInstruction() bool {
	return ins.Opcode().ToTyped().CodeOrMode() == BPF_IMM
}

func (ins *Instruction) AddPseudoInstruction(pseudoIns uint64) {
	ins.Pseudo = pseudoIns
	ins.Extended64 = true
}

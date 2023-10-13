package instruction

import (
	"fmt"
	"strings"
)

const (
	buggyCase = "this case should be impossible, this is a bug!"
)

func helperAssignment(operator string, op ArithmeticOpcode, ins Instruction) string {
	switch op.Source() {
	case BPF_K:
		return fmt.Sprintf("%s %s= %d", ins.Regs().DstReg(), operator, ins.Imm())
	case BPF_X:
		return fmt.Sprintf("%s %s= %s", ins.Regs().DstReg(), operator, ins.Regs().SrcReg())
	default:
		panic(buggyCase)
	}
}

func helperShiftmask(operator string, op ArithmeticOpcode, ins Instruction) string {
	var mask uint8
	switch ins.Opcode().Class() {
	case BPF_ALU:
		mask = 0x1F
	case BPF_ALU64:
		mask = 0x3F
	default:
		panic(buggyCase)
	}

	switch op.Source() {
	case BPF_K:
		return fmt.Sprintf("%s %s= %d", ins.Regs().DstReg(), operator, ins.Imm()&Imm(mask))
	case BPF_X:
		return fmt.Sprintf("%s %s= (%s & %d)", ins.Regs().DstReg(), operator, ins.Regs().SrcReg(), mask)
	default:
		panic(buggyCase)
	}
}

func disassembleArithmetic(op ArithmeticOpcode, ins Instruction) string {
	switch op.Code() {
	case BPF_ADD:
		return helperAssignment("+", op, ins)
	case BPF_SUB:
		return helperAssignment("-", op, ins)
	case BPF_MUL:
		return helperAssignment("*", op, ins)
	case BPF_DIV:
		if ins.Regs().SrcReg() == 0 {
			return fmt.Sprintf("%s = 0", ins.Regs().DstReg())
		}
		return helperAssignment("/", op, ins)
	case BPF_OR:
		return helperAssignment("|", op, ins)
	case BPF_AND:
		return helperAssignment("&", op, ins)
	case BPF_LSH:
		return helperShiftmask("<<", op, ins)
	case BPF_RSH:
		return helperShiftmask(">>", op, ins)
	case BPF_NEG:
		return fmt.Sprintf("%s = ~%s", ins.Regs().DstReg(), ins.Regs().SrcReg())
	case BPF_MOD:
		if ins.Regs().SrcReg() == 0 {
			return fmt.Sprintf("%s = %s", ins.Regs().DstReg(), ins.Regs().DstReg())
		}
		return helperAssignment("%", op, ins)
	case BPF_XOR:
		return helperAssignment("^", op, ins)
	case BPF_MOV:
		return helperAssignment("", op, ins)
	case BPF_ARSH:
		return helperShiftmask("s>>", op, ins)
	case BPF_END:
		return "byte swap! TODO"
	default:
		panic(buggyCase)
	}
}

func helperJumpConditional(operator string, op JumpOpcode, ins Instruction) string {
	switch op.Source() {
	case BPF_K:
		return fmt.Sprintf("if %s %s %d goto +%d", ins.Regs().DstReg(), operator, ins.Imm(), ins.Offset())
	case BPF_X:
		return fmt.Sprintf("if %s %s %s goto +%d", ins.Regs().DstReg(), operator, ins.Regs().SrcReg(), ins.Offset())
	default:
		panic(buggyCase)
	}
}

func disassembleJump(op JumpOpcode, ins Instruction) string {
	switch op.Code() {
	case BPF_JA:
		return fmt.Sprintf("goto +%d", ins.Offset())
	case BPF_JEQ:
		return helperJumpConditional("==", op, ins)
	case BPF_JGT:
		return helperJumpConditional(">", op, ins)
	case BPF_JGE:
		return helperJumpConditional(">=", op, ins)
	case BPF_JSET:
		return helperJumpConditional("&", op, ins)
	case BPF_JNE:
		return helperJumpConditional("!=", op, ins)
	case BPF_JSGT:
		return helperJumpConditional("s>", op, ins)
	case BPF_JSGE:
		return helperJumpConditional("s>=", op, ins)
	case BPF_CALL:
		switch op.Source() {
		case BPF_K, 0x2:
			return fmt.Sprintf("call %d", ins.Imm())
		case 0x1:
			return fmt.Sprintf("call +%d", ins.Offset())
		default:
			panic(buggyCase)
		}
	case BPF_EXIT:
		return "exit"
	case BPF_JLT:
		return helperJumpConditional("<", op, ins)
	case BPF_JLE:
		return helperJumpConditional("<=", op, ins)
	case BPF_JSLT:
		return helperJumpConditional("s<", op, ins)
	case BPF_JSLE:
		return helperJumpConditional("s<=", op, ins)
	default:
		panic(buggyCase)
	}
}

func disassembleImm(ins Instruction) string {
	switch ins.ImmSrc() {
	case BPF_IMM0:
		return fmt.Sprintf("%s = %d ll", ins.Regs().DstReg(), ins.Imm64())
	case BPF_IMM1:
		return fmt.Sprintf("%s = map_by_fd(%d)", ins.Regs().DstReg(), ins.Imm())
	case BPF_IMM2:
		return fmt.Sprintf("%s = map_val(map_by_fd(%d)) + %d", ins.Regs().DstReg(), ins.Imm(), ins.NextImm())
	case BPF_IMM3:
		return fmt.Sprintf("%s = var_addr(%d)", ins.Regs().DstReg(), ins.Imm())
	case BPF_IMM4:
		return fmt.Sprintf("%s = code_addr(%d)", ins.Regs().DstReg(), ins.Imm())
	case BPF_IMM5:
		return fmt.Sprintf("%s = map_by_idx(%d)", ins.Regs().DstReg(), ins.Imm())
	case BPF_IMM6:
		return fmt.Sprintf("%s = map_val(map_by_idx(%d)) + %d", ins.Regs().DstReg(), ins.Imm(), ins.NextImm())
	default:
		panic(buggyCase)
	}
}

// helperOffsetOperator helps to simplify the visual representation of
// operations such as (a + -b) to (a - b).
func helperOffsetOperator(off Offset) (offset Offset, operator string) {
	if off < 0 {
		operator = "-"
		offset = -off
	} else {
		operator = "+"
		offset = off
	}
	return offset, operator
}

func disassembleMem(op LoadAndStoreOpcode, ins Instruction) string {
	offset, operator := helperOffsetOperator(ins.Offset())

	// TODO: understand the difference between the size and unsigned size
	// "Where size is one of: BPF_B, BPF_H, BPF_W, or BPF_DW and 'unsigned size' is one of u8, u16, u32 or u64."
	switch ins.Opcode().Class() {
	case BPF_STX:
		return fmt.Sprintf("*(%s *)(%s %s %d) = %s", op.Size(), ins.Regs().DstReg(), operator, offset, ins.Regs().SrcReg())
	case BPF_ST:
		return fmt.Sprintf("*(%s *)(%s %s %d) = %d", op.Size(), ins.Regs().DstReg(), operator, offset, ins.Imm())
	case BPF_LDX:
		return fmt.Sprintf("%s = *(%s *)(%s + %d)", ins.Regs().DstReg(), op.Size(), ins.Regs().SrcReg(), ins.Offset())
	default:
		panic(buggyCase)
	}
}

func disassembleAtomic(op LoadAndStoreOpcode, ins Instruction) string {
	offset, offsetOperator := helperOffsetOperator(ins.Offset())

	switch op.Size() {
	case BPF_B, BPF_H:
		panic(fmt.Errorf("atomic operation on size %q is not supported", op.Size()))
	}

	var operator string
	switch ins.AtomicOperationImm() {
	case AtomicOperation(BPF_ADD):
		operator = "+"
	case AtomicOperation(BPF_OR):
		operator = "|"
	case AtomicOperation(BPF_AND):
		operator = "&"
	case AtomicOperation(BPF_XOR):
		operator = "^"
	case BPF_XCHG:
		panic("implement me")
	case BPF_CMPXCHG:
		panic("implement me")
	default:
		panic(fmt.Errorf("atomic operation does not support the %q operator", OpcodeArithmetic(ins.Imm())))
	}

	switch ins.Opcode().Class() {
	case BPF_STX:
		if ins.AtomicOperationImm()&AtomicOperation(BPF_FETCH) == 0x1 {
			panic("implement me")
		}
		return fmt.Sprintf("*(%s *)(%s %s %d) %s= %s", op.Size(), ins.Regs().DstReg(), offsetOperator, offset, operator, ins.Regs().SrcReg())
	default:
		panic(fmt.Errorf("atomic operation does not support the %q class", ins.Opcode().Class()))
	}
}

func disassembleLoadAndStore(op LoadAndStoreOpcode, ins Instruction) string {
	switch op.Mode() {
	case BPF_IMM:
		return disassembleImm(ins)
	case BPF_ABS:
		return "legacy BPF packet access (absolute)"
	case BPF_IND:
		return "legacy BPF packet access (indirect)"
	case BPF_MEM:
		return disassembleMem(op, ins)
	case BPF_MEMSX:
		return fmt.Sprintf("%s = *(s%s *)(%s + %d)", ins.Regs().DstReg(), strings.TrimPrefix(op.Size().String(), "u"), ins.Regs().SrcReg(), ins.Offset())
	case BPF_ATOMIC:
		return disassembleAtomic(op, ins)
	default:
		panic(buggyCase)
	}
}

func (ins Instruction) Disassemble() string {
	typedOpcode := ins.Opcode().ToTyped()
	switch op := typedOpcode.(type) {
	case ArithmeticOpcode:
		return disassembleArithmetic(op, ins)
	case JumpOpcode:
		return disassembleJump(op, ins)
	case LoadAndStoreOpcode:
		return disassembleLoadAndStore(op, ins)
	}
	return ""
}

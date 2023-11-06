package machine

import (
	"errors"
	"fmt"

	set "github.com/mtardy/mahebpf/pkg/instruction"
	"github.com/mtardy/mahebpf/pkg/program"
)

const (
	buggyCase = "this case should be impossible, this is a bug!"
)

var (
	ErrEndOfProgram = errors.New("end of program")
)

type VirtualMachine struct {
	// R0: return value from function calls, and exit value for eBPF programs
	// R1 - R5: arguments for function calls
	// R6 - R9: callee saved registers that function calls will preserve
	// R10: read-only frame pointer to access stack
	regs map[set.Register]int64

	// Program counter
	pc   int
	prog program.Program
}

func NewVirtualMachine() VirtualMachine {
	return VirtualMachine{
		regs: map[set.Register]int64{
			0:  0,
			1:  0,
			2:  0,
			3:  0,
			4:  0,
			5:  0,
			6:  0,
			7:  0,
			8:  0,
			9:  0,
			10: 0,
		},
	}
}

func (machine VirtualMachine) helperSource(source set.OpcodeSource, ins set.Instruction) int64 {
	switch source {
	case set.BPF_K:
		return int64(ins.Imm())
	case set.BPF_X:
		return machine.regs[ins.Regs().SrcReg()]
	default:
		panic(buggyCase)
	}
}

func (machine *VirtualMachine) executeArithmetic(op set.ArithmeticOpcode, ins set.Instruction) {
	switch op.Code() {
	case set.BPF_ADD:
		machine.regs[ins.Regs().DstReg()] += machine.helperSource(op.Source(), ins)
	case set.BPF_SUB:
		machine.regs[ins.Regs().DstReg()] -= machine.helperSource(op.Source(), ins)
	case set.BPF_MUL:
		machine.regs[ins.Regs().DstReg()] *= machine.helperSource(op.Source(), ins)
	// case set.BPF_DIV:
	// 	if set.Regs().SrcReg() == 0 {
	// 		return fmt.Sprintf("%s = 0", set.Regs().DstReg())
	// 	}
	// 	return helperAssignment("/", op, set)
	case set.BPF_OR:
		machine.regs[ins.Regs().DstReg()] |= machine.helperSource(op.Source(), ins)
	case set.BPF_AND:
		machine.regs[ins.Regs().DstReg()] &= machine.helperSource(op.Source(), ins)
	case set.BPF_LSH:
		machine.regs[ins.Regs().DstReg()] <<= machine.helperSource(op.Source(), ins)
	case set.BPF_RSH:
		machine.regs[ins.Regs().DstReg()] >>= machine.helperSource(op.Source(), ins)
	// 	return helperShiftmask(">>", op, set)
	// case set.BPF_NEG:
	// 	return fmt.Sprintf("%s = ~%s", set.Regs().DstReg(), set.Regs().SrcReg())
	// case set.BPF_MOD:
	// 	if set.Regs().SrcReg() == 0 {
	// 		return fmt.Sprintf("%s = %s", set.Regs().DstReg(), set.Regs().DstReg())
	// 	}
	// 	return helperAssignment("%", op, set)
	case set.BPF_XOR:
		machine.regs[ins.Regs().DstReg()] ^= machine.helperSource(op.Source(), ins)
	case set.BPF_MOV:
		machine.regs[ins.Regs().DstReg()] = machine.helperSource(op.Source(), ins)
	// case set.BPF_ARSH:
	// 	return helperShiftmask("s>>", op, set)
	// case set.BPF_END:
	// 	return "byte swap! TODO"
	default:
		// panic(buggyCase)
	}
}

func (machine *VirtualMachine) ExecuteStep() error {
	if machine.pc >= len(machine.prog) {
		return ErrEndOfProgram
	}
	ins := machine.prog[machine.pc]
	machine.pc++
	if ins.Extended64 {
		// jump over pseudo instruction
		machine.pc++
	}

	typedOpcode := ins.Opcode().ToTyped()
	switch op := typedOpcode.(type) {
	case set.ArithmeticOpcode:
		machine.executeArithmetic(op, ins)
	}
	return nil
}

func (machine VirtualMachine) NextInstruction() (*set.Instruction, error) {
	if machine.pc >= len(machine.prog) {
		return nil, ErrEndOfProgram
	}
	return &machine.prog[machine.pc], nil
}

func (machine *VirtualMachine) LoadProgram(prog program.Program) {
	machine.prog = prog
}

func (machine VirtualMachine) DumpRegs() string {
	return fmt.Sprint("Regs:", machine.regs)
}

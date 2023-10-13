package instruction

import (
	"testing"
)

const (
	// assembly
	// r1 += 0x11223344 // little
	//
	// opcode  src_Reg  dst_reg  offset  imm
	// 07      0        1        00 00   44 33 22 11
	exampleInstruction = 0x0701000044332211
)

func TestDecodeInstructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction uint64
		opcode      uint8
		regs        uint8
		srcReg      uint8
		dstReg      uint8
		offset      uint16
		imm         uint32
	}{
		{
			name:        "example",
			instruction: 0x0701000044332211,
			opcode:      0x07,
			regs:        0x01,
			srcReg:      0x0,
			dstReg:      0x1,
			offset:      0x0000,
			imm:         0x11223344,
		},
		{
			name:        "r1 = 0",
			instruction: 0xb701000000000000,
			opcode:      0xb7,
			regs:        0x01,
			srcReg:      0x0,
			dstReg:      0x1,
			offset:      0x0000,
			imm:         0x00000000,
		},
		{
			name:        "*(u32 *)(r10 - 4) = r1",
			instruction: 0x631afcff00000000,
			opcode:      0x63,
			regs:        0x1a,
			srcReg:      0x1,
			dstReg:      0xa,
			offset:      0xfffc,
			imm:         0x00000000,
		},
		{
			name:        "call 14",
			instruction: 0x850000000e000000,
			opcode:      0x85,
			regs:        0x00,
			srcReg:      0x0,
			dstReg:      0x0,
			offset:      0x0000,
			imm:         0x0000000e,
		},
		{
			name:        "r6 = r0",
			instruction: 0xbf06000000000000,
			opcode:      0xbf,
			regs:        0x06,
			srcReg:      0x0,
			dstReg:      0x6,
			offset:      0x0000,
			imm:         0x00000000,
		},
		{
			name:        "*(u32 *)(r10 - 8) = r6",
			instruction: 0x636af8ff00000000,
			opcode:      0x63,
			regs:        0x6a,
			srcReg:      0x6,
			dstReg:      0xa,
			offset:      0xfff8,
			imm:         0x00000000,
		},
		{
			name:        "r2 = r10",
			instruction: 0xbfa2000000000000,
			opcode:      0xbf,
			regs:        0xa2,
			srcReg:      0xa,
			dstReg:      0x2,
			offset:      0x0000,
			imm:         0x00000000,
		},
		{
			name:        "r2 += -4",
			instruction: 0x07020000fcffffff,
			opcode:      0x07,
			regs:        0x02,
			srcReg:      0x0,
			dstReg:      0x2,
			offset:      0x0000,
			imm:         0xfffffffc,
		},
		{
			name:        "r1 = 0 ll",
			instruction: 0x1801000000000000,
			opcode:      0x18,
			regs:        0x01,
			srcReg:      0x0,
			dstReg:      0x1,
			offset:      0x0000,
			imm:         0x00000000,
		},
		{
			name:        "r4 = 0",
			instruction: 0xb704000000000000,
			opcode:      0xb7,
			regs:        0x04,
			srcReg:      0x0,
			dstReg:      0x4,
			offset:      0x0000,
			imm:         0x00000000,
		},
		{
			name:        "call 2",
			instruction: 0x8500000002000000,
			opcode:      0x85,
			regs:        0x00,
			srcReg:      0x0,
			dstReg:      0x0,
			offset:      0x0000,
			imm:         0x00000002,
		},
		{
			name:        "goto +1",
			instruction: 0x0500010000000000,
			opcode:      0x05,
			regs:        0x00,
			srcReg:      0x0,
			dstReg:      0x0,
			offset:      0x00001,
			imm:         0x00000000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			instruction := NewInstruction(test.instruction)

			t.Run("opcode", func(t *testing.T) {
				opcode := instruction.Opcode()
				if opcode != Opcode(test.opcode) {
					t.Errorf("got %v, want %v", opcode, test.opcode)
				}
			})

			t.Run("regs", func(t *testing.T) {
				regs := instruction.Regs()
				if regs != Regs(test.regs) {
					t.Errorf("got %v, want %v", regs, test.regs)
				}

				t.Run("src", func(t *testing.T) {
					srcReg := instruction.Regs().SrcReg()
					if srcReg != Register(test.srcReg) {
						t.Errorf("got %v, want %v", srcReg, test.srcReg)
					}
				})

				t.Run("dst", func(t *testing.T) {
					dstReg := instruction.Regs().DstReg()
					if dstReg != Register(test.dstReg) {
						t.Errorf("got %v, want %v", dstReg, test.dstReg)
					}
				})
			})

			t.Run("offset", func(t *testing.T) {
				offset := instruction.Offset()
				if offset != Offset(test.offset) {
					t.Errorf("got %v, want %v", offset, test.offset)
				}
			})

			t.Run("imm", func(t *testing.T) {
				imm := instruction.Imm()
				if imm != Imm(test.imm) {
					t.Errorf("got 0x%x, want 0x%x", imm, test.imm)
				}
			})
		})
	}
}

func TestDecodeInstruction(t *testing.T) {
	t.Run("opcode", func(t *testing.T) {
		opcode := NewInstruction(exampleInstruction).Opcode()
		if opcode != 0x07 {
			t.Errorf("got %v, want %v", opcode, 0x07)
		}
	})

	t.Run("regs", func(t *testing.T) {
		regs := NewInstruction(exampleInstruction).Regs()
		if regs != 0x01 {
			t.Errorf("got %v, want %v", regs, 0x01)
		}

		t.Run("src", func(t *testing.T) {
			src_reg := NewInstruction(exampleInstruction).Regs().SrcReg()
			if src_reg != 0x0 {
				t.Errorf("got %v, want %v", src_reg, 0x0)
			}
		})

		t.Run("dst", func(t *testing.T) {
			dst_reg := NewInstruction(exampleInstruction).Regs().DstReg()
			if dst_reg != 0x1 {
				t.Errorf("got %v, want %v", dst_reg, 0x1)
			}
		})
	})

	t.Run("offset", func(t *testing.T) {
		offset := NewInstruction(exampleInstruction).Offset()
		if offset != 0x0000 {
			t.Errorf("got %v, want %v", offset, 0x0000)
		}
	})

	t.Run("imm", func(t *testing.T) {
		imm := NewInstruction(exampleInstruction).Imm()
		if imm != 0x11223344 {
			t.Errorf("got 0x%x, want 0x%x", imm, 0x11223344)
		}
	})
}

func TestStringer(t *testing.T) {
	ins := NewInstruction(exampleInstruction)
	t.Log(ins.Opcode().Code())
}

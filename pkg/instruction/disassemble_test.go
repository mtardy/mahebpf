package instruction

import (
	"testing"
)

func getTestingInstructions() []Instruction {
	rawInstructions := []uint64{
		0x0701000044332211,
		0xb701000000000000,
		0x631afcff00000000,
		0x850000000e000000,
		0xbf06000000000000,
		0x636af8ff00000000,
		0xbfa2000000000000,
		0x07020000fcffffff,
		0x1801000000000000,
		0xb704000000000000,
		0x8500000002000000,
		0x0500010000000000,
	}
	instructions := make([]Instruction, 0, len(rawInstructions))
	for _, rins := range rawInstructions {
		instructions = append(instructions, NewInstruction(rins))
	}
	return instructions
}

func BenchmarkDisassemble(b *testing.B) {
	instructions := getTestingInstructions()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		instructions[n%len(instructions)].Disassemble()
	}
}

func FuzzDisassemble(f *testing.F) {
	instructions := getTestingInstructions()
	for _, ins := range instructions {
		f.Add(ins.Basic)
	}
	f.Fuzz(func(t *testing.T, rawIns uint64) {
		NewInstruction(rawIns).Disassemble()
	})
}

func TestInstruction_Disassemble(t *testing.T) {
	type fields struct {
		Basic      uint64
		Pseudo     uint64
		Extended64 bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "r1 = 0",
			fields: fields{
				Basic: 0xb701000000000000,
			},
			want: "r1 = 0",
		},
		{
			name: "*(u32 *)(r10 - 4) = r1",
			fields: fields{
				Basic: 0x631afcff00000000,
			},
			want: "*(u32 *)(r10 - 4) = r1",
		},
		{
			name: "call 14",
			fields: fields{
				Basic: 0x850000000e000000,
			},
			want: "call 14",
		},
		{
			name: "r6 = r0",
			fields: fields{
				Basic: 0xbf06000000000000,
			},
			want: "r6 = r0",
		},
		{
			name: "*(u32 *)(r10 - 8) = r6",
			fields: fields{
				Basic: 0x636af8ff00000000,
			},
			want: "*(u32 *)(r10 - 8) = r6",
		},
		{
			name: "r2 = r10",
			fields: fields{
				Basic: 0xbfa2000000000000,
			},
			want: "r2 = r10",
		},
		{
			name: "r2 += -4",
			fields: fields{
				Basic: 0x07020000fcffffff,
			},
			want: "r2 += -4",
		},
		{
			name: "r1 = 0 ll",
			fields: fields{
				Basic:      0x1801000000000000,
				Pseudo:     0x0000000000000000,
				Extended64: true,
			},
			want: "r1 = 0 ll",
		},
		{
			name: "call 1",
			fields: fields{
				Basic: 0x8500000001000000,
			},
			want: "call 1",
		},
		{
			name: "if r0 != 0 goto +9",
			fields: fields{
				Basic: 0x5500090000000000,
			},
			want: "if r0 != 0 goto +9",
		},
		{
			name: "r2 = r10",
			fields: fields{
				Basic: 0xbfa2000000000000,
			},
			want: "r2 = r10",
		},
		{
			name: "r2 += -4",
			fields: fields{
				Basic: 0x07020000fcffffff,
			},
			want: "r2 += -4",
		},
		{
			name: "r3 = r10",
			fields: fields{
				Basic: 0xbfa3000000000000,
			},
			want: "r3 = r10",
		},
		{
			name: "r3 += -8",
			fields: fields{
				Basic: 0x07030000f8ffffff,
			},
			want: "r3 += -8",
		},
		{
			name: "r1 = 0 ll",
			fields: fields{
				Basic:      0x1801000000000000,
				Pseudo:     0x0000000000000000,
				Extended64: true,
			},
			want: "r1 = 0 ll",
		},
		{
			name: "r4 = 0",
			fields: fields{
				Basic: 0xb704000000000000,
			},
			want: "r4 = 0",
		},
		{
			name: "call 2",
			fields: fields{
				Basic: 0x8500000002000000,
			},
			want: "call 2",
		},
		{
			name: "goto +1",
			fields: fields{
				Basic: 0x0500010000000000,
			},
			want: "goto +1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins := Instruction{
				Basic:      tt.fields.Basic,
				Pseudo:     tt.fields.Pseudo,
				Extended64: tt.fields.Extended64,
			}
			if got := ins.Disassemble(); got != tt.want {
				t.Errorf("Instruction.Disassemble() = %q, want %q", got, tt.want)
			}
		})
	}
}

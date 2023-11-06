package program

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/mtardy/mahebpf/pkg/instruction"
)

type Program []instruction.Instruction

func NewProgram() Program {
	return Program{}
}

type DisassembledProgram struct {
	Instruction  instruction.Instruction
	Disassembled string
}

func (p Program) Disassemble() []DisassembledProgram {
	out := []DisassembledProgram{}
	for _, ins := range p {
		out = append(out, DisassembledProgram{
			Instruction:  ins,
			Disassembled: ins.Disassemble(),
		})
	}
	return out
}

func parseBytes(data []byte, width int, parser func(data []byte, index, width int) (uint64, error)) (*Program, error) {
	prog := NewProgram()
	for i := 0; i+width <= len(data); i = i + width {
		parsedInstruction, err := parser(data, i, width)
		if err != nil {
			return nil, err
		}
		ins := instruction.NewInstruction(parsedInstruction)

		var pseudoIns uint64
		if ins.NeedPseudoInstruction() {
			if i+width >= len(data) {
				return nil, fmt.Errorf("ins 0x%016x needs a pseudo instruction and it's not available", ins.Basic)
			}
			i = i + width
			pseudoIns, err = parser(data, i, width)
			if err != nil {
				return nil, err
			}
			ins.AddPseudoInstruction(pseudoIns)
		}

		prog = append(prog, ins)

		if ins.NeedPseudoInstruction() {
			prog = append(prog, instruction.NewPseudoInstruction(pseudoIns))
		}
	}
	return &prog, nil
}

func listELFSections(file *elf.File) []string {
	sections := make([]string, 0, len(file.Sections))
	for _, sec := range file.Sections {
		if sec.Name == "" {
			continue
		}
		sections = append(sections, sec.Name)
	}
	return sections
}

func ListELFSections(path string) ([]string, error) {
	file, err := elf.Open(path)
	if err != nil {
		return nil, err
	}
	return listELFSections(file), nil
}

func FromELF(path string, section string) (*Program, error) {
	file, err := elf.Open(path)
	if err != nil {
		return nil, err
	}

	sec := file.Section(section)

	if sec == nil {
		return nil, fmt.Errorf("section not found, available sections: %v", listELFSections(file))
	}

	byteCode, err := sec.Data()
	if err != nil {
		return nil, err
	}

	if len(byteCode)%8 != 0 {
		return nil, errors.New("section program len is not a multiple of 8")
	}

	return parseBytes(byteCode, 8, func(data []byte, index, width int) (uint64, error) {
		return binary.BigEndian.Uint64(data[index : index+width]), nil
	})
}

func FromASCII(path string) (*Program, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text = bytes.ReplaceAll(text, []byte(" "), []byte(""))
	text = bytes.ReplaceAll(text, []byte("\n"), []byte(""))

	if len(text)%16 != 0 {
		return nil, errors.New("text len is not a multiple of 16")
	}

	return parseBytes(text, 16, func(data []byte, index, width int) (uint64, error) {
		return strconv.ParseUint(string(text[index:index+width]), width, 64)
	})
}

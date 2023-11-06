package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/mtardy/mahebpf/pkg/machine"
	"github.com/mtardy/mahebpf/pkg/program"
)

var (
	fileTypeOption string
	bytesOption    bool
	numberOption   bool
	executeOption  bool
)

const usage = `Usage: dbpf [flags] file [section]

An educational eBPF disassembler

Flags:`

func init() {
	flag.StringVar(&fileTypeOption, "type", "elf", "type of the file to disassemble (elf or ascii)")
	flag.BoolVar(&bytesOption, "bytes", true, "print instruction bytes")
	flag.BoolVar(&numberOption, "number", true, "print line number")
	flag.BoolVar(&executeOption, "e", false, "execute program")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func printDisassembled(disassembled []program.DisassembledProgram) {
	width := len(strconv.Itoa(len(disassembled) - 1))
	for i, ins := range disassembled {
		if ins.Instruction.IsPseudo {
			i++
			continue
		}
		out := strings.Builder{}
		if numberOption {
			out.WriteString(fmt.Sprintf("%*d: ", width, i))
		}
		if bytesOption {
			out.WriteString(ins.Instruction.String() + " ")
		}
		out.WriteString(ins.Instruction.Disassemble())
		fmt.Println(out.String())
	}
}

func Execute() {
	flag.Parse()

	// handle SIGPIPE
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGPIPE)
	go func() {
		<-sigs
		os.Exit(0)
	}()

	if executeOption {
		prog, err := program.FromASCII("bin.txt")
		if err != nil {
			panic(err)
		}

		vm := machine.NewVirtualMachine()
		vm.LoadProgram(*prog)
		for {
			ins, err := vm.NextInstruction()
			if err != nil {
				break
			}
			fmt.Println(ins.Disassemble())

			err = vm.ExecuteStep()
			if err != nil {
				break
			}
			fmt.Println(vm.DumpRegs())
		}
		return
	}

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, usage)
		flag.PrintDefaults()
		os.Exit(2)
	}

	var prog *program.Program
	var err error
	switch strings.ToLower(fileTypeOption) {
	case "elf":
		if len(flag.Args()) < 2 {
			sections, err := program.ListELFSections(flag.Arg(0))
			if err != nil {
				fatal(err)
			}
			fatal(fmt.Errorf("please provide an ELF section to disassemble, available sections: %v", sections))
		}
		prog, err = program.FromELF(flag.Arg(0), flag.Arg(1))
		if err != nil {
			fatal(err)
		}
	case "ascii":
		prog, err = program.FromASCII(flag.Arg(0))
		if err != nil {
			fatal(err)
		}
	default:
		fatal(fmt.Errorf("invalid type %q, the only type available are elf or ascii", fileTypeOption))
	}

	if prog != nil {
		printDisassembled(prog.Disassemble())
	}
}

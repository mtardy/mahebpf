# maheBPF

maheBPF™ for my asinine holistic enterprise BPF pseudocode fragmentor.

## Installation

Install it with a Golang install, the project is nice enough to not have
external dependencies (for now).

```shell-session
go install github.com/mtardy/mahebpf@latest
```

## Usage

### 🧝🏻‍♀️ ELF 🧝🏻‍♂️

Let's say you have a BPF program in an ELF at the section kprobe/pizza
(little-endian only club 😎 thanks) and you want to disassemble it with a
military-grade™ dissasembler:

```shell-session
mahebpf prog.o kprobe/pizza
```

For my very useful little program, the output looks like this:

```text
 0: b701000000000000 r1 = 0
 1: 631afcff00000000 *(u32 *)(r10 - 4) = r1
 2: 850000000e000000 call 14
 3: bf06000000000000 r6 = r0
 4: 636af8ff00000000 *(u32 *)(r10 - 8) = r6
 5: bfa2000000000000 r2 = r10
 6: 07020000fcffffff r2 += -4
 7: 1801000000000000 0000000000000000 r1 = 0 ll
 9: 8500000001000000 call 1
10: 5500090000000000 if r0 != 0 goto +9
11: bfa2000000000000 r2 = r10
12: 07020000fcffffff r2 += -4
13: bfa3000000000000 r3 = r10
14: 07030000f8ffffff r3 += -8
15: 1801000000000000 0000000000000000 r1 = 0 ll
17: b704000000000000 r4 = 0
18: 8500000002000000 call 2
19: 0500010000000000 goto +1
20: 6360000000000000 *(u32 *)(r0 + 0) = r6
21: b700000000000000 r0 = 0
22: 9500000000000000 exit
```

Cool no? A bit like `llvm-objdump -S prog.o` but in bad.

### 🇺🇸 ASCII 🦅 

If you like to store your eBPF bytecode in ASCII in a text format like a person
of taste, I got you covered. Let's say you have a program in a `prog.txt` that
looks like this:

```text
b7 01 00 00 00 00 00 00
63 1a fc ff 00 00 00 00
85 00 00 00 0e 00 00 00
bf 06 00 00 00 00 00 00
63 6a f8 ff 00 00 00 00
bf a2 00 00 00 00 00 00
07 02 00 00 fc ff ff ff
18 01 00 00 00 00 00 00
00 00 00 00 00 00 00 00
85 00 00 00 01 00 00 00
55 00 09 00 00 00 00 00
bf a2 00 00 00 00 00 00
07 02 00 00 fc ff ff ff
bf a3 00 00 00 00 00 00
07 03 00 00 f8 ff ff ff
18 01 00 00 00 00 00 00
00 00 00 00 00 00 00 00
b7 04 00 00 00 00 00 00
85 00 00 00 02 00 00 00
05 00 01 00 00 00 00 00
```

To disassemble this hexabeauty:

```shell-session
mahebpf --type ascii prog.txt
```

Boom 💥🤯, same output as before!

## Contribute

Don't.

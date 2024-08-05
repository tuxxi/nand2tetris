package instruction

import (
	"log"
	"strconv"
	"strings"

	"github.com/tuxxi/nand2tetris/06/symbol"
)

func parseJumpType(jmp string) jumpType {
	switch jmp {
	case "JGT":
		return 0b001
	case "JEQ":
		return 0b010
	case "JGE":
		return 0b011
	case "JLT":
		return 0b100
	case "JNE":
		return 0b101
	case "JLE":
		return 0b110
	case "JMP":
		return 0b111
	default:
		// panic
		log.Fatalln("unknown jump type:", jmp)
	}
	return 0
}

func parseDestBits(dest string) destBits {
	var ret destBits
	for _, c := range dest {
		if c == 'A' {
			ret |= (1 << 2)
		}
		if c == 'D' {
			ret |= (1 << 1)
		}
		if c == 'M' {
			ret |= (1 << 0)
		}
	}
	return ret
}

// Parse the instructions
func Parse(lines []string, symbols *symbol.SymbolTable) []Instruction {
	var (
		instrs  []Instruction
		ROMAddr uint16
	)

	for i := range lines {
		// trim leading and trailing whitespace
		trimmed := strings.TrimSpace(lines[i])

		// skip comments
		if strings.HasPrefix(trimmed, "//") {
			continue
		}

		// skip newlines and blank lines
		if trimmed == "\n" || trimmed == "" {
			continue
		}

		// Label
		if trimmed[0] == '(' && trimmed[len(trimmed)-1] == ')' {
			sym := trimmed[1 : len(trimmed)-1]
			symbols.AddLabel(sym, ROMAddr)
			continue
		}

		// A instruction
		if trimmed[0] == '@' {
			log.Printf("Encountered A instruction at PC %v: '%s'\n", ROMAddr, trimmed)
			value := trimmed[1:]

			var symbol string

			// try to parse the instruction as a number
			v, err := strconv.Atoi(value)
			if err != nil {
				symbol = value
				// Don't resolve variable any symbols now, just record that we are using it.
			} else if v > 0x7FFF {
				panic("value out of range")
			}
			instrs = append(instrs, &AInstruction{val: uint16(v), symbol: symbol})
			ROMAddr++
			continue
		}

		// C instruction
		instr := CInstruction{raw: trimmed}
		log.Printf("Encountered C instruction at PC %v: '%s'\n", ROMAddr, trimmed)
		// has dest?
		if dest, rest, hasDest := strings.Cut(trimmed, "="); hasDest {
			instr.dest = parseDestBits(dest)
			trimmed = rest // slice up
		}
		if op, jump, hasJump := strings.Cut(trimmed, ";"); hasJump {
			instr.jmp = parseJumpType(jump)
			trimmed = op // slice
		}
		instr.operation = trimmed
		instrs = append(instrs, instr)
		ROMAddr++
	}
	return instrs
}

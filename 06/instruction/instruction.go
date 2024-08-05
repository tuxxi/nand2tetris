package instruction

import (
	"fmt"
	"log"
	"strings"

	"github.com/tuxxi/nand2tetris/06/symbol"
)

type Instruction interface {
	// implement String(), for debugging
	fmt.Stringer

	// Get the encoded value of an instruction
	Encoded() uint16

	// Link an instruction with the previously populated symbol table.
	// Returns a new Instruction when linked, or nil if no change.
	Link(table *symbol.SymbolTable) Instruction
}

type AInstruction struct {
	val    uint16 // raw value
	symbol string // symbol which needs resolving
}

func (a AInstruction) String() string {
	if a.symbol != "" {
		return fmt.Sprintf("@%s", a.symbol)
	}
	return fmt.Sprintf("@%d", a.val)
}

func (a AInstruction) Encoded() uint16 {
	if a.val > 0x7FFF {
		panic("bad a-instruction: cannot encode")
	}
	return uint16(a.val)
}

// Link a symbol with the instructions.
func (a AInstruction) Link(st *symbol.SymbolTable) Instruction {
	if a.symbol != "" {
		sym, ok := st.Get(a.symbol)
		if !ok {
			// Not found - it's an implicitly created variable.
			// We need to create it
			sym = st.AddVariable(a.symbol)
		}
		ret := AInstruction{val: sym.Location}
		log.Printf("LINKED: %s -> %s via %s\n", a, ret, sym)
		return ret
	}
	return nil
}

type jumpType uint8
type destBits uint8

type CInstruction struct {
	dest      destBits // d1-d3
	jmp       jumpType // j1-j3
	operation string   // c1-c6

	raw string
}

func (c CInstruction) String() string {
	return c.raw
}

func (c CInstruction) Encoded() uint16 {
	var ret uint16

	// set all the jmp bits
	ret |= uint16(c.jmp) << 0

	// set all the dst bits
	ret |= (uint16(c.dest << 3))

	// set all the comp bits
	op := c.operation
	if strings.Contains(c.operation, "M") {
		ret |= (1 << 12) // 'a' bit
		// make it easier to work with
		op = strings.ReplaceAll(c.operation, "M", "A")
	}

	// do this very dumb. there's probably a cleaner way to encode this
	comp := 0
	switch op {
	case "0":
		comp = 0b101010
	case "1":
		comp = 0b111111
	case "-1":
		comp = 0b111010
	case "D":
		comp = 0b001100
	case "A":
		comp = 0b110000
	case "!D":
		comp = 0b001101
	case "!A":
		comp = 0b110001
	case "-D":
		comp = 0b001111
	case "-A":
		comp = 0b110011
	case "D+1":
		comp = 0b011111
	case "A+1":
		comp = 0b110111
	case "D-1":
		comp = 0b001110
	case "A-1":
		comp = 0b110010
	case "D+A":
		comp = 0b000010
	case "D-A":
		comp = 0b010011
	case "A-D":
		comp = 0b000111
	case "D&A":
		comp = 0b000000
	case "D|A":
		comp = 0b010101
	}
	ret |= (uint16(comp) << 6)

	// C instructions always have last 3 bits set
	ret |= (uint16(0b111) << 13)
	return ret
}

// no-op, to implement interface : C-instructions do not have symbols
func (c CInstruction) Link(*symbol.SymbolTable) Instruction { return nil }

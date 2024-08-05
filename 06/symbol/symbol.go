package symbol

import (
	"fmt"
	"log"
	"strings"
)

// Symbols

type SymbolType byte

const (
	symbolTypeBuiltin  = 0x0
	symbolTypeLabel    = 0x1
	symbolTypeVariable = 0x2
)

func (typ SymbolType) String() string {
	switch typ {
	case symbolTypeBuiltin:
		return "BUILTIN"
	case symbolTypeLabel:
		return "LABEL"
	case symbolTypeVariable:
		return "VARIABLE"
	}
	return "(UNKNOWN)"
}

type Symbol struct {
	Location uint16

	typ    SymbolType
	symbol string
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s (%s) -> %v", s.typ, s.symbol, s.Location)
}

type SymbolTable struct {
	ramAddr uint16
	tab     map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		tab: make(map[string]Symbol),
		// RAM pseudo-variables start at 0x16
		ramAddr: 0x10,
	}

	// add global builtin variables
	st.addBuiltin("SP", 0x0000)
	st.addBuiltin("LCL", 0x0001)
	st.addBuiltin("ARG", 0x0002)
	st.addBuiltin("THIS", 0x0003)
	st.addBuiltin("THAT", 0x0004)
	st.addBuiltin("SCREEN", 0x4000)
	st.addBuiltin("KBD", 0x6000)

	// add pesudo-registers
	for i := uint16(0); i < 16; i++ {
		st.addBuiltin(fmt.Sprintf("R%d", i), i)
	}

	return st
}

func (st *SymbolTable) String() string {
	b := strings.Builder{}
	b.WriteRune('{')
	b.WriteRune('\n')
	for _, v := range st.tab {
		b.WriteString(fmt.Sprintf("\t%v\n", v))
	}
	b.WriteRune('}')
	return b.String()
}

func (st *SymbolTable) addBuiltin(symbol string, loc uint16) {
	st.tab[symbol] = Symbol{symbol: symbol, typ: symbolTypeBuiltin, Location: loc}
}

func (st *SymbolTable) AddLabel(symbol string, loc uint16) {
	st.tab[symbol] = Symbol{symbol: symbol, typ: symbolTypeLabel, Location: loc}
	log.Printf("Adding label: %s\n", st.tab[symbol])
}

func (st *SymbolTable) AddVariable(symbol string) Symbol {
	sym := Symbol{symbol: symbol, typ: symbolTypeLabel, Location: st.ramAddr}
	st.tab[symbol] = sym
	log.Printf("Adding variable: %s\n", st.tab[symbol])
	st.ramAddr++
	return sym
}

func (st *SymbolTable) Get(symbol string) (Symbol, bool) {
	s, ok := st.tab[symbol]
	return s, ok
}

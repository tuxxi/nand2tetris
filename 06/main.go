package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tuxxi/nand2tetris/06/instruction"
	"github.com/tuxxi/nand2tetris/06/symbol"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("usage: %s filename.asm", os.Args[0])
		os.Exit(1)
	}
	inFile := os.Args[1]
	b, err := os.ReadFile(inFile)
	if err != nil {
		panic(err)
	}
	log.Println(string(b))
	lines := strings.Split(string(b), "\n")

	// first pass, generate labels
	symbols := symbol.NewSymbolTable()
	log.Println("Initialized symbol table:", symbols)
	instrs := instruction.Parse(lines, symbols)

	// second pass, resolve RAM addresses and symbols.
	resolvedInstrs := linkSymbols(instrs, symbols)

	name := strings.Split(filepath.Base(inFile), ".")
	outName := name[0] + ".hack"
	f, err := os.OpenFile(outName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// encode the entire IR and stream it to the output file
	log.Println("Writing output to", outName)
	err = encode(resolvedInstrs, f)
	if err != nil {
		panic(err)
	}
	log.Println("Finished writing output:", outName)
}

func linkSymbols(instructions []instruction.Instruction, symbols *symbol.SymbolTable) []instruction.Instruction {
	var (
		resolvedInstructions []instruction.Instruction
	)
	for _, instr := range instructions {
		if newInstr := instr.Link(symbols); newInstr != nil {
			instr = newInstr
		}
		resolvedInstructions = append(resolvedInstructions, instr)
	}

	return resolvedInstructions
}

func encode(linked []instruction.Instruction, out io.Writer) error {
	// encode the final IR as binary
	for i, instr := range linked {
		// encode
		encoded := instr.Encoded()

		// Encode the instruction as ASCII-encoded binary.
		// Each instructio is separated with a newline, except for at EOF.
		var nl string
		if i < len(linked)-1 {
			nl = "\n"
		}
		s := fmt.Sprintf("%016b%s", encoded, nl)
		log.Printf("%d: %s -> %s", i, instr.String(), s)
		_, err := out.Write([]byte(s))
		if err != nil {
			return err
		}
	}
	return nil
}

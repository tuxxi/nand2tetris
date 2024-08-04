// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
// The algorithm is based on repetitive addition.

// a * b = (a+a+a+a+...+a)
// for 0 .. a: r2+=b


// Repeteadly add R1 into R2, R0 times.


// Set R2 to 0
@R2
M=0

// Copy R0 into R15 to use as our counter
@R0
D=M
@R15
M=D

(loop)
    // if R15 == 0; goto end
    @R15
    D=M
    @end
    D;JEQ

    @R1
    D=M   // D = RAM[1]

    @R2
    M=D+M // RAM[2] += D

    @R15
    M=M-1 // R15--

    // loop
    @loop
    0;JMP

(end)
@end
0;JMP

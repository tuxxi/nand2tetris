// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed,
// the screen should be cleared.

// @screen[0..31]       = row 0   (512 px)
// ...
// @screen[8160...8191] = row 255 (512 px)
//
// SetPixel(row, col, val):
//   addr = 16384 + (32*row + col/16)
//   word = RAM[addr]
//   new_word[col%16] = val
//   RAM[addr] = new_word

(START)
    @KBD
    D=M     // D = keycode

    @CLEAR_SCREEN
    D;JEQ  // if D == 0 goto CLEAR_SCREEN
           // else fallthrough to BLACK_SCEREN

(BLACK_SCREEN)
    // for i ... rows
    @8192
    D=A     // D = 8191

    (loop_1)
    @SCREEN
    A=D+A
    M=-1      // RAM[screen + D] = 0xFFFF

    @loop_1
    D=D-1;JGE  // D=D-1; if D >= 0 goto loop_1

    // reset
    @START
    0;JMP

(CLEAR_SCREEN)
    // for i ... rows
    @8192
    D=A     // D = 8191

    (loop_cls)
    @SCREEN
    A=D+A
    M=0      // RAM[screen + D] = 0x0

    @loop_cls
    D=D-1;JGE  // D=D-1; if D >= 0 goto loop_cls

    // reset
    @START
    0;JMP

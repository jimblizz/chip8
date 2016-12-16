package main

import (
    "fmt"
)

/*
	0NNN	Calls RCA 1802 program at address NNN.
	00E0	Clears the screen.
	00EE	Returns from a subroutine.
	1NNN	Jumps to address NNN.
	2NNN	Calls subroutine at NNN.
	3XNN	Skips the next instruction if VX equals NN.
	4XNN	Skips the next instruction if VX doesn't equal NN.
	5XY0	Skips the next instruction if VX equals VY.
	6XNN	Sets VX to NN.
	7XNN	Adds NN to VX.
	8XY0	Sets VX to the value of VY.
	8XY1	Sets VX to VX or VY.
	8XY2	Sets VX to VX and VY.
	8XY3	Sets VX to VX xor VY.
	8XY4	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
	8XY5	VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
	8XY6	Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.
	8XY7	Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
	8XYE	Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift.
	9XY0	Skips the next instruction if VX doesn't equal VY.
	ANNN	Sets I to the address NNN.
	BNNN	Jumps to the address NNN plus V0.
	CXNN	Sets VX to a random number, masked by NN.
	DXYN	Sprites stored in memory at location in index register (I), maximum 8bits wide. Wraps around the screen. If when drawn, clears a pixel, register VF is set to 1 otherwise it is zero. All drawing is XOR drawing (i.e. it toggles the screen pixels)
	EX9E	Skips the next instruction if the key stored in VX is pressed.
	EXA1	Skips the next instruction if the key stored in VX isn't pressed.
	FX07	Sets VX to the value of the delay timer.
	FX0A	A key press is awaited, and then stored in VX.
	FX15	Sets the delay timer to VX.
	FX18	Sets the sound timer to VX.
	FX1E	Adds VX to I.
	FX29	Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
	FX33	Stores the Binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)
	FX55	Stores V0 to VX in memory starting at address I.
	FX65	Fills V0 to VX with values from memory starting at address I.
*/

func (c* Chip8) Cycle () (error) {
    // Each opcode is two bytes, so we get two memory locations and increment the PC by 2 each cycle
    op := uint16(c.Memory[c.PC]) << 8 | uint16(c.Memory[c.PC+1])

    draw := false

    fmt.Printf("%X    %X       ", c.PC, op)

    switch op & 0xF000 { // Get the first char of the op

    case 0x0000:
        if op & 0xFF00 != 0x0000 {
            goto NOTIMPLEMENTED
        }

        switch op & 0x000F {
        case 0x0000:
            // 00E0	Clears the screen.
            fmt.Printf("00E0 Clear the screen")
            c.gfx = make([]byte, screenWidth * screenHeight)

        default:
            // 00EE	Returns from a subroutine.
            if c.SP <= 0 {
                return fmt.Errorf("Stack bottom")
            }
            c.PC = c.Stack[c.SP]
            c.SP--
        }

        break

    case 0x1000:
        // 1NNN	Jumps to address NNN.
        fmt.Printf("1NNN Jump to %X", op & 0x0FFF)

        c.PC = op & 0x0FFF // Jump to location
        // Mask: 0x0FFF
        // Uses: 0xNYYY (N = no, Y = Yes)

        goto SKIPADVANCE
        break

    case 0x2000:
        // 2NNN Calls subroutine at NNN.
        nnn := op & 0x0FFF
        fmt.Printf("2NNN Call subroutine at %X ", nnn)

        c.SP++                  // Increment the stack pointer
        c.Stack[c.SP] = c.PC    // Set the stack pointer current SP, this tis the return location
        c.PC = op & 0x0FFF      // Move the program counter

        goto SKIPADVANCE
        break

    case 0x6000:
        // 6XNN	Sets VX to NN
        x := (op & 0x0F00) >> 8
        nn := byte(op)

        fmt.Printf("6XNN Set V%X to %X", x, nn)

        c.V[x] = nn
        break

    case 0x7000:
        // 7XNN	Adds NN to VX.
        x := (op & 0x0F00) >> 8
        nn := byte(op & 0x00F)

        fmt.Printf("6XNN Set V%X to %X", x, nn)

        c.V[x] += nn;
        break

    case 0x8000:
        x := (op & 0x0F00) >> 8 // 2nd col
        y := (op & 0x00F0) >> 4 // 3rd col
        z := (op & 0x000F) // 3rd col

        // Switch based on col 4
        switch z {

        case 0x0000:
            // 8XY0	Sets VX to the value of VY.
            fmt.Printf("8XY0 Set V%X = V%X", x, y)
            c.V[x] = c.V[y]
            break

        case 0x0001:
            // 8XY1	Sets VX to VX or VY.
            fmt.Printf("8XY0 Set V%X = V%X or V%X", x, x, y)
            c.V[x] = c.V[y] | c.V[x]
            break

        case 0x0002:
            // 8XY2	Sets VX to VX and VY.
            fmt.Printf("8XY0 Set V%X = V%X and V%X", x, x, y)
            c.V[x] = c.V[y] & c.V[x]
            break

        case 0x0003:
            // 8XY3	Sets VX to VX xor VY.
            // Set Vx to the XOR of Vx and Vy
            // xor is ^
            fmt.Printf("8XY3 Set V%X = V%X xor V%X", x, x, y)
            c.V[x] = c.V[y] ^ c.V[x]
            break

        case 0x0004:
            // 8XY4	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
            fmt.Printf("8XY4 Set V%X = V%X + V%X, set VF = carry", x, x, y)

            // Add Vx and Vy together
            r := uint16(c.V[x]) + uint16(c.V[y])

            // If the result > 8 bits (>256) the set VF to 1
            var cf byte     // Defaults to 0
            if r > 0xFF {   // r > 255
                cf = 1
            }
            c.V[0xF] = cf

            c.V[x] = byte(r) // And then set Vx to the result of the addition
            break

        case 0x0005:
            // 8XY5	VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
            goto NOTIMPLEMENTED
            break

        case 0x0006:
            // 8XY6	Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.
            goto NOTIMPLEMENTED
            break

        case 0x0007:
            // 8XY7	Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
            goto NOTIMPLEMENTED
            break

        case 0x000E:
            // 8XYE	Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift.
            goto NOTIMPLEMENTED
            break

        default:
            goto NOTIMPLEMENTED
        }
        break

    case 0xA000:
        nnn := op & 0x0FFF
        fmt.Printf("ANNN Set I to %X", nnn)

        c.I = nnn
        break

    case 0xD000:
        // DXYN	Sprites stored in memory at location in index register (I), maximum 8bits wide.
        // Wraps around the screen. If when drawn, clears a pixel, register VF is set to 1 otherwise it is zero.
        // All drawing is XOR drawing (i.e. it toggles the screen pixels)

        // Height comes from the opcode, x,y from the stated register

        // Based on the implementation from (stolen from):
        // http://www.multigesture.net

        x := uint16(c.V[( op & 0x0F00 ) >> 8])
        y := uint16(c.V[( op & 0x00F0 ) >> 4])
        h := uint16(op & 0x000F) // Height

        fmt.Printf("Draw a sprite %x (from V%X) by %x (from V%X) and to height %X", x, c.V[x], y, c.V[y], h)

        var pixel byte
        c.V[0xF] = 0 // TODO: Why?

        // Y
        // Max h hight
        for yline := uint16(0); yline < h && yline+y < screenHeight; yline++ {
            pixel = c.Memory[c.I+yline]

            // X
            // Max 8 bits wide
            for xline := uint16(0); xline < 8; xline++ {
                // TODO: Why this check?
                if (pixel & (0x80 >> xline)) != 0 {
                    offset := (x + xline + ((y + yline) * screenWidth))

                    if c.gfx[offset] == 1 {
                        // VF is set to 1 if any screen pixels are flipped from
                        // set to unset when the sprite is drawn, and to 0 if
                        // that doesn't happen.
                        c.V[0xF] = 1
                    }
                    c.gfx[offset] ^= 1

                }
            }
        }
        draw = true
        break

    case 0xF000:
        // 2nd bit as register index
        x := (op & 0x0F00) >> 8

        // Switch on the last two bits
        switch op & 0x00FF {

        case 0x55:
            // FX55	Stores V0 to VX in memory starting at address I.
            // Copy V0...Vx to memory, in order, starting from memory address in I.
            fmt.Printf("FX55 Copy V0..V%X to memory, starting at I (%X)", x, c.I)

            for i := 0; uint16(i) <= x; i++ {
                c.Memory[c.I+uint16(i)] = c.V[i]
            }
            break

        case 0x1E:
            // FX1E Adds VX to I.
            // TODO: Why would you do that?

            v := uint16(c.V[x])

            fmt.Printf("FX1E Add V%X to I", v)

            c.I = c.I + v
            break

        case 0x65:
            // FX65	Fills V0 to VX with values from memory starting at address I.
            // So this is the inverse of FX55
            fmt.Printf("FX65 Copy memory to V0..V%X, starting at I (%X)", x, c.I)

            for i := uint16(0); i <= uint16(x); i++ { // For loop counting in hex :)
                c.V[i] = c.Memory[c.I+i] // Copy for each iteration
            }
            break

        default:
            goto NOTIMPLEMENTED
        }

        //FX07	Sets VX to the value of the delay timer.
        //FX0A	A key press is awaited, and then stored in VX.
        //FX15	Sets the delay timer to VX.
        //FX18	Sets the sound timer to VX.

        //FX29	Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
        //FX33	Stores the Binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)

        break

    default:
        goto NOTIMPLEMENTED
    }

    c.PC += 2

SKIPADVANCE:

    if draw {
        fmt.Print("     Draw!")
    }

    fmt.Println("")

    return nil


NOTIMPLEMENTED:
    return fmt.Errorf("%X not implemeted", op & 0xF000, op)
}
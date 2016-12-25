package main

import (
    "github.com/davecgh/go-spew/spew"
    "fmt"
    "io/ioutil"
    "bytes"
    "io"
    "time"
)

const (
    screenWidth     = 64
    screenHeight    = 32
    cpuFrequency    = 60 // To original was 60hz
)

var chip8Fontset = [80]byte{
    // Extra spooky magic numbers
    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
    0x20, 0x60, 0x20, 0x20, 0x70, // 1
    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
    0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8 struct {
    // The 4096 bytes of memory.
    //
    // Memory Map:
    // +---------------+= 0xFFF (4095) End of Chip-8 RAM
    // |               |
    // |               |
    // |               |
    // |               |
    // |               |
    // | 0x200 to 0xFFF|
    // |     Chip-8    |
    // | Program / Data|
    // |     Space     |
    // |               |
    // |               |
    // |               |
    // +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
    // |               |
    // |               |
    // |               |
    // +---------------+= 0x200 (512) Start of most Chip-8 programs
    // | 0x000 to 0x1FF|
    // | Reserved for  |
    // |  interpreter  |
    // +---------------+= 0x000 (0) Start of Chip-8 RAM
    Memory [4096]byte

    I uint16                // Address register

    V [16]byte              // 16x 8bit registers, V0 to VF

    PC uint16               // Program counter

    Stack [16]uint16        // Stack

    SP byte                 // Stack pointer

    Clock <-chan time.Time  // CHIP-8 is normally 60hz, this channel supplies ticks

    stop chan struct{}      // Shutdown the cpu

    TickCount int

    gfx []byte              // Graphics memory

    VideoDriver *Video            // Video driver
}

func NewCpu () *Chip8 {
    c := new(Chip8)

    // Initialize program counter
    c.PC = 0x200

    // Clear graphics
    c.gfx = make([]byte, screenWidth*screenHeight)

    // Set the fontset on each memory location
    // TODO: This is stolen, why do we need to do this? Cant it all be 0x0?
    for i := 0; i < 80; i++ {
        c.Memory[i] = chip8Fontset[i]
    }

    // Create the stop channel
    c.stop = make(chan struct{})

    c.Clock = time.Tick(time.Second / time.Duration(cpuFrequency))

    return c
}

// LoadBytes loads the bytes into memory.
func (c *Chip8) LoadBytes(p []byte) (int, error) {
    return c.load(0x200, bytes.NewReader(p))
}

func (c *Chip8) load(offset int, r io.Reader) (int, error) {
    return r.Read(c.Memory[offset:])
}

func (c* Chip8) Stop () {
    close(c.stop)
}

func (c* Chip8) Run () error {

    //go func () {
    //    for {
    //        select {
    //        case <-inputStream:
    //            fmt.Println("INPUT!")
    //        }
    //    }
    //}()

    c.VideoDriver = new(Video)
    go func (c *Chip8) {
        err := c.VideoDriver.Init(&c.gfx)
        if err != nil {
            panic(err)
        }
    }(c)

    for {
        select {
        case <-c.stop:
            // The cpu has stopped, exit the loop
            return nil

        case <-c.Clock:
            // Every 10 ticks we put in a header
            if c.TickCount == 0 {
                if !opcodesIdentical {
                    fmt.Println(fmt.Sprintf("%s     %s     %s", "PC", "Opcode", "Action"))
                }
            }
            if c.TickCount == 10 {
                c.TickCount = 0
            } else {
                c.TickCount += 1
            }

            // Do a cycle/step
            err := c.Cycle()
            if err != nil {
                // This shouldn't happen yet, but if the cycle errors, abort
                fmt.Println(err.Error())

                return nil
            }
        }
    }

    return nil
}

func main() {
    spew.Println("Starting")

    cpu := NewCpu()

    //program, err := ioutil.ReadFile("/Users/jamesblizzard/Downloads/c8games/BLINKY")
    //program, err := ioutil.ReadFile("/Users/jamesblizzard/Downloads/ibm")
    program, err := ioutil.ReadFile("/Users/jamesblizzard/Downloads/c8games/BLITZ")
    if err != nil {
        panic(err.Error())
    } else {

        //fmt.Println("--- DUMPING ROM ---")
        //spew.Dump(program)
        //fmt.Println("--- END OF ROM ---")

        _, err = cpu.LoadBytes(program)
        if err != nil {
            panic(err.Error())
        }

        //fmt.Println("--- MEMORY, AFTER ROM LOAD ---")
        //spew.Dump(cpu.Memory)
        //fmt.Println("--- END OF MEMORY, AFTER ROM LOAD ---")

        err = cpu.Run()
        if err != nil {
            fmt.Println(err.Error())
        } else {
            fmt.Println("Clean exit")
        }
    }
}
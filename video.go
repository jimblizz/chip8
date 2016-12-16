package main

import (
    "fmt"
    "github.com/go-gl/gl/v2.1/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
    "runtime"
)

const (
    // Multiplies screen size by 10.
    pixelSize = 10
)

type Video struct {
    Window *glfw.Window
}

func (v *Video) Init(gfxMemory *[]byte) error {

    // we need a parallel OS thread to avoid audio stuttering
    runtime.GOMAXPROCS(2)

    // we need to keep OpenGL calls on a single thread
    runtime.LockOSThread()

    if err := glfw.Init(); err != nil {
        return err
    }
    defer glfw.Terminate()

    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 2)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)

    var err error
    v.Window, err = glfw.CreateWindow(screenWidth*pixelSize, screenHeight*pixelSize, "CHIP-8 Emulator", nil, nil)
    if err != nil {
        return err
    }

    v.Window.MakeContextCurrent()

    // Enable vertical sync on cards that support it.
    glfw.SwapInterval(1)

    if err := gl.Init(); err != nil {
        return err
    }

    glerror := gl.GetError()
    if glerror != gl.NO_ERROR {
        fmt.Println("OPENGL ERROR")
        panic(glerror)
    }

    gl.ClearColor(0, 0, 0, 0)
    gl.MatrixMode(gl.PROJECTION)

    // Change coordinates to range from [0, 64] and [0,32].
    gl.Ortho(0, screenWidth, screenHeight, 0, 0, 1)


    // Unnecessary sanity check. :-P
    //if glfw.WindowParam(glfw.Opened) == 0 {
    //    return fmt.Errorf("No window opened")
    //}

    for !v.Window.ShouldClose() {
        pixels := *gfxMemory
        v.draw(pixels)
        glfw.PollEvents()
    }

    return nil
}

func (v *Video) quit() {
    gl.End()
    glfw.Terminate()
}

func (v *Video) close() {
    //glfw.CloseWindow()
}

func (v *Video) draw(pixels []byte) {
    // No need to clear the screen since I explicitly redraw all pixels, at
    // least currently.

    gl.MatrixMode(gl.POLYGON)

    for yline := 0; yline < screenHeight; yline++ {

        for xline := 0; xline < screenWidth; xline++ {

            x, y := float32(xline), float32(yline)
            if pixels[xline+yline*64] == 0 {
                gl.Color3f(0, 0, 0)
            } else {
                gl.Color3f(1, 1, 1)
            }
            gl.Rectf(x, y, x+1, y+1)
        }
    }
    v.Window.SwapBuffers()
}
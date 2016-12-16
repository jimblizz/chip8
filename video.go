package main

import "github.com/go-gl/glfw/v3.2/glfw"

const (
    pixelSize = 10 // Zoom!
)

type video struct {
}

func (v *video) init() error  {
    if err := glfw.Init(); err != nil {
        return err
    }

    return nil
}
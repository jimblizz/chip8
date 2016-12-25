package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

var keyMap = map[glfw.Key]byte{
	glfw.Key1:     0x1,
	glfw.Key2:     0x2,
	glfw.Key3:     0x3,
	glfw.Key4:     0xc,
	glfw.KeyQ:     0x4,
	glfw.KeyW:     0x5,
	glfw.KeyE:     0x6,
	glfw.KeyR:     0xd,
	glfw.KeyA:     0x7,
	glfw.KeyS:     0x8,
	glfw.KeyD:     0x9,
	glfw.KeyF:     0xe,
	glfw.KeyZ:     0xa,
	glfw.KeyX:     0x0,
	glfw.KeyC:     0xb,
	glfw.KeyV:     0xf,
	glfw.KeyUp:    0x2,
	glfw.KeyLeft:  0x4,
	glfw.KeyRight: 0x6,
}

var lastKeyPressed byte

func keyPressed(_ *glfw.Window, k glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	if action == glfw.Press {
		if val, ok := keyMap[k]; ok {
			setKey(val)
		}
	}
}

func setKey(key byte) {
	lastKeyPressed = key
}

func getKey() byte {
	old := lastKeyPressed
	lastKeyPressed = 0
	return old
}

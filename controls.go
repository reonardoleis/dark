package main

import (
	"github.com/tfriedel6/canvas/sdlcanvas"
)

var controls *Controls

type keysDown = map[Key]bool
type keysPressed = map[Key]bool

type Controls struct {
	keysDown    keysDown
	keysPressed keysPressed
}

type Key string

const (
	KeyW         = "KeyW"
	KeyA         = "KeyA"
	KeyS         = "KeyS"
	KeyD         = "KeyD"
	KeyQ         = "KeyQ"
	KeyE         = "KeyE"
	KeyM         = "KeyM"
	KeyLeftShift = "ShiftLeft"
	KeyU         = "KeyU"
	KeyJ         = "KeyJ"
	KeyF         = "KeyF"
	KeyR         = "KeyR"
)

func newControls() {
	keysDown := make(map[Key]bool)
	keysPressed := make(map[Key]bool)

	controls = &Controls{keysDown, keysPressed}
}

func GetControls() *Controls {
	return controls
}

func (c *Controls) Bind(window *sdlcanvas.Window) {
	window.KeyUp = c.handleKeyUp
	window.KeyDown = c.handleKeyDown
}

func (c *Controls) handleKeyDown(scancode int, rn rune, name string) {
	c.keysDown[Key(name)] = true
	c.keysPressed[Key(name)] = true
}

func (c *Controls) handleKeyUp(scancode int, rn rune, name string) {
	c.keysDown[Key(name)] = false
}

func (c *Controls) IsKeyDown(key Key) bool {
	down, ok := c.keysDown[key]
	return ok && down
}

func (c *Controls) IsKeyPressed(key Key) bool {
	pressed, ok := c.keysPressed[key]

	if ok && pressed {
		delete(c.keysPressed, key)
		return true
	}

	return false
}

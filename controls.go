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
	mouseDeltaX int
	mouseDeltaY int
	window      *sdlcanvas.Window
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
	KeyM1        = "KeyM1"
	Key1         = "Digit1"
	Key2         = "Digit2"
	Key3         = "Digit3"
	Key4         = "Digit4"
	KeyI         = "KeyI"
	KeyP         = "KeyP"
)

func newControls() {
	keysDown := make(map[Key]bool)
	keysPressed := make(map[Key]bool)
	controls = &Controls{keysDown, keysPressed, 0, 0, nil}
}

func GetControls() *Controls {
	return controls
}

func (c *Controls) Bind(window *sdlcanvas.Window) {
	window.KeyUp = c.handleKeyUp
	window.KeyDown = c.handleKeyDown
	window.MouseDown = c.handleMouseDown
	window.MouseUp = c.handleMouseUp
	window.MouseMove = c.handleMouseMove
	c.window = window
}

func (c *Controls) handleMouseMove(x, y int) {
	c.mouseDeltaX = x
	c.mouseDeltaY = y
}

func (c *Controls) handleMouseDown(button int, x, y int) {
	key := KeyM1
	c.keysPressed[Key(key)] = true
}

func (c *Controls) handleMouseUp(button int, x, y int) {
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

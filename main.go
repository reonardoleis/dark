package main

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

func main() {
	wnd, drawingCanvas, err := sdlcanvas.CreateWindow(1280, 720, "dark")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	newControls()
	GetControls().Bind(wnd)
	loadTextures(drawingCanvas)

	displayCanvas := canvas.New(wnd.Backend)

	drawingCanvas.LoadFont("./assets/font.ttf")
	manager := newManager(drawingCanvas, displayCanvas)

	wnd.MainLoop(manager.loop)
}

package main

import (
	"time"

	"github.com/tfriedel6/canvas"
)

var (
	DELTA_TIME = 0.0
)

type Manager struct {
	currentFrameTime int64
	lastFrameTime    int64
	canvas           *canvas.Canvas
	displayCanvas    *canvas.Canvas
	player           *Player
	dungeon          *Dungeon
	renderer3d       *Renderer3D
	renderer2d       *Renderer2D
}

func newManager(canvas, displayCanvas *canvas.Canvas) *Manager {
	initialPos := newVector2(22, 22)

	cameraPlane := newVector2(0, 0.66)
	cameraDirection := newVector2(-1, 0)
	camera := newCamera(cameraPlane, cameraDirection)

	player := newPlayer(initialPos, camera)

	dungeon := newDungeon(mapWidth, mapHeight)
	dungeon.generate(generateSimpleLevel)

	renderer3d := newRenderer3D(canvas, newRaycaster())
	renderer2d := newRenderer2D(canvas)

	return &Manager{
		currentFrameTime: 0,
		lastFrameTime:    0,
		canvas:           canvas,
		displayCanvas:    displayCanvas,
		player:           player,
		dungeon:          dungeon,
		renderer3d:       renderer3d,
		renderer2d:       renderer2d,
	}
}

func (m *Manager) loop() {
	m.lastFrameTime = m.currentFrameTime
	m.currentFrameTime = time.Now().UnixMilli()

	frameTime := m.currentFrameTime - m.lastFrameTime

	m.canvas.ClearRect(0, 0, screenWidth, screenHeight)
	m.displayCanvas.ClearRect(0, 0, 1280, 720)
	m.renderer3d.RenderFloor(m.player, m.player.camera)
	m.renderer3d.RenderWalls(m.player, m.player.camera, m.dungeon)
	//m.renderer2d.RenderMap(m.player, m.dungeon)
	m.renderer3d.RenderPlayer(m.player)
	m.player.Controls(m.dungeon, 1.0)
	m.canvas.SetFont(nil, 25)

	DELTA_TIME = float64(frameTime) / 1000

	frame := m.canvas.GetImageData(0.0, 0.0, m.canvas.Width(), m.canvas.Height())
	m.displayCanvas.DrawImage(
		frame,
		0.0, 0.0, screenWidth, screenHeight,
		0.0, 0.0, 1280, 720,
	)

}
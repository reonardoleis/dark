package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

var (
	DELTA_TIME = 0.0
)

type Manager struct {
	currentFrameTime     int64
	lastFrameTime        int64
	canvas               *canvas.Canvas
	displayCanvas        *canvas.Canvas
	player               *Player
	dungeon              *Dungeon
	renderer3d           *Renderer3D
	renderer2d           *Renderer2D
	window               *sdlcanvas.Window
	npcs                 []*NPC
	props                []*Prop
	spawner              Spawner
	informationPanelOpen bool
}

func generateProps(dungeon *Dungeon) []*Prop {
	props := make([]*Prop, 0)
	for y, line := range dungeon.grid {
		for x, entity := range line {
			if entity.isFloor() {
				random := rand.Intn(101)
				if random >= 90 {
					texture := scenarioPropTextures[rand.Intn(len(scenarioPropTextures))]
					propPosition := newVector2(float64(x)+0.5, float64(y)-0.5)
					divX, divY := getPropsDivs(texture)
					props = append(props, newProp(propPosition, texture, PropPositionCeiling, divX, divY))
				}
			}
		}
	}

	return props
}

func newManager(window *sdlcanvas.Window, canvas *canvas.Canvas, displayCanvas *canvas.Canvas) *Manager {
	cameraPlane := newVector2(0, 0.66)
	cameraDirection := newVector2(-1, 0)
	camera := newCamera(cameraPlane, cameraDirection)

	dungeon := newDungeon(mapWidth, mapHeight, 1)
	dungeon.generate(generateLevel)

	renderer3d := newRenderer3D(canvas, newRaycaster())
	renderer2d := newRenderer2D(canvas)

	emptyPosition := [2]int{0, 0}
	for y := range dungeon.grid {
		for x := range dungeon.grid[y] {
			if dungeon.At(x, y) != nil && dungeon.At(x, y).isFloor() {
				emptyPosition = [2]int{y, x}
				break
			}
		}
	}

	player := newPlayer(newVector2(float64(emptyPosition[1]), float64(emptyPosition[0])), camera)

	var npcs []*NPC
	spawner := newSpawner(dungeon, &npcs, player)
	spawner.initialSpawn()

	return &Manager{
		currentFrameTime: 0,
		lastFrameTime:    0,
		canvas:           canvas,
		displayCanvas:    displayCanvas,
		player:           player,
		dungeon:          dungeon,
		renderer3d:       renderer3d,
		renderer2d:       renderer2d,
		window:           window,
		npcs:             npcs,
		props:            generateProps(dungeon),
		spawner:          spawner,
	}
}

func (m *Manager) cleanupDeadNpcs() {
	for idx := range m.npcs {
		if !m.npcs[idx].alive {
			m.npcs = append(m.npcs[:idx], m.npcs[idx+1:]...)
		}
	}
}

func (m *Manager) loop() {
	m.lastFrameTime = m.currentFrameTime
	m.currentFrameTime = time.Now().UnixMilli()

	frameTime := m.currentFrameTime - m.lastFrameTime

	m.canvas.ClearRect(0, 0, screenWidth, screenHeight)
	m.displayCanvas.ClearRect(0, 0, 1280, 720)
	m.renderer3d.RenderFloorAndRoof(m.dungeon, m.player, m.player.camera)
	m.renderer3d.RenderWalls(m.player, m.player.camera, m.dungeon)
	m.renderer3d.RenderNpcs(m.player, m.npcs)
	m.renderer3d.RenderProps(m.player, m.props)
	//m.renderer2d.RenderMap(m.player, m.dungeon)
	m.renderer2d.RenderNpcsHealthbars(m.npcs, m.player)
	m.renderer3d.RenderPlayer(m.player)
	ComputeNpcEvents(m.npcs, m.player)
	m.renderer2d.RenderHealthbar(m.player)

	if m.player.alive {
		m.player.Controls(m.dungeon, 1.0)
		m.renderer2d.RenderCombatLog()
		m.player.SpendPoints()

		if m.player.points > 0 {
			m.renderer2d.RenderSpendPointsText(m.player)
		}

		if m.player.enemiesKilled >= minEnemiesToProgress {
			m.canvas.SetFillStyle(255, 255, 255)
			m.canvas.SetFont(nil, 15)
			m.canvas.FillText("Press R to progress to next level", screenWidth-300.0, 46.0)
			if GetControls().IsKeyDown(KeyR) {
				m.player.enemiesKilled = 0
				m.dungeon.level++
				m.dungeon.generate(generateLevel)
				emptyPosition := [2]int{0, 0}
				for y := range m.dungeon.grid {
					for x := range m.dungeon.grid[y] {
						if m.dungeon.At(x, y) != nil && m.dungeon.At(x, y).isFloor() {
							emptyPosition = [2]int{y, x}
							break
						}
					}
				}

				m.player.position.x = float64(emptyPosition[1])
				m.player.position.y = float64(emptyPosition[0])

				m.npcs = make([]*NPC, 0)
				spawner := newSpawner(m.dungeon, &m.npcs, m.player)
				spawner.initialSpawn()

			}
		}
	}

	m.renderer2d.RenderProgressText(m.player)

	if GetControls().IsKeyDown(KeyM) {
		os.Exit(0)
	}

	if GetControls().IsKeyPressed(KeyE) {
		m.informationPanelOpen = !m.informationPanelOpen
	}

	m.renderer2d.RenderExpBar(m.player)

	if !m.player.alive {
		m.canvas.SetFillStyle(50, 50, 50, 0.5)
		m.canvas.FillRect(0, 0, screenWidth, screenHeight)
		m.canvas.SetFillStyle(255, 255, 255)
		fontSize := 25.0
		m.canvas.SetFont(nil, fontSize)
		fontWidth := fontSize * (4. / 3.)
		text := "You are dead."
		textWidth := len(text) * int(fontWidth) / 2
		textStartX := screenWidth/2 - textWidth/2
		textStartY := screenHeight/2 - 10

		m.canvas.FillText(text, float64(textStartX), float64(textStartY))
	}

	if m.informationPanelOpen {
		m.renderer2d.RenderInformationPanel(m.player, m.dungeon)
	}

	m.cleanupDeadNpcs()
	m.player.currentEnemyNpcTarget = m.npcs[len(m.npcs)-1]

	for _, npc := range m.npcs {
		npc.ChasePlayer(m.player, m.dungeon)
	}

	ww, wh := m.window.Size()
	m.window.Window.WarpMouseInWindow(int32(ww/2), int32(wh/2))

	DELTA_TIME = float64(frameTime) / 1000

	frame := m.canvas.GetImageData(0.0, 0.0, m.canvas.Width(), m.canvas.Height())
	m.displayCanvas.DrawImage(
		frame,
		0.0, 0.0, screenWidth, screenHeight,
		0.0, 0.0, 1280, 720,
	)

}

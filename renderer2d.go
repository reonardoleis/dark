package main

import (
	"fmt"
	"math"

	"github.com/tfriedel6/canvas"
)

type Renderer2D struct {
	canvas *canvas.Canvas
}

func newRenderer2D(canvas *canvas.Canvas) *Renderer2D {
	return &Renderer2D{
		canvas: canvas,
	}
}

func (r *Renderer2D) RenderMap(player *Player, dungeon *Dungeon) {
	for y := 0; y < len(dungeon.grid); y++ {
		for x := 0; x < len(dungeon.grid[y]); x++ {
			if dungeon.grid[y][x] != nil {
				if dungeon.grid[y][x].isWall() {
					r.canvas.SetFillStyle(255, 0, 0)
				} else {
					r.canvas.SetFillStyle(0, 0, 255)
				}

				if dungeon.grid[y][x].debugId == 4 {
					r.canvas.SetFillStyle(255, 255, 0)
				} else if dungeon.grid[y][x].debugId == 5 {
					r.canvas.SetFillStyle(255, 0, 255)
				} else if dungeon.grid[y][x].debugId == 6 {
					r.canvas.SetFillStyle(80, 90, 150)
				}

				r.canvas.FillRect(float64(x)*cellScaleX, float64(y)*cellScaleY, cellScaleX, cellScaleY)
			}
		}
	}

	r.canvas.SetFillStyle(0, 255, 0)
	r.canvas.FillRect(player.position.x*cellScaleX, player.position.y*cellScaleY, cellScaleX, cellScaleY)

}

func (r *Renderer2D) RenderHealthbar(player *Player) {
	br, bg, bb := 255, 255, 255
	ir, ig, ib := 255, 0, 0

	width := 200.
	height := 16.
	borderRadius := 5.

	positionX := 20.
	positionY := float64(screenHeight) - height - 20.

	r.canvas.SetFillStyle(br, bg, bb)
	r.canvas.FillRect(positionX, positionY, width, height)

	r.canvas.SetFillStyle(0, 0, 0)
	r.canvas.FillRect(positionX+borderRadius/2, positionY+borderRadius/2, width-borderRadius, height-borderRadius)

	r.canvas.SetFillStyle(ir, ig, ib)
	percentage := float64(player.stats.Hp) / float64(player.stats.MaxHp)
	r.canvas.FillRect(positionX+borderRadius/2, positionY+borderRadius/2, (width-borderRadius)*percentage, height-borderRadius)
}

func (r *Renderer2D) RenderNpcsHealthbars(npcs []*NPC, player *Player) {
	for _, npc := range npcs {
		if npc.distance < enemyAttackDistance && npc.id == player.currentEnemyNpcTarget.id {
			r.canvas.SetFillStyle(0, 0, 0)
			r.canvas.FillRect((screenWidth/2)-100, 10, 200, 20)
			r.canvas.SetFillStyle(255, 0, 0)

			maxSizeX := 200.0
			sizeX := maxSizeX * (float64(npc.stats.Hp) / float64(npc.stats.MaxHp))
			r.canvas.FillRect((screenWidth/2)-maxSizeX/2, 10, sizeX, 20)

			fontSize := 12.0
			r.canvas.SetFont(nil, fontSize)
			fontWidth := fontSize * (4. / 3.)

			text := fmt.Sprintf("Lv. %d %s", npc.stats.Level, npc.name)
			textWidth := len(text) * int(fontWidth) / 2
			textStartX := screenWidth/2 - textWidth/2

			r.canvas.SetFillStyle(255, 255, 255)
			r.canvas.FillText(text, float64(textStartX), 22.5)

			return
		}
	}
}

func (r *Renderer2D) RenderCombatLog() {
	r.canvas.SetFillStyle(255, 255, 255)
	fontSize := 12.0

	y := float64(screenHeight) - 80

	r.canvas.SetFont(nil, fontSize)

	for _, log := range combatLog {
		fontWidth := fontSize * (4. / 3.)
		textWidth := len(log) * int(fontWidth) / 2
		textStartX := float64(screenWidth - textWidth - 32.0)
		r.canvas.FillText(log, textStartX, y)
		y += 10.0
	}
}

func (r *Renderer2D) RenderExpBar(player *Player) {
	br, bg, bb := 255, 255, 255
	ir, ig, ib := 255, 255, 0

	width := 200.
	height := 16.
	borderRadius := 5.

	positionX := screenWidth - width - 20.
	positionY := float64(screenHeight) - height - 20.

	r.canvas.SetFillStyle(br, bg, bb)
	r.canvas.FillRect(positionX, positionY, width, height)

	r.canvas.SetFillStyle(0, 0, 0)
	r.canvas.FillRect(positionX+borderRadius/2, positionY+borderRadius/2, width-borderRadius, height-borderRadius)

	r.canvas.SetFillStyle(ir, ig, ib)
	percentage := float64(player.stats.Exp) / float64(player.stats.NextLevel)
	r.canvas.FillRect(positionX+borderRadius/2, positionY+borderRadius/2, (width-borderRadius)*percentage, height-borderRadius)

}

func (r *Renderer2D) RenderInformationPanel(player *Player, dungeon *Dungeon) {
	startX := 25.0
	startY := 25.0
	endX := screenWidth - 25.0*2
	endY := screenHeight - 25.0*2

	r.canvas.SetFillStyle(255, 255, 255)
	r.canvas.FillRect(startX, startY, endX, endY)

	r.canvas.SetFillStyle(200, 200, 200)
	r.canvas.FillRect(startX+3.0, startY+3.0, endX-6.0, endY-6.0)

	r.canvas.SetFillStyle(0.0, 0.0, 0.0)
	r.canvas.SetFont(nil, 12)

	statsTexts := player.stats.String()

	startY = 40.0
	for _, statText := range statsTexts {
		r.canvas.FillText(statText, 32.0, startY)
		startY += 10.0
	}

	attributesText := player.attributes.String()
	for _, attributeText := range attributesText {
		r.canvas.FillText(attributeText, 32.0, startY)
		startY += 10.0
	}

	startY = 40.0
	powerupsTexts := player.powerups.String()
	for _, powerupsText := range powerupsTexts {
		r.canvas.FillText(powerupsText, endX-200, startY)
		startY += 10.0
	}

	r.renderMap(player, dungeon)
}

func (r *Renderer2D) RenderSpendPointsText(player *Player) {
	text1 := fmt.Sprintf("You have %d unspent points", player.points)

	r.canvas.SetFillStyle(255, 255, 255)
	r.canvas.SetFont(nil, 15)

	r.canvas.FillText(text1, 36, 36)
	r.canvas.FillText("Press 1 to upgrade STR", 46, 56)
	r.canvas.FillText("Press 2 to upgrade AGI", 46, 76)
	r.canvas.FillText("Press 3 to upgrade VIG", 46, 96)
	r.canvas.FillText("Press 4 to upgrade INT", 46, 116)
}

func (r *Renderer2D) renderMap(player *Player, dungeon *Dungeon) {
	playerX := int(math.Floor(player.position.x))
	playerY := int(math.Floor(player.position.y))

	startX := playerX - mapDisplaySize/2
	endX := playerX + mapDisplaySize/2

	if startX < 0 {
		startX = 0
		endX = mapDisplaySize
	} else if endX >= mapWidth {
		endX = mapWidth - 1
		startX = endX - mapDisplaySize
	}

	startY := playerY - mapDisplaySize/2
	endY := playerY + mapDisplaySize/2

	if startY < 0 {
		startY = 0
		endY = mapDisplaySize
	} else if endY >= mapHeight {
		endY = mapHeight - 1
		startY = endY - mapDisplaySize
	}

	offsetX := 200.0
	offsetY := 42.0

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			_, image := dungeon.At(x, y).getTexture1()

			drawX := offsetX + (float64(x-startX))*cellScaleX
			drawY := offsetY + (float64(y-startY))*cellScaleY
			r.canvas.DrawImage(image, 0.0, 0.0, float64(image.Bounds().Dx()), float64(image.Bounds().Dy()),
				drawX, drawY, cellScaleX, cellScaleY,
			)
		}
	}

	drawPlayerX := offsetX + (float64(mapDisplaySize/2))*cellScaleX
	drawPlayerY := offsetY + (float64(mapDisplaySize/2))*cellScaleY
	r.canvas.SetFillStyle(0, 255, 0)
	r.canvas.FillRect(drawPlayerX, drawPlayerY, cellScaleX, cellScaleY)
}

func (r *Renderer2D) RenderProgressText(player *Player) {
	r.canvas.SetFillStyle(255, 255, 255)
	r.canvas.SetFont(nil, 11)

	x := screenWidth / 2
	y := screenHeight - 32

	r.canvas.FillText(fmt.Sprintf("Kill %d more enemies to progress", max(0, minEnemiesToProgress-player.enemiesKilled)), float64(x)-95, float64(y)+7)

}

func (r *Renderer2D) RenderDungeonLevel(dungeon *Dungeon) {
	r.canvas.SetFillStyle(255, 255, 255)
	r.canvas.SetFont(nil, 13)

	r.canvas.FillText(fmt.Sprintf("Dungeon level %d", dungeon.level), 32.0, 12.0)
}

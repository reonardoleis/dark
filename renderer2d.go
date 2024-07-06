package main

import "github.com/tfriedel6/canvas"

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
				r.canvas.FillRect(float64(x)*cellScaleX, float64(y)*cellScaleY, cellScaleX, cellScaleY)
			}
		}
	}

	r.canvas.SetFillStyle(0, 255, 0)
	r.canvas.FillRect(player.position.x*cellScaleX, player.position.y*cellScaleY, cellScaleX, cellScaleY)

}

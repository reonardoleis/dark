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
	r.canvas.SetFillStyle(255, 0, 0)
	for y := 0; y < len(dungeon.grid); y++ {
		for x := 0; x < len(dungeon.grid[y]); x++ {
			if dungeon.grid[y][x] != nil {
				r.canvas.FillRect(float64(x)*cellScaleX, float64(y)*cellScaleY, cellScaleX, cellScaleY)
			}
		}
	}

	r.canvas.SetFillStyle(0, 255, 0)
	r.canvas.FillRect(player.position.x*cellScaleX, player.position.y*cellScaleY, cellScaleX, cellScaleY)

}

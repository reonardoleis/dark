package main

import (
	"image"
	"image/color"
	"math"

	"github.com/tfriedel6/canvas"
)

type Renderer3D struct {
	canvas    *canvas.Canvas
	raycaster *Raycaster
}

func newRenderer3D(canvas *canvas.Canvas, raycaster *Raycaster) *Renderer3D {
	return &Renderer3D{
		canvas:    canvas,
		raycaster: raycaster,
	}
}

func (r *Renderer3D) RenderWalls(player *Player, camera *Camera, dungeon *Dungeon) {
	rays := r.raycaster.ProcessConcurrently(dungeon, player, camera)
	for x, ray := range rays {

		texture := ray.entity.getTexture()
		wallHeight := float64(screenHeight) / ray.perpendicularWallDistace
		start := float64(screenHeight)/2 - float64(wallHeight)/2

		percentage := ray.euclideanDistance / float64(mapWidth) * darkenScale

		color := 255 * (percentage)
		if color < 0 {
			color = 0
		}
		if color > 255 {
			color = 255
		}
		var wallX float64
		if ray.side == 0 {
			wallX = player.position.y + ray.perpendicularWallDistace*ray.direction.y
		} else {
			wallX = player.position.x + ray.perpendicularWallDistace*ray.direction.x
		}
		wallX -= math.Floor(wallX)
		texX := int(wallX * float64(texture.Width()))
		if ray.side == 0 && ray.direction.x > 0 {
			texX = texture.Width() - texX - 1
		}
		if ray.side == 1 && ray.direction.y < 0 {
			texX = texture.Width() - texX - 1
		}
		if texX < 0 {
			texX = -texX
		}

		texX %= texture.Width()
		r.canvas.DrawImage(
			texture,
			float64(texX), 0.0, 1.0, float64(texture.Height()),
			float64(x), start, 1.0, wallHeight,
		)

		r.canvas.SetFillStyle(0, 0, 0, int(color))
		r.canvas.FillRect(
			float64(x), start, 1.0, wallHeight,
		)
	}

}

var (
	pLookup = make([]float64, screenHeight)
)

func init() {
	for idx := range pLookup {
		pLookup[idx] = -1
	}
}

func (r Renderer3D) RenderFloor(player *Player, camera *Camera) {
	floorTexture := textures[Wall].rgba
	ceilingTexture := textures[Wall].rgba
	floor := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))
	ceiling := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))
	floorShading := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))
	ceilingShading := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))

	posZ := 0.5 * float64(screenHeight)

	leftmostRay, rightmostRay := r.raycaster.GetRayDirections(camera)
	for y := screenHeight / 2; y < screenHeight; y++ {

		p := pLookup[y]
		if p == -1 {
			p = float64(y) - float64(screenHeight)/2
			pLookup[y] = p
		}

		rowDistance := posZ / p

		floorStepX := rowDistance * (rightmostRay.x - leftmostRay.x) / float64(screenWidth)
		floorStepY := rowDistance * (rightmostRay.y - leftmostRay.y) / float64(screenWidth)

		floorX := player.position.x + rowDistance*leftmostRay.x
		floorY := player.position.y + rowDistance*leftmostRay.y

		initialDistance := math.Sqrt(math.Pow(floorX-player.position.x, 2) + math.Pow(floorY-player.position.y, 2))
		for x := 0; x < screenWidth; x++ {
			texX := int(floorX*float64(floorTexture.Bounds().Dx())) % floorTexture.Bounds().Dy()
			texY := int(floorY*float64(floorTexture.Bounds().Dy())) % floorTexture.Bounds().Dy()

			if texX < 0 {
				texX += floorTexture.Bounds().Dx()
			}
			if texY < 0 {
				texY += floorTexture.Bounds().Dy()
			}

			floorColor := floorTexture.At(texX, texY)
			ceilingColor := ceilingTexture.At(texX, texY)

			if (int(floorX)+int(floorY))%2 == 0 {
				r, g, b, a := floorColor.RGBA()
				r = uint32(float64(r) * 0.5)
				g -= uint32(float64(g) * 0.5)
				b -= uint32(float64(b) * 0.5)
				floorColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			}
			if (int(floorX)+int(floorY))%2 == 1 {
				r, g, b, a := ceilingColor.RGBA()
				r = uint32(float64(r) * 0.5)
				g -= uint32(float64(g) * 0.5)
				b -= uint32(float64(b) * 0.5)
				ceilingColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			}

			ratio := float64(screenHeight) / float64(screenWidth)
			distance := initialDistance + math.Abs(float64(x)-float64(screenWidth/2))/float64(screenWidth)
			percentage := distance / float64(mapWidth) * darkenScale * ratio
			darkenFactor := percentage
			if percentage > 1 {
				darkenFactor = 1.0
			}
			floorX += floorStepX
			floorY += floorStepY

			r, g, b, a := floorColor.RGBA()
			a = uint32(float64(a) - (255 * darkenFactor))
			floorColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			floor.Set(x, y, floorColor)
			r, g, b, a = ceilingColor.RGBA()

			a = uint32(float64(a) - (255 * darkenFactor))
			ceilingColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			ceiling.Set(x, screenHeight-y, ceilingColor)
		}
	}
	_ = floorShading
	_ = ceilingShading
	r.canvas.DrawImage(floor)
	//r.canvas.DrawImage(floorShading)
	r.canvas.DrawImage(ceiling)
	//r.canvas.DrawImage(ceilingShading)
}

func (r *Renderer3D) RenderPlayer(player *Player) {
	texture := textures[Sword].image

	r.canvas.DrawImage(
		texture, 0.0, 0.0, float64(texture.Width()), float64(texture.Height()),
		screenWidth-float64(texture.Width())*1.5, screenHeight/2, float64(texture.Width()), float64(texture.Height()),
	)
}

//  DrawImage("image", dx, dy)
//  DrawImage("image", dx, dy, dw, dh)
//  DrawImage("image", sx, sy, sw, sh, dx, dy, dw, dh)
// Where dx/dy/dw/dh are the destination coordinates and sx/sy/sw/sh are the
// source coordinates

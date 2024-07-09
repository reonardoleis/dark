package main

import (
	"image"
	"image/color"
	"math"
	"sort"

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

		texture, _ := ray.entity.getTexture1()
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

type floorAndRoofWorkerIn struct {
	player         *Player
	camera         *Camera
	dungeon        *Dungeon
	floorTexture   *image.Image
	ceilingTexture *image.Image
	y              int
	leftmostRay    Vector2
	rightmostRay   Vector2
	floor          *image.RGBA
	ceiling        *image.RGBA
}

func (r Renderer3D) floorAndRoofWorker(in chan floorAndRoofWorkerIn, out chan int) {
	posZ := 0.5 * float64(screenHeight)

	for v := range in {
		y := v.y
		leftmostRay := v.leftmostRay
		rightmostRay := v.rightmostRay
		player := v.player
		floor := v.floor
		ceiling := v.ceiling
		floorTexture := *v.floorTexture
		ceilingTexture := *v.ceilingTexture

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
			gridX, gridY := int(floorX), int(floorY)
			if v.dungeon.At(gridX, gridY) != nil {
				_, floorTexture = v.dungeon.At(gridX, gridY).getTexture1()
				_, ceilingTexture = v.dungeon.At(gridX, gridY).getTexture2()
			}
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

		out <- 1
	}
}

func (r Renderer3D) RenderFloorAndRoof(dungeon *Dungeon, player *Player, camera *Camera) {
	floorTexture := textures[Dirtmud].rgba
	ceilingTexture := textures[Wall].rgba
	floor := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))
	ceiling := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))

	in := make(chan floorAndRoofWorkerIn, screenHeight/2)
	out := make(chan int, screenHeight/2)

	numWorkers := 8
	for range numWorkers {
		go r.floorAndRoofWorker(in, out)
	}

	leftmostRay, rightmostRay := r.raycaster.GetRayDirections(camera)
	for y := screenHeight / 2; y < screenHeight; y++ {
		in <- floorAndRoofWorkerIn{
			player:         player,
			camera:         camera,
			floorTexture:   &floorTexture,
			ceilingTexture: &ceilingTexture,
			y:              y,
			leftmostRay:    leftmostRay,
			rightmostRay:   rightmostRay,
			floor:          floor,
			ceiling:        ceiling,
			dungeon:        dungeon,
		}
	}

	close(in)
	counter := 0
	for range out {
		counter++
		if counter == screenHeight/2 {
			close(out)
		}
	}

	r.canvas.DrawImage(floor)
	r.canvas.DrawImage(ceiling)
}

func (r Renderer3D) RenderFloorAndRoofSequential(player *Player, camera *Camera) {
	floorTexture := textures[Wall].rgba
	ceilingTexture := textures[Wall].rgba
	floor := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))
	ceiling := image.NewRGBA(image.Rect(0.0, 0.0, screenWidth, screenHeight))

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
	r.canvas.DrawImage(floor)
	//r.canvas.DrawImage(floorShading)
	r.canvas.DrawImage(ceiling)
	//r.canvas.DrawImage(ceilingShading)
}

func (r *Renderer3D) RenderPlayer(player *Player) {
	texture := textures[Sword].image

	r.canvas.Save()

	r.canvas.Rotate(player.currentAttackStep)

	r.canvas.DrawImage(
		texture, 0.0, 0.0, float64(texture.Width()), float64(texture.Height()),
		screenWidth-float64(texture.Width())*4-(screenWidth*0.1), screenHeight/2, float64(texture.Width())*4, float64(texture.Height()*4),
	)

	r.canvas.Restore()
}

//  DrawImage("image", dx, dy)
//  DrawImage("image", dx, dy, dw, dh)
//  DrawImage("image", sx, sy, sw, sh, dx, dy, dw, dh)
// Where dx/dy/dw/dh are the destination coordinates and sx/sy/sw/sh are the
// source coordinates

func (r *Renderer3D) RenderNpcs(player *Player, npcs []*NPC) {

	darkenImage := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	image := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	for _, npc := range npcs {

		npc.distance = ((player.position.x - npc.position.x) * (player.position.x - npc.position.x)) +
			((player.position.y - npc.position.y) * (player.position.y - npc.position.y))

	}

	sort.SliceStable(npcs, func(i, j int) bool {
		return npcs[i].distance > npcs[j].distance
	})

	for _, npc := range npcs {
		spriteX := npc.position.x - player.position.x
		spriteY := npc.position.y - player.position.y

		invDet := 1.0 / (player.camera.plane.x*player.camera.direction.y - player.camera.plane.y*player.camera.direction.x)
		transformX := invDet * (player.camera.direction.y*spriteX - player.camera.direction.x*spriteY)
		transformY := invDet * (-player.camera.plane.y*spriteX + player.camera.plane.x*spriteY)

		spriteScreenX := int((float64(screenWidth) / 2) * (1 + transformX/transformY))
		tex := textures[npc.textureId]
		texWidth := tex.rgba.Bounds().Dx()
		texHeight := tex.rgba.Bounds().Dy()

		vMoveScreen := int(float64(texHeight) / transformY)
		spriteHeight := int(math.Abs(float64(screenHeight)/transformY)) * 2
		drawStartY := -spriteHeight/2 + screenHeight/2 + vMoveScreen
		if drawStartY < 0 {
			drawStartY = 0
		}
		drawEndY := spriteHeight/2 + screenHeight/2 + vMoveScreen
		if drawEndY >= screenHeight {
			drawEndY = screenHeight - 1
		}

		spriteWidth := int(math.Abs(float64(screenHeight)/transformY)) * 2

		drawStartX := -spriteWidth/2 + spriteScreenX
		if drawStartX < 0 {
			drawStartX = 0
		}
		drawEndX := spriteWidth/2 + spriteScreenX
		if drawEndX >= screenWidth {
			drawEndX = screenWidth - 1
		}

		for stripe := drawStartX; stripe < drawEndX; stripe++ {
			texX := int(256*(stripe-(-spriteWidth/2+spriteScreenX))*texWidth/spriteWidth) / 256

			if transformY > 0 && stripe > 0 && stripe < screenWidth && transformY < zBuffer[stripe] {
				for y := drawStartY; y < drawEndY; y++ {
					d := (y-vMoveScreen)*256 - screenHeight*128 + spriteHeight*128
					texY := ((d * texHeight / spriteHeight) / 256)
					c := tex.rgba.At(texX, texY)
					percentage := float64(npc.distance) / 4 / darkenScale
					darkenFactor := percentage
					if percentage > 1 {
						darkenFactor = 1.0
					}

					if r, g, b, a := c.RGBA(); r+g+b != 0 {
						if npc.hitHighlightTimeleft > 0 {
							r = 200
						}
						shade := 255 * darkenFactor
						darkenImage.Set(stripe, y, color.RGBA{uint8(0), uint8(0), uint8(0), uint8(shade)})
						color := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
						image.Set(stripe, y, color)
					}
				}
			}
		}
	}

	r.canvas.DrawImage(image)
	r.canvas.DrawImage(darkenImage)
}

func (r *Renderer3D) RenderProps(player *Player, props []*Prop) {
	image := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	for _, prop := range props {

		prop.distance = ((player.position.x - prop.position.x) * (player.position.x - prop.position.x)) +
			((player.position.y - prop.position.y) * (player.position.y - prop.position.y))

	}

	sort.SliceStable(props, func(i, j int) bool {
		return props[i].distance > props[j].distance
	})

	for _, prop := range props {
		spriteX := prop.position.x - player.position.x
		spriteY := prop.position.y - player.position.y

		invDet := 1.0 / (player.camera.plane.x*player.camera.direction.y - player.camera.plane.y*player.camera.direction.x)
		transformX := invDet * (player.camera.direction.y*spriteX - player.camera.direction.x*spriteY)
		transformY := invDet * (-player.camera.plane.y*spriteX + player.camera.plane.x*spriteY)

		spriteScreenX := int((float64(screenWidth) / 2) * (1 + transformX/transformY))
		tex := textures[prop.textureId]
		texWidth := tex.rgba.Bounds().Dx()
		texHeight := tex.rgba.Bounds().Dy()

		vMoveScreen := int(float64(texHeight)/transformY) * prop.divY
		if prop.propPosition == PropPositionCeiling {
			vMoveScreen = -vMoveScreen
		}
		spriteHeight := int(math.Abs(float64(screenHeight)/transformY)) / prop.divY
		drawStartY := -spriteHeight/2 + screenHeight/2 + vMoveScreen
		if drawStartY < 0 {
			drawStartY = 0
		}
		drawEndY := spriteHeight/2 + screenHeight/2 + vMoveScreen
		if drawEndY >= screenHeight {
			drawEndY = screenHeight - 1
		}

		spriteWidth := int(math.Abs(float64(screenHeight)/transformY)) / prop.divX

		drawStartX := -spriteWidth/2 + spriteScreenX
		if drawStartX < 0 {
			drawStartX = 0
		}
		drawEndX := spriteWidth/2 + spriteScreenX
		if drawEndX >= screenWidth {
			drawEndX = screenWidth - 1
		}

		for stripe := drawStartX; stripe < drawEndX; stripe++ {
			texX := int(256*(stripe-(-spriteWidth/2+spriteScreenX))*texWidth/spriteWidth) / 256

			if transformY > 0 && stripe > 0 && stripe < screenWidth && transformY < zBuffer[stripe] {
				for y := drawStartY; y < drawEndY; y++ {
					d := (y-vMoveScreen)*256 - screenHeight*128 + spriteHeight*128
					texY := ((d * texHeight / spriteHeight) / 256)
					c := tex.rgba.At(texX, texY)
					percentage := float64(prop.distance) / darkenScale * 0.85
					darkenFactor := percentage
					if percentage > 1 {
						darkenFactor = 1.0
					}

					if r, g, b, a := c.RGBA(); r+g+b != 0 {
						color := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a - uint32(255*darkenFactor))}
						image.Set(stripe, y, color)
					}
				}
			}
		}
	}

	r.canvas.DrawImage(image)
}

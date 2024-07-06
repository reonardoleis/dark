package main

import (
	"math"
	"sort"
)

type Raycaster struct{}

type Ray struct {
	direction                Vector2
	deltaDistance            Vector2
	perpendicularWallDistace float64
	hit                      bool
	side                     int
	euclideanDistance        float64
	entity                   *Entity
}

func newRaycaster() *Raycaster {
	return &Raycaster{}
}

func newRay(direction Vector2) Ray {
	return Ray{direction, newVector2(0, 0), 0, false, -1, -1, nil}
}

func (r Raycaster) cast(camera *Camera, cameraX float64) Ray {
	rayDirection := newVector2(
		camera.direction.x+camera.plane.x*cameraX,
		camera.direction.y+camera.plane.y*cameraX,
	)

	ray := newRay(rayDirection)

	var deltaDistanceX, deltaDistanceY float64
	if rayDirection.x != 0 {
		deltaDistanceX = math.Abs(1.0 / rayDirection.x)
	} else {
		deltaDistanceX = 1e30
	}

	if rayDirection.y != 0 {
		deltaDistanceY = math.Abs(1.0 / rayDirection.y)
	} else {
		deltaDistanceY = 1e30
	}

	ray.deltaDistance = newVector2(deltaDistanceX, deltaDistanceY)

	return ray
}

type processWorkerIn struct {
	dungeon *Dungeon
	player  *Player
	camera  *Camera
	x       int
}

type processWorkerOut struct {
	ray Ray
	x   int
}

func (r Raycaster) processWorker(in chan processWorkerIn, out chan processWorkerOut) {
	for v := range in {
		camera := v.camera
		x := v.x
		player := v.player
		dungeon := v.dungeon
		cameraX := float64(2*x)/float64(screenWidth) - 1
		ray := r.cast(camera, cameraX)

		mapX := int(player.position.x)
		mapY := int(player.position.y)

		var stepX, stepY int
		var sideDistanceX, sideDistanceY float64

		if ray.direction.x < 0 {
			stepX = -1
			sideDistanceX = (player.position.x - float64(mapX)) * ray.deltaDistance.x
		} else {
			stepX = 1
			sideDistanceX = (float64(mapX) + 1.0 - player.position.x) * ray.deltaDistance.x
		}

		if ray.direction.y < 0 {
			stepY = -1
			sideDistanceY = (player.position.y - float64(mapY)) * ray.deltaDistance.y
		} else {
			stepY = 1
			sideDistanceY = (float64(mapY) + 1.0 - player.position.y) * ray.deltaDistance.y
		}

		for !ray.hit {
			if sideDistanceX < sideDistanceY {
				sideDistanceX += ray.deltaDistance.x
				mapX += stepX
				ray.side = 0
			} else {
				sideDistanceY += ray.deltaDistance.y
				mapY += stepY
				ray.side = 1
			}

			if dungeon.At(mapX, mapY) != nil && dungeon.At(mapX, mapY).isWall() {
				ray.hit = true
				ray.entity = dungeon.At(mapX, mapY)
			}
		}

		if ray.side == 0 {
			ray.perpendicularWallDistace = sideDistanceX - ray.deltaDistance.x
		} else {
			ray.perpendicularWallDistace = sideDistanceY - ray.deltaDistance.y
		}

		hitX := player.position.x + (sideDistanceX-ray.deltaDistance.x)*ray.direction.x
		hitY := player.position.y + (sideDistanceY-ray.deltaDistance.y)*ray.direction.y
		ray.euclideanDistance = euclideanDistance(player.position.x, player.position.y, hitX, hitY)

		out <- processWorkerOut{
			ray: ray,
			x:   x,
		}

	}
}

func (r Raycaster) ProcessConcurrently(dungeon *Dungeon, player *Player, camera *Camera) []Ray {
	rays := make([]Ray, screenWidth)

	in := make(chan processWorkerIn, screenWidth)
	out := make(chan processWorkerOut, screenWidth)

	numWorkers := 32
	for i := 0; i < numWorkers; i++ {
		go r.processWorker(in, out)
	}

	for x := 0; x < screenWidth; x++ {
		in <- processWorkerIn{
			camera:  camera,
			dungeon: dungeon,
			player:  player,
			x:       x,
		}
	}

	close(in)

	infos := []processWorkerOut{}
	for v := range out {
		infos = append(infos, v)
		if len(infos) == screenWidth {
			close(out)
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].x < infos[j].x
	})

	for idx, info := range infos {
		rays[idx] = info.ray
	}

	return rays
}

func (r Raycaster) Process(dungeon *Dungeon, player *Player, camera *Camera) []Ray {
	rays := make([]Ray, screenWidth)

	for x := 0; x < screenWidth; x++ {
		cameraX := float64(2*x)/float64(screenWidth) - 1
		ray := r.cast(camera, cameraX)

		mapX := int(player.position.x)
		mapY := int(player.position.y)

		var stepX, stepY int
		var sideDistanceX, sideDistanceY float64

		if ray.direction.x < 0 {
			stepX = -1
			sideDistanceX = (player.position.x - float64(mapX)) * ray.deltaDistance.x
		} else {
			stepX = 1
			sideDistanceX = (float64(mapX) + 1.0 - player.position.x) * ray.deltaDistance.x
		}

		if ray.direction.y < 0 {
			stepY = -1
			sideDistanceY = (player.position.y - float64(mapY)) * ray.deltaDistance.y
		} else {
			stepY = 1
			sideDistanceY = (float64(mapY) + 1.0 - player.position.y) * ray.deltaDistance.y
		}

		for !ray.hit {
			if sideDistanceX < sideDistanceY {
				sideDistanceX += ray.deltaDistance.x
				mapX += stepX
				ray.side = 0
			} else {
				sideDistanceY += ray.deltaDistance.y
				mapY += stepY
				ray.side = 1
			}

			if dungeon.At(mapX, mapY) != nil {
				ray.hit = true
				ray.entity = dungeon.At(mapX, mapY)
			}
		}

		if ray.side == 0 {
			ray.perpendicularWallDistace = sideDistanceX - ray.deltaDistance.x
		} else {
			ray.perpendicularWallDistace = sideDistanceY - ray.deltaDistance.y
		}

		hitX := player.position.x + (sideDistanceX-ray.deltaDistance.x)*ray.direction.x
		hitY := player.position.y + (sideDistanceY-ray.deltaDistance.y)*ray.direction.y
		ray.euclideanDistance = euclideanDistance(player.position.x, player.position.y, hitX, hitY)
		rays[x] = ray
	}

	return rays
}

func (r Raycaster) GetRayDirections(camera *Camera) (Vector2, Vector2) {
	return newVector2(camera.direction.x-camera.plane.x, camera.direction.y-camera.plane.y),
		newVector2(camera.direction.x+camera.plane.x, camera.direction.y+camera.plane.y)
}

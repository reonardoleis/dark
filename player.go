package main

import "math"

type Player struct {
	position      Vector2
	camera        *Camera
	moveSpeed     float64
	rotationSpeed float64
}

func newPlayer(position Vector2, camera *Camera) *Player {
	return &Player{
		position:      position,
		camera:        camera,
		moveSpeed:     3.0,
		rotationSpeed: 2.25,
	}
}

func (p *Player) Controls(dungeon *Dungeon, deltaTime float64) {
	if GetControls().IsKeyDown(KeyW) {
		empty := dungeon.grid[int(p.position.y)][int(p.position.x+p.camera.direction.x*p.moveSpeed*DELTA_TIME)].isFloor()
		if empty {
			p.position.x += p.camera.direction.x * p.moveSpeed * DELTA_TIME
		}

		empty = dungeon.grid[int(p.position.y+p.camera.direction.y*p.moveSpeed*DELTA_TIME)][int(p.position.x)].isFloor()
		if empty {
			p.position.y += p.camera.direction.y * p.moveSpeed * DELTA_TIME
		}
	}

	if GetControls().IsKeyDown(KeyS) {
		empty := dungeon.grid[int(p.position.y)][int(p.position.x-p.camera.direction.x*p.moveSpeed*DELTA_TIME)].isFloor()
		if empty {
			p.position.x -= p.camera.direction.x * p.moveSpeed * DELTA_TIME
		}

		empty = dungeon.grid[int(p.position.y-p.camera.direction.y*p.moveSpeed*DELTA_TIME)][int(p.position.x)].isFloor()
		if empty {
			p.position.y -= p.camera.direction.y * p.moveSpeed * DELTA_TIME
		}
	}

	if GetControls().IsKeyDown(KeyD) {
		oldCameraDirectionX := p.camera.direction.x
		p.camera.direction.x = p.camera.direction.x*math.Cos(-p.rotationSpeed*DELTA_TIME) - p.camera.direction.y*math.Sin(-p.rotationSpeed*(DELTA_TIME))
		p.camera.direction.y = oldCameraDirectionX*math.Sin(-p.rotationSpeed*DELTA_TIME) + p.camera.direction.y*math.Cos(-p.rotationSpeed*(DELTA_TIME))

		oldCameraPlaneX := p.camera.plane.x
		p.camera.plane.x = p.camera.plane.x*math.Cos(-p.rotationSpeed*DELTA_TIME) - p.camera.plane.y*math.Sin(-p.rotationSpeed*(DELTA_TIME))
		p.camera.plane.y = oldCameraPlaneX*math.Sin(-p.rotationSpeed*DELTA_TIME) + p.camera.plane.y*math.Cos(-p.rotationSpeed*(DELTA_TIME))
	}

	if GetControls().IsKeyDown(KeyA) {
		oldCameraDirectionX := p.camera.direction.x
		p.camera.direction.x = p.camera.direction.x*math.Cos(p.rotationSpeed*DELTA_TIME) - p.camera.direction.y*math.Sin(p.rotationSpeed*(DELTA_TIME))
		p.camera.direction.y = oldCameraDirectionX*math.Sin(p.rotationSpeed*DELTA_TIME) + p.camera.direction.y*math.Cos(p.rotationSpeed*(DELTA_TIME))

		oldCameraPlaneX := p.camera.plane.x
		p.camera.plane.x = p.camera.plane.x*math.Cos(p.rotationSpeed*DELTA_TIME) - p.camera.plane.y*math.Sin(p.rotationSpeed*(DELTA_TIME))
		p.camera.plane.y = oldCameraPlaneX*math.Sin(p.rotationSpeed*DELTA_TIME) + p.camera.plane.y*math.Cos(p.rotationSpeed*(DELTA_TIME))
	}

}

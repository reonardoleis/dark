package main

import (
	"math"
	"math/rand"
)

type Player struct {
	position              Vector2
	camera                *Camera
	moveSpeed             float64
	rotationSpeed         float64
	attackDuration        float64
	isAttacking           bool
	currentAttackStep     float64
	isWalking             bool
	stats                 *Stats
	currentEnemyNpcTarget *NPC
	attributes            *Attributes
	alive                 bool
	points                int
	enemiesKilled         int
	totalEnemiesKilled    int
}

func newPlayer(position Vector2, camera *Camera) *Player {
	player := &Player{
		position:       position,
		camera:         camera,
		moveSpeed:      3.0,
		rotationSpeed:  0.4,
		attackDuration: 0.25,
		stats:          newStats(),
		attributes:     newAttributes(5, 5, 5, 5),
		alive:          true,
	}

	go playFootsteps(player)

	return player
}

func (p *Player) Controls(dungeon *Dungeon, deltaTime float64) {
	debuffH := 1.0
	debuffV := 1.0

	if GetControls().IsKeyDown(KeyW) || GetControls().IsKeyDown(KeyS) {
		debuffH = 0.5
	}
	if GetControls().IsKeyDown(KeyA) || GetControls().IsKeyDown(KeyD) {
		debuffV = 0.5
	}

	if GetControls().IsKeyDown(KeyW) || GetControls().IsKeyDown(KeyS) || GetControls().IsKeyDown(KeyA) || GetControls().IsKeyDown(KeyD) {
		p.isWalking = true
	} else {
		p.isWalking = false
	}

	if GetControls().IsKeyDown(KeyW) {
		empty := dungeon.grid[int(p.position.y)][int(p.position.x+p.camera.direction.x*p.moveSpeed*DELTA_TIME*debuffV)].isFloor()
		if empty {
			p.position.x += p.camera.direction.x * p.moveSpeed * DELTA_TIME * debuffV
		}

		empty = dungeon.grid[int(p.position.y+p.camera.direction.y*p.moveSpeed*DELTA_TIME*debuffV)][int(p.position.x)].isFloor()
		if empty {
			p.position.y += p.camera.direction.y * p.moveSpeed * DELTA_TIME * debuffV
		}
	}

	if GetControls().IsKeyDown(KeyS) {
		empty := dungeon.grid[int(p.position.y)][int(p.position.x-p.camera.direction.x*p.moveSpeed*DELTA_TIME*debuffV)].isFloor()
		if empty {
			p.position.x -= p.camera.direction.x * p.moveSpeed * DELTA_TIME * debuffV
		}

		empty = dungeon.grid[int(p.position.y-p.camera.direction.y*p.moveSpeed*DELTA_TIME*debuffV)][int(p.position.x)].isFloor()
		if empty {
			p.position.y -= p.camera.direction.y * p.moveSpeed * DELTA_TIME * debuffV
		}
	}

	if GetControls().IsKeyDown(KeyD) {
		perpendicularDirectionX := p.camera.direction.y
		perpendicularDirectionY := -p.camera.direction.x

		empty := dungeon.grid[int(p.position.y)][int(p.position.x+perpendicularDirectionX*p.moveSpeed*DELTA_TIME*debuffH)].isFloor()
		if empty {
			p.position.x += perpendicularDirectionX * p.moveSpeed * DELTA_TIME * debuffH
		}

		empty = dungeon.grid[int(p.position.y+perpendicularDirectionY*p.moveSpeed*DELTA_TIME*debuffH)][int(p.position.x)].isFloor()
		if empty {
			p.position.y += perpendicularDirectionY * p.moveSpeed * DELTA_TIME * debuffH
		}
	}

	if GetControls().IsKeyDown(KeyA) {
		perpendicularDirectionX := -p.camera.direction.y
		perpendicularDirectionY := p.camera.direction.x

		empty := dungeon.grid[int(p.position.y)][int(p.position.x+perpendicularDirectionX*p.moveSpeed*DELTA_TIME*debuffH)].isFloor()
		if empty {
			p.position.x += perpendicularDirectionX * p.moveSpeed * DELTA_TIME * debuffH
		}

		empty = dungeon.grid[int(p.position.y+perpendicularDirectionY*p.moveSpeed*DELTA_TIME*debuffH)][int(p.position.x)].isFloor()
		if empty {
			p.position.y += perpendicularDirectionY * p.moveSpeed * DELTA_TIME * debuffH
		}
	}

	rotation := 0.5
	increments := rotation / p.attackDuration
	if GetControls().IsKeyPressed(KeyM1) && !p.isAttacking && p.currentAttackStep <= 0.0 {
		p.Attack()
		go playWeaponSwing1()
		p.isAttacking = true
	}

	if p.isAttacking {
		p.currentAttackStep += increments * DELTA_TIME
		if p.currentAttackStep >= rotation {
			p.isAttacking = false
		}
	} else if p.currentAttackStep > 0.0 {
		p.currentAttackStep -= increments * DELTA_TIME
	}
	ww, _ := GetControls().window.Size()
	mouseX := float64(GetControls().mouseDeltaX)
	diff := mouseX / float64(ww)
	rotationSpeed := p.rotationSpeed

	diff = 1 + (diff - 0.5)

	if int(mouseX) < ww/2 {
		diff = 1 / diff
	}

	diff = math.Abs(1-diff) * 100

	p.rotationSpeed = rotationSpeed * diff

	if GetControls().mouseDeltaX > ww/2 {
		oldCameraDirectionX := p.camera.direction.x
		p.camera.direction.x = p.camera.direction.x*math.Cos(-p.rotationSpeed*DELTA_TIME) - p.camera.direction.y*math.Sin(-p.rotationSpeed*(DELTA_TIME))
		p.camera.direction.y = oldCameraDirectionX*math.Sin(-p.rotationSpeed*DELTA_TIME) + p.camera.direction.y*math.Cos(-p.rotationSpeed*(DELTA_TIME))

		oldCameraPlaneX := p.camera.plane.x
		p.camera.plane.x = p.camera.plane.x*math.Cos(-p.rotationSpeed*DELTA_TIME) - p.camera.plane.y*math.Sin(-p.rotationSpeed*(DELTA_TIME))
		p.camera.plane.y = oldCameraPlaneX*math.Sin(-p.rotationSpeed*DELTA_TIME) + p.camera.plane.y*math.Cos(-p.rotationSpeed*(DELTA_TIME))

	} else if GetControls().mouseDeltaX < ww/2 {
		oldCameraDirectionX := p.camera.direction.x
		p.camera.direction.x = p.camera.direction.x*math.Cos(p.rotationSpeed*DELTA_TIME) - p.camera.direction.y*math.Sin(p.rotationSpeed*(DELTA_TIME))
		p.camera.direction.y = oldCameraDirectionX*math.Sin(p.rotationSpeed*DELTA_TIME) + p.camera.direction.y*math.Cos(p.rotationSpeed*(DELTA_TIME))

		oldCameraPlaneX := p.camera.plane.x
		p.camera.plane.x = p.camera.plane.x*math.Cos(p.rotationSpeed*DELTA_TIME) - p.camera.plane.y*math.Sin(p.rotationSpeed*(DELTA_TIME))
		p.camera.plane.y = oldCameraPlaneX*math.Sin(p.rotationSpeed*DELTA_TIME) + p.camera.plane.y*math.Cos(p.rotationSpeed*(DELTA_TIME))
	}

	p.rotationSpeed = rotationSpeed
}

func (p *Player) Attack() {
	if p.currentEnemyNpcTarget.distance > enemyAttackDistance {
		return
	}

	if p.currentEnemyNpcTarget.Dodges() {
		addDamageDone(p.currentEnemyNpcTarget.name, 0, false, true)
		return
	}

	p.currentEnemyNpcTarget.hitHighlightTimeleft = hitHighlightTime
	damage, crit := damage(p.attributes)

	if p.currentEnemyNpcTarget.stats.Hp <= damage {
		p.currentEnemyNpcTarget.stats.Hp = 0
		p.currentEnemyNpcTarget.alive = false
		p.handleExp(p.currentEnemyNpcTarget)
		p.enemiesKilled++
		p.totalEnemiesKilled++
	}

	p.currentEnemyNpcTarget.stats.Hp -= damage
	addDamageDone(p.currentEnemyNpcTarget.name, damage, crit, false)
}

func (p *Player) handleExp(npc *NPC) {
	gainedExp := expFactorPerEnemyLevel * float64(npc.stats.Level) * float64(npc.stats.Level)

	p.stats.Exp += int(gainedExp)
	for p.stats.Exp >= p.stats.NextLevel {
		diff := math.Max(0, float64(p.stats.Exp-p.stats.NextLevel))

		p.stats.Level++
		p.stats.Exp = int(diff)
		p.points += 1
		p.stats.NextLevel += int(float64(p.stats.NextLevel) * 1.1)
		p.stats.Hp = p.stats.MaxHp
	}
}
func (p *Player) Dodges() bool {
	chance := dodgeChance(p.attributes)

	random := rand.Intn(101)

	return float64(random)/100 <= chance
}

func (p *Player) SpendPoints() {
	if p.points <= 0 {
		return
	}

	if GetControls().IsKeyPressed(Key1) {
		p.points--
		p.attributes.Str++
	}
	if GetControls().IsKeyPressed(Key2) {
		p.points--
		p.attributes.Agi++
	}
	if GetControls().IsKeyPressed(Key3) {
		p.points--
		p.attributes.Vig++
		p.stats.MaxHp += p.attributes.Vig * p.attributes.Vig
	}
	if GetControls().IsKeyPressed(Key4) {
		p.points--
		p.attributes.Int++
	}
}

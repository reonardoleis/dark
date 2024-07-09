package main

import (
	"math/rand"
)

type NpcName string

const (
	GOBLIN = "Goblin"
)

type EnemyConfig struct {
	attackMinDelaySeconds float64
}

var (
	enemyConfigs = map[NpcName]EnemyConfig{
		GOBLIN: {attackMinDelaySeconds: 0.8},
	}
)

var (
	npcSerialId = 1
)

type NPC struct {
	id                          int
	name                        NpcName
	position                    Vector2
	textureId                   TextureId
	distance                    float64
	attributes                  *Attributes
	stats                       *Stats
	alive                       bool
	timeSinceLastAttack         float64
	randomAttackDelayMultiplier float64
	hitHighlightTimeleft        float64
}

func newNpc(position Vector2, textureId TextureId, attributes *Attributes, stats *Stats) *NPC {
	stats.MaxHp = maxHp(attributes)
	stats.Hp = stats.MaxHp
	npcSerialId++
	return &NPC{
		id:                          npcSerialId - 1,
		position:                    position,
		textureId:                   textureId,
		attributes:                  attributes,
		stats:                       stats,
		alive:                       true,
		randomAttackDelayMultiplier: 1.0 + rand.Float64()*1.2,
	}
}

func getNpcName(npc *NPC) NpcName {
	switch npc.textureId {
	case Sprite1:
		return GOBLIN
	}

	return NpcName("")
}

func (n *NPC) calcMoveSpeed() float64 {
	return float64(n.attributes.Agi) * 0.5
}

func (n *NPC) ChasePlayer(player *Player, dungeon *Dungeon) {
	if n.distance < enemyAttackDistance || n.distance > enemyMaxDistanceMovement {
		return
	}

	dirX := 0
	dirY := 0

	if n.position.x > player.position.x {
		dirX = -1
	} else if n.position.x < player.position.x {
		dirX = 1
	}

	if n.position.y > player.position.y {
		dirY = -1
	} else if n.position.y < player.position.y {
		dirY = 1
	}

	if DELTA_TIME > 100. {
		return
	}

	incX := float64(dirX) * n.calcMoveSpeed() * DELTA_TIME
	incY := float64(dirY) * n.calcMoveSpeed() * DELTA_TIME

	emptyX := dungeon.grid[int(n.position.y)][int(n.position.x+incX)].isFloor()
	if !emptyX {
		incX = 0
	}
	empty := dungeon.grid[int(n.position.y+incY)][int(n.position.x)].isFloor()
	if !empty {
		incY = 0
	}

	n.position.x += incX
	n.position.y += incY
}

func (n *NPC) Dodges() bool {
	chance := dodgeChance(n.attributes)

	random := rand.Intn(101)

	return float64(random)/100 <= chance
}

func (n *NPC) Attack(player *Player) {
	if n.timeSinceLastAttack < enemyConfigs[n.name].attackMinDelaySeconds*n.randomAttackDelayMultiplier {
		n.timeSinceLastAttack += DELTA_TIME
		return
	}

	n.timeSinceLastAttack = 0
	dodges := player.Dodges()
	if dodges {
		addDamageTaken(n.name, 0, false, true)
	}

	damage, crit := damage(n.attributes)
	addDamageTaken(n.name, damage, crit, false)

	if player.stats.Hp-damage <= 0 {
		player.stats.Hp = 0
		player.alive = false
		return
	}

	n.randomAttackDelayMultiplier = 1.0 + rand.Float64()*1.2
	player.stats.Hp -= damage
}

func ComputeNpcEvents(npcs []*NPC, player *Player) {
	for _, npc := range npcs {
		if npc.distance > enemyAttackDistance {
			continue
		}

		if npc.hitHighlightTimeleft > 0 {
			npc.hitHighlightTimeleft -= DELTA_TIME
			if npc.hitHighlightTimeleft < 0 {
				npc.hitHighlightTimeleft = 0
			}
		}

		npc.Attack(player)
	}
}

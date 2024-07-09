package main

import (
	"math/rand"
)

type Spawner struct {
	dungeon *Dungeon
	npcs    *[]*NPC
	player  *Player
}

func newSpawner(dungeon *Dungeon, npcs *[]*NPC, player *Player) Spawner {
	return Spawner{
		dungeon: dungeon,
		npcs:    npcs,
		player:  player,
	}
}

func (s Spawner) initialSpawn() {
	for y, line := range s.dungeon.grid {
		for x, cell := range line {
			random := rand.Intn(101)
			if random < 95 || !cell.isFloor() {
				continue
			}

			position := newVector2(float64(x), float64(y))
			stats := newStats()
			stats.Level = s.dungeon.level + rand.Intn(4)
			npc := newNpc(position, Sprite1, randomAttributes(s.dungeon.level), stats)
			npc.name = getNpcName(npc)

			*s.npcs = append(*s.npcs, npc)
		}
	}

}

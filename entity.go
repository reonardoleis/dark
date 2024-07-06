package main

import "github.com/tfriedel6/canvas"

type EntityType int

const (
	Wall = iota
	Enemy
	Chest
	Door
	Lever
	Bookshelf
	Snowbricks
	Floorbrick
	Hexblue
	Sword
)

type Entity struct {
	t        EntityType
	position Vector2
}

func newEntity(t EntityType, p Vector2) *Entity {
	return &Entity{
		t:        t,
		position: p,
	}
}

func (e Entity) getTexture() *canvas.Image {
	return textures[e.t].image
}

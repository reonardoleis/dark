package main

import (
	"image"

	"github.com/tfriedel6/canvas"
)

type TextureId int

const (
	Bricks1 = iota
	Enemy
	Chest
	Door
	Lever
	Bookshelf
	Snowbricks
	Floorbrick
	Hexblue
	Sword
	Dirtmud
	Bricks1Window1
	Bricks1Window2
	Stonefloor1
	Ground1_0
	Ground1_1
	Ground2_0
	Ground2_1
	Ground3_0
	Ground3_1
	Wall1_0
	Wall1_1
	Wall2_0
	Wall2_1
	Wall2_2
	Wall3_0
	Wall3_1
)

type TextureVariations struct {
	variations         []TextureId
	variationThreshold int
}

var (
	varies = map[TextureId]TextureVariations{
		Wall1_0:   {[]TextureId{Wall1_1}, 80},
		Wall2_0:   {[]TextureId{Wall2_1, Wall2_2}, 70},
		Wall3_0:   {[]TextureId{Wall3_1}, 50},
		Ground1_0: {[]TextureId{Ground1_1}, 50},
		Ground2_0: {[]TextureId{Ground2_1}, 98},
		Ground3_0: {[]TextureId{Ground3_1}, 80},
	}
)

type EntityType int

const (
	Wall = iota
	FloorCeiling
)

type Entity struct {
	t        EntityType
	texture1 TextureId
	texture2 TextureId
	position Vector2
}

func newEntity(texture1, texture2 TextureId, t EntityType, p Vector2) *Entity {
	return &Entity{
		t:        t,
		texture1: texture1,
		texture2: texture2,
		position: p,
	}
}

func (e Entity) getTexture1() (*canvas.Image, image.Image) {
	return textures[e.texture1].image, textures[e.texture1].rgba
}

func (e Entity) getTexture2() (*canvas.Image, image.Image) {
	return textures[e.texture2].image, textures[e.texture2].rgba
}

func (e Entity) isFloor() bool {
	return e.t != Wall
}

func (e Entity) isWall() bool {
	return e.t == Wall
}

func (t TextureId) HasVariations() bool {
	return len(varies[t].variations) > 0
}

func (t TextureId) Variations() TextureVariations {
	return varies[t]
}

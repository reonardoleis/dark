package main

import (
	"image"
	"math/rand"

	"github.com/tfriedel6/canvas"
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
		Wall5_0:   {[]TextureId{Wall5_1}, 70},
		Wall7_0:   {[]TextureId{Wall7_1}, 95},
		Wall8_0:   {[]TextureId{Wall8_1}, 92},
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
	debugId  int
	isRoom   bool
	isSafe   bool
}

func newEntity(texture1, texture2 TextureId, t EntityType, p Vector2, isRoom, isSafe bool, debugId int) *Entity {
	return &Entity{
		t:        t,
		texture1: texture1,
		texture2: texture2,
		position: p,
		isRoom:   isRoom,
		isSafe:   isSafe,
		debugId:  debugId,
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

func (t TextureId) GetRandomVariation() TextureId {
	if !t.HasVariations() {
		return t
	}

	variations := t.Variations()
	variate := rand.Intn(101) >= variations.variationThreshold

	if !variate {
		return t
	}

	return variations.variations[rand.Intn(len(variations.variations))]
}

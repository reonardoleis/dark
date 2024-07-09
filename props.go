package main

type PropPosition int

const (
	PropPositionCeiling PropPosition = iota
	PropPositionFloor   PropPosition = iota
)

type Prop struct {
	textureId    TextureId
	position     Vector2
	distance     float64
	propPosition PropPosition
	divX         int
	divY         int
}

func newProp(position Vector2, textureId TextureId, propPosition PropPosition, divX, divY int) *Prop {
	return &Prop{textureId, position, 0.0, propPosition, divX, divY}
}

func getPropsDivs(textureId TextureId) (int, int) {
	switch textureId {
	case Chains1:
		return 2, 2
	case Clothes1:
		return 1, 1
	case Barrel:
		return 2, 2
	}

	return 0, 0
}

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

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
	//
	Ground4_0
	Ground5_0
	Ground6_0
	Ground6_1
	//
	Wall1_0
	Wall1_1
	Wall2_0
	Wall2_1
	Wall2_2
	Wall3_0
	Wall3_1
	//
	Wall4_0
	Wall5_0
	Wall5_1
	Wall6_0
	Wall7_0
	Wall7_1
	Wall8_0
	Wall8_1
	//
	Sprite1
	Chains1
	Clothes1
	Barrel
)

var (
	menuImage *canvas.Image
)

type Texture struct {
	rgba  image.Image
	image *canvas.Image
}

func newTexture(rgba image.Image, image *canvas.Image) *Texture {
	return &Texture{rgba, image}
}

var (
	texturesPath = "./assets"
	texturesName = []string{"sword1.png", "bricks1.png", "snowbrick1.png", "bookshelf1.png", "ceilinghexblue.png", "floorbrick1.png", "dirtmud.png",
		"bricks1window1.png", "bricks1window2.png", "stonefloor1.png",
		"ground1_0.png", "ground1_1.png", "ground2_0.png", "ground2_1.png", "ground3_0.png", "ground3_1.png",
		"wall1_0.png", "wall1_1.png", "wall2_0.png", "wall2_1.png", "wall2_2.png", "wall3_0.png", "wall3_1.png",
		"ground4_0.png", "ground5_0.png", "ground6_0.png", "ground6_1.png",
		"wall4_0.png", "wall5_0.png", "wall5_1.png", "wall6_0.png", "wall7_0.png", "wall7_1.png", "wall8_0.png", "wall8_1.png",
		"sprite1.png",
		"chains1.png", "clothes1.png",
		"barrel.png",
	}
	texturesType = []TextureId{
		Sword, Bricks1, Snowbricks, Bookshelf, Hexblue, Floorbrick, Dirtmud, Bricks1Window1, Bricks1Window2, Stonefloor1,
		Ground1_0, Ground1_1, Ground2_0, Ground2_1, Ground3_0, Ground3_1,
		Wall1_0, Wall1_1, Wall2_0, Wall2_1, Wall2_2, Wall3_0, Wall3_1,
		Ground4_0, Ground5_0, Ground6_0, Ground6_1,
		Wall4_0, Wall5_0, Wall5_1, Wall6_0, Wall7_0, Wall7_1, Wall8_0, Wall8_1,
		Sprite1,
		Chains1,
		Clothes1,
		Barrel,
	}
	textures             = make(map[TextureId]*Texture)
	floorTextures        = []TextureId{}
	ceilingTextures      = []TextureId{}
	wallTextures         = []TextureId{}
	roomFloorTextures    = []TextureId{}
	roomCeilingTextures  = []TextureId{}
	roomWallTextures     = []TextureId{}
	scenarioPropTextures = []TextureId{}
)

func loadTextures(canvas *canvas.Canvas) {
	for idx, textureName := range texturesName {
		path := fmt.Sprintf("%s/%s", texturesPath, textureName)

		image, err := canvas.LoadImage(path)
		if err != nil {
			panic(err)
		}

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		png, err := png.Decode(file)
		if err != nil {
			panic(err)
		}

		textures[texturesType[idx]] = newTexture(png, image)
		file.Close()
	}

	floorTextures = []TextureId{Bricks1, Stonefloor1, Dirtmud, Ground1_0, Ground2_0, Ground3_0, Ground6_0}
	ceilingTextures = []TextureId{Bricks1, Ground1_0, Wall2_0, Ground5_0}
	wallTextures = []TextureId{Bricks1, Snowbricks, Wall1_0, Wall2_0, Wall3_0, Wall7_0}

	roomFloorTextures = []TextureId{Floorbrick, Ground4_0, Ground5_0}
	roomCeilingTextures = []TextureId{Ground4_0, Ground5_0, Ground6_0}
	roomWallTextures = []TextureId{Wall4_0, Wall5_0, Wall6_0, Wall8_0}
	scenarioPropTextures = []TextureId{Clothes1, Chains1}

	path := fmt.Sprintf("%s/%s", texturesPath, "menu.png")

	image, err := canvas.LoadImage(path)
	if err != nil {
		panic(err)
	}

	menuImage = image

}

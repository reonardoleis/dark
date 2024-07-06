package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/tfriedel6/canvas"
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
	texturesName = []string{"sword1.png", "bricks1.png", "snowbrick1.png", "bookshelf1.png", "ceilinghexblue.png", "floorbrick1.png"}
	texturesType = []EntityType{Sword, Wall, Snowbricks, Bookshelf, Hexblue, Floorbrick}
	textures     = make(map[EntityType]*Texture)
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
}

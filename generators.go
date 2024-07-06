package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/ojrac/opensimplex-go"
)

const (
	TOP    = 0
	BOTTOM = 1
	LEFT   = 2
	RIGHT  = 3
)

func generateLevel(grid [][]*Entity) {
	tempGrid := make([][]int, len(grid))
	for y := 0; y < len(tempGrid); y++ {
		tempGrid[y] = make([]int, len(grid[0]))
	}

	paths := 0

	for paths < int(math.Ceil(float64(mapWidth)/7.0)) {
		currentX := 1
		currentY := 1
		dirX := 0
		dirY := 0
		goal := 0
		startSide := rand.Intn(4)
		switch startSide {
		case TOP:
			currentY = 1
			currentX = rand.Intn(len(grid[0]))
			dirY = 1
			goal = len(tempGrid) - 2
		case BOTTOM:
			currentY = len(tempGrid) - 2
			currentX = rand.Intn(len(grid[0]))
			dirY = -1
			goal = len(tempGrid) - 2
		case LEFT:
			currentX = 1
			currentY = rand.Intn(len(grid))
			dirX = 1
			goal = len(tempGrid[0]) - 2
		case RIGHT:
			currentX = len(tempGrid[0]) - 2
			currentY = rand.Intn(len(grid))
			dirX = -1
			goal = len(tempGrid[0]) - 2
		}

		counter := 0
		for counter < goal+2 {
			tempGrid[currentY][currentX] = 1
			flip := rand.Intn(3)
			if startSide == TOP || startSide == BOTTOM {
				if flip == 0 {
					dirX = -1
				} else if flip == 1 {
					dirX = 1
				} else {
					dirX = 0
				}
			} else {
				if flip == 0 {
					dirY = -1
				} else if flip == 1 {
					dirY = 1
				} else {
					dirY = 0
				}
			}

			if currentY+dirY < 1 {
				dirY = 1
			}

			if currentY+dirY > len(tempGrid)-2 {
				dirY = -1
			}

			if currentX+dirX < 1 {
				dirX = 1
			}

			if currentX+dirX > len(tempGrid[0])-2 {
				dirX = -1
			}

			currentY += dirY
			currentX += dirX
			switch startSide {
			case TOP, BOTTOM:
				tempGrid[currentY+dirY][currentX] = 1
			case LEFT, RIGHT:
				tempGrid[currentY][currentX+dirX] = 1
			}
			counter++
		}

		paths++
	}

	rooms := [][2]int{}
	for y := maxRoomSize + 1; y < len(grid)-maxRoomSize-1; y++ {
		for x := maxRoomSize + 1; x < len(grid[y])-maxRoomSize-1; x++ {
			random := rand.Intn(101)
			if tempGrid[y][x] == 0 {
				continue
			}

			possible := true
			if random > roomThreshold {
				for _, room := range rooms {
					roomX := float64(room[0])
					roomY := float64(room[1])
					candidateX := float64(x)
					candidateY := float64(y)

					distance := euclideanDistance(roomX, roomY, candidateX, candidateY)
					if distance < minRoomDistance {
						possible = false
						break
					}
				}

				if possible {
					sizeX := rand.Intn(maxRoomSize-minRoomSize) + minRoomSize
					sizeY := rand.Intn(maxRoomSize-minRoomSize) + minRoomSize

					for oy := 0; oy < sizeY; oy++ {
						for ox := 0; ox < sizeX; ox++ {
							if tempGrid[y+oy][x+ox] == 4 {
								continue
							}
							tempGrid[y+oy][x+ox] = 2
							if oy == 0 || oy == sizeY-1 || ox == 0 || ox == sizeX-1 {
								if tempGrid[y+oy][x+ox] == 0 {
									tempGrid[y+oy][x+ox] = 3

								} else {
									tempGrid[y+oy][x+ox] = 2
								}
							}

						}
					}

					rooms = append(rooms, [2]int{x, y})
				}
			}
		}
	}

	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			if tempGrid[y][x] == 1 {
				if tempGrid[y+1][x] == 0 && tempGrid[y][x+1] == 0 && tempGrid[y-1][x] == 0 && tempGrid[y][x-1] == 0 {
					tempGrid[y][x] = 0
				}
			}
		}
	}

	noise1 := opensimplex.NewNormalized(time.Now().UnixNano())
	scale := .03
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			v1 := noise1.Eval2(float64(x)*scale, float64(y)*scale)
			index1 := int(v1 * float64(len(floorTextures)))
			index2 := int(v1 * float64(len(ceilingTextures)))
			index3 := int(v1 * float64(len(wallTextures)))

			floorTexture := floorTextures[index1]
			ceilingTexture := ceilingTextures[index2]
			wallTexture := wallTextures[index3]

			random := rand.Intn(101)
			if wallTexture.HasVariations() {
				variations := wallTexture.Variations()
				if random >= variations.variationThreshold {
					random = rand.Intn(len(variations.variations))
					wallTexture = variations.variations[random]
				}
			}

			if y == 0 || y == len(grid)-1 || x == 0 || x == len(grid[y])-1 {
				grid[y][x] = newEntity(wallTexture, wallTexture, Wall, newVector2(float64(x), float64(y)))
			} else if tempGrid[y][x] == 0 {
				grid[y][x] = newEntity(wallTexture, wallTexture, Wall, newVector2(float64(x), float64(y)))
			} else if tempGrid[y][x] == 1 {
				grid[y][x] = newEntity(floorTexture, ceilingTexture, FloorCeiling, newVector2(float64(x), float64(y)))
			} else if tempGrid[y][x] == 2 {
				grid[y][x] = newEntity(Bookshelf, Bookshelf, FloorCeiling, newVector2(float64(x), float64(y)))
			} else if tempGrid[y][x] == 4 {
				grid[y][x] = newEntity(Bookshelf, Bookshelf, Wall, newVector2(float64(x), float64(y)))
			} else {
				grid[y][x] = newEntity(Bookshelf, Bookshelf, Wall, newVector2(float64(x), float64(y)))
			}
		}
	}

}

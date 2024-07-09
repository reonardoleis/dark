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

func isRoomCorner(ox, oy, sizeX, sizeY int) bool {
	return (ox == 0 && oy == 0) || (ox == sizeX-1 && oy == sizeY-1) ||
		(ox == 0 && oy == sizeY-1) || (ox == sizeX-1 && oy == 0)
}

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

			switch startSide {
			case TOP, BOTTOM:
				tempGrid[currentY+dirY][currentX] = 1
			case LEFT, RIGHT:
				tempGrid[currentY][currentX+dirX] = 1
			}
			currentY += dirY

			currentX += dirX

			counter++
		}

		paths++
	}

	rooms := [][2]int{}

	type roomData struct {
		floorTexture   TextureId
		ceilingTexture TextureId
		wallTexture    TextureId
	}

	roomsData := make([][]roomData, len(tempGrid))
	for y := 0; y < len(roomsData); y++ {
		roomsData[y] = make([]roomData, len(tempGrid[0]))
	}

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
					floorTextureId := roomFloorTextures[rand.Intn(len(roomFloorTextures))]
					ceilingTextureId := roomCeilingTextures[rand.Intn(len(roomCeilingTextures))]
					wallTextureId := roomWallTextures[rand.Intn(len(roomWallTextures))]

					currentRoomData := roomData{
						floorTexture:   floorTextureId,
						ceilingTexture: ceilingTextureId,
						wallTexture:    wallTextureId,
					}

					sizeX := rand.Intn(maxRoomSize-minRoomSize) + minRoomSize
					sizeY := rand.Intn(maxRoomSize-minRoomSize) + minRoomSize

					openedSides := map[int]bool{}
					for oy := 0; oy < sizeY; oy++ {
						for ox := 0; ox < sizeX; ox++ {
							tempGrid[y+oy][x+ox] = 4
							roomsData[y+oy][x+ox] = currentRoomData
							if ox == 0 || oy == 0 || ox == sizeX-1 || oy == sizeY-1 {
								tempGrid[y+oy][x+ox] = 5

								if !isRoomCorner(ox, oy, sizeX, sizeY) {
									if ox == 0 && !openedSides[0] && tempGrid[y+oy][x+ox-1] == 1 {
										openedSides[0] = true
										tempGrid[y+oy][x+ox] = 4
									} else if ox == sizeX-1 && !openedSides[1] && tempGrid[y+oy][x+ox+1] == 1 {
										openedSides[1] = true
										tempGrid[y+oy][x+ox] = 4
									} else if oy == 0 && !openedSides[2] && tempGrid[y+oy-1][x+ox] == 1 {
										openedSides[2] = true
										tempGrid[y+oy][x+ox] = 4
									} else if oy == sizeY-1 && !openedSides[3] && tempGrid[y+oy+1][x+ox] == 1 {
										openedSides[3] = true
										tempGrid[y+oy][x+ox] = 4
									}
								}
							}

						}
					}

					for oy := 0; oy < sizeY; oy++ {
						for ox := 0; ox < sizeX; ox++ {

							if isRoomCorner(ox, oy, sizeX, sizeY) {

								if ox == 0 && !openedSides[0] && tempGrid[y+oy][x+ox-1] == 1 {
									openedSides[0] = true
									tempGrid[y+oy][x+ox] = 4
									tempGrid[y+oy][x+ox-1] = 6
									tempGrid[y+oy][x+ox+1] = 6
								} else if ox == sizeX-1 && !openedSides[1] && tempGrid[y+oy][x+ox+1] == 1 {
									openedSides[1] = true
									tempGrid[y+oy][x+ox] = 4
									tempGrid[y+oy][x+ox+1] = 6
									tempGrid[y+oy][x+ox-1] = 6
								} else if oy == 0 && !openedSides[2] && tempGrid[y+oy-1][x+ox] == 1 {
									openedSides[2] = true
									tempGrid[y+oy][x+ox] = 4
									tempGrid[y+oy-1][x+ox] = 6
									tempGrid[y+oy+1][x+ox] = 6
								} else if oy == sizeY-1 && !openedSides[3] && tempGrid[y+oy+1][x+ox] == 1 {
									openedSides[3] = true
									tempGrid[y+oy][x+ox] = 4
									tempGrid[y+oy+1][x+ox] = 6
									tempGrid[y+oy-1][x+ox] = 6
								}

							}

						}
					}

					rooms = append(rooms, [2]int{x, y})
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

			floorTexture := floorTextures[index1].GetRandomVariation()
			ceilingTexture := ceilingTextures[index2].GetRandomVariation()
			wallTexture := wallTextures[index3].GetRandomVariation()

			if y == 0 || y == len(grid)-1 || x == 0 || x == len(grid[y])-1 {
				grid[y][x] = newEntity(wallTexture, wallTexture, Wall, newVector2(float64(x), float64(y)), false, false, tempGrid[y][x])
			} else if tempGrid[y][x] == 0 {
				grid[y][x] = newEntity(wallTexture, wallTexture, Wall, newVector2(float64(x), float64(y)), false, false, tempGrid[y][x])
			} else if tempGrid[y][x] == 1 {
				grid[y][x] = newEntity(floorTexture, ceilingTexture, FloorCeiling, newVector2(float64(x), float64(y)), false, false, tempGrid[y][x])
			} else if tempGrid[y][x] == 2 {
				grid[y][x] = newEntity(Bookshelf, Bookshelf, FloorCeiling, newVector2(float64(x), float64(y)), false, false, tempGrid[y][x])
			} else if tempGrid[y][x] == 4 || tempGrid[y][x] == 6 {
				grid[y][x] = newEntity(roomsData[y][x].floorTexture.GetRandomVariation(), roomsData[y][x].ceilingTexture.GetRandomVariation(), FloorCeiling, newVector2(float64(x), float64(y)), true, false, tempGrid[y][x])
			} else if tempGrid[y][x] == 5 {
				grid[y][x] = newEntity(roomsData[y][x].wallTexture.GetRandomVariation(), roomsData[y][x].wallTexture.GetRandomVariation(), Wall, newVector2(float64(x), float64(y)), true, false, tempGrid[y][x])
			}
		}
	}

}

package main

import "math/rand"

func generateSimpleLevel(grid [][]*Entity) {
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if y == 0 || y == len(grid)-1 || x == 0 || x == len(grid[y])-1 {
				random := rand.Intn(101)
				if random < 90 {
					grid[y][x] = newEntity(Wall, newVector2(float64(x), float64(y)))
				} else {
					grid[y][x] = newEntity(Wall, newVector2(float64(x), float64(y)))
				}
			} else if random := rand.Intn(101); random > 90 {
				grid[y][x] = newEntity(Wall, newVector2(float64(x), float64(y)))
			}
		}
	}
}

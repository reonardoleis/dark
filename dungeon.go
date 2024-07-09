package main

type Dungeon struct {
	grid  [][]*Entity
	level int
}

type LevelGenerator = func(grid [][]*Entity)

func newDungeon(w, h, level int) *Dungeon {
	grid := make([][]*Entity, h)
	for i := 0; i < len(grid); i++ {
		grid[i] = make([]*Entity, w)
	}

	return &Dungeon{
		level: level,
		grid:  grid,
	}
}

func (d *Dungeon) generate(lg LevelGenerator) {
	lg(d.grid)
}

func (d *Dungeon) At(x, y int) *Entity {
	if y < 0 || y > len(d.grid)-1 || x < 0 || x > len(d.grid[y])-1 {
		return nil
	}
	return d.grid[y][x]
}

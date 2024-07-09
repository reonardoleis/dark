package main

import (
	"fmt"
	"strings"
)

type Stats struct {
	Level     int
	Exp       int
	NextLevel int
	Hp        int
	MaxHp     int
}

func newStats() *Stats {
	return &Stats{
		Level:     1,
		Exp:       0,
		NextLevel: 500,
		Hp:        100,
		MaxHp:     100,
	}
}

func (s Stats) String() []string {
	fmt.Printf("%+v\n", s)
	text := fmt.
		Sprintf("Level: %d\nExp: %d\nNext level: %d\nHP: %d\nMax. HP: %d", s.Level, s.Exp, s.NextLevel, s.Hp, s.MaxHp)
	return strings.Split(text, "\n")
}

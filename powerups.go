package main

import (
	"fmt"
	"math/rand"
)

var (
	powerupsNames = []string{"Lifesteal", "Move speed multiplier"}
)

type Powerups struct {
	Lifesteal      float64
	MoveSpeedMulti float64
}

func newPowerups() *Powerups {
	return &Powerups{}
}

func (r *Powerups) Random() (string, float64) {
	var increment float64
	random := rand.Intn(2)

	switch random {
	case 0:
		increment = rand.Float64() / 14
		r.Lifesteal += increment
	default:
		increment = rand.Float64() / 10
		r.MoveSpeedMulti += increment
	}

	return powerupsNames[random], increment
}

func (r *Powerups) String() []string {
	return []string{
		fmt.Sprintf("Lifesteal: %.2f%%", r.Lifesteal*100),
		fmt.Sprintf("Move speed multiplier: %.2f%%", r.MoveSpeedMulti*100),
	}
}

package main

import (
	"math"
	"math/rand"
)

func maxAttribute(dungeonLevel int) int {
	maxAttr := int(math.Floor(float64(dungeonLevel) * 1.66))
	return maxAttr
}

func randomAttribute(playerLevel int) int {
	maxAttr := maxAttribute(playerLevel)
	r := rand.Intn(maxAttr) + 1

	return r
}

func randomAttributes(playerLevel int) *Attributes {
	str := randomAttribute(playerLevel)
	agi := randomAttribute(playerLevel)
	vig := randomAttribute(playerLevel)
	int := randomAttribute(playerLevel)

	return newAttributes(str, agi, vig, int)
}

func maxHp(attributes *Attributes) int {
	return int(float64(attributes.Vig) * 50)
}

func dodgeChance(attributes *Attributes) float64 {
	chance := (float64(attributes.Agi) * 0.5) / 100

	return math.Min(chance, 0.2)
}

func critMultiplier(attributes *Attributes) float64 {
	return 1 + float64(attributes.Agi)*0.08
}

func critChance(attributes *Attributes) float64 {
	chance := (float64(attributes.Agi)) * 0.03
	return math.Min(chance, 1.0)
}

func damage(attributes *Attributes) (int, bool) {
	damage := float64(attributes.Str) * 2
	critChance := critChance(attributes)

	random := rand.Float64()
	crit := random <= critChance
	if crit {
		damage *= critMultiplier(attributes)
	}

	return int(damage), crit
}

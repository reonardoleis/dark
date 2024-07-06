package main

import "math"

type Vector2 struct {
	x float64
	y float64
}

func newVector2(x, y float64) Vector2 {
	return Vector2{x, y}
}

func (v Vector2) Add(vec Vector2) Vector2 {
	return newVector2(v.x+vec.x, v.y+vec.y)
}

func (v Vector2) Sub(vec Vector2) Vector2 {
	return newVector2(v.x-vec.x, v.y-vec.y)
}

func (v Vector2) Dot(vec Vector2) float64 {
	return v.x*vec.x + v.y*vec.y
}

func (v Vector2) Rotate(rads float64) Vector2 {
	x := v.x*math.Cos(rads) - v.y*math.Sin(rads)
	y := v.x*math.Sin(rads) + v.y*math.Cos(rads)

	return newVector2(x, y)
}

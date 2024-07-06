package main

type Camera struct {
	plane     Vector2
	direction Vector2
}

func newCamera(plane, direction Vector2) *Camera {
	return &Camera{plane, direction}
}

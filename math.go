package main

import "math"

func euclideanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func perpendicularToEuclidean(perpendicularDistance float64, playerPosition Vector2, rayDirection Vector2) float64 {
	euclideanDistance := math.Sqrt(perpendicularDistance*perpendicularDistance +
		(playerPosition.x-rayDirection.x)*(playerPosition.x-rayDirection.x) +
		(playerPosition.y-rayDirection.y)*(playerPosition.y-rayDirection.y))

	return euclideanDistance
}

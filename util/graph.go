package util

import (
	"fmt"
	"math"
)

type Grid map[string]Coordinate

type Coordinate struct {
	X int
	Y int
}

func (coordinate *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", coordinate.X, coordinate.Y)
}

func (coordinate *Coordinate) GetClosest(others map[string]Coordinate) Coordinate {

	currentSmallest := 999999.0 // this is bad by I am lazy right now
	var currentSmallestAt Coordinate

	for _, other := range others {
		distance := coordinate.Distance(other)
		if distance < currentSmallest {
			currentSmallestAt = other
			currentSmallest = distance
		}
	}

	return currentSmallestAt
}

// Distance finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func (coordinate *Coordinate) Distance(other Coordinate) float64 {
	first := math.Pow(float64(other.X-coordinate.X), 2)
	second := math.Pow(float64(other.Y-coordinate.Y), 2)
	return math.Sqrt(first + second)
}

func (coordinate *Coordinate) GetPathPoints(end Coordinate) []Coordinate {

	coordinates := []Coordinate{}

	if coordinate.X == end.X && coordinate.Y == end.Y {
		return coordinates
	}

	change := Coordinate{end.X - coordinate.X, end.Y - coordinate.Y}

	gcd := Abs(Gcd(change.X, change.Y))
	change.X = change.X / gcd
	change.Y = change.Y / gcd

	at := Coordinate{coordinate.X + change.X, coordinate.Y + change.Y}

	for at.X != end.X || at.Y != end.Y {
		coordinates = append(coordinates, Coordinate{at.X, at.Y})
		at.X += change.X
		at.Y += change.Y
	}

	return coordinates
}

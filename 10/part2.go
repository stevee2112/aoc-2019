package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	"aoc-2019/util"
	"github.com/golang/geo/r3"
)

func main() {

	// Get Data
	_, file, _, _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	asteroidField := make(util.Grid)
	y := 0
	for scanner.Scan() {
		x := 0
		chars := strings.Split(scanner.Text(), "")

		for _, char := range chars {
			if char == "#" {
				asteroid := util.Coordinate{x, y}
				asteroidField[asteroid.String()] = asteroid
			}

			x++
		}
		y++
	}

	most := 0
	mostAt := util.Coordinate{0, 0}
	angles := []float64{}
	asteroidsByAngle := map[float64]map[string]util.Coordinate{}

	for _, asteroid := range asteroidField {
		if canSee := asteroidsInView(asteroid, asteroidField); canSee > most {
			most = canSee
			mostAt = asteroid
		}
	}

	fmt.Println(most, "can be seen from", mostAt.String())

	// remove asteroid from asteroidField now that we have a monitor on it
	delete(asteroidField, mostAt.String())

	for _, asteroid := range asteroidField {
		if mostAt.X == asteroid.X && mostAt.Y == asteroid.Y {
			continue
		}
		angle := getAngleToAsteroid(mostAt, asteroid)

		if _, ok := asteroidsByAngle[angle]; !ok {
			asteroidsByAngle[angle] = map[string]util.Coordinate{}
		}

		asteroidsByAngle[angle][asteroid.String()] = asteroid
	}

	// sort angles
	for angle, _ := range asteroidsByAngle {
		angles = append(angles, angle)
	}
	sort.Float64s(angles)

	destroyCounter := 1
	destroyedAt := map[int]util.Coordinate{}

	// rotation
	for len(asteroidField) > 0 {
		for _, angle := range angles {
			if len(asteroidsByAngle[angle]) > 0 {
				// we need to destroy an asteroid
				asteroids := asteroidsByAngle[angle]

				// get closest asteroid
				asteroid := mostAt.GetClosest(asteroids)

				delete(asteroidField, asteroid.String())
				delete(asteroidsByAngle[angle], asteroid.String())

				destroyedAt[destroyCounter] = asteroid
				destroyCounter++
			}
		}
	}

	bet := destroyedAt[200]
	fmt.Println((bet.X * 100) + bet.Y)
}

func getAngleToAsteroid(at util.Coordinate, asteroid util.Coordinate) float64 {

	vectorUp := r3.Vector{0, -1, 0}
	vector := r3.Vector{float64(asteroid.X - at.X), float64(asteroid.Y - at.Y), 0}

	if asteroid.X < at.X {
		return 360.0 - vectorUp.Angle(vector).Degrees()
	} else {
		return vectorUp.Angle(vector).Degrees()
	}
}

func asteroidsInView(at util.Coordinate, asteroids util.Grid) int {

	canSee := 0

OtherAsteroid:
	for _, otherAsteroid := range asteroids {

		if at.String() == otherAsteroid.String() {
			continue
		}

		pathPoints := at.GetPathPoints(otherAsteroid)

		if len(pathPoints) > 0 {
			for _, pathPoint := range pathPoints {
				// asteroid in the way
				if _, ok := asteroids[pathPoint.String()]; ok {
					continue OtherAsteroid
				}
			}
		}

		canSee++

	}

	return canSee
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"aoc-2019/util"
)

func main() {

	// Get Data
	_, file, _, _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	asteroids := make(util.Grid)
	y := 0
	for scanner.Scan() {
		x := 0
		chars := strings.Split(scanner.Text(), "")

		for _, char := range chars {
			if char == "#" {
				asteroid := util.Coordinate{x, y}
				asteroids[asteroid.String()] = asteroid
			}

			x++
		}
		y++
	}

	most := 0
	mostAt := util.Coordinate{0, 0}

	for _, asteroid := range asteroids {
		if canSee := asteroidsInView(asteroid, asteroids); canSee > most {
			most = canSee
			mostAt = asteroid
		}
	}

	fmt.Println(most, "can be seen from", mostAt.String())
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

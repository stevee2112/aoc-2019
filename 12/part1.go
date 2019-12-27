package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"strconv"
	"github.com/golang/geo/r3"
	"regexp"
)

type Motion struct {
	Index int;
	Position r3.Vector
	Velocity r3.Vector
}

type Moon Motion

func NewMoon(index int, x float64, y float64, z float64) Moon{
	return Moon{index, r3.Vector{x, y, z}, r3.Vector{}}
}

func main() {

	// Get Data
	_, file, _, _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	moons := []Moon{}

	i := 0
	for scanner.Scan() {
		reg := regexp.MustCompile("[^0-9,-]")
		rawMoonData :=  strings.Split(reg.ReplaceAllString(scanner.Text(), ""), ",")

		x,_ := strconv.ParseFloat(rawMoonData[0], 64)
		y,_ := strconv.ParseFloat(rawMoonData[1], 64)
		z,_ := strconv.ParseFloat(rawMoonData[2], 64)
		moons = append(moons, NewMoon(i, x, y, z))
		i++

	}

	steps := 1000
	for i := 0;i < steps;i++ {
		updateOrbitalMotion(&moons)
	}

	fmt.Println(computeEnergy(moons))
}

func computeEnergy(moons []Moon) float64 {

	energy := 0.0

		for _,moon := range moons {
			energy += (SumAbsCoordinates(moon.Position) * SumAbsCoordinates(moon.Velocity))
		}


	return energy
}

func SumAbsCoordinates(vector r3.Vector) float64 {

	absVector := vector.Abs()
	return absVector.X + absVector.Y + absVector.Z
}

func updateOrbitalMotion(moons *[]Moon) {

	currentMoons := make([]Moon, 4)
	copy(currentMoons, *moons)

	for i,moon := range *moons {
		changeVector := r3.Vector{}
		for _,otherMoon := range currentMoons {

			if moon.Index == otherMoon.Index {
				continue
			}

			if moon.Position.X > otherMoon.Position.X {
				changeVector.X--
			} else if moon.Position.X < otherMoon.Position.X {
				changeVector.X++
			}

			if moon.Position.Y > otherMoon.Position.Y {
				changeVector.Y--
			} else if moon.Position.Y < otherMoon.Position.Y {
				changeVector.Y++
			}

			if moon.Position.Z > otherMoon.Position.Z {
				changeVector.Z--
			} else if moon.Position.Z < otherMoon.Position.Z {
				changeVector.Z++
			}
		}

		// apply gravity
		(*moons)[i].Velocity.X += changeVector.X
		(*moons)[i].Velocity.Y += changeVector.Y
		(*moons)[i].Velocity.Z += changeVector.Z

		// apply velocity
		(*moons)[i].Position.X += (*moons)[i].Velocity.X
		(*moons)[i].Position.Y += (*moons)[i].Velocity.Y
		(*moons)[i].Position.Z += (*moons)[i].Velocity.Z
	}
}

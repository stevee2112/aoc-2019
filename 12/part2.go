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
	"encoding/gob"
	"bytes"
	"crypto/md5"
	"aoc-2019/util"
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

	stepsTillSame := util.Lcm(
		stepsTillSameByAxis(moons, "x"),
		stepsTillSameByAxis(moons, "y"),
		stepsTillSameByAxis(moons, "z"),
	)

	fmt.Println(stepsTillSame)
}

func stepsTillSameByAxis(moons []Moon, axis string) int {

	axisOnlyMoons := []Moon{}

	for i,moon := range moons {

		if axis == "x" {
			axisOnlyMoons = append(axisOnlyMoons, NewMoon(i,moon.Position.X,0,0))
		}

		if axis == "y" {
			axisOnlyMoons = append(axisOnlyMoons, NewMoon(i,0,moon.Position.Y,0))
		}

		if axis == "z" {
			axisOnlyMoons = append(axisOnlyMoons, NewMoon(i,0,0,moon.Position.Z))
		}

	}

	stateOfUniverse := map[string]bool{}

	originalUniverseInBytes,_ := GetBytes(axisOnlyMoons)
	originalUniverseMd5 := fmt.Sprintf("%x", md5.Sum(originalUniverseInBytes))
	stateOfUniverse[originalUniverseMd5] = true

	steps := 0
	found := false
	for !found {
		updateOrbitalMotion(&axisOnlyMoons)
		universeInBytes,_ := GetBytes(axisOnlyMoons)
		universeMd5 := fmt.Sprintf("%x", md5.Sum(universeInBytes))

		if _,ok := stateOfUniverse[universeMd5]; ok {
			found = true
		}

		stateOfUniverse[universeMd5] = true
		steps++
	}

	return steps
}

func printMoons(moons []Moon, step int) {

	fmt.Println("STEP", step);
	fmt.Println("--------------");

	for _,moon := range moons {
		fmt.Printf(
			"pos=<x= %d, y=  %d, z= %d>, vel=<x= %d, y= %d, z= %d>\n",
			int(moon.Position.X),
			int(moon.Position.Y),
			int(moon.Position.Z),
			int(moon.Velocity.X),
			int(moon.Velocity.Y),
			int(moon.Velocity.Z),
		)
	}
}

func GetBytes(key interface{}) ([]byte, error) {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
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

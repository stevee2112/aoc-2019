package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
    "bufio"
	"strings"
)

type SpaceObject struct {
	Code string
	InOrbit []SpaceObject
}

func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	spaceObjects := make(map[string][]SpaceObject)

	for scanner.Scan() {
		orbits := strings.Split(scanner.Text(),")")

		if _, ok := spaceObjects[orbits[0]]; !ok {
			spaceObjects[orbits[0]] = []SpaceObject{}
		}

		spaceObjects[orbits[0]] = append(spaceObjects[orbits[0]], SpaceObject{orbits[1], []SpaceObject{}})
	}

	COM := buildOrbits(spaceObjects, &SpaceObject{"COM", []SpaceObject{}})

	fmt.Println(getOrbitCount(COM, 0))
}

func buildOrbits(spaceObjects map[string][]SpaceObject, spaceObject *SpaceObject) *SpaceObject {

	for _,orbiter := range spaceObjects[spaceObject.Code] {
		spaceObject.InOrbit = append(spaceObject.InOrbit, *buildOrbits(spaceObjects, &orbiter))
	}

	return  spaceObject
}

func getOrbitCount(spaceObject *SpaceObject, depth int) int {

	totalOrbits := depth

	for _, orbiter := range spaceObject.InOrbit {
		totalOrbits += getOrbitCount(&orbiter, depth + 1)
	}

	return totalOrbits
}

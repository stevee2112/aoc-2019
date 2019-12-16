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
	Orbiting *SpaceObject
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

		spaceObjects[orbits[0]] = append(spaceObjects[orbits[0]], SpaceObject{orbits[1], []SpaceObject{}, nil})
	}

	COM := buildOrbits(spaceObjects, &SpaceObject{"COM", []SpaceObject{}, nil}, nil)

	YOU := findSpaceObject("YOU", COM)
	SAN := findSpaceObject("SAN", COM)

	commonAncestor := getCommonAncestor(getPathToAncestor(YOU, COM), getPathToAncestor(SAN, COM))

	fmt.Println(len(getPathToAncestor(YOU, commonAncestor)) + len(getPathToAncestor(SAN, commonAncestor)))
}

func getCommonAncestor(path1 []SpaceObject, path2 []SpaceObject) *SpaceObject {

	for _,pathObj1 := range path1 {
		for _,pathObj2 := range path2 {
			if pathObj1.Code == pathObj2.Code {
				return &pathObj1
			}
		}
	}

	return nil
}

func getPathToAncestor(spaceObject *SpaceObject, ancestor *SpaceObject) []SpaceObject {

	path := []SpaceObject{}

	if (spaceObject.Orbiting.Code != ancestor.Code) {
		path = append(path, *spaceObject)
		path = append(path, getPathToAncestor(spaceObject.Orbiting, ancestor)...)
	}

	return path
}

func findSpaceObject(code string, root *SpaceObject) *SpaceObject {

	if root != nil {
		if root.Code == code {
			return root
		} else {
			for _,orbiter := range root.InOrbit {
				if matchingSpaceObject := findSpaceObject(code, &orbiter); matchingSpaceObject != nil {
					return matchingSpaceObject
				}
			}
		}
	}

	return nil
}

func buildOrbits(spaceObjects map[string][]SpaceObject, spaceObject *SpaceObject, orbiting *SpaceObject) *SpaceObject {

	spaceObject.Orbiting = orbiting

	for i,_ := range spaceObjects[spaceObject.Code] {
		spaceObject.InOrbit = append(spaceObject.InOrbit, *buildOrbits(spaceObjects, &spaceObjects[spaceObject.Code][i], spaceObject))
	}

	return spaceObject
}

func getOrbitCount(spaceObject SpaceObject, depth int) int {

	fmt.Println(strings.Repeat("-", (depth * 5)), spaceObject.Code)
	totalOrbits := depth

	for _, orbiter := range spaceObject.InOrbit {
		totalOrbits += getOrbitCount(orbiter, depth + 1)
	}

	return totalOrbits
}

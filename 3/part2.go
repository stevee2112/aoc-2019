package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
    "bufio"
	"strings"
	"strconv"
)

type Coordinate struct {
	X int
	Y int
}

func (coordinate *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", coordinate.X, coordinate.Y)
}

type Wire struct {
	End Coordinate
	Path map[string]PathEl
	Steps int
}

type PathEl struct {
	Coordinate Coordinate
	Steps int
}

func NewWire() *Wire {
	wire := &Wire{}
	wire.Path = make(map[string]PathEl)
	coordinate := Coordinate{0,0}
	wire.Path[coordinate.String()] = PathEl{coordinate, 0}
	wire.End = coordinate
	return wire
}

func (wire *Wire) AddCoordinate(coordinate Coordinate) {
	wire.Steps++
	wire.Path[coordinate.String()] = PathEl{coordinate, wire.Steps}
	wire.End = coordinate
}

func (wire *Wire) ProcessPathData(pathData []string) {
	for _, instruction := range pathData {
		direction := instruction[0:1]
		length, _ := strconv.Atoi(instruction[1:])

		for at := 0;at < length;at++ {
			currentX := wire.End.X
			currentY := wire.End.Y

			if (direction == "U") {
				wire.AddCoordinate(Coordinate{currentX, currentY + 1})
			}
			if (direction == "D") {
				wire.AddCoordinate(Coordinate{currentX, currentY - 1})
			}
			if (direction == "L") {
				wire.AddCoordinate(Coordinate{currentX - 1, currentY})
			}
			if (direction == "R") {
				wire.AddCoordinate(Coordinate{currentX + 1, currentY})
			}
		}
	}
}


func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	scanner.Scan()
	pathData1 := strings.Split(scanner.Text(), ",")

	scanner.Scan()
	pathData2 := strings.Split(scanner.Text(), ",")

	wire1 := NewWire()
	wire1.ProcessPathData(pathData1)

	wire2 := NewWire()
	wire2.ProcessPathData(pathData2)

	intersections := wireIntersections(wire1, wire2)

	fewestSteps := 0;
	for _, intersection := range intersections {
		wire1Steps := wire1.Path[intersection.String()].Steps
		wire2Steps := wire2.Path[intersection.String()].Steps

		totalSteps := wire1Steps + wire2Steps

		if fewestSteps == 0 || totalSteps < fewestSteps {
			fewestSteps = totalSteps
		}

	}

	fmt.Println(fewestSteps)
}

func wireIntersections(wire1 *Wire, wire2 *Wire) []Coordinate {

	var intersections []Coordinate

	for _, pathEl1 := range wire1.Path {
		if _, ok := wire2.Path[pathEl1.Coordinate.String()]; ok && pathEl1.Coordinate.String() != "0,0" {
				intersections = append(intersections, pathEl1.Coordinate)
		}
	}

	return intersections
}

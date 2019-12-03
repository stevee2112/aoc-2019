package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
    "bufio"
	"strings"
	"strconv"
	"math"
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
	Path map[string]Coordinate
}

func NewWire() *Wire {
	wire := &Wire{}
	wire.Path = make(map[string]Coordinate)
	wire.AddCoordinate(Coordinate{0,0})
	return wire
}

func (wire *Wire) AddCoordinate(coordinate Coordinate) {
	wire.Path[coordinate.String()] = coordinate
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

	shortestFromCenter := 0;
	for _, intersection := range intersections {
		distance := int(math.Abs(float64(intersection.X)) + math.Abs(float64(intersection.Y)))
		if shortestFromCenter == 0 || distance < shortestFromCenter {
			shortestFromCenter = distance
		}
	}

	fmt.Println(shortestFromCenter)
}

func wireIntersections(wire1 *Wire, wire2 *Wire) []Coordinate {

	var intersections []Coordinate

	for _, coordinate1 := range wire1.Path {
		if _, ok := wire2.Path[coordinate1.String()]; ok && coordinate1.String() != "0,0" {
				intersections = append(intersections, coordinate1)
		}
	}

	return intersections
}

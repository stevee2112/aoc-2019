package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strings"
	"strconv"
	"aoc-2019/util"
	"math/rand"
	"time"
	"github.com/gdamore/tcell"
)

type Tile struct {
	coordinate util.Coordinate
	value string
}

type Coordinate struct {
	X int
	Y int
	Tile string
}

func (coordinate *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", coordinate.X, coordinate.Y)
}

type Grid map[string]Coordinate

type Screen struct {
	Grid Grid
}

func NewScreen() *Screen {
	screen := &Screen{}
	screen.Grid = make(Grid)
	return screen
}

func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	rawInput, _ := os.Open(path.Dir(file) + "/input")

	defer rawInput.Close()
	scanner := bufio.NewScanner(rawInput)

	scanner.Scan()
	rawData := scanner.Text()

	intCodeStrings := strings.Split(rawData, ",")
	var intCodes []int
	for _, intCode := range intCodeStrings {
		intCodeInt, _ := strconv.Atoi(intCode)
		intCodes = append(intCodes, intCodeInt)
	}

	run(intCodes)
}
func run(intCodes []int) {

	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)

	rand.Seed(1) // TODO eventually do not do this
	direction := 1
	grid := make(Grid)

	at := Coordinate{0,0,"D"}

	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

programRun:
	for {
		select {
		case out := <-output:
			switch out {
			case 0:
				wallAt := getCoordinateByDirection(at, direction)
				wallAt.Tile = "#"
				grid[wallAt.String()] = wallAt;
			case 1:
				updateCoordinateByDirection(&at, direction)
			case 2:
				o2At := getCoordinateByDirection(at, direction)
				o2At.Tile = "o"
				grid[o2At.String()] = o2At;
				break programRun
			}
		case <-needsInput:
			direction = rand.Intn(4) + 1;
			input <- direction
		case <- done:
			break programRun
		}
	}

	grid["0,0"] = Coordinate{0,0,"D"}
	normalized := getNormalizedGrid(grid);

	scn, _ := tcell.NewScreen()
	scn.Init()
	scn.Clear()

	for _,coordinate := range normalized {
		scn.SetContent(coordinate.X, coordinate.Y, rune(coordinate.Tile[0]), []rune(""), tcell.StyleDefault)
	}

	scn.Show()

	time.Sleep(time.Second * 40)
	scn.Fini()

}

func getNormalizedGrid(grid Grid) Grid {

	normalized := make(Grid)
	min := getGridMinX(grid)
	max := getGridMaxY(grid)

	for _,coordinate := range grid {
		newCoordinate := Coordinate{coordinate.X + (min * -1), util.Abs(coordinate.Y + (max * -1)), coordinate.Tile}
		normalized[newCoordinate.String()] = newCoordinate
	}

	return normalized
}

func getGridMinX(grid Grid) int {

	min := 99999999999999 // not great but lazy

	for _,coordinate := range grid {
		if coordinate.X < min {
			min = coordinate.X
		}
	}

	return min
}

func getGridMaxY(grid Grid) int {

	max := -999999999 // not great but lazy

	for _,coordinate := range grid {
		if coordinate.Y > max {
			max = coordinate.Y
		}
	}

	return max
}

func updateCoordinateByDirection(coordinate *Coordinate, direction int) {
	switch direction {
	case 1: // North
		coordinate.Y++
	case 2: // South
		coordinate.Y--
	case 3: // West
		coordinate.X--
	case 4: // East
		coordinate.X++
	}
}

func getCoordinateByDirection(coordinate Coordinate, direction int) Coordinate {
	switch direction {
	case 1: // North
		coordinate.Y++
	case 2: // South
		coordinate.Y--
	case 3: // West
		coordinate.X--
	case 4: // East
		coordinate.X++
	}

	return coordinate
}

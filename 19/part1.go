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
)

type Path struct {
	Coordinate util.Coordinate
	Direction int
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

	rows := 50;
	cols := 50;
	grid := util.TileGrid{}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tile := run(intCodes, x, y)
			grid[tile.Coordinate.String()] = tile
		}
	}

	util.PrintTileGrid(grid, 5)

	affected := 0
	for _,tile := range grid {
		if tile.Value == "1" {
			affected++
		}
	}

	fmt.Println(affected)
}

func run(intCodes []int, x int, y int) util.Tile {

	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)

	spaceTile := util.Tile{util.Coordinate{x,y}, nil}

	go func() {
		input <- spaceTile.Coordinate.X
		input <- spaceTile.Coordinate.Y
	}();

	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

programRun:
	for {
		select {
		case out := <-output:
			spaceTile.Value = strconv.Itoa(out)
		case <-needsInput:
		case <- done:
			break programRun
		}
	}

	return spaceTile
}

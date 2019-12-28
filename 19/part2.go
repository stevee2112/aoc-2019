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

	startAtX := 1308
	startAtY := 1049
	rows := 100;
	cols := 100;
	grid := util.TileGrid{}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tile := run(intCodes, x + startAtX, y + startAtY)
			grid[tile.Coordinate.String()] = tile
		}
	}

	squareSize := 100
	affected := 0
	smallestFoundAtSum := 9999999999999
	smallestFoundAt := ""

	count := 0
	for _,tile := range grid {
		if tile.Value == "1" {
			affected++

			rightCornerCoor := util.Coordinate{(tile.Coordinate.X + squareSize) - 1, tile.Coordinate.Y}
			bottomCornerCoor := util.Coordinate{tile.Coordinate.X, (tile.Coordinate.Y + squareSize) - 1}

			rightCorner, rightOk := grid[rightCornerCoor.String()]
			bottomCorner, bottomOk := grid[bottomCornerCoor.String()]

			if rightOk && bottomOk {
				if rightCorner.Value == "1" && bottomCorner.Value == "1" {
					tile.Value = "+"

					if (tile.Coordinate.X + tile.Coordinate.Y) < smallestFoundAtSum {
						smallestFoundAt = tile.Coordinate.String()
						smallestFoundAtSum = tile.Coordinate.X + tile.Coordinate.Y
					}
				}
			}
		}

		count++
	}

	smallestAtTile := grid[smallestFoundAt]
	smallestAtTile.Value = "+"
	grid[smallestFoundAt] = smallestAtTile

	fmt.Println(smallestFoundAt)
	fmt.Println(affected)

	if affected == (rows * cols) {
		fmt.Println("possible solution at", smallestAtTile.Coordinate.String())
		fmt.Println("puzzle answer", smallestAtTile.Coordinate.X * 10000 + smallestAtTile.Coordinate.Y)
	}
	
	//util.PrintTileGridShifted(grid, 5, -startAtX, -startAtY)
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

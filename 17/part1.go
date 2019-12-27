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

	run(intCodes)
}

func run(intCodes []int) {

	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)

	outputAscii := []int{}

	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

programRun:
	for {
		select {
		case out := <-output:
			outputAscii = append(outputAscii, out)
		case <-needsInput:
		case <- done:
			break programRun
		}
	}

	rowAt := 0
	colAt := 0

	grid := util.TileGrid{}

	for _,ascii := range outputAscii {
		if ascii == 10 {
			rowAt++
			colAt = 0
		} else {
			char := fmt.Sprintf("%c", rune(ascii))
			at := util.Tile{util.Coordinate{colAt, rowAt}, char}
			grid[at.Coordinate.String()] = at
			colAt++
		}
	}

	alignmentSum := 0
	intersections := getIntersections(grid)

	for _,intersection := range intersections {
		alignment := intersection.Coordinate.X * intersection.Coordinate.Y
		alignmentSum += alignment
	}

	fmt.Println(alignmentSum)
	util.PrintTileGrid(grid, 10)
}

func getIntersections(grid util.TileGrid) []util.Tile {

	tiles := []util.Tile{}


	for _,tile := range grid {

		adj := getAdjacent(tile, grid)

		intersection := true

		if len(adj) != 4 {
			intersection = false
		}

		for _,tile := range adj {
			if tile.Value != "#" {
				intersection = false
			}
		}

		if intersection {
			tiles = append(tiles, tile)
		}
	}

	return tiles
}

func getAdjacent(at util.Tile, grid util.TileGrid) []util.Tile {
	tiles := []util.Tile{}

	//north
	north := util.Coordinate{at.Coordinate.X, at.Coordinate.Y - 1}
	if val,ok := grid[north.String()]; ok {
		tiles = append(tiles, val)
	}

	//south
	south := util.Coordinate{at.Coordinate.X, at.Coordinate.Y + 1}
	if val,ok := grid[south.String()]; ok {
		tiles = append(tiles, val)
	}

	//west
	west := util.Coordinate{at.Coordinate.X - 1, at.Coordinate.Y}
	if val,ok := grid[west.String()]; ok {
		tiles = append(tiles, val)
	}

	//east
	east := util.Coordinate{at.Coordinate.X + 1, at.Coordinate.Y}
	if val,ok := grid[east.String()]; ok {
		tiles = append(tiles, val)
	}

	return tiles
}

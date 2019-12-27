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
	"container/list"
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

	rand.Seed(1) // TODO eventually do not do this
	var direction int
	grid := make(util.TileGrid)

	startAt := util.Tile{util.Coordinate{0,0},"D"}
	at := util.Tile{startAt.Coordinate, "D"}
	o2At := util.Tile{util.Coordinate{0,0},"o"}
	count := 0

	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

programRun:
	for {
		select {
		case out := <-output:
			switch out {
			case 0:
				wallAtCoor := getCoordinateByDirection(at.Coordinate, direction)
				wallAt := util.Tile{wallAtCoor, "#"}
				grid[wallAt.Coordinate.String()] = wallAt;
			case 1:
				updateCoordinateByDirection(&at.Coordinate, direction)
			case 2:
				o2AtCoor := getCoordinateByDirection(at.Coordinate, direction)
				updateCoordinateByDirection(&at.Coordinate, direction)
				o2At = util.Tile{o2AtCoor, "o"}
				grid[o2At.Coordinate.String()] = o2At;
			}
		case <-needsInput:
			if count > 1000000 {
				break programRun
			}
			direction = rand.Intn(4) + 1;
			count++
			input <- direction
		case <- done:
			break programRun
		}
	}

	// add some o2
	fmt.Println(addO2(o2At.Coordinate, grid, 0))
}

func addO2(at util.Coordinate, grid util.TileGrid, timeToFill int) int {
	totalTimeToFill := timeToFill
	freeAdj := getFreeAdjacent(at, grid)
	timeToFill++

	for _,adj := range freeAdj {
		grid[adj.String()] = util.Tile{adj, strconv.Itoa(timeToFill)}
		timeToFillAdj := addO2(adj, grid, timeToFill)

		if timeToFillAdj > totalTimeToFill {
			totalTimeToFill = timeToFillAdj
		}
	}

	return totalTimeToFill
}

func getFreeAdjacent(at util.Coordinate, grid util.TileGrid) []util.Coordinate {
	freeAdj := []util.Coordinate{}

	for i := 1; i <= 4; i++ {

		newCoor := getCoordinateByDirection(at, i)

		pathClear := false
		_,ok := grid[newCoor.String()]
		if !ok {
			pathClear = true
		}

		if pathClear {
			freeAdj = append(freeAdj, newCoor)
		}
	}

	return freeAdj
}

type TileDistance struct {
	Tile util.Tile
	Distance int
}

func shortestDistance(grid util.TileGrid, src util.Tile, dest util.Tile) int {

    // check source and destination cell
    // of the matrix have value 1
	if (grid[src.Coordinate.String()] == grid[dest.Coordinate.String()]) {
		return -1
	}

	visited := map[string]bool{}

    // Mark the source cell as visited
	visited[src.Coordinate.String()] = true

    // Create a queue for BFS
	queue := list.New()

    // Distance of source cell is 0
    queue.PushBack(TileDistance{src, 0})  // Enqueue source cell

    // Do a BFS starting from source cell
    for queue.Len() > 0 {

        curr := queue.Front();
		tileAt := curr.Value.(TileDistance)

        // If we have reached the destination cell,
        // we are done
		if tileAt.Tile.Coordinate.String() == dest.Coordinate.String() {
            return tileAt.Distance
		}

        // Otherwise dequeue the front cell in the queue
        // and enqueue its adjacent cells
        queue.Remove(curr);

        for i := 1; i <= 4; i++ {

			newCoor := getCoordinateByDirection(tileAt.Tile.Coordinate, i)

			pathClear := false
			val,ok := grid[newCoor.String()]
			if !ok || val.Value == "o" {
				pathClear = true
			}

			haveVisited := false
			if _,ok := visited[newCoor.String()]; ok {
				haveVisited = true
			}

			if pathClear && !haveVisited {
				visited[newCoor.String()] = true
				queue.PushBack(TileDistance{util.Tile{newCoor, nil}, tileAt.Distance + 1})
			}
        }
    }

    // Return -1 if destination cannot be reached
    return -1;
}

func updateCoordinateByDirection(coordinate *util.Coordinate, direction int) {
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

func getCoordinateByDirection(coordinate util.Coordinate, direction int) util.Coordinate {
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

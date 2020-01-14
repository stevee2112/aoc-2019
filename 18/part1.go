package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strings"
	"aoc-2019/util"
	"sort"
)

func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	rawInput, _ := os.Open(path.Dir(file) + "/example")

	defer rawInput.Close()
	scanner := bufio.NewScanner(rawInput)

	maze := util.TileGrid{}
	rowAt := 0;
	colAt := 0;

	keys := util.TileGrid{}
	doors := util.TileGrid{}
	at := util.Tile{}

	for scanner.Scan() {
		colAt = 0;
		row  := scanner.Text()
		for _,char := range strings.Split(row, "") {
			tile := util.Tile{util.Coordinate{colAt, rowAt}, char}
			maze[tile.Coordinate.String()] = tile

			if strings.ToLower(char) == char && char != "#" && char != "." && char != "@" {
				keys[char] = tile
			}

			if strings.ToUpper(char) == char && char != "#" && char != "." && char != "@" {
				doors[char] = tile
			}

			if char == "@" {
				at = util.Tile{util.Coordinate{colAt, rowAt}, char}
			}
			colAt++
		}

		rowAt++
	}

	fmt.Println(move(maze, at, keys, doors, 0))
}

func move(maze util.TileGrid, at util.Tile, keys util.TileGrid, doors util.TileGrid, depth int) int {

	depth++
	steps := 0

	// we have got all the keys
	if len(keys) == 0 {
		return 0;
	}

	// Get all reachable keys
	keysReachable := map[string]int{}
	for key,keyTile := range keys {
		stepsToKey := util.ShortestPath(maze, at, keyTile, func(tile util.Tile) bool {
			valString := tile.Value.(string)
			return (strings.ToLower(valString) == valString && valString != "#")
		})

		if stepsToKey > 0 {
			keysReachable[key] = stepsToKey
		}
	}

	//fmt.Println(keysReachable)

	shortestKey,_ := getShortest(maze, at, keys, doors, 0)
	fmt.Println(shortestKey)

    //util.PrintTileGridTerminal(maze)
    //shortestKey = util.GetFromStdin();
	keyTile := keys[shortestKey];

	// steps to key
	steps += util.ShortestPath(maze, at, keyTile, func(tile util.Tile) bool {
		valString := tile.Value.(string)
		return (strings.ToLower(valString) == valString && valString != "#")
	})

	// Clear where current at and move to key location
	maze[at.Coordinate.String()] = util.Tile{util.Coordinate{at.Coordinate.X, at.Coordinate.Y}, "."}
	maze[keyTile.Coordinate.String()] =
		util.Tile{util.Coordinate{keyTile.Coordinate.X, keyTile.Coordinate.Y}, "@"}
	at.Coordinate = keyTile.Coordinate

	// Clear key
	delete(keys, keyTile.Value.(string))

	// Clear Door
	doorTile,ok := doors[strings.ToUpper(keyTile.Value.(string))]

	if ok {
		maze[doorTile.Coordinate.String()] =
			util.Tile{util.Coordinate{doorTile.Coordinate.X, doorTile.Coordinate.Y}, "."}
		delete(doors, doorTile.Value.(string))
	}

	steps += move(maze, at, keys, doors, depth)

	return steps
}

func getShortest(maze util.TileGrid, at util.Tile, keys util.TileGrid, doors util.TileGrid, depth int) (string, int) {

	depth++

	// we have got all the keys
	if len(keys) == 0 {
		return "", 0
	}

	// Get all reachable keys
	keysReachable := map[string]int{}
	for key,keyTile := range keys {
		stepsToKey := util.ShortestPath(maze, at, keyTile, func(tile util.Tile) bool {
			valString := tile.Value.(string)
			return (strings.ToLower(valString) == valString && valString != "#")
		})

		if stepsToKey > 0 {
			keysReachable[key] = stepsToKey
		}
	}

	shortest := 99999999999 // lazy
	shortestKey := ""

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv

	for key,steps := range keysReachable {
		ss = append(ss, kv{key,steps})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	max := len(ss)

	// if max > 2 {
	// 	max = 2
	// }

	// For each rechable key clone stuff then use shortest value for return of moved
	for _,kv := range ss[:max] {

		keyStr := kv.Key

		keySteps := 0;

		tempMaze := util.CloneTileGrid(maze)
		tempKeys := util.CloneTileGrid(keys)
		tempDoors := util.CloneTileGrid(doors)
		tempAt := tempMaze[at.Coordinate.String()]

		keyTile := tempKeys[keyStr];

		// steps to key
		keySteps += util.ShortestPath(tempMaze, tempAt, keyTile, func(tile util.Tile) bool {
			valString := tile.Value.(string)
			return (strings.ToLower(valString) == valString && valString != "#")
		})

		// Clear where current at and move to key location
		tempMaze[tempAt.Coordinate.String()] = util.Tile{util.Coordinate{tempAt.Coordinate.X, tempAt.Coordinate.Y}, "."}
		tempMaze[keyTile.Coordinate.String()] =
			util.Tile{util.Coordinate{keyTile.Coordinate.X, keyTile.Coordinate.Y}, "@"}
		tempAt.Coordinate = keyTile.Coordinate

		// Clear key
		delete(tempKeys, keyTile.Value.(string))

		// Clear Door
		doorTile,ok := tempDoors[strings.ToUpper(keyTile.Value.(string))]

		if ok {
			tempMaze[doorTile.Coordinate.String()] =
				util.Tile{util.Coordinate{doorTile.Coordinate.X, doorTile.Coordinate.Y}, "."}
			delete(tempDoors, doorTile.Value.(string))
		}

		_, stepRet := getShortest(tempMaze, tempAt, tempKeys, tempDoors, depth)
		keySteps += stepRet

		if keySteps <= shortest {
			shortest = keySteps
			shortestKey = keyStr
		}
	}

	return shortestKey, shortest
}

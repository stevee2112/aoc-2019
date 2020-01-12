package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strings"
	"aoc-2019/util"
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
	at := &util.Tile{}

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
				at = &tile
			}
			colAt++
		}

		rowAt++
	}

	steps := 0

	util.PrintTileGrid(maze, 2)

	for len(keys) > 0 {
		steps += move(&maze, at, &keys, &doors)
	}

	fmt.Println(steps)
}

func move(maze *util.TileGrid, at *util.Tile, keys *util.TileGrid, doors *util.TileGrid) int {

	// find closest key
	keyDistance := map[string]int{}


	for key,keyTile := range *keys {
		keyDistance[key] = util.ShortestPath(*maze, *at, keyTile, func(tile util.Tile) bool {
			valString := tile.Value.(string)
			return (strings.ToLower(valString) == valString && valString != "#")
		})
	}

	var closestKey string
	closestSteps := 999999999999 // again lazy

	for key,distance := range keyDistance {
		if distance > 0 && distance < closestSteps {
			closestSteps = distance
			closestKey = key
		}
	}

	keyTile := (*keys)[closestKey];

	// Clear where current at and move to key location
	(*maze)[at.Coordinate.String()] = util.Tile{util.Coordinate{at.Coordinate.X, at.Coordinate.Y}, "."}
	(*maze)[keyTile.Coordinate.String()] =
		util.Tile{util.Coordinate{keyTile.Coordinate.X, keyTile.Coordinate.Y}, "@"}
	at.Coordinate = keyTile.Coordinate

	// Clear key
	delete(*keys, keyTile.Value.(string))

	// Clear Door
	doorTile,ok := (*doors)[strings.ToUpper(keyTile.Value.(string))]

	if ok {
		(*maze)[doorTile.Coordinate.String()] =
			util.Tile{util.Coordinate{doorTile.Coordinate.X, doorTile.Coordinate.Y}, "."}
		delete(*doors, doorTile.Value.(string))
	}

	util.PrintTileGrid(*maze, 2)

	return closestSteps
}

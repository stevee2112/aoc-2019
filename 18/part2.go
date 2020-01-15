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
	"crypto/md5"
	"encoding/json"
	"unicode"
)

type kvr struct {
	Key   string
	Value int
	Robot string
}


var globalCache map[string]kvr

func main() {

	// initial cache
	globalCache = make(map[string]kvr)

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	rawInput, _ := os.Open(path.Dir(file) + "/inputpart2")

	defer rawInput.Close()
	scanner := bufio.NewScanner(rawInput)

	maze := util.TileGrid{}
	rowAt := 0;
	colAt := 0;

	keys := util.TileGrid{}
	doors := util.TileGrid{}
	robots := util.TileGrid{}

	for scanner.Scan() {
		colAt = 0;
			row  := scanner.Text()
		for _,char := range strings.Split(row, "") {
			tile := util.Tile{util.Coordinate{colAt, rowAt}, char}
			maze[tile.Coordinate.String()] = tile

			if strings.ToLower(char) == char && char != "#" && char != "." && !unicode.IsNumber(rune(char[0])) {
				keys[char] = tile
			}

			if strings.ToUpper(char) == char && char != "#" && char != "." && !unicode.IsNumber(rune(char[0])) {
				doors[char] = tile
			}

			if unicode.IsNumber(rune(char[0])) {
				robots[char] = tile
			}
			colAt++
		}

		rowAt++
	}

	fmt.Println(move(maze, robots, keys, doors))
}

func move(maze util.TileGrid, robots util.TileGrid, keys util.TileGrid, doors util.TileGrid) int {

	steps := 0

	// we have got all the keys
	if len(keys) == 0 {
		return 0;
	}

    //util.PrintTileGridTerminal(maze)

	// UNCOMMENT WHEN READY TO DO WITH ALGO
	shortestKVR := getShortest(maze, robots, keys, doors, 0)
	// fmt.Println(shortestKVR.Robot, shortestKVR.Key)

	shortestKey := shortestKVR.Key
	robotToMove := shortestKVR.Robot

	fmt.Println(robotToMove, shortestKey)
	// fmt.Print("Robot to Move: ");
	// robotToMove = util.GetFromStdin();
	// fmt.Print("Key to get: ");
    // shortestKey = util.GetFromStdin();


	keyTile := keys[shortestKey];
	movingRobot := robots[robotToMove]
	at := maze[movingRobot.Coordinate.String()]

	// steps to key
	moves := util.ShortestPath(maze, at, keyTile, func(tile util.Tile) bool {
		valString := tile.Value.(string)
		return (strings.ToLower(valString) == valString && valString != "#")
	})

	steps += moves

	// Clear where current at and move to key location
	maze[at.Coordinate.String()] = util.Tile{util.Coordinate{at.Coordinate.X, at.Coordinate.Y}, "."}
	maze[keyTile.Coordinate.String()] =
		util.Tile{util.Coordinate{keyTile.Coordinate.X, keyTile.Coordinate.Y}, robotToMove}
	at.Coordinate = keyTile.Coordinate
	robots[robotToMove] = util.Tile{util.Coordinate{keyTile.Coordinate.X, keyTile.Coordinate.Y}, robotToMove}

	// Clear key
	delete(keys, keyTile.Value.(string))

	// Clear Door
	doorTile,ok := doors[strings.ToUpper(keyTile.Value.(string))]

	if ok {
		maze[doorTile.Coordinate.String()] =
			util.Tile{util.Coordinate{doorTile.Coordinate.X, doorTile.Coordinate.Y}, "."}
		delete(doors, doorTile.Value.(string))
	}

	steps += move(maze, robots, keys, doors)

	return steps
}

func getShortest(maze util.TileGrid, robots util.TileGrid, keys util.TileGrid, doors util.TileGrid, depth int) kvr {

	// we have got all the keys
	if len(keys) == 0 {
		return kvr{}
	}

	//checkCache
	if val,ok := globalCache[getCacheKey(maze,robots,keys,doors)]; ok {
		return val
	}

	// Get all reachable keys
	keysReachable := map[string]kvr{}
	for robot, at := range robots {
		for key,keyTile := range keys {
			stepsToKey := util.ShortestPath(maze, at, keyTile, func(tile util.Tile) bool {
				valString := tile.Value.(string)
				return (strings.ToLower(valString) == valString && valString != "#")
			})

			if stepsToKey > 0 {
				keysReachable[key] = kvr{key, stepsToKey, robot}
			}
		}
	}

	shortest := 99999999999 // lazy
	shortestKey := ""
	robotToMove := ""

	var ss []kvr

	for _,kvr := range keysReachable {
		ss = append(ss, kvr)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	max := len(ss)

	if max > 1 {
		max = 1
	}

	// For each rechable key clone stuff then use shortest value for return of moved
	for _,kvr := range ss[:max] {

		keyStr := kvr.Key

		keySteps := 0;

		tempMaze := util.CloneTileGrid(maze)
		tempKeys := util.CloneTileGrid(keys)
		tempDoors := util.CloneTileGrid(doors)
		tempRobots := util.CloneTileGrid(robots)

		keyTile := tempKeys[keyStr];
		movingRobot := tempRobots[kvr.Robot]
		tempAt := tempMaze[movingRobot.Coordinate.String()]

		// steps to key
		keySteps += util.ShortestPath(tempMaze, tempAt, keyTile, func(tile util.Tile) bool {
			valString := tile.Value.(string)
			return (strings.ToLower(valString) == valString && valString != "#")
		})
		keySteps += 1

		// Clear where current at and move to key location
		tempMaze[tempAt.Coordinate.String()] = util.Tile{util.Coordinate{tempAt.Coordinate.X, tempAt.Coordinate.Y}, "."}
		tempMaze[keyTile.Coordinate.String()] =
			util.Tile{util.Coordinate{keyTile.Coordinate.X, keyTile.Coordinate.Y}, kvr.Robot}
		tempAt.Coordinate = keyTile.Coordinate
		tempRobots[robotToMove] = util.Tile{util.Coordinate{keyTile.Coordinate.X, keyTile.Coordinate.Y}, robotToMove}

		// Clear key
		delete(tempKeys, keyTile.Value.(string))

		// Clear Door
		doorTile,ok := tempDoors[strings.ToUpper(keyTile.Value.(string))]

		if ok {
			tempMaze[doorTile.Coordinate.String()] =
				util.Tile{util.Coordinate{doorTile.Coordinate.X, doorTile.Coordinate.Y}, "."}
			delete(tempDoors, doorTile.Value.(string))
		}

		kvrRet := getShortest(tempMaze, tempRobots, tempKeys, tempDoors, depth)
		keySteps += kvrRet.Value

		if keySteps <= shortest {
			shortest = keySteps
			shortestKey = keyStr
			robotToMove = kvr.Robot
		}
	}

	saveToCache(maze, robots, keys, doors, kvr{shortestKey, shortest, robotToMove})

	return kvr{shortestKey, shortest, robotToMove}
}

func getCacheKey(maze util.TileGrid, robots util.TileGrid, keys util.TileGrid, doors util.TileGrid) string {
	type CacheKeyData struct {
		Maze util.TileGrid
		Robots util.TileGrid
		Keys util.TileGrid
		Doors util.TileGrid
	}

	data, _ := json.Marshal(CacheKeyData{maze, robots, keys, doors})
	hash := fmt.Sprintf("%x", md5.Sum(data))

	return hash

}

func saveToCache(maze util.TileGrid, robots util.TileGrid, keys util.TileGrid, doors util.TileGrid, data kvr) {
	globalCache[getCacheKey(maze, robots, keys, doors)] = data
}

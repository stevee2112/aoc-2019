package util

import (
	"fmt"
	"math"
	"container/list"
)

type Grid map[string]Coordinate

type TileGrid map[string]Tile

type Tile struct {
	Coordinate Coordinate
	Value interface{}
}

type Coordinate struct {
	X int
	Y int
}

type TileDistance struct {
	Tile Tile
	Distance int
}

type FreeSpaceFunc func (tile Tile) bool

func (coordinate *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", coordinate.X, coordinate.Y)
}

func (coordinate *Coordinate) GetClosest(others map[string]Coordinate) Coordinate {

	currentSmallest := 999999.0 // this is bad by I am lazy right now
	var currentSmallestAt Coordinate

	for _, other := range others {
		distance := coordinate.Distance(other)
		if distance < currentSmallest {
			currentSmallestAt = other
			currentSmallest = distance
		}
	}

	return currentSmallestAt
}

// Distance finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func (coordinate *Coordinate) Distance(other Coordinate) float64 {
	first := math.Pow(float64(other.X-coordinate.X), 2)
	second := math.Pow(float64(other.Y-coordinate.Y), 2)
	return math.Sqrt(first + second)
}

func (coordinate *Coordinate) GetPathPoints(end Coordinate) []Coordinate {

	coordinates := []Coordinate{}

	if coordinate.X == end.X && coordinate.Y == end.Y {
		return coordinates
	}

	change := Coordinate{end.X - coordinate.X, end.Y - coordinate.Y}

	gcd := Abs(Gcd(change.X, change.Y))
	change.X = change.X / gcd
	change.Y = change.Y / gcd

	at := Coordinate{coordinate.X + change.X, coordinate.Y + change.Y}

	for at.X != end.X || at.Y != end.Y {
		coordinates = append(coordinates, Coordinate{at.X, at.Y})
		at.X += change.X
		at.Y += change.Y
	}

	return coordinates
}

func GetNormalizedGrid(grid TileGrid) TileGrid {

	normalized := make(TileGrid)
	min := getGridMinX(grid)
	max := getGridMaxY(grid)

	for _,tile := range grid {
		newTile := Tile{Coordinate{tile.Coordinate.X + (min * -1), Abs(tile.Coordinate.Y + (max * -1))}, tile.Value}
		normalized[newTile.Coordinate.String()] = newTile
	}

	return normalized
}

func getGridMinX(grid TileGrid) int {

	min := 99999999999999 // not great but lazy

	for _,tile := range grid {
		if tile.Coordinate.X < min {
			min = tile.Coordinate.X
		}
	}

	return min
}

func getGridMaxX(grid TileGrid) int {

	max := -99999999999999 // not great but lazy

	for _,tile := range grid {
		if tile.Coordinate.X > max {
			max = tile.Coordinate.X
		}
	}

	return max
}

func getGridMaxY(grid TileGrid) int {

	max := -999999999 // not great but lazy

	for _,tile := range grid {
		if tile.Coordinate.Y > max {
			max = tile.Coordinate.Y
		}
	}

	return max
}

func ShortestPath(grid TileGrid, src Tile, dest Tile, freeSpaceFunc FreeSpaceFunc) int {

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
			freeSpace := freeSpaceFunc(val)
			if ok && freeSpace {
				pathClear = true
			}

			haveVisited := false
			if _,ok := visited[newCoor.String()]; ok {
				haveVisited = true
			}

			if pathClear && !haveVisited {
				visited[newCoor.String()] = true
				queue.PushBack(TileDistance{Tile{newCoor, nil}, tileAt.Distance + 1})
			}
        }
    }

    // Return -1 if destination cannot be reached
    return -1;
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

func CloneTileGrid(tileGrid TileGrid) TileGrid {
	newTileGrid := TileGrid{}

	for key,tile := range tileGrid {
		newTileGrid[key] = Tile{Coordinate{tile.Coordinate.X, tile.Coordinate.Y}, tile.Value}
	}

	return newTileGrid
}

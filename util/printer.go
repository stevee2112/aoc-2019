package util

import (
	"time"
	"github.com/gdamore/tcell"
)

func PrintTileGrid(grid TileGrid, duration int) {

	scn, _ := tcell.NewScreen()
	scn.Init()
	scn.Clear()

	for _,tile := range grid {
		scn.SetContent(tile.Coordinate.X, tile.Coordinate.Y, rune(tile.Value.(string)[0]), []rune(""), tcell.StyleDefault)
	}

	scn.Show()
	time.Sleep(time.Second * time.Duration(duration))
	scn.Fini()
}

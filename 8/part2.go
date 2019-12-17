package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strconv"
	"github.com/gdamore/tcell"
	"time"
)

type Image struct {
	Width int
	Height int
	Layers []Layer
}

func (image *Image) Rendered() map[string]Coordinate {

	renderedImage := make(map[string]Coordinate)

	for _,layer := range image.Layers {
		for _,coordinate := range layer {

			if _, ok := renderedImage[coordinate.String()]; !ok {
				renderedImage[coordinate.String()] = Coordinate{coordinate.X, coordinate.Y, 2}
			}

			if coordinate.Pixel != 2 && renderedImage[coordinate.String()].Pixel == 2 {
				renderedImage[coordinate.String()] = coordinate
			}
		}
	}

	return renderedImage
}

type Layer map[string]Coordinate

func (layer *Layer) GetPixelTypeCount(pixelType int) int {

	count := 0
	for _,coordinate := range *layer {
		if coordinate.Pixel == pixelType {
			count++
		}
	}

	return count
}

type Coordinate struct {
	X int
	Y int
	Pixel int
}

func (coordinate *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", coordinate.X, coordinate.Y)
}

func NewImage(width int, height int, imageData []int) *Image {
	image := &Image{}
	image.Width = width
	image.Height = height
	image.Layers = []Layer{}

	colAt := 0;
	rowAt := 0;
	pixelsPerLayer := image.Height * image.Width
	pixelsLeftInLayer := 0
	for _,pixel := range imageData {
		if (pixelsLeftInLayer == 0) {
			layer := make(map[string]Coordinate)
			image.Layers = append(image.Layers, layer)
			pixelsLeftInLayer = pixelsPerLayer
			colAt = 0
			rowAt = 0
		}

		currentLayer := image.Layers[len(image.Layers) - 1]
		newPixelCoor := Coordinate{colAt,rowAt,pixel}
		currentLayer[newPixelCoor.String()] = newPixelCoor

		colAt++
		if colAt >= image.Width {
			colAt = 0
			rowAt++
		}

		pixelsLeftInLayer--
	}

	return image
}

func displayRenderedImage(renderedImage map[string]Coordinate) {
	scn, _ := tcell.NewScreen()
	scn.Init()
	scn.Clear()

	for _,coordinate := range renderedImage {

		if coordinate.Pixel == 1 {
			scn.SetContent(coordinate.X, coordinate.Y, rune(strconv.Itoa(coordinate.Pixel)[0]), []rune(""), tcell.StyleDefault)
		}
	}

	scn.Show()

	quit := make(chan struct{})
	go func() {
		for {
			ev := scn.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					scn.Sync()
				}
			case *tcell.EventResize:
				scn.Sync()
			}
		}
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
	}

	scn.Fini()
}

func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	rawInput, _ := os.Open(path.Dir(file) + "/input")

	defer rawInput.Close()
	scanner := bufio.NewScanner(rawInput)
	scanner.Split(bufio.ScanRunes)

	rawImageData := []int{}

    for scanner.Scan() {
		pixel,_ := strconv.Atoi(scanner.Text())
        rawImageData = append(rawImageData, pixel)
    }

	image := NewImage(25, 6, rawImageData)

	displayRenderedImage(image.Rendered())
}

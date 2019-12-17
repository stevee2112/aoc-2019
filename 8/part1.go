package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strconv"
)

type Image struct {
	Width int
	Height int
	Layers []Layer
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

	fewest0Layer := 0
	fewest0 := image.Layers[0].GetPixelTypeCount(0)
	for i,layer := range image.Layers {
		zeros := layer.GetPixelTypeCount(0);
		if zeros < fewest0 {
			fewest0Layer = i
			fewest0 = zeros
		}
	}
	fmt.Println(image.Layers[fewest0Layer].GetPixelTypeCount(1) * image.Layers[fewest0Layer].GetPixelTypeCount(2))

}

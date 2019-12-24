package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strings"
	"strconv"
	"github.com/gdamore/tcell"
	"time"
	"aoc-2019/util"
)

type Coordinate struct {
	X int
	Y int
	Tile int
}

func (coordinate *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", coordinate.X, coordinate.Y)
}

type Screen struct {
	Grid map[string]Coordinate
}

func NewScreen() *Screen {
	screen := &Screen{}
	screen.Grid = make(map[string]Coordinate)
	return screen
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

	intCodes[0] = 2 // Free play

	run(intCodes)
}
func run(intCodes []int) {

	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)


	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

	inputCounter := 0;
	coordinate := Coordinate{0,0,0}
	score := 0
	ballXAt := 0
	paddleXAt := 0

	scn, _ := tcell.NewScreen()
	scn.Init()
	scn.Clear()
	scn.Show()

programRun:
	for {
		select {
		case out := <-output:
			if inputCounter == 0 {
				coordinate.X = out
				inputCounter++
			} else if inputCounter == 1 {
				coordinate.Y = out
				inputCounter++
			} else if inputCounter == 2 {
				coordinate.Tile = out
				inputCounter = 0
				if coordinate.String() == "-1,0" {
					score = out
				} else {

					if (coordinate.Tile == 4) {
						ballXAt = coordinate.X
					}

					if (coordinate.Tile == 3) {
						paddleXAt = coordinate.X
					}

					char := strconv.Itoa(coordinate.Tile)[0];

					if char == '0' {
						char = ' '
					}

					scn.SetContent(coordinate.X, coordinate.Y, rune(char), []rune(""), tcell.StyleDefault)
					scn.Sync()
				}
			}
		case <-needsInput:
			time.Sleep(time.Millisecond)

			if paddleXAt == ballXAt {
				input <- 0
			}

			if paddleXAt < ballXAt {
				input <- 1
			}

			if paddleXAt > ballXAt {
				input <- -1
			}
		case <- done:
			break programRun
		}
	}

	scn.Fini()
	fmt.Println(score)
}

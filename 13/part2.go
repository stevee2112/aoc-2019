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


	go processIntCodes(intCodes, input, output, needsInput, done)

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
			time.Sleep(time.Millisecond * 50)

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

func processIntCodes(originalIntCodes []int, input chan int, output chan int, needsInput chan bool, done chan bool) []int {

	intCodes := make([]int, len(originalIntCodes))
	copy(intCodes, originalIntCodes)

	relativeBase := 0
	at := 0
	step := map[string]int{
		"01": 4,
		"02": 4,
		"03": 2,
		"04": 2,
		"05": 3,
		"06": 3,
		"07": 4,
		"08": 4,
		"09": 2,
	}

	halt := false
	for halt == false {
		movePointer := true
		instruction := fmt.Sprintf("%05d", intCodes[at])

		opCode := instruction[3:]

		if opCode == "99" {
			done <- true
			break
		}

		if opCode == "01" {

			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at + 3]
			param3Mode := instruction[0:1]

			if param3Mode == "2" {
				param3 = param3 + relativeBase
			}

			if param3 >= len(intCodes) {
				intCodes = (growMemory(intCodes, param3 + 1))
			}

			intCodes[param3] = getValue(intCodes, param1, param1Mode, relativeBase) +
				getValue(intCodes, param2, param2Mode, relativeBase)
		}

		if opCode == "02" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at + 3]
			param3Mode := instruction[0:1]

			if param3Mode == "2" {
				param3 = param3 + relativeBase
			}

			if param3 >= len(intCodes) {
				intCodes = (growMemory(intCodes, param3 + 1))
			}

			intCodes[param3] = getValue(intCodes, param1, param1Mode, relativeBase) *
				getValue(intCodes, param2, param2Mode, relativeBase)
		}

		if opCode == "03" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]

			if param1Mode == "2" {
				param1 = param1 + relativeBase
			}

			if param1 >= len(intCodes) {
				intCodes = (growMemory(intCodes, param1 + 1))
			}

			needsInput <- true
			intCodes[param1] = <-input
		}

		if opCode == "04" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			output <- getValue(intCodes, param1, param1Mode, relativeBase)
		}

		if opCode == "05" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]

			if getValue(intCodes, param1, param1Mode, relativeBase) != 0 {
				at = getValue(intCodes, param2, param2Mode, relativeBase)
				movePointer = false
			}
		}

		if opCode == "06" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]

			if getValue(intCodes, param1, param1Mode, relativeBase) == 0 {
				at = getValue(intCodes, param2, param2Mode, relativeBase)
				movePointer = false
			}
		}

		if opCode == "07" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at + 3]
			param3Mode := instruction[0:1]

			if param3Mode == "2" {
				param3 = param3 + relativeBase
			}

			if param3 >= len(intCodes) {
				intCodes = (growMemory(intCodes, param3 + 1))
			}

			if getValue(intCodes, param1, param1Mode, relativeBase) <
				getValue(intCodes, param2, param2Mode, relativeBase) {
				intCodes[param3] = 1
			} else {
				intCodes[param3] = 0
			}
		}

		if opCode == "08" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at + 3]
			param3Mode := instruction[0:1]

			if param3Mode == "2" {
				param3 = param3 + relativeBase
			}

			if param3 >= len(intCodes) {
				intCodes = (growMemory(intCodes, param3 + 1))
			}

			if getValue(intCodes, param1, param1Mode, relativeBase) ==
				getValue(intCodes, param2, param2Mode, relativeBase) {
				intCodes[param3] = 1
			} else {
				intCodes[param3] = 0
			}
		}

		if opCode == "09" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			relativeBase += getValue(intCodes, param1, param1Mode, relativeBase)
		}

		if movePointer {
			at += step[opCode]
		}
	}

	return intCodes
}

func getValue(intCodes []int, parameter int, mode string, relativeBase int) int {

	var value int

	if mode == "0" {

		if parameter >= len(intCodes) {
			value = 0;
		} else {
			value = intCodes[parameter]
		}
	}

	if mode == "1" {
		value = parameter
	}

	if mode == "2" {
		if parameter + relativeBase >= len(intCodes) {
			value = 0;
		} else {
			value = intCodes[parameter + relativeBase]
		}
	}

	return value
}

func growMemory(intCodes []int, size int) []int {
	newMemory := make([]int, size, size)
	copy(newMemory, intCodes)

	return newMemory
}

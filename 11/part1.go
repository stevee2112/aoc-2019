package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strings"
	"strconv"
)

type Coordinate struct {
	X int
	Y int
	Color int
}

func (coordinate *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", coordinate.X, coordinate.Y)
}

type Robot struct {
	At Coordinate
	Painted map[string]Coordinate
	Direction string
}

func NewRobot() *Robot {
	robot := &Robot{}
	robot.Painted = make(map[string]Coordinate)
	coordinate := Coordinate{0,0,0}
	robot.At = coordinate
	robot.Direction = "N"
	return robot
}

func (robot *Robot) GetCurrentPanelColor() int {
	if val, ok := robot.Painted[robot.At.String()]; ok {
		return val.Color
	} else {
		return 0 // Assume panel is black if not in path
	}
}

func (robot *Robot) PaintPanel(color int) {
	robot.At.Color = color
	robot.Painted[robot.At.String()] = robot.At
}

func (robot *Robot) Move(rotation int) {

	currentX := robot.At.X
	currentY := robot.At.Y

	switch robot.Direction {
	case "N":
		if rotation == 0 {
			robot.MoveToCoordinate(Coordinate{currentX, currentY - 1, 0})
			robot.Direction = "W"
		}

		if rotation == 1 {
			robot.MoveToCoordinate(Coordinate{currentX, currentY + 1, 0})
			robot.Direction = "E"
		}
	case "S":
		if rotation == 0 {
			robot.MoveToCoordinate(Coordinate{currentX, currentY + 1, 0})
			robot.Direction = "E"
		}

		if rotation == 1 {
			robot.MoveToCoordinate(Coordinate{currentX, currentY - 1, 0})
			robot.Direction = "W"
		}
	case "E":
		if rotation == 0 {
			robot.MoveToCoordinate(Coordinate{currentX + 1, currentY, 0})
			robot.Direction = "N"
		}

		if rotation == 1 {
			robot.MoveToCoordinate(Coordinate{currentX - 1, currentY, 0})
			robot.Direction = "S"
		}
	case "W":
		if rotation == 0 {
			robot.MoveToCoordinate(Coordinate{currentX - 1, currentY, 0})
			robot.Direction = "S"
		}

		if rotation == 1 {
			robot.MoveToCoordinate(Coordinate{currentX + 1, currentY, 0})
			robot.Direction = "N"
		}
	}
}

func (robot *Robot) MoveToCoordinate(coordinate Coordinate) {
	if val, ok := robot.Painted[coordinate.String()]; ok {
		robot.At = val
	} else {
		robot.At = coordinate
	}
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

	robot := paintRobot(intCodes)
	fmt.Println(len(robot.Painted))

}

func paintRobot(intCodes []int) (*Robot) {

	robot := NewRobot()
	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)


	go processIntCodes(intCodes, input, output, needsInput, done)

	paintOrMove := 0

	for {
		select {
		case out := <-output:
			if paintOrMove == 0 {
				robot.PaintPanel(out)
				//fmt.Println("paint", out)
			}

			if paintOrMove == 1 {
				robot.Move(out)
				//fmt.Println("move", out)
			}
			paintOrMove = 1 - paintOrMove
		case <-needsInput:
			input <- robot.GetCurrentPanelColor()
		case <- done:
			return robot
		}
	}
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

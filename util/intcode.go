package util

import (
	"fmt"
)

func ProcessIntCodes(originalIntCodes []int, input chan int, output chan int, needsInput chan bool, done chan bool) []int {

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
				intCodes = (growIntCodeMemory(intCodes, param3 + 1))
			}

			intCodes[param3] = getIntCodeValue(intCodes, param1, param1Mode, relativeBase) +
				getIntCodeValue(intCodes, param2, param2Mode, relativeBase)
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
				intCodes = (growIntCodeMemory(intCodes, param3 + 1))
			}

			intCodes[param3] = getIntCodeValue(intCodes, param1, param1Mode, relativeBase) *
				getIntCodeValue(intCodes, param2, param2Mode, relativeBase)
		}

		if opCode == "03" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]

			if param1Mode == "2" {
				param1 = param1 + relativeBase
			}

			if param1 >= len(intCodes) {
				intCodes = (growIntCodeMemory(intCodes, param1 + 1))
			}

			needsInput <- true
			intCodes[param1] = <-input
		}

		if opCode == "04" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			output <- getIntCodeValue(intCodes, param1, param1Mode, relativeBase)
		}

		if opCode == "05" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]

			if getIntCodeValue(intCodes, param1, param1Mode, relativeBase) != 0 {
				at = getIntCodeValue(intCodes, param2, param2Mode, relativeBase)
				movePointer = false
			}
		}

		if opCode == "06" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]

			if getIntCodeValue(intCodes, param1, param1Mode, relativeBase) == 0 {
				at = getIntCodeValue(intCodes, param2, param2Mode, relativeBase)
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
				intCodes = (growIntCodeMemory(intCodes, param3 + 1))
			}

			if getIntCodeValue(intCodes, param1, param1Mode, relativeBase) <
				getIntCodeValue(intCodes, param2, param2Mode, relativeBase) {
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
				intCodes = (growIntCodeMemory(intCodes, param3 + 1))
			}

			if getIntCodeValue(intCodes, param1, param1Mode, relativeBase) ==
				getIntCodeValue(intCodes, param2, param2Mode, relativeBase) {
				intCodes[param3] = 1
			} else {
				intCodes[param3] = 0
			}
		}

		if opCode == "09" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			relativeBase += getIntCodeValue(intCodes, param1, param1Mode, relativeBase)
		}

		if movePointer {
			at += step[opCode]
		}
	}

	return intCodes
}

func getIntCodeValue(intCodes []int, parameter int, mode string, relativeBase int) int {

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

func growIntCodeMemory(intCodes []int, size int) []int {
	newMemory := make([]int, size, size)
	copy(newMemory, intCodes)

	return newMemory
}

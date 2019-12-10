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

func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	scanner.Scan()
	rawData := scanner.Text()

	intCodeStrings := strings.Split(rawData, ",")
	var intCodes []int
	for _, intCode := range intCodeStrings {
		intCodeInt, _ := strconv.Atoi(intCode)
		intCodes = append(intCodes, intCodeInt)
	}

	programInput := 1
	intCodes = processIntCodes(intCodes, programInput)
}

func processIntCodes(intCodes []int, input int) []int {

	at := 0
	step := map[string]int{
		"01": 4,
		"02": 4,
		"03": 2,
		"04": 2,
	}
	halt := false
	for halt == false {

		instruction := fmt.Sprintf("%05d", intCodes[at])
		opCode := instruction[3:]

		if opCode == "99" {
			break
		}

		if opCode == "01" {

			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at + 3]

			intCodes[param3] = getValue(intCodes, param1, param1Mode) +
				getValue(intCodes, param2, param2Mode)
		}

		if opCode == "02" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at + 2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at + 3]

			intCodes[param3] = getValue(intCodes, param1, param1Mode) *
				getValue(intCodes, param2, param2Mode)
		}

		if opCode == "03" {
			param1 := intCodes[at + 1]
			intCodes[param1] = input
		}

		if opCode == "04" {
			param1 := intCodes[at + 1]
			param1Mode := instruction[2:3]
			fmt.Println("Output:", getValue(intCodes, param1, param1Mode))
		}

		at += step[opCode]

	}

	return intCodes
}

func getValue(intCodes []int, parameter int, mode string) int {
	value := parameter

	if mode == "0" {
		value = intCodes[parameter];
	}

	return value
}

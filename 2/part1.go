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

	// restore gravity assist
	intCodes[1] = 12
	intCodes[2] = 2

	intCodes = processIntCodes(intCodes)

	fmt.Println(intCodes)
}

func processIntCodes(intCodes []int) []int {

	at := 0
	step := 4
	halt := false
	for halt == false {

		opCode := intCodes[at]

		if opCode == 99 {
			break
		}

		input1 := intCodes[at + 1]
		input2 := intCodes[at + 2]
		position := intCodes[at + 3]

		if opCode == 1 {
			intCodes[position] = intCodes[input1] + intCodes[input2]
		}

		if opCode == 2 {
			intCodes[position] = intCodes[input1] * intCodes[input2]
		}

		at += step

	}

	return intCodes
}

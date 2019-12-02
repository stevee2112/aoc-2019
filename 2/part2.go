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
	var originalIntCodes []int
	for _, intCode := range intCodeStrings {
		intCodeInt, _ := strconv.Atoi(intCode)
		originalIntCodes = append(originalIntCodes, intCodeInt)
	}

	for pos1 := 0;pos1 < 100;pos1++ {
		for pos2 := 0;pos2 < 100;pos2++ {

			intCodes := make([]int, len(originalIntCodes))
			copy(intCodes, originalIntCodes)

			output := processIntCodes(intCodes, pos1, pos2)

			if (output == 19690720) {
				fmt.Println(100 * pos1 + pos2)
			}
		}
	}
}

func processIntCodes(intCodes []int, pos1 int, pos2 int) int {

	// restore gravity assist
	intCodes[1] = pos1
	intCodes[2] = pos2

	at := 0
	step := 4
	halt := false
	for halt == false {

		opCode := intCodes[at]

		if opCode == 99 || (opCode != 1 && opCode != 2) {
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

	return intCodes[0]
}

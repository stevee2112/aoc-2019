package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func main() {

	// Get Data
	_, file, _, _ := runtime.Caller(0)

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

	allSeq := permutation([]int{0, 1, 2, 3, 4})

	maxOutput := 0

	for _, seq := range allSeq {
		output := getThrusterOutput(intCodes, seq)
		if output > maxOutput {
			maxOutput = output
		}
	}

	fmt.Println(maxOutput)
}

func getThrusterOutput(intCodes []int, seq []int) int {
	var amps [5]chan int
	for i := range amps {
		amps[i] = make(chan int)
	}

	for i := 0; i < len(amps); i++ {
		prevValue := i
		if i > 0 {
			prevValue = <-amps[i-1]
		}
		go processIntCodes(intCodes, []int{seq[i], prevValue}, amps[i])
	}

	return <-amps[4]

}

func processIntCodes(originalIntCodes []int, inputs []int, output chan int) []int {

	intCodes := make([]int, len(originalIntCodes))
	copy(intCodes, originalIntCodes)

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
	}

	halt := false
	for halt == false {

		movePointer := true
		instruction := fmt.Sprintf("%05d", intCodes[at])
		opCode := instruction[3:]

		if opCode == "99" {
			break
		}

		if opCode == "01" {

			param1 := intCodes[at+1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at+2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at+3]

			intCodes[param3] = getValue(intCodes, param1, param1Mode) +
				getValue(intCodes, param2, param2Mode)
		}

		if opCode == "02" {
			param1 := intCodes[at+1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at+2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at+3]

			intCodes[param3] = getValue(intCodes, param1, param1Mode) *
				getValue(intCodes, param2, param2Mode)
		}

		if opCode == "03" {
			param1 := intCodes[at+1]
			input := inputs[0]
			inputs = inputs[1:]
			intCodes[param1] = input
		}

		if opCode == "04" {
			param1 := intCodes[at+1]
			param1Mode := instruction[2:3]
			outputValue := getValue(intCodes, param1, param1Mode)
			output <- outputValue
		}

		if opCode == "05" {
			param1 := intCodes[at+1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at+2]
			param2Mode := instruction[1:2]

			if getValue(intCodes, param1, param1Mode) != 0 {
				at = getValue(intCodes, param2, param2Mode)
				movePointer = false
			}
		}

		if opCode == "06" {
			param1 := intCodes[at+1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at+2]
			param2Mode := instruction[1:2]

			if getValue(intCodes, param1, param1Mode) == 0 {
				at = getValue(intCodes, param2, param2Mode)
				movePointer = false
			}
		}

		if opCode == "07" {
			param1 := intCodes[at+1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at+2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at+3]

			if getValue(intCodes, param1, param1Mode) <
				getValue(intCodes, param2, param2Mode) {
				intCodes[param3] = 1
			} else {
				intCodes[param3] = 0
			}
		}

		if opCode == "08" {
			param1 := intCodes[at+1]
			param1Mode := instruction[2:3]
			param2 := intCodes[at+2]
			param2Mode := instruction[1:2]
			param3 := intCodes[at+3]

			if getValue(intCodes, param1, param1Mode) ==
				getValue(intCodes, param2, param2Mode) {
				intCodes[param3] = 1
			} else {
				intCodes[param3] = 0
			}
		}

		if movePointer {
			at += step[opCode]
		}

	}

	return intCodes
}

func getValue(intCodes []int, parameter int, mode string) int {
	value := parameter

	if mode == "0" {
		value = intCodes[parameter]
	}

	return value
}

func permutation(xs []int) (permuts [][]int) {
	var rc func([]int, int)
	rc = func(a []int, k int) {
		if k == len(a) {
			permuts = append(permuts, append([]int{}, a...))
		} else {
			for i := k; i < len(xs); i++ {
				a[k], a[i] = a[i], a[k]
				rc(a, k+1)
				a[k], a[i] = a[i], a[k]
			}
		}
	}
	rc(xs, 0)

	return permuts
}

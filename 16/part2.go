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
	var rawData strings.Builder
	rawData.WriteString(scanner.Text())

	// multipler
	multiple := 10000

	for i := 1;i < multiple;i++ {
		rawData.WriteString(scanner.Text())
	}

	offset,_  := strconv.Atoi(rawData.String()[:7])
	sequenceDigits := strings.Split(rawData.String(), "")

	var sequence []int
	for _,digit  := range sequenceDigits {
		intDigit,_ := strconv.Atoi(digit)
		sequence = append(sequence, intDigit)
	}

	fmt.Println(sequenceAsString(FFT(sequence[offset:], 100), 8))
}

func FFT (sequence []int, phases int) []int {

	output := sequence
	for i := 0; i < phases;i++ {
		output = FFTPhase(output)
	}

	return output
}

func FFTPhase(sequence []int) []int {

	newSequence := make([]int, len(sequence))

	lastVal := 0

	for i := len(sequence) - 1; i >= 0;i-- {
		current := sequence[i]
		lastVal += current
		newSequence[i] = lastVal % 10
	}

	return newSequence;
}

func sequenceAsString(sequence []int, size int) string {
	sequenceString := ""

	at := 0
	for _,digit := range sequence {
		char := strconv.Itoa(digit)
		sequenceString += char
		at++

		if at >= size {
			break;
		}
	}

	return sequenceString
}

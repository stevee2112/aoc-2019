package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
	"bufio"
	"strings"
	"strconv"
	"aoc-2019/util"
)

func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	scanner.Scan()
	rawData := scanner.Text()

	sequenceDigits := strings.Split(rawData, "")
	var sequence []int
	for _,digit  := range sequenceDigits {
		intDigit,_ := strconv.Atoi(digit)
		sequence = append(sequence, intDigit)
	}

	fmt.Println(sequenceAsString(FFT(sequence, 100), 8))
}

func FFT (sequence []int, phases int) []int {

	output := sequence
	for i := 0; i < phases;i++ {
		output = FFTPhase(output)
		fmt.Println(i)
	}

	return output
}

func FFTPhase(sequence []int) []int {

	newSequence := []int{}
	pattern := []int{0, 1, 0, -1}

	for i := 0; i < len(sequence);i++ {
		newSequence = append(newSequence, FFTDigit(sequence, pattern, i))
	}

	return newSequence;
}

func FFTDigit(sequence []int, pattern []int, position int) int {

	sum := 0

	for i,digit := range sequence {
		val := digit * FFTPatternValue(pattern, position, i)
		sum += val;
	}

	return util.Abs(sum) % 10;
}

// This function good probably changed to use actual math to speed things up.....
func FFTPatternValue(pattern []int, position int, at int) int{

	expandedPattern := []int{}
	multiple := position + 1

	for _,patternDigit := range pattern {
		for i := 0;i < multiple; i++ {
			expandedPattern = append(expandedPattern, patternDigit)
		}
	}

	// Shift first digit
	expandedPattern = append(expandedPattern, expandedPattern[0])
	expandedPattern = expandedPattern[1:]

	if at >= len(expandedPattern) {
		return expandedPattern[(at % len(expandedPattern))]
	} else {
		return expandedPattern[at]
	}
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

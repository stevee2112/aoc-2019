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
	"container/list"
)

type Path struct {
	Coordinate util.Coordinate
	Direction int
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

	// wake up
	intCodes[0] = 2

	run(intCodes)
}

func run(intCodes []int) {

	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)

	list := list.New()

	commandArray := []string{
		"A,B,A,C,A,B,C,B,C,B", // Main
		"L,10,R,8,L,6,R,6", // A
		"L,8,L,8,R,8", // B
		"R,8,L,6,L,10,L,10", // C
		"n", // video feed
	}

	for _, commandString := range commandArray {
		for _,command := range strings.Split(commandString, "") {
			list.PushBack(command)
		}

		list.PushBack("\n") // newline
	}

	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

programRun:
	for {
		select {
		case out := <-output:

			if out > 127 {
				fmt.Println("Dust vacuumed:", out)
			} else {
				fmt.Print(string(out)) // convert ascii to string
			}
		case <-needsInput:
			if list.Len() > 0 {
				next := list.Front()
				char := next.Value.(string)
				list.Remove(next)
				fmt.Print(char)
				input <- toAscii(char)
			}

			// if we wanted to manually enter values
			//input <- toAscii(util.GetFromStdin())
		case <- done:
			break programRun
		}
	}
}

func toAscii(char string) int {
	charRune := rune(char[0])
	return int(charRune)
}

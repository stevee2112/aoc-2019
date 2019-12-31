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

	run(intCodes)
}

func run(intCodes []int) {

	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)

	list := list.New()

	commandArray := []string{
		"NOT A J", // if a is empty
		"NOT B T",
		"OR T J", // or b is empty
		"NOT C T",
		"OR T J", // or c is empty
		"AND D J", // then jump
		"OR J T", // set t to jump
		"AND E T", // and e is safe
		"OR H T", // or h is safe
		"AND T J",
		"RUN",
	}

	for _, commandString := range commandArray {
		for _,ascii := range stringToAsciiArray(commandString) {
			list.PushBack(ascii)
		}
	}

	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

programRun:
	for {
		select {
		case out := <-output:

			if out > 127 {
				fmt.Println("Damage:", out)
			} else {
				fmt.Print(string(out)) // convert ascii to string
			}
		case <-needsInput:

			// no queued input get from stdin and and to queue
			if list.Len() < 1 {
				for _,ascii := range stringToAsciiArray(util.GetFromStdin()) {
					list.PushBack(ascii)
				}
			}

			if list.Len() > 0 {
				next := list.Front()
				ascii := next.Value.(int)
				list.Remove(next)
				input <- ascii
			}
		case <- done:
			break programRun
		}
	}
}

func toAscii(char string) int {
	charRune := rune(char[0])
	return int(charRune)
}

func stringToAsciiArray(string string) []int {

	ascii := []int{}

	for _, char := range strings.Split(string, "") {
		ascii = append(ascii, toAscii(char))
	}

	// add new line
	ascii = append(ascii, 10)

	return ascii
}

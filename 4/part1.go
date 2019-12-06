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
	boundry := strings.Split(scanner.Text(), "-")

	min, _ := strconv.Atoi(boundry[0])
	max, _ := strconv.Atoi(boundry[1])

	matches := 0

	for at := min; at <= max; at++ {
		asString := strconv.Itoa(at)
		chars := strings.Split(asString, "")
	if match(chars) {
			matches++
		}
	}

	fmt.Println(matches)
}

func match(chars []string) bool {

	hasTwoAdj := false;
	lastChar, chars := chars[0], chars[1:]
	lastInt,_ := strconv.Atoi(lastChar)

	for _,char := range chars {
		charAsInt,_ := strconv.Atoi(char)

		if(charAsInt < lastInt) {
			return false;
		}

		if (charAsInt == lastInt) {
			hasTwoAdj = true
		}
		lastInt = charAsInt;
	}

	return hasTwoAdj
}

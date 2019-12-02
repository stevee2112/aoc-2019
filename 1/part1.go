package main

import (
	"fmt"
	"os"
	"runtime"
	"path"
    "bufio"
	"strconv"
)

func main() {

	// Get Data
	_, file, _,  _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	fuelSum := 0

	for scanner.Scan() {
		mass, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		fuelSum += getFuelNeeded(int(mass))
	}

	fmt.Println(int(fuelSum))
}

func getFuelNeeded(massOrFuel int) int {

	fuelSum := 0
	fuel := (massOrFuel / 3) - 2
	fuelSum += int(fuel)

	return fuelSum
}

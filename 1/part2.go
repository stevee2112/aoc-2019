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
		fuelSum += getFuelNeeded(mass)
	}

	fmt.Println(fuelSum)
}

func getFuelNeeded(massOrFuel int64) int {

	fuelSum := 0
	fuel := (massOrFuel / 3) - 2
	if fuel > 0 {
		fuelSum += getFuelNeeded(fuel)
		fuelSum += int(fuel)
	}

	return fuelSum
}

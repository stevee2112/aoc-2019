package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"strconv"
	"math"
)

type Chemical struct {
	Amount float64
	Name string
}

type ChemicalReaction struct {
	Output Chemical
	Input []Chemical
}

func main() {

	// Get Data
	_, file, _, _ := runtime.Caller(0)

	input, _ := os.Open(path.Dir(file) + "/input")

	defer input.Close()
	scanner := bufio.NewScanner(input)

	reactions := map[string]ChemicalReaction{}

	for scanner.Scan() {
		rawReaction := strings.Split(scanner.Text(), "=>")
		rawInput := strings.Split(strings.Trim(rawReaction[0], " "), ",")
		rawOutput := strings.Trim(rawReaction[1], " ")

		outputAmount,_ := strconv.ParseFloat(strings.Split(rawOutput, " ")[0], 64)
		outputChemical := Chemical{outputAmount, strings.Split(rawOutput, " ")[1]}

		chemicalReaction := ChemicalReaction{outputChemical, []Chemical{}}

		for _,rawInputChemical := range rawInput {
			rawInputChemicalParts :=  strings.Split(strings.Trim(rawInputChemical, " "), " ")
			inputAmount,_ := strconv.ParseFloat(rawInputChemicalParts[0], 64)
			inputChemical := Chemical{inputAmount, rawInputChemicalParts[1]}

			chemicalReaction.Input = append(chemicalReaction.Input, inputChemical)
		}

		reactions[chemicalReaction.Output.Name] = chemicalReaction
	}

	refuse := map[string]float64{}

	oreRemaining := 1000000000000;
	fuel := 0;

	for oreRemaining > 0 {
		oreUsed := int(getOreCost(Chemical{100, "FUEL"}, reactions, refuse))
		oreRemaining -= oreUsed
		fuel += 100

		if oreRemaining < 1000000000 {
			break
		}
	}

	for oreRemaining > 0 {
		oreUsed := int(getOreCost(Chemical{1, "FUEL"}, reactions, refuse))
		oreRemaining -= oreUsed
		fuel += 1
	}


	fmt.Println(fuel - 1)
	fmt.Println(oreRemaining)
	fmt.Println(refuse)
}

func getOreCost(
	chemical Chemical,
	reactions map[string]ChemicalReaction,
	refuse map[string]float64) float64 {

	if chemical.Name == "ORE" {
		return chemical.Amount
	}

	amount := 0.0

	if refuse[chemical.Name] > 0 {
		if refuse[chemical.Name] >= chemical.Amount { // there is enough in excess
			refuse[chemical.Name] -= chemical.Amount
			return 0.0
		} else {
			chemical.Amount -= refuse[chemical.Name]
			refuse[chemical.Name] = 0.0
			return getOreCost(chemical, reactions, refuse)
		}
	}

	reaction := reactions[chemical.Name]
	multiple := math.Ceil(chemical.Amount / reaction.Output.Amount)

	for _,input := range reaction.Input {
		amount += getOreCost(Chemical{multiple * input.Amount, input.Name}, reactions, refuse)
	}

    //Here we put the excess in the refuse map
	surplus := multiple * reaction.Output.Amount - chemical.Amount
    if surplus > 0 {
        refuse[chemical.Name] += surplus
    }

	return amount
}

func isMadeByOre(chemical Chemical, reactions map[string]ChemicalReaction) bool {

	reaction := reactions[chemical.Name]

	if len(reaction.Input) == 1 && reaction.Input[0].Name == "ORE" {
		return true
	}

	return false
}

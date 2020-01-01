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

type Packet util.Coordinate

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

	startAt := 0;
	networkDeviceCount := 50;
	network := map[int]chan Packet{}

	//make network
	for i := startAt; i < networkDeviceCount; i++ {
		network[i] = make(chan Packet)
	}

	for i := startAt; i < networkDeviceCount; i++ {
		go runNetworkComputer(intCodes, i, network, true)
	}

	// select not for - SUPER IMPORTANT
	select{}

}

func runNetworkComputer(intCodes []int, address int, network map[int]chan Packet, logging bool) {

	input := make(chan int)
	output := make(chan int)
	needsInput := make(chan bool)
	done := make(chan bool)

	list := list.New()

	list.PushBack(address) // address

	go util.ProcessIntCodes(intCodes, input, output, needsInput, done)

	outputAt := 0
	addressToSend := 0
	packet := Packet{}
programRun:
	for {
		select {
		case out := <-output:
			if outputAt == 0 {
				addressToSend = out
			} else if outputAt == 1 {
				packet.X = out
			} else if outputAt == 2 {
				packet.Y = out
				outputAt = -1 // restart

				if logging {
					fmt.Println(address, "sending packet to", addressToSend, packet)
				}

				network[addressToSend] <- packet
			}
			outputAt++
		case packet := <- network[address]:
			if logging {
				fmt.Println(address, "recieved packet", packet)
			}
			list.PushBack(packet.X)
			list.PushBack(packet.Y)
		case <-needsInput:
			if list.Len() > 0 {
				next := list.Front()
				intVal := next.Value.(int)
				list.Remove(next)
				input <- intVal
			} else {
				input <- -1
			}
		case <- done:
			break programRun
		}
	}
}

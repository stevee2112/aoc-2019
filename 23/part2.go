package main

import (
	"aoc-2019/util"
	"bufio"
	"container/list"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type Packet util.Coordinate

type IdleStatus struct {
	Address int
	Idle    bool
}

func main() {

	// Get Data
	_, file, _, _ := runtime.Caller(0)

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

	startAt := 0
	networkDeviceCount := 50
	network := map[int]chan Packet{}
	idleChannel := make(chan IdleStatus)
	idleCheck := map[int]bool{}
	natPacket := Packet{}

	//make network
	for i := startAt; i < networkDeviceCount; i++ {
		network[i] = make(chan Packet)
		idleCheck[i] = false
	}

	// add nat to network
	network[255] = make(chan Packet)

	for i := startAt; i < networkDeviceCount; i++ {
		go runNetworkComputer(intCodes, i, network, idleChannel, false)
	}

	for {
		select {
		case packet := <-network[255]:
			natPacket = packet
		case idleStatus := <-idleChannel:
			idleCheck[idleStatus.Address] = idleStatus.Idle

			if checkIfAllIdle(idleCheck) == true && natPacket.X != 0 {
				fmt.Println("All idle sending packet to address 0", natPacket)
				restartPacket := Packet{natPacket.X, natPacket.Y}
				natPacket.X = 0
				go func() {
					network[0] <- restartPacket
				}()
			}
		}
	}
}

func checkIfAllIdle(idleSet map[int]bool) bool {

	allIdle := true

	for _, idleStatus := range idleSet {
		if idleStatus == false {
			allIdle = false
			break
		}
	}

	return allIdle
}

func runNetworkComputer(intCodes []int, address int, network map[int]chan Packet, idleChannel chan IdleStatus, logging bool) {

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
		case packet := <-network[address]:
			if logging {
				fmt.Println(address, "recieved packet", packet)
			}
			list.PushBack(packet.X)
			list.PushBack(packet.Y)
		case <-needsInput:
			if list.Len() > 0 {
				idleChannel <- IdleStatus{address, false}
				next := list.Front()
				intVal := next.Value.(int)
				list.Remove(next)
				input <- intVal
			} else {
				idleChannel <- IdleStatus{address, true}
				input <- -1
			}
		case <-done:
			fmt.Println("here")
			break programRun
		}
	}
}

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/mpps/utils"
)

var numPackets = 100000000
var numThreads = 8

func sender(addr *net.UDPAddr, done chan<- bool, threadIndex, packetCount int) {
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	data := []byte("1234fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc")
	packetIndexStart := threadIndex * packetCount
	packetIndexEnd := packetIndexStart + packetCount
	for i := packetIndexStart; i < packetIndexEnd; i++ {
		utils.AddSequence(data, i)
		utils.AddCheckSum(data)
		_, e := conn.Write(data)
		checkError(e)
	}
	done <- true
}

func main() {
	serverAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 22222, Zone: ""}
	done := make(chan bool, numThreads)
	for i := 0; i < numThreads; i++ {
		go sender(&serverAddr, done, i, numPackets/numThreads)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

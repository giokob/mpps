package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mpps/utils"
)

var byteCount = 122

func dataReceiver(conn *net.UDPConn, dataChannel chan<- []byte) {
	data := make([]byte, byteCount)
	for {
		_, _, err := conn.ReadFromUDP(data)
		checkError(err)
		dataChannel <- data
	}
}

func dataAnalyzer(dataChannel chan []byte) {
	totalCounter := 0
	lastSecondCounter := 0
	invalid := 0
	ticket := time.NewTicker(time.Second)
	for {
		select {
		case <-ticket.C:
			fmt.Printf("total: %d, lastSec: %d, invalid: %d\n", totalCounter, lastSecondCounter, invalid)
			lastSecondCounter = 0
			invalid = 0
		case data, ok := <-dataChannel:
			if !ok {
				fmt.Println("closing channel")
				close(dataChannel)
			}
			if !utils.CheckSum(data) {
				invalid++
			}

			lastSecondCounter++
			totalCounter++
		}
	}
}

func main() {
	addr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 22222, Zone: ""}
	conn, err := net.ListenUDP("udp", &addr)
	checkError(err)
	defer conn.Close()

	dataChannel := make(chan []byte, 10)
	go dataReceiver(conn, dataChannel)
	dataAnalyzer(dataChannel)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

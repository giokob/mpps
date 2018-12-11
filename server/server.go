package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mpps/utils"
)

func dataReciver(conn *net.UDPConn, dataChannel chan<- []byte) {
	data := make([]byte, 122)
	for {
		_, _, err := conn.ReadFromUDP(data)
		checkError(err)
		dataChannel <- data
	}
}

func main() {
	addr := net.UDPAddr{net.ParseIP("127.0.0.1"), 22222, ""}
	conn, err := net.ListenUDP("udp", &addr)
	checkError(err)
	defer conn.Close()

	totalCounter := 0
	lastSecondCounter := 0
	invalid := 0
	ticket := time.NewTicker(time.Second)
	dataChannel := make(chan []byte, 10)
	go dataReciver(conn, dataChannel)

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

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

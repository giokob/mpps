package main

import (
	"fmt"
	"net"
	"os"

	"github.com/mpps/utils"
)

var NUMPACKETS int = 100000000
var NUMTHREADS int = 8

func sender(addr *net.UDPAddr, done chan<- int) {
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	data := []byte("1234fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc")
	for i := 0; i < NUMPACKETS/NUMTHREADS; i++ {
		utils.AddSequence(data, i)
		utils.AddCheckSum(data)
		// fmt.Println(data)
		// fmt.Println(readSequence(data))
		// fmt.Println(checkSum(data))
		_, e := conn.Write(data)
		checkError(e)
	}
	done <- 1
}

func main() {
	serverAddr := net.UDPAddr{net.ParseIP("127.0.0.1"), 22222, ""}
	done := make(chan int, NUMTHREADS)
	for i := 0; i < NUMTHREADS; i++ {
		go sender(&serverAddr, done)
	}

	for i := 0; i < NUMTHREADS; i++ {
		<-done
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

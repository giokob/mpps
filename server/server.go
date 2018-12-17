package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"github.com/mpps/utils"
)

var byteCount = 122

func dataReceiver(dataChannel chan<- []byte) {
	lc := net.ListenConfig{
		Control: listenCtrl,
	}
	conn, err := lc.ListenPacket(context.Background(), "udp4", "127.0.0.1:22222")
	checkError(err)
	defer conn.Close()

	data := make([]byte, byteCount)
	for {
		// _, _, err := conn.ReadFromUDP(data)
		_, _, err := conn.ReadFrom(data)
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

func listenCtrl(network string, address string, c syscall.RawConn) error {
	var operr error
	var fn = func(s uintptr) {
		operr = unix.SetsockoptInt(int(s), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	}

	err := c.Control(fn)
	if err != nil {
		return err
	}
	if operr != nil {
		return operr
	}
	return nil
}

func main() {

	// addr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 22222, Zone: ""}
	// conn, err := net.ListenUDP("udp", &addr)
	// checkError(err)
	// defer conn.Close()

	dataChannel := make(chan []byte, 10)
	go dataReceiver(dataChannel)
	go dataReceiver(dataChannel)
	go dataReceiver(dataChannel)
	go dataReceiver(dataChannel)
	dataAnalyzer(dataChannel)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

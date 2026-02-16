package main

import (
	"fmt"
	"net"
	"netcat/utils"
	"os"
	"strconv"
)

func main() {

	// Default server port
	var port string = ":8989"
	// check command-line arguments
	if len(os.Args) > 2 {
		fmt.Print("usage: go run . port")
		return
	}

	// if a port is provided, use it
	if len(os.Args) == 2 {
		// convert the port from string to integer
		p, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Print("usage: go run . port")
			return
		}
		// verify that the port is within the allowed range (1024 - 49151)
		// 0–1023 are reserved (system ports)
		// 49152+ are dynamic/ephemeral ports
		if p >= 1024 && p <= 49151 {
			port = ":" + os.Args[1] // format the port for net.Listen (":8080")
		} else {
			fmt.Print("usage: go run . port")
			return
		}
	}

	fmt.Println("Listening on the port", port)
	// creat communication channels
	clientChannel := make(chan utils.Client)
	messageChannel := make(chan utils.Message)
	validNameChannel := make(chan bool)
	// limit maximum connected clients to 10 using buffered channel
	slots := make(chan struct{}, 10)
	// start TCP server listner
	ln, _ := net.Listen("tcp", port)
	// start ChatManager in a separete goroutine
	go utils.ChatManager(clientChannel, messageChannel, validNameChannel)
	// main loop to accept incoming connections
	for {
		Conn, err := ln.Accept()
		if err != nil {
			continue
		}
		// check if there is room, reserve a slot and handle connection concurently

		select {
		case slots <- struct{}{}:
			go utils.HandleConn(Conn, clientChannel, messageChannel, validNameChannel, slots)

		// if server is full, reject connection
		default:
			fmt.Fprint(Conn, "room is full, try later\n")
			Conn.Close()
		}
	}

}

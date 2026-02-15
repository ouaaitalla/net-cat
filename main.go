package main

import (
	"fmt"
	"net"
	"netcat/utils"
	"os"
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
		port = ":" + os.Args[1]
	}

	fmt.Println("server started in port", port)
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

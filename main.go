package main

import (
	"fmt"
	"net"
	"os"
	"netcat/utils"
)

func main() {
	var port string = ":8989"
	if len(os.Args) > 2 {
		fmt.Print("usage: go run . port")
		return
	}
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	clientChannel := make(chan utils.Client)
	messageChannel := make(chan utils.Message)
	validNameChannel := make(chan bool)
	mainChannel := make(chan bool)

	ln, err := net.Listen("tcp", port)
if err != nil {
	fmt.Println("Failed to listen:", err)
	return
}

	go utils.ChatManager(clientChannel, messageChannel, validNameChannel, mainChannel)
	for {
		Conn, err := ln.Accept()
		if err != nil {
			continue
		}

		checkClientmax := utils.Client{
			Name: "",
			Conn: nil,
		}
		clientChannel <- checkClientmax
		msg := <-mainChannel

		if msg {
			fmt.Fprint(Conn, "room is full, try later\n")
			Conn.Close()
			continue
		}
		go utils.HandleConn(Conn, clientChannel, messageChannel, validNameChannel, mainChannel)
	}
}

package main

import (
	"fmt"
	"net"
	"netcat/utils"
	"os"
)

func main() {
	var port string = ":8989"
	if len(os.Args) > 2 {
		fmt.Print("usage: go run . port")
		return
	}
	if len(os.Args) == 2 {
		port = ":" +os.Args[1]
	}
	fmt.Println("server started in port", port)
	clientChannel := make(chan utils.Client)
	messageChannel := make(chan utils.Message)
	validNameChannel := make(chan bool)
	ln, _ := net.Listen("tcp", port)
	go utils.ChatManager(clientChannel, messageChannel, validNameChannel)
	var conn net.Conn
	for {
		conn, _ = ln.Accept()
		go utils.HandleConn(conn, clientChannel, messageChannel,validNameChannel)
	}
}

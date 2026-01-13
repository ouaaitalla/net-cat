package main

import (
	"net"
	"netcat/utils"
	)

func main() {
	clientChannel := make(chan utils.Client)
	messageChannel := make(chan utils.Message)
	ln, _ := net.Listen("tcp", ":8080")
	go utils.ChatManager(clientChannel, messageChannel)
	var conn net.Conn
	for {
		conn, _ = ln.Accept()
		go utils.HandleConn(conn, clientChannel, messageChannel)
	}
}

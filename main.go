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
		port = ":" + os.Args[1]
	}
	clientChannel := make(chan utils.Client)
	messageChannel := make(chan utils.Message)
	validNameChannel := make(chan bool)
	slots := make(chan struct{}, 3)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	go utils.ChatManager(clientChannel, messageChannel, validNameChannel)
for {
    Conn, err := ln.Accept()
    if err != nil {
        continue
    }

    select {
    case slots <- struct{}{}:
        go func() {
            utils.HandleConn(Conn, clientChannel, messageChannel, validNameChannel, slots)
        }()

    default:
        fmt.Fprint(Conn, "room is full, try later\n")
        Conn.Close()
    }
}
}

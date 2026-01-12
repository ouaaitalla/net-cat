package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type client struct {
	name string
	conn net.Conn
}

type Message struct {
	textMessage string
	conn        net.Conn
}

var frr string

func formatMessage(name, message string) string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s", currentTime, name, message)
}

func main() {

	clientChannel := make(chan client)
	messageChannel := make(chan Message)

	ln, _ := net.Listen("tcp", ":8080")

	go chatManager(clientChannel, messageChannel)

	for {
		conn, _ := ln.Accept()
		go handleConn(conn, clientChannel, messageChannel)
	}
}

func handleConn(conn net.Conn, clientChannel chan client, messageChannel chan Message) {
	var clientName string
	reader := bufio.NewReader(conn)
	logo, err := os.ReadFile("logolinux.txt")
	if err != nil {

	}
	for {
		if clientName != "" {
			frr = formatMessage(clientName, "")
			fmt.Fprint(conn, frr)
			message, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF{
					cl := client{
						name: clientName,
						conn : nil,
					}
					clientChannel <- cl
					return
				}
			}
			message = formatMessage(clientName, message)
			messageStruct := Message{
				textMessage: message,
				conn:        conn,
			}
			messageChannel <- messageStruct

		} else {
			fmt.Fprint(conn, "welcom to tcp chat\n")
			fmt.Fprint(conn, string(logo))
			fmt.Fprint(conn, "Enter your name : ")
			message, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			message = strings.TrimSpace(message)
			if message != "" {
				clientName = message
				if len(clientName) <= 25 {
					cl := client{
						name: clientName,
						conn: conn,
					}
					clientChannel <- cl
				} else {
					fmt.Fprint(conn, "cannot use name longer then 25 caracter\n")
				}
			} else {
				fmt.Fprint(conn, "cannot use an empty name\n")
			}
		}
	}
}

func chatManager(clientChannel chan client, messageChannel chan Message) {
	clients := make([]client, 0)
	var chatHistory string
	for {
		select {
		case clientInfo := <-clientChannel:
			if clientInfo.conn != nil {
				clients = append(clients, clientInfo)
				fmt.Fprint(clientInfo.conn, chatHistory)
			} else {
				removeClient(&clients, clientInfo)
				leftMessage := clientInfo.name + "has left chat\n"
				leftMsgStruct := Message{
					textMessage: leftMessage,
					conn : nil,
				}
				broadCast(leftMsgStruct,clients)
			}
		case clientMessage := <-messageChannel:
			broadCast(clientMessage, clients)
			chatHistory += clientMessage.textMessage + "\n"
			
		}
	}
}

func broadCast(clientMessage Message, clients []client) {
	var msg string
	for _, client := range clients {
		if client.conn != clientMessage.conn {
			msg = formatMessage(client.name, "")
			fmt.Fprint(client.conn, "\n"+clientMessage.textMessage)
			fmt.Fprint(client.conn, msg)
		}
	}
}

func removeClient(clients *[]client, clientInfo client) {
	for i, c := range *clients {
		if c.name == clientInfo.name {
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
			return
		}
	}
}

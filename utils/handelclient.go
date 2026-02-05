package utils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var clients []Client

func HandleConn(Conn net.Conn, clientChannel chan Client, messageChannel chan Message, validNameChannel chan bool, mainChannel chan bool) {
	var clientName string
	reader := bufio.NewReader(Conn)
	logo, _ := os.ReadFile("logolinux.txt")
	if len(clients) > 9 {
		fmt.Fprint(Conn, "room chat is full try later")
		Conn.Close()
	}
	for {
		Form = formatMessage(clientName, "")
		if clientName != "" {
			message, err := reader.ReadString('\n')
			if message == "\n"{
				fmt.Fprint(Conn, Form)
				continue
			}
			if err != nil {
				if err == io.EOF {
					cl := Client{
						Name: clientName,
						Conn: nil,
					}
					clientChannel <- cl
					return
				}
			}
			Form = formatMessage(clientName, "")
			fmt.Fprint(Conn, Form)
			message = formatMessage(clientName, message)
			messageStruct := Message{
				textMessage: message,
				Conn:        Conn,
			}
			messageChannel <- messageStruct
		} else {
			fmt.Fprint(Conn, "Welcome to TCP-Chat!\n")
			fmt.Fprint(Conn, string(logo))
			fmt.Fprint(Conn, "[ENTER YOUR Name]:")
			message, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			message = strings.TrimSpace(message)
			if message != "" {
				if len(message) <= 25 {
					cl := Client{
						Name: message,
						Conn: Conn,
					}
					clientChannel <- cl
					v := <-validNameChannel
					if v {
						clientName = message
					}
				} else {
					fmt.Fprint(Conn, "cannot use Name longer then 25 caracter\n")
				}
			} else {
				fmt.Fprint(Conn, "cannot use an empty Name\n")
			}
		}
	}
}

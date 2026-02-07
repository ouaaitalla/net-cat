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

func HandleConn(conn net.Conn, clientChannel chan Client, messageChannel chan Message, validNameChannel chan bool) {
	var clientName string
	reader := bufio.NewReader(conn)
	logo, _ := os.ReadFile("logolinux.txt")
	for {
		var Form string
		Form = formatMessage(clientName, "")
		if clientName != "" {
			message, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					cl := Client{
						name: clientName,
						conn: nil,
					}
					clientChannel <- cl
					return
				}
			}
			if message == "\n" {
				fmt.Fprint(conn, Form)
				continue
			}
			
			fmt.Fprint(conn, Form)
			message = formatMessage(clientName, message)
			messageStruct := Message{
				textMessage: message,
				conn:        conn,
			}
			messageChannel <- messageStruct
		} else {
			fmt.Fprint(conn, "Welcome to TCP-Chat!\n")
			fmt.Fprint(conn, string(logo))
			fmt.Fprint(conn, "[ENTER YOUR NAME]:")
			message, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			message = strings.TrimSpace(message)
			if message != "" {
				if len(message) <= 25 {
					cl := Client{
						name: message,
						conn: conn,
					}
					clientChannel <- cl
					v := <-validNameChannel
					if v {
						clientName = message
					}
				} else {
					fmt.Fprint(conn, "cannot use name longer then 25 caracter\n")
				}
			} else {
				fmt.Fprint(conn, "cannot use an empty name\n")
			}
		}
	}
}

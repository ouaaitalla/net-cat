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
	var Form string
	reader := bufio.NewReader(conn)
	logo, _ := os.ReadFile("logolinux.txt")
	for {

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
			if !isValidASCII(message) || strings.TrimSpace(message) == "" || message == "\n" {
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
				if len(message) <= 25 && isValidASCII(message) {
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
					fmt.Fprint(conn, "cannot use name longer then 15 caracter or caracter not printable\n")
				}
			} else {
				fmt.Fprint(conn, "cannot use an empty name\n")
			}
		}
	}
}

func isValidASCII(s string) bool {
	for _, r := range s {
		if r == '\n' {
			continue
		}
		if r >= 32 && r <= 126 {
			continue
		}
		return false
	}
	return true
}

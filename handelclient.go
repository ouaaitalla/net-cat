// package main
// import (
//     "bufio"
//     "fmt"
//     "net"
//     "strings"
// )

// func handleConn(client *client, ch chan string) {
//     reader := bufio.NewReader(client.conn)
//     var chatHistory string
//     go func() {
//         for msg := range client.ch{
//             if client.name != ""{
//                 if !strings.HasPrefix(msg,client.name) {
//                     fmt.Fprint(client.conn,"\n" + msg + client.name + " : ")
//                 }
//             } else{
//                 chatHistory += msg
//             }
//         }
//     }()
//     for {
//         if client.name != ""{
//             fmt.Fprint(client.conn,client.name + " : ")
//             message, err := reader.ReadString('\n')
//             if err != nil {
//                 return
//             }
//             message = client.name + " : "+ message
//             ch <- message
//         } else{
//             fmt.Fprint(client.conn,"Enter your name : ")
//             message, err := reader.ReadString('\n')
//             message = strings.TrimSpace(message)
//             if err != nil {
//                 return
//             }
//             if message != ""{
//                 client.name = message
//                 fmt.Fprint(client.conn,chatHistory)
//             }
//         }
//     }
// }

// type client struct {
//     name string
//     conn net.Conn
//     ch   chan string
// }

// func main() {
//     var chatHistory string
//     var clients     []*client
//     ch := make(chan string)
//     historyChan := make(chan string)
//     newClientChan := make(chan string)
//     ln,_ := net.Listen("tcp", ":8080")
//     go func(){
//         for msg := range ch{
//             for _,cl := range clients{
//                 cl.ch <- msg
//             }
//         }
//     }()
//     go func(){
//         var ChatHistory string
//         for{
//             select{
//             case msg := <- historyChan:
//                 ChatHistory += msg
//             case notification := <- newClientChan:
//                 if notification != ""{
//                     historyChan <- ChatHistory
//                 }
//             }
//         }
//     }()
//     for {
//         conn, _ := ln.Accept()
//         c := &client{
//         conn: conn,
//         ch:   make(chan string),
//         }
//         clients = append(clients,c)
//         go handleConn(c,ch)
//         newClientChan <- "new"
//         chatHistory = <- historyChan
//         c.ch <- chatHistory
//     }
// }
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type client struct {
	name string
	conn net.Conn
	ch   chan string
}
type Message struct {
	text string;
	sender *client
}
	var frr string
func handleConn(client *client, ch chan Message, chatHistory *[]string) {
	reader := bufio.NewReader(client.conn)
	go func() {
		for msg := range client.ch {
			if client.name != "" {
				if !strings.HasPrefix(msg, client.name) {
					fmt.Fprint(client.conn, "\n"+msg)
					fmt.Fprint(client.conn, frr)
				}
			} else {
				*chatHistory = append(*chatHistory, msg)
			}
		}
	}()

	for {
		if client.name != "" {

			frr = formatMessage(client.name, "")
			fmt.Fprint(client.conn, frr)
			message, err := reader.ReadString('\n')
			// fmt.Fprint(client.conn, frr)
			if err != nil {
				return
			}
			message = strings.TrimSpace(message)
			if message != "" {
				message += "\n"
				formatted := formatMessage(client.name, message)
				ch <- Message{
					text:   formatted,
					sender: client,
				}
				bin := strings.TrimSuffix(formatted, "\n")
				*chatHistory = append(*chatHistory,(bin))
			}
		} else {
			logo, err := os.ReadFile("logolinux.txt")
			if err != nil {
				continue
			}
			fmt.Fprint(client.conn, "welcom to tcp chat\n")
			fmt.Fprint(client.conn, string(logo))
			fmt.Fprint(client.conn, "Enter your name : ")
			message, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			message = strings.TrimSpace(message)
			if message != "" {
				client.name = message

				for _, msg := range *chatHistory {
					fmt.Fprintln(client.conn, msg)
				}
			}
		}
	}
}

func formatMessage(name, message string) string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s", currentTime, name, message)
}

func main() {
	var chatHistory []string
	var clients []*client

	ch := make(chan Message)

	ln, _ := net.Listen("tcp", ":8080")

	go func() {
		for msg := range ch {
			for _, cl := range clients {
				if cl != msg.sender {
					cl.ch <- msg.text
				}
			}
		}
	}()

	for {
		conn, _ := ln.Accept()

		c := &client{
			conn: conn,
			ch:   make(chan string, 10),
		}
		clients = append(clients, c)

		go handleConn(c, ch, &chatHistory)
	}
}

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

// HandleConn handles a single client TCP connection.
// Parameters:
// - conn: the TCP connection
// - clientChannel: used to notify ChatManager about new/disconnected clients
// - messageChannel: used to send messages to ChatManager
// - validNameChannel: used to receive validation result for client name
func HandleConn(conn net.Conn, clientChannel chan Client, messageChannel chan Message, validNameChannel chan bool, slots chan struct{}) {

	// This defer runs when the function exits.
	// It closes the connection and releases a slot (for limiting max clients).
	defer func() {
		conn.Close()
		<-slots
	}()

	var clientName string                   // stores the client's name
	var Form string                         // stores the formatted prompt
	reader := bufio.NewReader(conn)         // buffered reader to read input line by line
	logo, _ := os.ReadFile("logolinux.txt") // read logo from file

	for {

		// Prepare the prompt format using the client name
		Form = formatMessage(clientName, "")
		// If the client already has a valid name
		if clientName != "" {

			// Read a message from the client
			message, err := reader.ReadString('\n')
			if err != nil {
				//if client disconnected , notify ChatManager
				if err == io.EOF {
					cl := Client{
						name: clientName,
						conn: nil, // nil means client disconnected
					}
					clientChannel <- cl
					return
				}
			}
			// valid message (must be pritable ASCII and not empty)
			if !isValidASCII(message) || strings.TrimSpace(message) == "" || message == "\n" {
				fmt.Fprint(conn, Form) // resend prompt
				continue
			}
			// send prompt back to client
			fmt.Fprint(conn, Form)
			// Format message before broadcasting
			message = formatMessage(clientName, message)
			// creat message struct and send before broadcasting
			messageStruct := Message{
				textMessage: message,
				conn:        conn,
			}
			messageChannel <- messageStruct
		} else {
			// if client does not yet have a name, ask for it
			fmt.Fprint(conn, "Welcome to TCP-Chat!\n")
			fmt.Fprint(conn, string(logo))
			fmt.Fprint(conn, "[ENTER YOUR NAME]:")
			// Read name input
			message, err := reader.ReadString('\n')
			if err != nil {
				return // connection closed
			}
			message = strings.TrimSpace(message)
			// validate name (not empty , <= 25 chars, printable ASCII)
			if message != "" {
				if len(message) <= 25 && isValidASCII(message) {
					// send proposed name to ChatManager
					cl := Client{
						name: message,
						conn: conn,
					}
					clientChannel <- cl
					// wait for validation result (true = accepted)
					v := <-validNameChannel
					if v {
						clientName = message // set client name if valid
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

// isValidASCII returns true if the string contains only printable ASCII characters (32–126)
// and newline,
// otherwise false.
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

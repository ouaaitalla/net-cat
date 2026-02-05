package utils

import "fmt"

func broadCast(clientMessage Message, clients []Client) {
	var msg string
	for _, client := range clients {
		if client.Conn != clientMessage.Conn {
			msg = formatMessage(client.Name, "")
			fmt.Fprint(client.Conn, "\n"+clientMessage.textMessage)
			fmt.Fprint(client.Conn, msg)
		}
	}
}

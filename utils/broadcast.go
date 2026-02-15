package utils

import "fmt"


// broadCast sends a message to all connected clients
// except the sender of the message.
func broadCast(clientMessage Message, clients []Client) {
	var msg string
	for _, client := range clients {
		// do not send the messge back to the sender 
		if client.conn != clientMessage.conn {
			// prepare the prompt format for the receiving client 
			msg = formatMessage(client.name, "")

			// snd the message and then resend the prompt 
			fmt.Fprint(client.conn, "\n"+clientMessage.textMessage)
			fmt.Fprint(client.conn, msg)
		}
	}
}

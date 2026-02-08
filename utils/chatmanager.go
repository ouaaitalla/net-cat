package utils

import "fmt"

func ChatManager(clientChannel chan Client, messageChannel chan Message, validNameChannel chan bool) {
	var form string
	var chatHistory string
	var clients []Client
	for {
		select {
		case clientInfo := <-clientChannel:
			if clientInfo.Conn != nil {
				if IsUnicName(clientInfo.Name, clients) {
					fmt.Fprint(clientInfo.Conn, "this Name already exist\n")
					validNameChannel <- false
					break
				} else {
					validNameChannel <- true
				}
				clients = append(clients, clientInfo)
				jM := clientInfo.Name + " has joined a chat \n"
				joinMessage := Message{
					textMessage: jM,
					Conn:        clientInfo.Conn,
				}

				broadCast(joinMessage, clients)
				fmt.Fprint(clientInfo.Conn, chatHistory)
				form = formatMessage(clientInfo.Name, "")
				fmt.Fprint(clientInfo.Conn, form)
			} else {
				removeClient(&clients, clientInfo)
				leftMessage := clientInfo.Name + " has left chat\n"
				leftMsgStruct := Message{
					textMessage: leftMessage,
					Conn:        nil,
				}
				broadCast(leftMsgStruct, clients)
			}
		case clientMessage := <-messageChannel:
			broadCast(clientMessage, clients)
			chatHistory += clientMessage.textMessage

		}
	}
}

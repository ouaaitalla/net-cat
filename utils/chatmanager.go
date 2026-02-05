package utils

import "fmt"

func ChatManager(clientChannel chan Client, messageChannel chan Message, validNameChannel chan bool, mainChannel chan bool) {
	var Form string
	var chatHistory string
	for {
		select {
		case clientInfo := <-clientChannel:
			if clientInfo.Name == "" && clientInfo.Conn == nil {
				if len(clients) <= 10 {
					mainChannel <- false
				} else {
					mainChannel <- true
				}
				continue
			}
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
				Form = formatMessage(clientInfo.Name, "")
				fmt.Fprint(clientInfo.Conn, Form)
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

func IsUnicName(Name string, clients []Client) bool {
	for _, client := range clients {
		if client.Name == Name {
			return true
		}
	}
	return false
}

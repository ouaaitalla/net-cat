package utils

import "fmt"

func ChatManager(clientChannel chan Client, messageChannel chan Message, validNameChannel chan bool) {
	var chatHistory string
	for {
		select {
		case clientInfo := <-clientChannel:
			if clientInfo.conn != nil {
				if IsUnicName(clientInfo.name, clients){
					fmt.Fprint(clientInfo.conn,"this name already exist\n")
					validNameChannel <- false
					break
				} else{
					validNameChannel <- true
				}
				clients = append(clients, clientInfo)
				jM := clientInfo.name + " has joined a chat \n"
				joinMessage := Message{
					textMessage: jM,
					conn: clientInfo.conn,
				}
				
				broadCast(joinMessage, clients)
				fmt.Fprint(clientInfo.conn, chatHistory)
				var form string
				form = formatMessage(clientInfo.name, "")
				fmt.Fprint(clientInfo.conn, form)
			} else {
				removeClient(&clients, clientInfo)
				leftMessage := clientInfo.name + " has left chat\n"
				leftMsgStruct := Message{
					textMessage: leftMessage,
					conn:        nil,
				}
				broadCast(leftMsgStruct, clients)
			}
		case clientMessage := <-messageChannel:
			broadCast(clientMessage, clients)
			chatHistory += clientMessage.textMessage

		}
	}
}

func IsUnicName(name string,clients []Client)bool{
	for _, client := range clients{
		if client.name ==  name {
			return true
		}
	}
	return false
}

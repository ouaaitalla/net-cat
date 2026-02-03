package utils

import "fmt"

func ChatManager(clientChannel chan Client, messageChannel chan Message) {
	var chatHistory string
	for {
		select {
		case clientInfo := <-clientChannel:
			if clientInfo.conn != nil {
				if IsUnicName(clientInfo.name){
					fmt.Fprint(clientInfo.conn,"this name already exist")
					clientInfo.conn.Close()
				}
				clients = append(clients, clientInfo)
				LNC = append(LNC, clientInfo.name)
				jM := clientInfo.name + " has joined a chat \n"
				joinMessage := Message{
					textMessage: jM,
					conn: clientInfo.conn,
				}
				fmt.Fprint(clientInfo.conn, chatHistory)
				broadCast(joinMessage, clients)
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

func IsUnicName(name string)bool{
	for _, ls := range LNC {
		if ls ==  name {
			return true
		}
	}
	return false
}

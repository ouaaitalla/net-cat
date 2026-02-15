package utils

import "fmt"


// ChatManager is responsible for managing all connected clients,
// handling join/leave events, and broadcasting messages.
func ChatManager(clientChannel chan Client, messageChannel chan Message, validNameChannel chan bool) {
	var chatHistory string // stores all previous chat messages 
	var form string // stores formatted prompt  
	for {
		select {
		// handel new client connection or client disconnection 
		case clientInfo := <-clientChannel:
			// if conn is not nil client is is joining 
			if clientInfo.conn != nil {
				// check if the username already exists   
				if IsUnicName(clientInfo.name, clients) {
					fmt.Fprint(clientInfo.conn, "this name already exist\n")
					validNameChannel <- false
					break
				} else {
					validNameChannel <- true
				}
				// add client to the lisgt 
				clients = append(clients, clientInfo)
				// creat join notification message 
				jM := clientInfo.name + " has joined a chat \n"
				joinMessage := Message{
					textMessage: jM,
					conn:        clientInfo.conn,
				}
				// broadcast join message to all clients 
				broadCast(joinMessage, clients)
				// send previous cat history to the client 
				fmt.Fprint(clientInfo.conn, chatHistory)
				// send prompt format to the new client 
				form = formatMessage(clientInfo.name, "")
				fmt.Fprint(clientInfo.conn, form)
			} else {
				// if conn is nil client has disconnected 
				// remove client from the list 
				removeClient(&clients, clientInfo)
				// creat left notification message 
				leftMessage := clientInfo.name + " has left chat\n"
				leftMsgStruct := Message{
					textMessage: leftMessage,
					conn:        nil,
				}
				// broadcast leave message 
				broadCast(leftMsgStruct, clients)
			}
		// handle incoming chat messages 
		case clientMessage := <-messageChannel:
			// Broadcast message to all connected clients 
			broadCast(clientMessage, clients)
			// append messsage to chat history 
			chatHistory += clientMessage.textMessage

		}
	}
}


// IsUnicName checks  if a user name  alrady exists in the clients slice.
// return true if the name exists, otherwise false 
func IsUnicName(name string, clients []Client) bool {
	for _, client := range clients {
		if client.name == name {
			return true
		}
	}
	return false
}

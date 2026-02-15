package utils

// removeClient removes a client from the clients slice by matching the client name.
func removeClient(clients *[]Client, clientInfo Client) {
	for i, c := range *clients {
		if c.name == clientInfo.name {
			// Remove the client at index i by slicing and appending
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
			return
		}
	}
}

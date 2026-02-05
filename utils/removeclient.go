package utils

func removeClient(clients *[]Client, clientInfo Client) {
	for i, c := range *clients {
		if c.Name == clientInfo.Name {
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
			return
		}
	}
}

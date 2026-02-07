package utils

import "net"

type Client struct {
	Name string
	Conn net.Conn
}

type Message struct {
	textMessage string
	Conn        net.Conn
}


func IsUnicName(Name string, clients []Client) bool {
	for _, client := range clients {
		if client.Name == Name {
			return true
		}
	}
	return false
}

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

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

var Form string

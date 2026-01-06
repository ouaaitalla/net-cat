package netcat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
    Name string
    Conn net.Conn
}

func HandelClient(conn net.Conn) {
	defer conn.Close()
	linuxLogo, err := os.ReadFile("logolinux.txt")
	if err != nil {
		return
	}
	fmt.Fprint(conn, "Welcome to TCP-Chat!\n")
	fmt.Fprint(conn, linuxLogo)
	fmt.Fprint(conn, "[ENTER YOUR NAME]: ")
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Fprintln(conn, "name cannot be empty")
		return
	}
	client := &Client{
    Name: name,
    Conn: conn,
	}
	AddClient(client)
	SendHistory(conn)

	joinMsg := fmt.Sprintf("%s has joined our chat...\n", client.Name)
	Broadcast(joinMsg, conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		if strings.HasPrefix(msg, "/name ") {
    	newName := strings.TrimSpace(strings.TrimPrefix(msg, "/name "))

    	if newName == "" {
        	fmt.Fprintln(c.Conn, "Name cannot be empty")
        	continue
    	}

    	oldName := c.Name
    	c.Name = newName

    	notify := fmt.Sprintf("%s changed name to %s\n", oldName, newName)
    	Broadcast(notify, c.Conn)

    	continue
	}

		formatted := formatMessage(client.Name, msg)
		SaveMessage(formatted)
		Broadcast(formatted, nil)
		RemoveClient(conn)
		leaveMsg := fmt.Sprintf("%s has left our chat...\n", client.Name)
		Broadcast(leaveMsg, conn)
	}
}

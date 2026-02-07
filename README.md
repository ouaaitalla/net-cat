
# TCP Chat – NetCat Style Group Chat (Go)

## Description

This project is a simplified reimplementation of NetCat as a TCP group chat server written in Go.
It allows multiple clients to connect to a server, choose a username, and exchange messages in real time.

The server manages connections concurrently using goroutines and channels, broadcasts messages to all connected clients, and keeps a history of the chat for new users joining later.

The behavior is similar to `nc` but focused on group chat features.

---

## Features

* TCP server with multi-client support (1 → many)
* Concurrent connection handling with goroutines
* Communication using channels
* Username required on connection
* Username must be:

  * non-empty
  * ASCII printable
  * unique
  * limited length
* Maximum 10 connections supported
* Chat message format with timestamp and username:

```
[YYYY-MM-DD HH:MM:SS][username]:message
```

* Join notification broadcast
* Leave notification broadcast
* Chat history sent to newly connected clients
* Empty messages are ignored
* Non-ASCII messages are rejected
* Default port = 8989
* Custom port supported
* Linux ASCII logo displayed on connection

---

## Project Structure

```
.
├── main.go
├── utils/
|   ├── broadcast.go
|   ├── chatmanager.go
|   ├── formatmessage.go
|   ├── handleclient.go
|   ├── logolinux.txt
|   ├── removeclient.go
|   └── helpers.go
├── go.mod
└── README.md
```

### Main Components

**main.go**

* Parses port argument
* Starts TCP listener
* Creates channels
* Launches ChatManager goroutine
* Accepts connections concurrently using goroutines
* Validates port usage

**ChatManager**

* Central coordinator
* Manages:

  * connected clients
  * name validation
  * join/leave events
  * message broadcasting
  * chat history

**HandleConn**

* Handles each client in a goroutine
* Shows Linux logo + welcome message
* Requests username
* Reads client messages
* Sends messages to manager via channel

**Broadcast**

* Sends messages to all clients except sender
* Adds prompt format after each message

---

## Usage

### Start server (default port)

```bash
go run .
```

Output:

```
server started in port :8989
```

### Start server with custom port

```bash
go run . 2525
```

Output:

```
server started in port :2525
```

### Wrong usage

```bash
go run . 2525 localhost
```

Output:

```
usage: go run . port
```

---

## Client Connection (using nc)

Connect using NetCat:

```bash
nc localhost 8989
```

Client receives:

```
Welcome to TCP-Chat!
<linux logo>
[ENTER YOUR NAME]:
```

After entering a valid name, the client receives:

* Chat history
* Prompt line
* Future messages from other clients

---

## Concurrency Model

The project uses Go concurrency primitives:

### Goroutines

* One goroutine per client connection
* One goroutine for ChatManager

### Channels

```
clientChannel      → client join/leave events
messageChannel     → chat messages
validNameChannel   → username validation result
```

This avoids race conditions and centralizes state changes inside ChatManager.

---

## Message Flow

```
Client → HandleConn → messageChannel → ChatManager → broadcast → Clients
```

Steps:

1. Client sends message
2. HandleConn validates and formats it
3. Message sent to messageChannel
4. ChatManager receives it
5. Message added to history
6. Broadcast to all other clients

---

## Client Rules

* Name must not be empty
* Must be ASCII printable
* Must be unique
* Length limited
* Duplicate names are rejected

---

## Chat History

The server stores chat history in memory:

```
chatHistory string
```

When a new client joins:

* Full history is sent before prompt appears

---

## Error Handling

Handled cases:

* Client disconnect (EOF)
* Empty messages ignored
* Invalid ASCII rejected
* Duplicate usernames rejected
* Empty name rejected
* Maximum connections enforcement (10 clients)
---

## Allowed Packages Used

* net
* fmt
* os
* bufio
* io
* strings
* time

(All compliant with project constraints)

---

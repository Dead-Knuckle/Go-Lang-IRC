package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents a connected client
type Client struct {
	conn      net.Conn
	nickname  string
	lastHeart time.Time
}

type Message struct {
	Msg      string `json:"msg"`
	Username string `json:"username"`
}

// ChatRoom represents the chat room
type ChatRoom struct {
	clients map[net.Conn]*Client
}

func main() {
	// Start the server
	chatRoom := &ChatRoom{
		clients: make(map[net.Conn]*Client),
	}
	listener, err := net.Listen("tcp", ":6667")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on port 6667.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleClient(conn, chatRoom)
	}
}

func handleClient(conn net.Conn, chatRoom *ChatRoom) {
	defer conn.Close()

	client := &Client{
		conn:      conn,
		lastHeart: time.Now(),
	}

	// Prompt the client to choose a nickname
	sendMessage(conn, "SERVER", "Welcome to the IRC chat. Please enter a nickname:")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		nickname := scanner.Text()

		if strings.TrimSpace(nickname) == "" {
			continue
		}

		// Check if the chosen nickname is already taken
		taken := false
		for _, c := range chatRoom.clients {
			if c.nickname == nickname {
				taken = true
				break
			}
		}

		if taken {
			sendMessage(conn, "SERVER", "Sorry, that nickname is already taken. Please choose another one:")
		} else {
			client.nickname = nickname
			break
		}
	}

	if client.nickname == "" {
		fmt.Println("Error: Client disconnected without choosing a nickname.")
		return
	}

	// Add client to chat room
	chatRoom.clients[conn] = client
	fmt.Println("Client", client.nickname, "Joined!")

	// Welcome message
	sendMessage(conn, "SERVER", "Welcome, "+client.nickname+"! You are now connected to the chat.")
	sendMessageToAll(chatRoom, client.nickname, " has joined the chat.")

	go manageHeartbeat(conn, chatRoom)

	for scanner.Scan() {
		message := scanner.Text()

		if strings.TrimSpace(message) == "" {
			continue
		}

		if message == "HEARTBEAT" {
			fmt.Println("MSG:HEARTBEAT")
			client.lastHeart = time.Now() // Update the lastHeart time on receiving a PONG
			continue
		}

		// Parse message
		parts := strings.SplitN(strings.TrimSpace(message), " ", 3)
		command := parts[0]

		switch command {
		case "NICK":
			if len(parts) != 2 {
				sendMessage(conn, "SERVER", "Usage: NICK <user_name>")
				continue
			}

			nameToBe := parts[1]

			for _, c := range chatRoom.clients {
				if c.nickname == nameToBe {
					sendMessage(client.conn, "SERVER", "That name is already taken.")
					continue
				}
			}

			client.nickname = nameToBe

			sendMessage(conn, "SERVER", "Nickname set to "+nameToBe)
		case "JOIN":
			sendMessageToAll(chatRoom, client.nickname, " has joined the chat.")
		case "QUIT":
			delete(chatRoom.clients, conn)
			sendMessageToAll(chatRoom, client.nickname, " has left the chat.")
			return
		case "MSG":
			if len(parts) < 3 {
				sendMessage(conn, "SERVER", "Usage: MSG <nickname> <message>")
				continue
			}

			recipientNick := parts[1]
			message := parts[2]

			sendPrivateMessage(chatRoom, client.conn, client.nickname, recipientNick, message)
		case "LST":
			var nicknames []string
			for _, c := range chatRoom.clients {
				nicknames = append(nicknames, c.nickname)
			}
			sendMessage(conn, "SERVER", strings.Join(nicknames, ", "))
		case "HEARTBEAT":
			continue
		case "HELP":
			sendMessage(conn, "SERVER", "'LST' - LIST IRC NAMES : 'MSG <username>' PRIVATELY MSG ANOTHER USER : 'QUIT' SAFEFULY LEAVE CHATROOM 'PING' RETURNS THE MILISECOND PING")
		case "PING":
			startTime := time.Now()
			sendMessage(conn, "SERVER", "PONG")
			endTime := time.Now()
			pingTime := endTime.Sub(startTime).Milliseconds()
			sendMessage(conn, "SERVER", fmt.Sprintf("Ping: %d ms", pingTime))
		default:
			// Broadcast the message to the chat room
			sendMessageToGroup(chatRoom, client.conn, client.nickname, message)
		}

	}
}

func manageHeartbeat(conn net.Conn, chatRoom *ChatRoom) {
	client := chatRoom.clients[conn]
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println(client.lastHeart)

		if time.Since(client.lastHeart) > 15*time.Second {
			fmt.Println("Heartbeat not received, closing connection for", client.nickname)
			sendMessageToAll(chatRoom, client.nickname, " has been disconnected.")
			delete(chatRoom.clients, conn)
			conn.Close()
			return
		}

		fmt.Println("Sending heart beat to ", conn.RemoteAddr().String())
		sendMessage(conn, "SERVER", "HEARTBEAT")
	}
}

func sendMessage(conn net.Conn, username string, message string) {
	send := Message{
		Msg:      message,
		Username: username,
	}

	bytes, _ := json.Marshal(send)
	fmt.Println(string(bytes))
	fmt.Fprintf(conn, "%s\n", string(bytes))
}

func sendMessageToAll(chatRoom *ChatRoom, username string, message string) {
	for _, client := range chatRoom.clients {
		sendMessage(client.conn, username, message)
	}
}

func sendMessageToGroup(chatRoom *ChatRoom, sendClient net.Conn, username string, message string) {
	for _, client := range chatRoom.clients {

		if client.conn != sendClient {
			sendMessage(client.conn, username, message)
		}
	}
}

func sendPrivateMessage(chatRoom *ChatRoom, conn net.Conn, senderNick, recipientNick, message string) {
	for _, client := range chatRoom.clients {
		if client.nickname == recipientNick {
			sendMessage(client.conn, senderNick+"(private) ", message)
			return
		}
	}
	// If recipient not found
	sendMessage(conn, "SERVER", "User "+recipientNick+" not found")
}

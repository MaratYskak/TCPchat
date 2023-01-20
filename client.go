package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func RegisterNewUser(conn net.Conn, server Server, deadConn chan net.Conn) {
	conn.Write(greetingMessage())
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	buf := make([]byte, 1024)
	nbyte, err := conn.Read(buf)

	previousMessages := loadPreviousMessages()
	conn.Write(previousMessages)

	if err != nil {
		deadConn <- conn
	} else {

		username := make([]byte, nbyte-1)
		copy(username, buf[:nbyte])

		server.Clients[conn] = Client{
			id:   1,
			name: string(username),
		}
		for conn := range server.Clients {
			text := fmt.Sprintf("\n%s has joined our chat...\n", username)
			conn.Write([]byte(text))
			time := getCurrentTime()
			text = fmt.Sprintf("[%s][%s]:", time, server.Clients[conn].name)
			conn.Write([]byte(text))
		}

	}

}

func AcceptMessages(conn net.Conn, messages chan Message, server Server, deadConn chan net.Conn) {

	client := server.Clients[conn]
	if client.id == 0 {
		RegisterNewUser(conn, server, deadConn)
	}

	buf := make([]byte, 1024)
	for {
		nbyte, err := conn.Read(buf)
		if err != nil {
			deadConn <- conn
			fmt.Println("connection terminated")

			break
		} else {
			message := make([]byte, nbyte)
			copy(message, buf[:nbyte])

			clientName := server.Clients[conn].name

			readyMessage := Message{
				Text:       string(message),
				SenderName: clientName,
				Sender:     conn,
			}
			messages <- readyMessage
		}
	}
}

func getCurrentTime() string {
	currentTime := time.Now()
	time := currentTime.Format("2006-01-02 15:04:05")
	return time
}

func greetingMessage() []byte {
	data, err := os.ReadFile("./files/greeting.txt")
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func logConfig() {
	file, err := os.Create("logs.txt")
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(file)
	log.SetFlags(0)
}

func loadPreviousMessages() []byte {
	file, err := os.Open("logs.txt")
	if err != nil {
		fmt.Println(err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	return text
}

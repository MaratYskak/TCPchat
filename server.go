package main

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	name string
	id   int
}

type Server struct {
	Clients     map[net.Conn]Client
	ClientCount int
}

type Message struct {
	Text       string
	SenderName string
	Sender     net.Conn
}

func Serve(server *Server, port string) {

	newConnections := make(chan net.Conn)
	messages := make(chan Message)
	deadConnections := make(chan net.Conn)
	logConfig()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				deadConnections <- conn
			}
			newConnections <- conn

		}
	}()

	for {
		var logtext string
		select {
		case conn := <-newConnections:
			fmt.Println("new connection")
			server.ClientCount++
			if server.ClientCount > 10 {
				fmt.Println("Client max count limit exceed")
			} else {
				go AcceptMessages(conn, messages, *server, deadConnections)
			}
		case message := <-messages:
			for conn := range server.Clients {
				time := getCurrentTime()
				logtext = fmt.Sprintf("[%s][%s]:%s", time, message.SenderName, message.Text)

				if conn == message.Sender {
					text := fmt.Sprintf("[%s][%s]:", time, message.SenderName)
					message.Sender.Write([]byte(text))
				} else {
					if len(message.Text) == 1 && message.Text[0] == 10 {
						continue
					}
					text := fmt.Sprintf("\n[%s][%s]:%s", time, message.SenderName, message.Text)
					conn.Write([]byte(text))

					text = fmt.Sprintf("[%s][%s]:", time, server.Clients[conn].name)
					conn.Write([]byte(text))
				}

			}
			log.Print(logtext)
		case deadConnection := <-deadConnections:
			clientName := server.Clients[deadConnection].name
			delete(server.Clients, deadConnection)
			for conn, client := range server.Clients {
				server.ClientCount--

				text := fmt.Sprintf("\n%s has left our chat...\n", clientName)
				conn.Write([]byte(text))

				time := getCurrentTime()
				text = fmt.Sprintf("[%s][%s]:", time, client.name)
				conn.Write([]byte(text))

			}
		}
	}
}

package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	args := os.Args[1:]
	port := ":8989"
	if len(args) == 1 {
		port = ":" + args[0]
	} else if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	fmt.Println("Server is listening on port", port, "using tcp")

	server := new(Server)
	server.Clients = make(map[net.Conn]Client)

	Serve(server, port)
}

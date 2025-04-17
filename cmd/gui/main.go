package main

import (
	"net"
)

func main() {
	clientConn, serverConn := net.Pipe()
	go server(serverConn)
	client(clientConn)
}

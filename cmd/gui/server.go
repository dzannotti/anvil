package main

import (
	"fmt"
	"io"
	"net"
)

func server(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Server: Connection established")

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Server: Client disconnected")
			} else {
				fmt.Printf("Server: Read error: %v", err)
			}
			return
		}

		message := string(buffer[:n])
		fmt.Printf("Server: Received message: %s", message)

		response := fmt.Sprintf("Server received: %s", message)
		if _, err := conn.Write([]byte(response)); err != nil {
			fmt.Printf("Server: Write error: %v", err)
			return
		}
	}
}

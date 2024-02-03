package server

import (
	"fmt"
	"net"
	"strings"
)

func Server() (string, error) {
	// listen to port 8080
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return "", err
	}

	// after the connection close the port
	defer listener.Close()

	// accept the connection when there is a connection
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection", err)
		return "", err
	}
	// close the connection when it is ended
	defer conn.Close()

	// open a new buffer that receives bytes
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data", err)
		return "", err
	}

	// Trim null bytes from the buffer
	trimmedBuffer := strings.TrimRight(string(buffer[:bytesRead]), "\x00")

	fmt.Println("", trimmedBuffer)
	return trimmedBuffer, nil
}
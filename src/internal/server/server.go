package server

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn){
	fmt.Println("handle connection")
}
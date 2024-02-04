package query

import (
	"fmt"
	"nestdb/internal/config"
	"nestdb/internal/server"
)

//Input: ptr of inputbuffer
//output:
func Read_input(input_Buffer *config.InputBuffer){
	/*
    reader := bufio.NewReader(os.Stdin)
    line , err := reader.ReadString('\n')


    if err != nil{
        fmt.Println("Error Reading Input")
    }
    bytesRead := len(line) -1
    */

    line, err := server.Server()
    if err != nil {
        fmt.Printf("server error")
    }
    bytesRead := len(line) 
    // Fix the input_length, buffer_length, buffer
    input_Buffer.Input_length = uint32(bytesRead)
    input_Buffer.Buffer_length = uint32(bytesRead)
    
    // The new pointer is representing the type of byte, make the buffer with bytesRead
    input_Buffer.Buffer = new([]byte)
    *input_Buffer.Buffer = make([]byte, bytesRead)
    
    // Copy the modified byte slice into the buffer
    copy(*input_Buffer.Buffer, []byte(line)[:bytesRead])
}
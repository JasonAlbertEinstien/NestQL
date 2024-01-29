package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

func sizeOfAttribute(Struct interface{} , Attribute interface{})(uint32){
    return uint32(unsafe.Sizeof(Attribute))
}

//just for current testing purpose
var ID_SIZE uint32 = sizeOfAttribute(Row{} , Row{}.id) 
var USERNAME_SIZE uint32 = sizeOfAttribute(Row{} , Row{}.username)
var EMAIL_SIZE uint32 = sizeOfAttribute(Row{} , Row{}.email)
var ID_OFFSET uint32 = 0
var USERNAME_OFFSET = ID_OFFSET + ID_SIZE
var EMAIL_OFFSET = USERNAME_OFFSET + USERNAME_SIZE
var ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE

//config seting
const (
    PAGE_SIZE uint32 = 4096
    TABLE_MAX_PAGES uint32 = 100
) 

var ROWS_PER_PAGE uint32 = PAGE_SIZE/ROW_SIZE
var TABLE_MAX_ROWS uint32 = ROWS_PER_PAGE*TABLE_MAX_PAGES

type Page [4096]byte

type Table struct{
    num_rows uint32
    pages   []* Page
}

//this part will create a new table and return a pointer point to table
func new_table()(*Table){
    table := &Table{
        pages: make([] *Page , TABLE_MAX_PAGES),
    }

    //set every place of the table allocate an address to avoid nil pointer problem
    for i := range table.pages {
		page := &Page{}
		table.pages[i] = page
	}

    return table
}

//This is a metacommand restul which is int 
type MetaCommandResult int
//The constant is either 0 or 1
const (
    META_COMMAND_SUCCESS MetaCommandResult= 0
    META_COMMAND_UNRECOGNIZE_COMMAND MetaCommandResult= 1
)
//PrepareResult is an integer
type PrepareResult int
//Either 0 or 1
const (
    PREPARE_SUCCESS PrepareResult= 0
    PREPARE_UNRECOGNIZED_STATEMENT PrepareResult= 1
    PREPARE_SYNTAX_ERROR PrepareResult =2 
)

type StatementType int

type Statement struct{
    _type StatementType
    row_to_insert Row 
}

const (
    STATEMENT_INSERT StatementType= 0
    STATEMENT_SELECT StatementType= 1
)

type Row struct{
   id uint32
   username [COLUMN_USERNMAE_SIZE]byte
   email [COLUMN_EMAIL_SIZE]byte
}

//serializeRow and derserializeRow with binary encoding
//the return byte is a binary code now
func serializeRow(source* Row)([]byte , error){
    buf := new(bytes.Buffer)
    err := binary.Write(buf , binary.LittleEndian , source.id)
    if err != nil {
        return nil , err
    }
    err = binary.Write(buf , binary.LittleEndian , source.username)
    if err != nil {
        return nil , err
    }
    err = binary.Write(buf , binary.LittleEndian , source.email)
    if err != nil{
        return nil , err
    }
    return buf.Bytes() , nil
}

func deserializeRow(data []byte)(*Row, error){
    buf := bytes.NewReader(data)
    destination := &Row{}
    err := binary.Read(buf , binary.LittleEndian , &destination.id)
    if err != nil {
        return nil , err
    }

    err = binary.Read(buf , binary.LittleEndian , &destination.username)
    if err != nil {
        return nil , err
    }

    err = binary.Read(buf , binary.LittleEndian , &destination.email)
    if err != nil{
        return nil , err
    }
    return destination, err
}

const (
    COLUMN_USERNMAE_SIZE uint8 = 32
    COLUMN_EMAIL_SIZE uint8 = 255
)
//so this function will be return 0 if the command is .exit  else return 1
func do_meta_command(input_buffer *InputBuffer)(MetaCommandResult){
    if(string(*input_buffer.buffer) == ".exit"){
        println("Thank you for your contribution to nestspace")
        os.Exit(0)
        return META_COMMAND_SUCCESS
    }else{
        return META_COMMAND_UNRECOGNIZE_COMMAND
    }
}
//this is for prepare statment
//after running this function your statement will be a number representing which kidn of statment
//return a value if it's already prepared
func prepare_statement(input_buffer *InputBuffer , statement *Statement)(PrepareResult){
    //compare two string even if the first one is nil it will return false instead of fatal erro , also allow later to read other args
    if strings.HasPrefix(string(*input_buffer.buffer), "insert") { 
        statement._type = STATEMENT_INSERT
        args_assigned , _ := fmt.Sscanf(string(*input_buffer.buffer) , "insert %d %s %s" , &(statement.row_to_insert.id), &(statement.row_to_insert.username) , &(statement.row_to_insert.email))
        if(args_assigned <3){
            return PREPARE_SYNTAX_ERROR
        }
        return PREPARE_SUCCESS
    }
    if string(*input_buffer.buffer) == "select"{
        statement._type = STATEMENT_SELECT
        return PREPARE_SUCCESS
    }
    return PREPARE_UNRECOGNIZED_STATEMENT
}

//exec statement function
func execute_statement(statement *Statement){
    switch(statement._type){
        case (STATEMENT_INSERT):{
            fmt.Println("this is the place for insert statment")
            break
        }
        case (STATEMENT_SELECT):{
            fmt.Println("this is a select statment")
            break
        }
    }
}

type ExecuteResult int

const (
    EXECUTE_SUCCESS ExecuteResult = 0
    EXECUTE_FAIL ExecuteResult = 1
)


// if the table current num row is alredy greater then return an error
//logic: figure out where to save (which page to save the data) 
//in the page (find the current page number)
// then serialize the data based on the statement
// and then save the data 
func execute_insert(statement *Statement , table *Table)(ExecuteResult){
    if(table.num_rows >= TABLE_MAX_ROWS){
        return EXECUTE_FAIL
    }

    row_to_insert := &(statement.row_to_insert)
    data , _ := serializeRow(row_to_insert)

    page_number := table.num_rows/ROWS_PER_PAGE
    page := table.pages[page_number]

    row_offset := table.num_rows %ROWS_PER_PAGE
    row_address := page[row_offset*ROW_SIZE:]

    copy(row_address , data)

    table.num_rows += 1

    return EXECUTE_SUCCESS  
}





//this is the structure of the input buffer
//It contain 
//buffer: saving the string
//buffer length : how long the length of the buffer
//input_length is for saving the user inptut 
type InputBuffer struct{
    buffer *[]byte
    buffer_length uint32
    input_length uint32
}

//create a new inputbuffer
// input : nil 
// return: pointer of an inputbuffer
func new_input_buffer()(*InputBuffer){
    return &InputBuffer{
        buffer: nil,
        buffer_length: 0,
        input_length: 0,
    }
}

func print_pomt(){
    fmt.Printf("nestspaceDB >")
}

//Input: ptr of inputbuffer
//output:
func read_input(input_Buffer *InputBuffer){
    reader := bufio.NewReader(os.Stdin)
    line , err := reader.ReadString('\n')

    if err != nil{
        fmt.Println("Error Reading Input")
    }
    bytesRead := len(line) -1

    //fix the input_length , buffer_length , buffer 
    input_Buffer.input_length = uint32(bytesRead)
	input_Buffer.buffer_length = uint32(len(line))

    //The new pointer is representing the type of byte and then make the buffer with byteread
    input_Buffer.buffer = new([]byte)
    *input_Buffer.buffer = make([]byte, bytesRead)
	copy(*input_Buffer.buffer, []byte(line)[:bytesRead])
}

/*
    variables: inputbuffer (ptr) 
    logic:
        print prompt 
        read input
        
        for everyloop you check 1. if the
    
*/
func main(){
    input_buffer := new_input_buffer()


    for{
            print_pomt()
            read_input(input_buffer)
            
            if(string(string(*input_buffer.buffer)[0]) == "."){
                switch(do_meta_command(input_buffer)){
                    case(META_COMMAND_SUCCESS):{
                        continue
                    }
                    case(META_COMMAND_UNRECOGNIZE_COMMAND):{
                        println("FUCK OFF WRONG COMMAND")
                        continue
                    }
                }
            }
            var statement Statement
            switch(prepare_statement(input_buffer , &statement)){
                case(PREPARE_SUCCESS):{
                    break
                }
                case(PREPARE_UNRECOGNIZED_STATEMENT):{
                    println("Bruh your sql is an unrecognized COMMAND")
                    continue
                }
            case (PREPARE_SYNTAX_ERROR):{
                    println("UMM can you please input correct type and amount of arguments. NESTSPACE IS CRYING :(")
                    continue
                }
            }
            execute_statement(&statement)
            fmt.Printf("Executed \n")
        }    
}



package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	config "nestdb/internal/config"
	server "nestdb/internal/server"
	"nestdb/internal/table"
	"nestdb/internal/utils"
)

//Implementing B-Tree Structure
//Node Header Format
//format of a leaf node should be
/*
   node_type 1 byte
   is roote 1 byte
   parent_pointer 2 - 5 byte
   byte 6 - 9 num_cells
   byte 14- 306 value 0
   vyte 307 -310 key1
   byte page xszie : value 1 ...
   each leaf node will have 12 values and key pary
   there are some byte is wasted

*/
type NodeType int 

const(
    NODE_INTERNAL NodeType = 0
    NODE_LEAF NodeType = 1
)

//NODE HEADER LAYOUT
const(
    NODE_TYPE_SIZE = 1
    NODE_TYPE_OFFSET = 0
    IS_ROOT_SIZE = 1
    IS_ROOT_OFFSET = NODE_TYPE_SIZE
    PARENT_POINTER_SIZE = 4
    PARENT_POINTER_OFFSET = IS_ROOT_OFFSET +IS_ROOT_SIZE
    COMMONE_NODE_HEADER_SIZE = NODE_TYPE_SIZE + IS_ROOT_SIZE + PARENT_POINTER_SIZE
)

//LEAF NODE FORMAT
const(
    LEAF_NODE_KEY_SIZE = 4
    LEAF_NODE_KEY_OFFSET = 0
    LEAF_NODE_HEADER_SIZE = LEAF_NODE_KEY_SIZE + COMMONE_NODE_HEADER_SIZE
)

var LEAF_NODE_VALUE_SZIZE uint32 = config.ROW_SIZE
var LEAF_NODE_VALUE_OFFSET uint32 = LEAF_NODE_KEY_OFFSET +LEAF_NODE_KEY_SIZE
var LEAF_NODE_CELL_SIZE = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SZIZE
var LEAF_NODE_SPACE_FOR_CELLS uint32 = config.PAGE_SIZE - LEAF_NODE_HEADER_SIZE
var LEAF_NODE_MAX_CELL uint32 = LEAF_NODE_SPACE_FOR_CELLS / LEAF_NODE_CELL_SIZE

//how to access the key with node 
//how insert is done ?
//How to tree is form ? 
type NodeTree struct{
    node *Node
}

type Node struct{
    is_leaf bool
}
//insertion of b tree
func insertNode(){}

func leftRotation(){}

func rightRotation(){}
//delete of b tree
func deleteNode(){}
//search in b tree 
// the algorithm that is used is simiplar to binary tree
func searchNode(){}




//so this function will be return 0 if the command is .exit  else return 1
func do_meta_command(input_buffer *config.InputBuffer , table *config.Table)(config.MetaCommandResult){
    if(string(*input_buffer.Buffer) == ".exit"){
        db_close(table)
        println("Thank you for your contribution")
        os.Exit(0)
        return config.META_COMMAND_SUCCESS
    }else{
        return config.META_COMMAND_UNRECOGNIZE_COMMAND
    }
}
//this is for prepare statment
//after running this function your statement will be a number representing which kidn of statment
//return a value if it's already prepared
func prepare_statement(input_buffer *config.InputBuffer , statement *config.Statement)(config.PrepareResult){
    //compare two string even if the first one is nil it will return false instead of fatal erro , also allow later to read other args
    if strings.HasPrefix(string(*input_buffer.Buffer), "insert") { 
        statement.Type = config.STATEMENT_INSERT

        //I LOVE CHATGPT I HATE TO WRITE THIS
        var usernameBytes []byte
        var emailBytes []byte
        args_assigned, _ := fmt.Sscanf(string(*input_buffer.Buffer), "insert %d %s %s", &(statement.Row_to_insert.Id), &usernameBytes, &emailBytes)
        if(args_assigned <3){
            return config.PREPARE_SYNTAX_ERROR
        }
        copy(statement.Row_to_insert.Username[:], usernameBytes)
        copy(statement.Row_to_insert.Email[:], emailBytes)
        return config.PREPARE_SUCCESS
    }
    if string(*input_buffer.Buffer) == "select"{
        statement.Type = config.STATEMENT_SELECT
        return config.PREPARE_SUCCESS
    }
    return config.PREPARE_UNRECOGNIZED_STATEMENT
}

//exec statement function
func execute_statement(statement *config.Statement , table *config.Table)(config.ExecuteResult){
    switch(statement.Type){
        case (config.STATEMENT_INSERT):{
            execute_insert(statement , table)
            break
        }
        case (config.STATEMENT_SELECT):{
            execute_select(statement , table)
            break
        }
    }
    return config.EXECUTE_FAIL
}
// if the table current num row is alredy greater then return an error
//logic: figure out where to save (which page to save the data) 
//in the page (find the current page number)
// then serialize the data based on the statement
// and then save the data 
func execute_insert(statement *config.Statement , table *(config.Table))(config.ExecuteResult){
	if table.Num_rows >= config.TABLE_MAX_ROWS {
		return config.EXECUTE_FAIL
	}

    row_to_insert := &(statement.Row_to_insert)
	data, _ := utils.Serialize(row_to_insert)

	// Get the address of the row
    /*
	// Find the page number and the specific position of the row

    */
    //page_number := table.Num_rows / ROWS_PER_PAGE
    //page := get_page(table.Pager , page_number)
    

	// Find the remainder by calculating the row size and get the address of that row
	//row_offset := table.Num_rows % ROWS_PER_PAGE
    cursor := table_end(table)
	row_address := cursor_value(cursor)
	copy(row_address, data)
	// Increment the number of rows
	table.Num_rows += 1
	// Flush the changes to disk
	//pager_flush(*table.Pager, page_number, PAGE_SIZE)
	return config.EXECUTE_SUCCESS
}

//This function is used for select operation (currently you can only select all from the database)
//the logic is: 
/*
    //the cursor will first point to the value and ten deserialize the row
    parameter: statement , table
    return EXECUTE RESULT
    what i want is to get the address of place i want to read (in this case every page)
    currently the table.page is like this
    table:  --> page1 --> [4096]byte
            --> page2 
            --> page3
    so inside the for loop i have to get the position i need for every row 
    for example table[page1][rowoffset*row_size: (rowoffset*row_size)+row_size]
    
*/
func execute_select(statement *config.Statement , table *config.Table)(config.ExecuteResult){
    cursor := table_stat(table)
    for cursor.end_of_table != true {
        data, _ := utils.Deserialize(cursor_value(cursor))
        cursor_advance(cursor)
        print_row(data)
    }
    return config.EXECUTE_SUCCESS
}
//function print table
func print_row(row *config.Row) {
	fmt.Printf("ID: %d, Username: %s, Email: %s\n", row.Id, row.Username, row.Email)
} 

//create a new inputbuffer
// input : nil 
// return: pointer of an inputbuffer
func new_input_buffer()(*config.InputBuffer){
    return &config.InputBuffer{
        Buffer: nil,
        Buffer_length: 0,
        Input_length: 0,
    }
}

func print_pomt(){
    fmt.Printf("nestspaceDB >")
}

//Input: ptr of inputbuffer
//output:
func read_input(input_Buffer *config.InputBuffer){
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

/*
    open database
        
    close database
 
    create a new database 

    get page function --> set the pages pointer to 
*/

//get that specific postion from the page
/*
    This section have bug ! Please be aware 
*/
func get_page(pager *config.Pager ,page_num uint32)*config.Page{
    //if page number > the size --> exit the os.exit()
    if(page_num > config.TABLE_MAX_PAGES){
        fmt.Println("TRY TO ACCESS MORE THAN TLBLE MAXPAGE")
        os.Exit(1)
    }   
//    fmt.Printf("INLOOP")
    //if that page is not being fetched before 
    if(pager.Pages[page_num] == nil){
        //create a new pointer for one page 
        page := new(config.Page)
        num_pages := pager.File_length/config.PAGE_SIZE

        if (pager.File_length % config.PAGE_SIZE != 0){
            num_pages += 1
        }
        if(page_num <= num_pages){

            //bug exit here : 
            _, err := syscall.Seek(pager.File_descriptor , int64(page_num*config.PAGE_SIZE) , 0)
            if err != nil{
                fmt.Printf("ERROR SEEKING FILE \n")
                os.Exit(1)
            }
            //the read and return to the byte
            bytesRead , err := syscall.Read(pager.File_descriptor , page[:])
            if err != nil {
                fmt.Printf("ERROR READING FILE \n")
                os.Exit(1)
            }
            if bytesRead == -1{
                fmt.Printf("ERROR READING FILE \n")
                os.Exit(1)
            }
        }
        pager.Pages[page_num] = page
    }
    return pager.Pages[page_num]
}

//set the file closed and save it   
//only save the datat when the dataase is closed 
//save the data to the database (from in memory to disk and then close the database)
func db_close(table *config.Table){
    pager := table.Pager
    num_full_pages := table.Num_rows/config.ROWS_PER_PAGE

    for i := 0 ; i < int(num_full_pages) ; i++ {
        if pager.Pages[i] == nil{
     
            continue
        }
        pager_flush(*pager , uint32(i), config.PAGE_SIZE)
        pager.Pages[i] = nil
    }
    num_additional_rows := table.Num_rows % config.ROWS_PER_PAGE
    if(num_additional_rows >0){
        page_num := num_full_pages
        if(pager.Pages[page_num] != nil){
            pager_flush(*pager , page_num , num_additional_rows*config.ROW_SIZE)
            pager.Pages[page_num] = nil
        }
    }
    err := syscall.Close(pager.File_descriptor)
    if err != nil {
        os.Exit(1)
    }
    for i:=0 ; i <int(config.TABLE_MAX_PAGES); i++{
        page := pager.Pages[i]
        if(page != nil){
            pager.Pages[i] = nil
        }
    }
}

//just a simple that accept a new table
//logic : if that page is nil --> that means the page does not exit so it will be exit 1 
//if not that try to move the pointer to that page and then 
func pager_flush(pager config.Pager , page_num uint32 ,size uint32){
    if(pager.Pages[page_num] == nil){
        fmt.Printf("Tried to flush null page \n")
        os.Exit(1)
    }
    offset , _:= syscall.Seek(pager.File_descriptor , int64(page_num)*int64(config.PAGE_SIZE) , 0)
    if (offset == -1){
        fmt.Printf("ERROR SEEKING \n")
        os.Exit(1)
    }
    bytes_written , err := syscall.Write(pager.File_descriptor , pager.Pages[page_num][:size])
    if err != nil {
        fmt.Printf("Error writing: %v\n", err)
		os.Exit(1)
    }
    if bytes_written != int(size) {
		fmt.Println("Incomplete write")
	}
}

/*
    cursor function : represent the location in the table
    cursor at the beginning of the table
    cursor at the end of the table 
    access to the place where your cursor is pointing to 
    access the cursor to the next row (dynamic pointer)

    purpose: delete the row pointed to by a cursor 
    Modify a row pointed by a cursor
    Search a table by a given id    
*/
type Cursor struct{
    table *config.Table
    row_num uint32
    end_of_table bool
}

//return a cursor that the where 
func table_stat(table *config.Table)(*Cursor){
    cursor := &Cursor{
        table: table,
        row_num: 0,
        end_of_table: (table.Num_rows == 0),
    }
    return cursor
}

//return a cursor to the table end
func table_end(table *config.Table)(*Cursor){
    cursor := &Cursor{
        table: table,
        row_num: table.Num_rows,   
        end_of_table: true,
    }
    return cursor
}

//calculate the place for the curosr
func cursor_value(cursor *Cursor)([]byte){
    row_num := cursor.row_num
    page_num := row_num/config.ROWS_PER_PAGE

    page := get_page(cursor.table.Pager , page_num)

    row_offset := row_num%config.ROWS_PER_PAGE
    bytes_offset := row_offset*config.ROW_SIZE

    return page[bytes_offset:row_offset*config.ROW_SIZE+config.ROW_SIZE-1]
}

func cursor_advance(cursor *Cursor){
    cursor.row_num += 1 
    if(cursor.row_num >= cursor.table.Num_rows){
        cursor.end_of_table = true
    }
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
    if(len(os.Args)) <2 {
        fmt.Println("Please supply a database filename uwu")
        os.Exit(1)
    }
    filename := os.Args[1]
    //a bug in here 
    
    //run the tcp server
    time.Sleep(time.Second)
    table , _ := table.Db_open(filename)
    for{
            print_pomt()
            //replace this read input function with reading the input from server
            read_input(input_buffer)
            //*input_buffer.buffer = <- input_server
            if len(*input_buffer.Buffer) > 0 && (string(string(*input_buffer.Buffer)[0]) == "."){
                switch(do_meta_command(input_buffer , table)){
                    case(config.META_COMMAND_SUCCESS):{
                        continue
                    }
                    case(config.META_COMMAND_UNRECOGNIZE_COMMAND):{
                        println("WRONG META COMMAND")
                        continue
                    }
                }
            }
            var statement config.Statement
            switch(prepare_statement(input_buffer , &statement)){
                case(config.PREPARE_SUCCESS):{
                    break
                }
                case(config.PREPARE_UNRECOGNIZED_STATEMENT):{
                    println("Unrecognized COMMAND")
                    break
                }
                case (config.PREPARE_SYNTAX_ERROR):{
                    println("Wrong Args or TYPES")
                    break
                }
            }
            execute_statement(&statement , table)
            fmt.Printf("Executed \n")
        }    
}
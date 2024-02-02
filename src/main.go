package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"syscall"

	tableds "nestdb/internal/table"
)

//open the file (if not then create one) --> return a file pointer
// logic:
/*
   input the file name
   using openfile(filename , os.xxx , ox . xxx , if it's allowing user to delete than it should always be
   0777 if not it's better to be 0666)
   save the file pointer to at the off set end of the file
   the the pager should have 1. a fd , 2. file length (if new then 0 ) 3. a pages which is 4096 byte set as the max page size
*/
//
func pager_open(Filename string)(*tableds.Pager){
    f, err := os.OpenFile(Filename , os.O_RDWR|os.O_CREATE, syscall.S_IWUSR | syscall.S_IRUSR)   
    if err != nil{
        fmt.Println("ERROR OF OPEN THE DB FILE OR CREATING THE FILE")
    }
    fileInfo , err := f.Stat()
    if err != nil {
        fmt.Println("ERROR CANNOT READ FILE DATA")
    }
    fileLength := fileInfo.Size()
    pager := &tableds.Pager{
        File_descriptor: int(f.Fd()),
        File_length: uint32(fileLength),
        //remember to make enough space for it if not it will have error
        Pages: make([]*tableds.Page, TABLE_MAX_PAGES),
    }
    for i := 0 ; i < int(TABLE_MAX_PAGES) ; i++{
        pager.Pages[i] = nil
    }
    return pager
}


//this part will create a new table and return a pointer point to table

func db_open(filename string) *tableds.Table {
    pager := pager_open(filename)
    //this pager should be nil
    var num_rows uint32 = pager.File_length/ROW_SIZE
	table := &tableds.Table{
        Pager: pager,
        Num_rows: num_rows,
	}
	return table
}

//replace this with db_opne

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

//deserializeRow with the given data --> 
// it will retrive it to  
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


//so this function will be return 0 if the command is .exit  else return 1
func do_meta_command(input_buffer *InputBuffer , table *tableds.Table)(MetaCommandResult){
    if(string(*input_buffer.buffer) == ".exit"){
        db_close(table)
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

        //I LOVE CHATGPT I HATE TO WRITE THIS
        var usernameBytes []byte
        var emailBytes []byte
        args_assigned, _ := fmt.Sscanf(string(*input_buffer.buffer), "insert %d %s %s", &(statement.row_to_insert.id), &usernameBytes, &emailBytes)
        if(args_assigned <3){
            return PREPARE_SYNTAX_ERROR
        }
        copy(statement.row_to_insert.username[:], usernameBytes)
        copy(statement.row_to_insert.email[:], emailBytes)
        return PREPARE_SUCCESS
    }
    if string(*input_buffer.buffer) == "select"{
        statement._type = STATEMENT_SELECT
        return PREPARE_SUCCESS
    }
    return PREPARE_UNRECOGNIZED_STATEMENT
}

//exec statement function
func execute_statement(statement *Statement , table *tableds.Table)(ExecuteResult){
    switch(statement._type){
        case (STATEMENT_INSERT):{
            execute_insert(statement , table)
            break
        }
        case (STATEMENT_SELECT):{
            execute_select(statement , table)
            break
        }
    }
    return EXECUTE_FAIL
}

// if the table current num row is alredy greater then return an error
//logic: figure out where to save (which page to save the data) 
//in the page (find the current page number)
// then serialize the data based on the statement
// and then save the data 
func execute_insert(statement *Statement , table *(tableds.Table))(ExecuteResult){
	if table.Num_rows >= TABLE_MAX_ROWS {
		return EXECUTE_FAIL
	}

    row_to_insert := &(statement.row_to_insert)
	data, _ := serializeRow(row_to_insert)

	// Get the address of the row
    /*
	// Find the page number and the specific position of the row
	page_number := table.Num_rows / ROWS_PER_PAGE
	page := table.Pager.Pages[page_number]
    */
    page_number := table.Num_rows / ROWS_PER_PAGE
    page := get_page(table.Pager , page_number)
    

	// Find the remainder by calculating the row size and get the address of that row
	row_offset := table.Num_rows % ROWS_PER_PAGE
	row_address := page[row_offset*ROW_SIZE:]
	copy(row_address, data)
	// Increment the number of rows
	table.Num_rows += 1
	// Flush the changes to disk
	//pager_flush(*table.Pager, page_number, PAGE_SIZE)
	return EXECUTE_SUCCESS
}

//This function is used for select operation (currently you can only select all from the database)
//the logic is: 
/*
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
func execute_select(statement *Statement , table *tableds.Table)(ExecuteResult){
    //fmt.Printf("" , table.Num_rows)
    for i := 0; i < int(table.Num_rows); i++ {
        page_number := table.Num_rows / ROWS_PER_PAGE
        page := get_page(table.Pager , page_number)

        row_offset := uint32(i) % ROWS_PER_PAGE
        row_address := page[row_offset*ROW_SIZE : row_offset*ROW_SIZE+ROW_SIZE-1]

        data, _ := deserializeRow(row_address)
        print_row(data)
    }
    return EXECUTE_SUCCESS
}
//function print table
func print_row(row *Row) {
	fmt.Printf("ID: %d, Username: %s, Email: %s\n", row.id, row.username, row.email)
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
    open database
        
    close database
 
    create a new database 

    get page function --> set the pages pointer to 
*/

//get that specific postion from the page
/*
    This section have bug ! Please be aware 
*/
func get_page(pager *tableds.Pager ,page_num uint32)*tableds.Page{
    //if page number > the size --> exit the os.exit()
    if(page_num > TABLE_MAX_PAGES){
        fmt.Println("EPIC FAIL YOU TRY TO ACCESS MORE THAN TLBLE MAXPAGE")
        os.Exit(1)
    }   
//    fmt.Printf("INLOOP")
    //if that page is not being fetched before 
    if(pager.Pages[page_num] == nil){
        //create a new pointer for one page 
        page := new(tableds.Page)
        num_pages := pager.File_length/PAGE_SIZE

        if (pager.File_length % PAGE_SIZE != 0){
            num_pages += 1
        }
        if(page_num <= num_pages){

            //bug exit here : 
            _, err := syscall.Seek(pager.File_descriptor , int64(page_num*PAGE_SIZE) , 0)
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
func db_close(table *tableds.Table){
    pager := table.Pager
    num_full_pages := table.Num_rows/ROWS_PER_PAGE

    for i := 0 ; i < int(num_full_pages) ; i++ {
        if pager.Pages[i] == nil{
     
            continue
        }
        pager_flush(*pager , uint32(i), PAGE_SIZE)
        pager.Pages[i] = nil
    }
    num_additional_rows := table.Num_rows %ROWS_PER_PAGE
    if(num_additional_rows >0){
        page_num := num_full_pages
        if(pager.Pages[page_num] != nil){
            pager_flush(*pager , page_num , num_additional_rows*ROW_SIZE)
            pager.Pages[page_num] = nil
        }
    }
    err := syscall.Close(pager.File_descriptor)
    if err != nil {
        os.Exit(1)
    }
    for i:=0 ; i <int(TABLE_MAX_PAGES); i++{
        page := pager.Pages[i]
        if(page != nil){
            pager.Pages[i] = nil
        }
    }
}

//just a simple that accept a new table
//logic : if that page is nil --> that means the page does not exit so it will be exit 1 
//if not that try to move the pointer to that page and then 
func pager_flush(pager tableds.Pager , page_num uint32 ,size uint32){
    if(pager.Pages[page_num] == nil){
        fmt.Printf("Tried to flush null page \n")
        os.Exit(1)
    }
    offset , _:= syscall.Seek(pager.File_descriptor , int64(page_num)*int64(PAGE_SIZE) , 0)
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
    cursor function 
    cursor at the beginning of the table
    cursor at the end of the table 
    access to the place where your cursor is pointing to 
    access the cursor to the next row (dynamic pointer)
*/
type Cursor struct{

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
        fmt.Println("Please supply a database filename")
        os.Exit(1)
    }
    filename := os.Args[1]
    //a bug in here 
    table := db_open(filename)
    for{
            print_pomt()
            read_input(input_buffer)
            if len(*input_buffer.buffer) > 0 && (string(string(*input_buffer.buffer)[0]) == "."){
                switch(do_meta_command(input_buffer , table)){
                    case(META_COMMAND_SUCCESS):{
                        continue
                    }
                    case(META_COMMAND_UNRECOGNIZE_COMMAND):{
                        println("OPPS WRONG META COMMAND")
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
                    break
                }
                case (PREPARE_SYNTAX_ERROR):{
                    println("UMM can you please input correct type and amount of arguments. NESTSPACE IS CRYING :(")
                    break
                }
            }
            execute_statement(&statement , table)
            fmt.Printf("Executed \n")
        }    
}


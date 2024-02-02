package main

import "unsafe"

func sizeOfAttribute(Struct interface{} , Attribute interface{})(uint32){
    return uint32(unsafe.Sizeof(Attribute))
}

const (
	PAGE_SIZE uint32 = 4096
    TABLE_MAX_PAGES uint32 = 100
)

const (
    META_COMMAND_SUCCESS MetaCommandResult= 0
    META_COMMAND_UNRECOGNIZE_COMMAND MetaCommandResult= 1
)
const (
    PREPARE_SUCCESS PrepareResult= 0
    PREPARE_UNRECOGNIZED_STATEMENT PrepareResult= 1
    PREPARE_SYNTAX_ERROR PrepareResult =2 
)

const (
    STATEMENT_INSERT StatementType= 0
    STATEMENT_SELECT StatementType= 1
)

const (
    COLUMN_USERNMAE_SIZE uint8 = 32
    COLUMN_EMAIL_SIZE uint8 = 255
)
const (
    EXECUTE_SUCCESS ExecuteResult = 0
    EXECUTE_FAIL ExecuteResult = 1
)

var ID_SIZE uint32 = sizeOfAttribute(Row{} , Row{}.id) 
var USERNAME_SIZE uint32 = 32
var EMAIL_SIZE uint32 = 255
var ID_OFFSET uint32 = 0
var USERNAME_OFFSET = ID_OFFSET + ID_SIZE
var EMAIL_OFFSET = USERNAME_OFFSET + USERNAME_SIZE
var ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
var ROWS_PER_PAGE uint32 = PAGE_SIZE/ROW_SIZE
var TABLE_MAX_ROWS uint32 = ROWS_PER_PAGE*TABLE_MAX_PAGES


//This is a metacommand restul which is int 
type MetaCommandResult int

//PrepareResult is an integer
type PrepareResult int
//Either 0 or 1

type StatementType int

type Statement struct{
    _type StatementType
    row_to_insert Row 
}

type Row struct{
   id uint32
   username [COLUMN_USERNMAE_SIZE]byte
   email [COLUMN_EMAIL_SIZE]byte
}
type ExecuteResult int

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


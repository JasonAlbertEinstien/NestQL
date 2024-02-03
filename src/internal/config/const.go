package config

const (
	PAGE_SIZE uint32 = 4096
    TABLE_MAX_PAGES uint32 = 100
    COLUMN_USERNMAE_SIZE uint8 = 32

    COLUMN_EMAIL_SIZE uint8 = 255

	META_COMMAND_SUCCESS MetaCommandResult= 0
    META_COMMAND_UNRECOGNIZE_COMMAND MetaCommandResult= 1

	PREPARE_SUCCESS PrepareResult= 0
    PREPARE_UNRECOGNIZED_STATEMENT PrepareResult= 1
    PREPARE_SYNTAX_ERROR PrepareResult =2 

	STATEMENT_INSERT StatementType= 0
    STATEMENT_SELECT StatementType= 1

	EXECUTE_SUCCESS ExecuteResult = 0
    EXECUTE_FAIL ExecuteResult = 1
)

var ID_SIZE uint32 = 5
var USERNAME_SIZE uint32 = 32
var EMAIL_SIZE uint32 = 255
var ID_OFFSET uint32 = 0
var USERNAME_OFFSET = ID_OFFSET + ID_SIZE
var EMAIL_OFFSET = USERNAME_OFFSET + USERNAME_SIZE
var ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
var ROWS_PER_PAGE uint32 = PAGE_SIZE/ROW_SIZE
var TABLE_MAX_ROWS uint32 = ROWS_PER_PAGE*TABLE_MAX_PAGES

type MetaCommandResult int
type PrepareResult int
type StatementType int 
type ExecuteResult int

type Statement struct{
	Type StatementType
    Row_to_insert Row 
}

type InputBuffer struct{
	Buffer *[]byte
	Buffer_length uint32
	Input_length uint32
}

type Row struct{
	Id uint32
	Username [COLUMN_USERNMAE_SIZE]byte
	Email [COLUMN_EMAIL_SIZE]byte
 }

 type Pager struct{
	File_descriptor int
	File_length uint32 
	Pages []*Page
}
type Page [4096]byte

type Table struct{
	Num_rows uint32
	Pager *Pager
}
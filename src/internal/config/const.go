package config

const (
	PAGE_SIZE uint32 = 4096
    TABLE_MAX_PAGES uint32 = 100
    COLUMN_USERNMAE_SIZE uint8 = 32

    COLUMN_EMAIL_SIZE uint8 = 255
)

var ID_SIZE uint32 = 4
var USERNAME_SIZE uint32 = 32
var EMAIL_SIZE uint32 = 255
var ID_OFFSET uint32 = 0
var USERNAME_OFFSET = ID_OFFSET + ID_SIZE
var EMAIL_OFFSET = USERNAME_OFFSET + USERNAME_SIZE
var ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
var ROWS_PER_PAGE uint32 = PAGE_SIZE/ROW_SIZE
var TABLE_MAX_ROWS uint32 = ROWS_PER_PAGE*TABLE_MAX_PAGES


type Row struct{
	Id uint32
	Username [COLUMN_USERNMAE_SIZE]byte
	Email [COLUMN_EMAIL_SIZE]byte
 }
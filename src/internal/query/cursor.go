package query

import (
	"nestdb/internal/config"
	"nestdb/internal/table"
)

//return a cursor that the where
func Table_stat(table *config.Table)(*config.Cursor){
    cursor := &config.Cursor{
        Table: table,
        Row_num: 0,
        End_of_table: (table.Num_rows == 0),
    }
    return cursor
}

//return a cursor to the table end
func Table_end(table *config.Table)(*config.Cursor){
    cursor := &config.Cursor{
        Table: table,
        Row_num: table.Num_rows,   
        End_of_table: true,
    }
    return cursor
}

//calculate the place for the curosr
func Cursor_value(cursor *config.Cursor)([]byte){
    row_num := cursor.Row_num
    page_num := row_num/config.ROWS_PER_PAGE

    page := table.Get_page(cursor.Table.Pager , page_num)

    row_offset := row_num%config.ROWS_PER_PAGE
    bytes_offset := row_offset*config.ROW_SIZE

    return page[bytes_offset:row_offset*config.ROW_SIZE+config.ROW_SIZE-1]
}

func Cursor_advance(cursor *config.Cursor){
    cursor.Row_num += 1 
    if(cursor.Row_num >= cursor.Table.Num_rows){
        cursor.End_of_table = true
    }
}
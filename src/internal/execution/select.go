package execution

import (
	"nestdb/internal/config"
	"nestdb/internal/query"
	"nestdb/internal/utils"
)

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
    cursor := query.Table_stat(table)
    for cursor.End_of_table != true {
        data, _ := utils.Deserialize(query.Cursor_value(cursor))
        query.Cursor_advance(cursor)
        utils.Print_row(data)
    }
    return config.EXECUTE_SUCCESS
}
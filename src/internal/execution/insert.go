package execution

import (
	"nestdb/internal/config"
	"nestdb/internal/query"
	"nestdb/internal/utils"
)

// if the table current num row is alredy greater then return an error
//logic: figure out where to save (which page to save the data)
//in the page (find the current page number)
// then serialize the data based on the statement
// and then save the data
func Execute_insert(statement *config.Statement , table *(config.Table))(config.ExecuteResult){
	if table.Num_rows >= config.TABLE_MAX_ROWS {
		return config.EXECUTE_FAIL
	}
	

    row_to_insert := &(statement.Row_to_insert)
	data, _ := utils.Serialize(row_to_insert)
	//currently the data is set to bee serialize already

	// Get the address of the row
    /*
	// Find the page number and the specific position of the row

    */
    //page_number := table.Num_rows / ROWS_PER_PAGE
    //page := get_page(table.Pager , page_number)
    

	// Find the remainder by calculating the row size and get the address of that row
	//row_offset := table.Num_rows % ROWS_PER_PAGE
    cursor := query.Table_end(table)
	row_address := query.Cursor_value(cursor)
	copy(row_address, data)
	// Increment the number of rows
	table.Num_rows += 1
	// Flush the changes to disk
	//pager_flush(*table.Pager, page_number, PAGE_SIZE)
	return config.EXECUTE_SUCCESS
}
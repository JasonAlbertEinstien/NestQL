package table

import (
	"fmt"
	"nestdb/internal/config"
	"os"
	"syscall"
)

//this part will create a new table and return a pointer point to table
func Db_open(filename string)(*config.Table , error){
	pager , err := Pager_open(filename)
	if err != nil{
		return nil , err
	}
	num_rows :=  pager.File_length/config.ROW_SIZE
	table := &config.Table{
		Pager: pager,
		Num_rows: num_rows,
	}
	return table , nil
} 

//set the file closed and save it   
//only save the datat when the dataase is closed 
//save the data to the database (from in memory to disk and then close the database)
func Db_close(table *config.Table){
    pager := table.Pager
    num_full_pages := table.Num_rows/config.ROWS_PER_PAGE

    for i := 0 ; i < int(num_full_pages) ; i++ {
        if pager.Pages[i] == nil{
     
            continue
        }
        Pager_flush(*pager , uint32(i), config.PAGE_SIZE)
        pager.Pages[i] = nil
    }
    num_additional_rows := table.Num_rows % config.ROWS_PER_PAGE
    if(num_additional_rows >0){
        page_num := num_full_pages
        if(pager.Pages[page_num] != nil){
            Pager_flush(*pager , page_num , num_additional_rows*config.ROW_SIZE)
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
func Pager_flush(pager config.Pager , page_num uint32 ,size uint32){
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
package table

import "nestdb/internal/config"

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

func Db_close(){}
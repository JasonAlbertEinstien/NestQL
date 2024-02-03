package table

import (
	"fmt"
	config "nestdb/internal/config"
	"os"
	"syscall"
)

//open the file (if not then create one) --> return a file pointer
// logic:
/*=
  input the file name
  using openfile(filename , os.xxx , ox . xxx , if it's allowing user to delete than it should always be
  0777 if not it's better to be 0666)
  save the file pointer to at the off set end of the file
  the the pager should have 1. a fd , 2. file length (if new then 0 ) 3. a pages which is 4096 byte set as the max page size
*/
//
func Pager_open(Filename string)(*config.Pager , error){
	f, err := os.OpenFile(Filename , os.O_RDWR|os.O_CREATE , syscall.S_IWUSR|syscall.S_IRUSR)
	if err != nil {
		fmt.Println("ERROR OF OPEN THE DB FILE OR CREATING THE FILE")
		return nil , err
	}
	fileInfo , err := f.Stat()
	if err != nil {
		fmt.Println("ERROR CANNOT READ FILE DATA")
		return nil , err
	}

	fileLength := fileInfo.Size()
	pager := &config.Pager{
		File_descriptor: int(f.Fd()),
        File_length: uint32(fileLength),
        //remember to make enough space for it if not it will have error
        Pages: make([]*config.Page, config.TABLE_MAX_PAGES),
	}

	for i := 0 ; i<int(config.TABLE_MAX_PAGES);i++{
		pager.Pages[i] = nil
	}

	return pager , nil
}
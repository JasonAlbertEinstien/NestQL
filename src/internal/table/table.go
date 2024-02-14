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
        Num_pages: uint32(fileLength)/config.PAGE_SIZE,
	}

    if uint32(fileLength) % config.PAGE_SIZE != 0{
        fmt.Println("CORRUPTED FILE")
        os.Exit(1)
    }

	for i := 0 ; i<int(config.TABLE_MAX_PAGES);i++{
		pager.Pages[i] = nil
	}

	return pager , nil
}

//get that specific postion from the page
/*
    This section have bug ! Please be aware 
*/
func Get_page(pager *config.Pager ,page_num uint32)*config.Page{
    //if page number > the size --> exit the os.exit()


    if(page_num > config.TABLE_MAX_PAGES){
        fmt.Println("TRY TO ACCESS MORE THAN TLBLE MAXPAGE")
        os.Exit(1)
    }   
//    fmt.Printf("INLOOP")
    //if that page is not being fetched before 
    if(pager.Pages[page_num] == nil){
        //create a new pointer for one page 
        page := new(config.Page)
        num_pages := pager.File_length/config.PAGE_SIZE

        if (pager.File_length % config.PAGE_SIZE != 0){
            num_pages += 1
        }
        if(page_num <= num_pages){

            //bug exit here : 
            _, err := syscall.Seek(pager.File_descriptor , int64(page_num*config.PAGE_SIZE) , 0)
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
        if (page_num >= pager.Num_pages){
            pager.Num_pages = page_num + 1
        }
    }
    return pager.Pages[page_num]
}
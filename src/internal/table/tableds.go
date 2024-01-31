package tableds

//this pacakge is used for saving the datastructre of a table

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
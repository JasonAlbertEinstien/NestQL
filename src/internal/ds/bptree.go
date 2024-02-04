// package ds
// for testing purpose
package main

//rewrite whole bptree concept

//parent class for internal page and leaf page
//Resources: https://15445.courses.cs.cmu.edu/spring2023/project2/#b+tree-structure-2a

type BPTreePage struct{ 
    page_type uint32 //either internal or leaf
    size_ uint32
    max_size_ uint32
}

type BPTreeInternalPage struct{
    keys []int
    pages []*BPTreePage
}

type BPTreeLeafPage struct{
    keys []int
    //use interger as value for testing purpose/ replace this section with the data once finish this section.
    value []byte
}

//insertation deletion 
func InsertNode(id int , root *BPTreeInternalPage)(){}

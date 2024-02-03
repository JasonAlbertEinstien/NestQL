// package ds
// for testing purpose
package main

//
type BPTreeNode struct{
	isLeaf bool
	keys []int
	children []*BPTreeNode
}
//create a new node
func NewBPTreeNode(leaf bool)(*BPTreeNode){
	return &BPTreeNode{
		isLeaf: leaf,
		keys: []int{},
		children: []*BPTreeNode{},
	}
}

type BPTree struct{
	root *BPTreeNode
}

func NewBPTree() (*BPTree){
	return &BPTree{
		root: NewBPTreeNode(true),
	}
}

func (t *BPTree) Insert(key int){
	node := t.root
}

func main(){

}
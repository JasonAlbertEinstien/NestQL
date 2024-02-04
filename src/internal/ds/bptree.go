// package ds
// for testing purpose
package main

import "nestdb/internal/config"

//Implementing B-Tree Structure
//Node Header Format
//format of a leaf node should be
/*
   node_type 1 byte
   is roote 1 byte
   parent_pointer 2 - 5 byte
   byte 6 - 9 num_cells
   byte 14- 306 value 0
   vyte 307 -310 key1
   byte page xszie : value 1 ...
   each leaf node will have 12 values and key pary
   there are some byte is wasted

*/
type NodeType int 

const(
    NODE_INTERNAL NodeType = 0
    NODE_LEAF NodeType = 1
)

//NODE HEADER LAYOUT
const(
    NODE_TYPE_SIZE = 1
    NODE_TYPE_OFFSET = 0
    IS_ROOT_SIZE = 1
    IS_ROOT_OFFSET = NODE_TYPE_SIZE
    PARENT_POINTER_SIZE = 4
    PARENT_POINTER_OFFSET = IS_ROOT_OFFSET +IS_ROOT_SIZE
    COMMONE_NODE_HEADER_SIZE = NODE_TYPE_SIZE + IS_ROOT_SIZE + PARENT_POINTER_SIZE
)

//LEAF NODE FORMAT
const(
    LEAF_NODE_KEY_SIZE = 4
    LEAF_NODE_KEY_OFFSET = 0
    LEAF_NODE_HEADER_SIZE = LEAF_NODE_KEY_SIZE + COMMONE_NODE_HEADER_SIZE
)

var LEAF_NODE_VALUE_SZIZE uint32 = config.ROW_SIZE
var LEAF_NODE_VALUE_OFFSET uint32 = LEAF_NODE_KEY_OFFSET +LEAF_NODE_KEY_SIZE
var LEAF_NODE_CELL_SIZE = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SZIZE
var LEAF_NODE_SPACE_FOR_CELLS uint32 = config.PAGE_SIZE - LEAF_NODE_HEADER_SIZE
var LEAF_NODE_MAX_CELL uint32 = LEAF_NODE_SPACE_FOR_CELLS / LEAF_NODE_CELL_SIZE

//how to access the key with node 
//how insert is done ?
//How to tree is form ? 
type NodeTree struct{
    node *Node
}

type Node struct{
    is_leaf bool
}

//currently testing in ds bptreee.go file
//insertion of b tree
func insertNode(){}
func leftRotation(){}
func rightRotation(){}
//delete of b tree
func deleteNode(){}
//search in b tree 
// the algorithm that is used is simiplar to binary tree
func searchNode(){}


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
	if len(node.keys) == 0 {
		node.keys = append(node.keys, key)
		return
	}

	if node.isLeaf{
		//insert and then sort it based on the key
	}

}

func main(){

}
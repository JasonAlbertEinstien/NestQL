// package ds
// for testing purpose
package tree

import "nestdb/internal/config"

type NodeType uint8
const(
    NODE_INTERNAL NodeType = 0
    NODE_LEAF NodeType = 1
    NODE_TYPE_SIZE uint32 = 1
    NODE_TYPE_OFFSET uint32 = 0
    IS_ROOT_SIZE uint32 = 1
    IS_ROOT_OFFSET uint32 = NODE_TYPE_SIZE
    PARENT_POINTER_SIZE uint32 = 4
    PARENT_POINTER_OFFSET = IS_ROOT_OFFSET + IS_ROOT_SIZE
    COMMON_NODE_HEADER_SIZE = NODE_TYPE_SIZE + IS_ROOT_SIZE + PARENT_POINTER_SIZE
)


/*
    Leaf node header 
*/
const(
    LEAF_NODE_NUM_CELLS_SIZE = 4
    LEAF_NODE_NUM_CELLS_OFFSET = COMMON_NODE_HEADER_SIZE
    LEAF_NODE_HEADER_SIZE = COMMON_NODE_HEADER_SIZE + LEAF_NODE_NUM_CELLS_SIZE
)

/*
    LEAF BODY LAYOUT
*/

const(
    LEAF_NODE_KEY_SIZE = 4 
    LEAF_NODE_KEY_OFFSET = 0 
)

var LEAF_NODE_VALUE_SIZE uint32= config.ROW_SIZE
var LEAF_NODE_VALUE_OFFSET uint32 = LEAF_NODE_KEY_OFFSET + LEAF_NODE_KEY_SIZE
var LEAF_NODE_CELL_SIZE uint32 = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SIZE
var LEAF_NODE_SPACE_FOR_CELL uint32 = config.PAGE_SIZE - LEAF_NODE_HEADER_SIZE
var LEAF_NODE_MAX_CELLS uint32 = LEAF_NODE_SPACE_FOR_CELL / LEAF_NODE_CELL_SIZE

func Leaf_node_num_cells (node *[]byte)([]byte){
    return (*node)[LEAF_NODE_NUM_CELLS_OFFSET:LEAF_NODE_NUM_CELLS_OFFSET+LEAF_NODE_NUM_CELLS_OFFSET]
}

func Leaf_node_cell(node *[]byte , cell_num uint32)([]byte){
    return (*node)[LEAF_NODE_HEADER_SIZE + cell_num *LEAF_NODE_CELL_SIZE: LEAF_NODE_HEADER_SIZE + (cell_num+1) *LEAF_NODE_CELL_SIZE -1]
}

func Leaf_node_key(node *[]byte , cell_num uint32)([]byte){
    return Leaf_node_cell(node , cell_num)
}  

func Initialize_leaf_node(node *[]byte){
    Leaf_node_num_cells(node)[0] = 0
}
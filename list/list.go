package list

import "errors"

/*
 Doubly LinkedList

 - List of nodes
 - each node is going to reference next and prev


 API:
 - InsertRoot(val)
 - Insert(ele, pos)
 - Remove(pos)
 - Replace(value, pos)
 - GetRootVal

 Node: object
 - next
 - prev
 - value

 List: object
 - len
 - root (head, with access to tail)
*/

var (
	// ErrInvalidIndex occurs when the specified index is out of range
	ErrInvalidIndex = errors.New("invalid positional index")
)

// List is the LinkedList struct
type List struct {
	len  int
	root *Node
}

// Node represents a node in the linked list
type Node struct {
	prev  *Node
	next  *Node
	value interface{}
}

// New creates a new LinkedList
func New() *List {
	return &List{
		root: &Node{
			next: new(Node),
			prev: new(Node),
		},
		len: 0,
	}
}

// Length returns the length of the LinkedList
func (l *List) Length() int {
	return l.len
}

// InsertRoot inserts a value into the LinkedList
func (l *List) InsertRoot(val interface{}) *List {
	if l.len == 0 {
		l.root.value = val
		l.len++
	}

	return l
}

// GetRootVal returns the value at the root node
func (l *List) GetRootVal() interface{} {
	return l.root.value
}

// Insert inserts a value at the positional index
func (l *List) Insert(val interface{}, pos int) (*List, error) {
	if pos < 0 || pos > l.len {
		return l, ErrInvalidIndex
	}

	currentIdx := 0
	currentNode := l.root

	// start at 0
	// iterate over nodes until we get to position
	for currentIdx != pos-1 {
		currentNode = currentNode.next
		currentIdx++
	}

	// save the prev of the current node as temp var
	temp := currentNode.prev

	// create new node with val
	// make the next of the new node be the current node
	newNode := &Node{
		value: val,
	}
	newNode.next = currentNode

	// put prev of new node as the temp var
	newNode.prev = temp
	// put the next of the temp var as the new node
	temp.next = newNode

	l.len++

	return l, nil
}

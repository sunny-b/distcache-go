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
	tail *Node
}

// Node represents a node in the linked list
type Node struct {
	prev  *Node
	next  *Node
	value interface{}
}

// New creates a new LinkedList
func New() *List {
	root := new(Node)
	return &List{
		root: root,
		tail: root,
		len:  0,
	}
}

// Length returns the length of the Linked List
func (l *List) Length() int {
	return l.len
}

// Append appends a value to the list at the end
func (l *List) Append(val interface{}) (err error) {
	defer func() {
		if err == nil {
			l.len++
		}
	}()
	if l.len == 0 {
		l.root.value = val
		return nil
	}

	currentNode, err := l.getNodeAtPos(pos)
	if err != nil {
		return err
	}

	// save the prev of the current node as temp var
	temp := currentNode.prev

	// create new node with val
	// make the next of the new node be the current node
	newNode := &Node{
		value: val,
		next:  currentNode,
		prev:  temp,
	}

	// put the next of the temp var as the new node
	temp.next = newNode

	l.len++

	return nil
}

// Get returns the node at a given position
func (l *List) Get(pos int) *Node {
	node, err := l.getNodeAtPos(pos)
	if err != nil {
		return nil
	}

	return node
}

func (l *List) getNodeAtPos(pos int) (*Node, error) {
	if pos < 0 || pos > l.len {
		return nil, ErrInvalidIndex
	}

	currentIdx := 0
	currentNode := l.root

	// start at 0
	// iterate over nodes until we get to position
	for currentIdx < pos-1 {
		currentNode = currentNode.next
		currentIdx++
	}

	return currentNode, nil
}

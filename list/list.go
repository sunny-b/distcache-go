package list

/*
 Doubly LinkedList

 - List of nodes
 - each node is going to reference next and prev


 API:
 - PushFront(val) - X
 - Remove(node) - X
 - MoveFront(node)
 - Head() - X
 - Tail() - X

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
	ErrInvalidIndex = "invalid positional index"
)

// List is the LinkedList struct
type List struct {
	len  int
	root Node
}

// Node represents a node in the linked list
type Node struct {
	prev  *Node
	next  *Node
	Value interface{}
}

// New creates a new LinkedList
func New() *List {
	l := new(List)
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0

	return l
}

// Length returns the length of the Linked List
func (l *List) Length() int {
	return l.len
}

// Append appends a value to the list at the end
func (l *List) PushFront(val interface{}) *Node {
	defer func() {
		l.len++
	}()

	// create new node with val
	// make the next of the new node be the current node
	newNode := &Node{
		Value: val,
		next:  l.root.next,
		prev:  &l.root,
	}

	l.root.next = newNode
	newNode.next.prev = newNode

	return newNode
}

// Head returns the head node of the LinkedList
func (l *List) Head() *Node {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Tail returns the tail node of the LinkedList
func (l *List) Tail() *Node {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// Remove is going to remove a Node from the LinkedList
func (l *List) Remove(node *Node) {
	node.next.prev = node.prev
	node.prev.next = node.next
	node.next = nil
	node.prev = nil
	l.len--
}

// Movefront is going to move the passed in Node to the front of the Linked List
func (l *List) MoveFront(node *Node) {
	currentFront := l.root.next

	node.prev = &l.root
	l.root.next = node

	currentFront.prev = node
	node.next = currentFront
}

// Get returns the node at a given position
//func (l *List) Get(pos int) *Node {
//node := l.getNodeAtPos(pos)

//return node
//}

//func (l *List) getNodeAtPos(pos int) *Node {
//if pos < 0 || pos > l.len {
//panic(ErrInvalidIndex)
//}

//currentIdx := 0
//currentNode := l.root

//// start at 0
//// iterate over nodes until we get to position
//for currentIdx < pos-1 {
//currentNode = currentNode.next
//currentIdx++
//}

//return currentNode
//}

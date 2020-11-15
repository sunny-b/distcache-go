package list

/*
 Doubly LinkedList

 - List of nodes
 - each node is going to reference next and prev


 API:
 - InsertRoot(val)
 - Insert(ele, pos)
 - Remove(pos)
 - Replace(value, pos)

 Node: object
 - next
 - prev
 - value

 List: object
 - len
 - root (head, with access to tail)
*/

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
		root: new(Node),
		len:  0,
	}
}

// InsertRoot inserts a value into the LinkedList
func (l *List) InsertRoot(val interface{}) bool {
	if l.len == 0 {
		l.root.value = val
		l.len++
	}

	return true
}

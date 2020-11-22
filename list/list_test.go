package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushFront(t *testing.T) {
	// when the list is empty, PushFront should set the root value as the passed in value
	val := 5
	t.Run("empty list", func(t *testing.T) {
		l := New()
		node := l.PushFront(5)

		assert.Equal(t, val, node.Value.(int))
		assert.Equal(t, val, l.Head().Value.(int))
	})

	// when the list is not empty, PushFront creates a new node and sets it as the head value
	t.Run("non-empty list", func(t *testing.T) {
		foobar := "foobar"
		l := New()

		l.PushFront(5)
		node := l.PushFront(foobar)

		assert.Equal(t, foobar, node.Value.(string))
		assert.Equal(t, foobar, l.Head().Value.(string))
		assert.Equal(t, 5, l.Tail().Value.(int))
	})
}

func TestRemove(t *testing.T) {
	t.Run("remove node", func(t *testing.T) {
		l := New()
		node := l.PushFront("foo")

		l.Remove(node)

		assert.Equal(t, 0, l.Length())
		assert.Nil(t, l.Head())
	})

	t.Run("remove tail", func(t *testing.T) {
		l := New()
		node1 := l.PushFront("foo")
		node2 := l.PushFront(5)

		assert.Equal(t, node1, l.Tail())
		assert.Equal(t, 2, l.Length())

		l.Remove(l.Tail())

		assert.Equal(t, node2, l.Tail())
		assert.Equal(t, 1, l.Length())
	})
}

func TestMoveFront(t *testing.T) {
	t.Run("move node to the front", func(t *testing.T) {
		l := New()
		node1 := l.PushFront(5)
		node2 := l.PushFront("foobar")

		assert.Equal(t, node2, l.Head())

		l.MoveFront(node1)

		assert.Equal(t, node1, l.Head())
	})
}

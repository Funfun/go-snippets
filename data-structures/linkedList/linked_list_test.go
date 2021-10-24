package linked_list

import (
	"fmt"
	"testing"
)

func TestLinkedListAppend(t *testing.T) {
	t.Run("append to linked list", func(t *testing.T) {
		l := LinkedList{}
		l.Append(&Node{Value: 77})

		if l.Length != 1 {
			t.Error("want length eq to 1")
		}
	})
}

func TestLinkedListString(t *testing.T) {
	t.Run("linkedList satisfies Stringer interface", func(t *testing.T) {
		var _ fmt.Stringer = (*LinkedList)(nil)
	})
}

func TestLinkedListDeleteWithValue(t *testing.T) {
	t.Run("remove node with search value", func(t *testing.T) {
		l := LinkedList{Head: &Node{Value: 11}, Length: 1}
		l.Append(&Node{Value: 12})
		l.Append(&Node{Value: 21})

		l.DeleteWithValue(12)

		if l.Length != 2 {
			t.Error("want length to be eq 2, got: ", l.Length)
		}
	})

	t.Run("when linkedList has 2 nodes with search values", func(t *testing.T) {
		l := LinkedList{Head: &Node{Value: 11}, Length: 1}
		l.Append(&Node{Value: 11})
		l.Append(&Node{Value: 10})

		l.DeleteWithValue(11)

		if l.Length != 1 {
			t.Error("want length to be eq 2, got: ", l.Length)
		}
	})

	t.Run("remove node should remain unchanged when search value is not found", func(t *testing.T) {
		l := LinkedList{Head: &Node{Value: 11}, Length: 1}
		l.Append(&Node{Value: 12})
		l.Append(&Node{Value: 21})

		l.DeleteWithValue(99)

		if l.Length != 3 {
			t.Error("want length to be eq 3, got: ", l.Length)
		}
	})
}

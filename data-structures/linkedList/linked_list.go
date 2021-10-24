package linked_list

import (
	"fmt"
	"strings"
)

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	Head   *Node
	Length int
}

// Append takes O(1) speed
func (l *LinkedList) Append(node *Node) {
	t := l.Head
	l.Head = node
	l.Head.Next = t
	l.Length++
}

func (l LinkedList) String() string {
	var out strings.Builder
	node := l.Head
	for {
		if node == nil {
			break
		}
		fmt.Fprintf(&out, "Node = %d", node.Value)
		node = node.Next

	}

	return out.String()
}

func (l *LinkedList) DeleteWithValue(value int) {
	node := l.Head
	if node.Value == value {
		l.Head = l.Head.Next
		l.Length--
		return
	}

	for {
		if node.Next == nil {
			return
		}

		if node.Next.Value == value {
			node.Next = node.Next.Next
			l.Length--
		} else {
			node = node.Next
		}
	}
}

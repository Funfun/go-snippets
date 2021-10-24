package tree

import (
	"fmt"
	"io"
)

type Tree struct {
	Left  *Tree
	Right *Tree
	Value int
}

func (t Tree) Print(w io.Writer) {
	if t.Left != nil {
		t.Left.Print(w)
	}

	fmt.Fprintf(w, "Value = %d\n", t.Value)

	if t.Right != nil {
		t.Right.Print(w)
	}
}

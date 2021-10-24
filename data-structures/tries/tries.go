package tries

import (
	"fmt"
	"io"
	"net/url"
)

type Page struct {
	Title string
	Url   url.URL
}

func (p Page) String() string {
	return fmt.Sprintf("Title = %s, URL = %s", p.Title, p.Url.String())
}

type Tries struct {
	Nodes map[rune]*Tries
	Value *Page
}

func (t *Tries) Insert(key string, page *Page) {
	node := t
	for _, c := range key {
		if _, ok := node.Nodes[c]; !ok {
			node.Nodes[c] = &Tries{Nodes: make(map[rune]*Tries), Value: nil}
		}
		node = node.Nodes[c]
	}

	node.Value = page
}

func (t *Tries) Print(w io.Writer) {
	fmt.Fprintf(w, "Value = %s\n", t.Value)

	for key, node := range t.Nodes {
		fmt.Fprintf(w, "Key = %s%s", string(key), " ")
		node.Print(w)
	}
}

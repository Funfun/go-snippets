package tries

import (
	"net/url"
	"os"
	"testing"
)

func TestTries(t *testing.T) {
	t.Run("insert", func(t *testing.T) {
		trie := Tries{Nodes: make(map[rune]*Tries)}
		trie.Insert("abc", &Page{Title: "ABC", Url: url.URL{Host: "example.com", Scheme: "https", Path: "ABC"}})
		trie.Insert("abd", &Page{Title: "ABD", Url: url.URL{Host: "example.com", Scheme: "https", Path: "ABD"}})
		trie.Insert("aab", &Page{Title: "AAC", Url: url.URL{Host: "example.com", Scheme: "https", Path: "AAB"}})
		trie.Print(os.Stdout)

		t.Error(1)
	})
}

package tree

import (
	"os"
	"testing"
)

func TestTree(t *testing.T) {
	t.Run("print test", func(t *testing.T) {
		tree := Tree{Left: &Tree{Value: 99}, Right: &Tree{Value: 100}, Value: 11}
		tree.Print(os.Stderr)
	})
}

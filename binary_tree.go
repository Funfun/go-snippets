func kthSmallest(root *TreeNode, k int) int {
    t := Constructor(root)
    c := 0
    for t.HasNext() {
        c++
        v := t.Next() 
        if k == c {
            return v
        }      
    }

    return -1
}

type BSTIterator struct {
	stack []*TreeNode
}

func Constructor(root *TreeNode) BSTIterator {
	st := []*TreeNode{}
	cur := root
	for cur != nil {
		st = append(st, cur)
		cur = cur.Left
	}
	return BSTIterator{
		stack: st,
	}
}

/** @return the next smallest number */
func (this *BSTIterator) Next() int {
	cur := this.stack[len(this.stack)-1]
	this.stack = this.stack[:len(this.stack)-1]
	ret := cur.Val
	cur = cur.Right
	for cur != nil {
		this.stack = append(this.stack, cur)
		cur = cur.Left
	}

	return ret
}

/** @return whether we have a next smallest number */
func (this *BSTIterator) HasNext() bool {
	return len(this.stack) > 0
}

package btree

import "testing"

func TestBTree(t *testing.T) {
	bt := New(5)
	for i := 0; i < 20; i++ {
		bt.Insert(i, nil)
	}
	bt.Print()
	bt.Remove(5)
	bt.Print()
}

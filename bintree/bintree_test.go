package bintree

import (
	"testing"

	"github.com/mooncaker816/gostructure/queue"
)

func TestInsertAsRoot(t *testing.T) {
	q := queue.NewQueue()
	bt := new(BinTree)
	root := bt.InsertAsRoot(0, nil)
	q.Enqueue(root)
	key := 1
	for l := 0; l < 4; l++ {
		for i := 0; i < 1<<uint(l); i++ {
			ni, _ := q.Dequeue()
			n := ni.(*Node)
			q.Enqueue(bt.InsertAsLChild(n, key, nil))
			key++
			q.Enqueue(bt.InsertAsRChild(n, key, nil))
			key++
		}
	}
	if bt.Root.Height != 4 {
		t.Errorf("Got %v expected %v", bt.Root.Height, 4)
	}
	if bt.Size != 1<<uint(bt.Root.Height+1)-1 {
		t.Errorf("Got %v expected %v", bt.Size, 1<<uint(bt.Root.Height+1)-1)
	}
	bt.Print()
	bt2 := new(BinTree)
	root2 := bt2.InsertAsRoot("k", nil)
	i := bt2.InsertAsLChild(root2, "i", nil)
	bt2.InsertAsRChild(root2, "j", nil)
	h := bt2.InsertAsRChild(i, "h", nil)
	b := bt2.InsertAsLChild(h, "b", nil)
	g := bt2.InsertAsRChild(h, "g", nil)
	bt2.InsertAsRChild(b, "a", nil)
	e := bt2.InsertAsLChild(g, "e", nil)
	bt2.InsertAsRChild(g, "f", nil)
	bt2.InsertAsLChild(e, "c", nil)
	bt2.InsertAsRChild(e, "d", nil)
	if bt2.Root.Height != 5 {
		t.Errorf("Got %v expected %v", bt2.Root.Height, 5)
	}
	if bt2.Size != 11 {
		t.Errorf("Got %v expected %v", bt2.Size, 11)
	}
	bt2.Print()
}

func TestIsRoot(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if !a.IsRoot() {
		t.Errorf("Got %v expected %v", a.IsRoot(), true)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if b.IsRoot() {
		t.Errorf("Got %v expected %v", b.IsRoot(), false)
	}
}

func TestIsLeftChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.IsLChild() {
		t.Errorf("Got %v expected %v", a.IsLChild(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if !b.IsLChild() {
		t.Errorf("Got %v expected %v", b.IsLChild(), true)
	}
	c := bt.InsertAsRChild(a, 3, nil)
	if c.IsLChild() {
		t.Errorf("Got %v expected %v", c.IsLChild(), false)
	}
}
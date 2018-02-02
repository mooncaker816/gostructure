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

func TestIsLChild(t *testing.T) {
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

func TestIsRChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.IsRChild() {
		t.Errorf("Got %v expected %v", a.IsRChild(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if b.IsRChild() {
		t.Errorf("Got %v expected %v", b.IsRChild(), false)
	}
	c := bt.InsertAsRChild(a, 3, nil)
	if !c.IsRChild() {
		t.Errorf("Got %v expected %v", c.IsRChild(), true)
	}
}

func TestHasParent(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasParent() {
		t.Errorf("Got %v expected %v", a.HasParent(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if !b.HasParent() {
		t.Errorf("Got %v expected %v", b.HasParent(), true)
	}
}

func TestHasLChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasLChild() {
		t.Errorf("Got %v expected %v", a.HasLChild(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if !a.HasLChild() {
		t.Errorf("Got %v expected %v", a.HasLChild(), true)
	}
	if b.HasLChild() {
		t.Errorf("Got %v expected %v", b.HasLChild(), false)
	}
}

func TestHasRChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasRChild() {
		t.Errorf("Got %v expected %v", a.HasRChild(), false)
	}
	b := bt.InsertAsRChild(a, 2, nil)
	if !a.HasRChild() {
		t.Errorf("Got %v expected %v", a.HasRChild(), true)
	}
	if b.HasRChild() {
		t.Errorf("Got %v expected %v", b.HasRChild(), false)
	}
}

func TestHasChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasChild() {
		t.Errorf("Got %v expected %v", a.HasChild(), false)
	}
	b := bt.InsertAsRChild(a, 2, nil)
	if !a.HasChild() {
		t.Errorf("Got %v expected %v", a.HasChild(), true)
	}
	if b.HasChild() {
		t.Errorf("Got %v expected %v", b.HasChild(), false)
	}
}

func TestHasBothChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasBothChild() {
		t.Errorf("Got %v expected %v", a.HasBothChild(), false)
	}
	bt.InsertAsRChild(a, 2, nil)
	if a.HasBothChild() {
		t.Errorf("Got %v expected %v", a.HasBothChild(), false)
	}
	bt.InsertAsLChild(a, 3, nil)
	if !a.HasBothChild() {
		t.Errorf("Got %v expected %v", a.HasBothChild(), true)
	}
}

func TestIsLeaf(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if !a.IsLeaf() {
		t.Errorf("Got %v expected %v", a.IsLeaf(), true)
	}
	b := bt.InsertAsRChild(a, 2, nil)
	if !b.IsLeaf() {
		t.Errorf("Got %v expected %v", b.IsLeaf(), true)
	}
	if a.IsLeaf() {
		t.Errorf("Got %v expected %v", a.IsLeaf(), false)
	}
}

func TestSibling(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.Sibling() != nil {
		t.Errorf("Got %v expected %v", a.Sibling(), nil)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	c := bt.InsertAsRChild(a, 3, nil)
	if b.Sibling() == nil || b.Sibling() != c {
		t.Errorf("Got %v expected %v", b.Sibling(), c)
	}
	if c.Sibling() == nil || c.Sibling() != b {
		t.Errorf("Got %v expected %v", c.Sibling(), b)
	}
}

func TestUncle(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	b := bt.InsertAsLChild(a, 2, nil)
	c := bt.InsertAsRChild(a, 3, nil)
	d := bt.InsertAsLChild(b, 4, nil)
	e := bt.InsertAsRChild(b, 5, nil)
	f := bt.InsertAsLChild(c, 6, nil)
	g := bt.InsertAsRChild(c, 7, nil)
	if d.Uncle() == nil || d.Uncle() != c || d.Uncle() != e.Uncle() {
		t.Errorf("Got %v expected %v", d.Uncle(), c)
	}
	if f.Uncle() == nil || f.Uncle() != b || f.Uncle() != g.Uncle() {
		t.Errorf("Got %v expected %v", f.Uncle(), b)
	}
}

func TestSize(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	b := bt.InsertAsLChild(a, 2, nil)
	c := bt.InsertAsRChild(a, 3, nil)
	d := bt.InsertAsLChild(b, 4, nil)
	f := bt.InsertAsLChild(c, 6, nil)
	g := bt.InsertAsRChild(c, 7, nil)
	if a.Size() != 6 {
		t.Errorf("Got %v expected %v", a.Size(), 6)
	}
	if b.Size() != 2 {
		t.Errorf("Got %v expected %v", b.Size(), 2)
	}
	if c.Size() != 3 {
		t.Errorf("Got %v expected %v", c.Size(), 3)
	}
	if d.Size() != 1 || d.Size() != f.Size() || d.Size() != g.Size() {
		t.Errorf("Got %v expected %v", d.Size(), 1)
	}
}

package avl

import (
	"fmt"
	"testing"

	"github.com/mooncaker816/gostructure/bintree"
)

func TestInsert(t *testing.T) {
	at := NewAvl()
	root := at.InsertAsRoot(16, nil)
	lc := at.InsertAsLChild(root, 10, nil)
	rc := at.InsertAsRChild(root, 25, nil)
	llc := at.InsertAsLChild(lc, 5, nil)
	rlc := at.InsertAsRChild(lc, 11, nil)
	lrc := at.InsertAsLChild(rc, 19, nil)
	rrc := at.InsertAsRChild(rc, 28, nil)
	lllc := at.InsertAsLChild(llc, 2, nil)
	at.InsertAsRChild(llc, 8, nil)
	//lrlc := at.InsertAsLChild(rlc, 2, nil)
	at.InsertAsRChild(rlc, 15, nil)
	at.InsertAsLChild(lrc, 17, nil)
	at.InsertAsRChild(lrc, 22, nil)
	at.InsertAsLChild(rrc, 27, nil)
	rrrc := at.InsertAsRChild(rrc, 37, nil)
	at.InsertAsRChild(lllc, 4, nil)
	//at.InsertAsLChild(rrlc, 13, nil)
	at.InsertAsLChild(rrrc, 33, nil)
	at.Print()
	at.Insert(13, nil)
	at.Print()
	at.Insert(30, nil)
	at.Print()
	at2 := NewAvl()
	at2.InsertAsRoot(100, nil)
	at2.Print()
	at2.Insert(95, nil)
	at2.Print()
	at2.Insert(85, nil)
	at2.Print()
	at2.Insert(75, nil)
	at2.Print()
	at2.Insert(65, nil)
	at2.Print()
	at2.Insert(55, nil)
	at2.Print()
	at2.Insert(45, nil)
	at2.Print()
	at2.Insert(35, nil)
	at2.Print()
	at2.TravLevel(withprintheight())
}

func TestRemove(t *testing.T) {
	at := NewAvl()
	root := at.InsertAsRoot(16, nil)
	lc := at.InsertAsLChild(root, 10, nil)
	rc := at.InsertAsRChild(root, 25, nil)
	llc := at.InsertAsLChild(lc, 5, nil)
	rlc := at.InsertAsRChild(lc, 11, nil)
	lrc := at.InsertAsLChild(rc, 19, nil)
	rrc := at.InsertAsRChild(rc, 28, nil)
	lllc := at.InsertAsLChild(llc, 2, nil)
	at.InsertAsRChild(llc, 8, nil)
	//lrlc := at.InsertAsLChild(rlc, 2, nil)
	at.InsertAsRChild(rlc, 15, nil)
	at.InsertAsLChild(lrc, 17, nil)
	at.InsertAsRChild(lrc, 22, nil)
	at.InsertAsLChild(rrc, 27, nil)
	rrrc := at.InsertAsRChild(rrc, 37, nil)
	at.InsertAsRChild(lllc, 4, nil)
	//at.InsertAsLChild(rrlc, 13, nil)
	at.InsertAsLChild(rrrc, 33, nil)
	at.Print()
	at.Remove(8)
	at.Print()
	at.Remove(16)
	at.Print()
	at.Remove(19)
	at.Print()
	at.TravLevel(withprintheight())
	at2 := NewAvl()
	r2 := at2.InsertAsRoot(5, nil)
	at2.InsertAsLChild(r2, 1, nil)
	x := at2.InsertAsRChild(r2, 10, nil)
	at2.InsertAsRChild(x, 11, nil)
	at2.Print()
	at2.Remove(1)
	at2.Print()
	if at2.Root != x {
		t.Errorf("root got %v expected %v\n", at2.Root, x)
	}
	at2.TravLevel(withprintheight())
}

func withprintheight() func(n *bintree.Node) {
	return func(n *bintree.Node) {
		fmt.Printf("%d ", n.Height)
	}
}

package bst

import (
	"fmt"
	"os"
	"testing"

	"github.com/mooncaker816/gostructure/bintree"
)

func TestSearch(t *testing.T) {
	bst := NewBst()
	root := bst.InsertAsRoot(16, nil)
	lc := bst.InsertAsLChild(root, 10, nil)
	rc := bst.InsertAsRChild(root, 25, nil)
	llc := bst.InsertAsLChild(lc, 5, nil)
	rlc := bst.InsertAsRChild(lc, 11, nil)
	lrc := bst.InsertAsLChild(rc, 19, nil)
	rrc := bst.InsertAsRChild(rc, 28, nil)
	lllc := bst.InsertAsLChild(llc, 2, nil)
	bst.InsertAsRChild(llc, 8, nil)
	//lrlc := bst.InsertAsLChild(rlc, 2, nil)
	rrlc := bst.InsertAsRChild(rlc, 15, nil)
	bst.InsertAsLChild(lrc, 17, nil)
	bst.InsertAsRChild(lrc, 22, nil)
	bst.InsertAsLChild(rrc, 27, nil)
	rrrc := bst.InsertAsRChild(rrc, 37, nil)
	bst.InsertAsRChild(lllc, 4, nil)
	bst.InsertAsLChild(rrlc, 13, nil)
	bst.InsertAsLChild(rrrc, 33, nil)
	bst.Print()
	n := bst.Search(15)
	if n != rrlc {
		t.Errorf("Search got %v expected %v", n, rrlc)
	}
	if bst.Hot != rrlc.Parent {
		t.Errorf("Search Hot got %v expected %v", bst.Hot, rrlc.Parent)
	}
	n = bst.Search(14)
	if n != nil {
		t.Errorf("Search got %v expected %v", n, nil)
	}
}

func TestInsert(t *testing.T) {
	bst := NewBst()
	root := bst.InsertAsRoot(16, nil)
	lc := bst.InsertAsLChild(root, 10, nil)
	rc := bst.InsertAsRChild(root, 25, nil)
	llc := bst.InsertAsLChild(lc, 5, nil)
	rlc := bst.InsertAsRChild(lc, 11, nil)
	lrc := bst.InsertAsLChild(rc, 19, nil)
	rrc := bst.InsertAsRChild(rc, 28, nil)
	lllc := bst.InsertAsLChild(llc, 2, nil)
	bst.InsertAsRChild(llc, 8, nil)
	//lrlc := bst.InsertAsLChild(rlc, 2, nil)
	rrlc := bst.InsertAsRChild(rlc, 15, nil)
	bst.InsertAsLChild(lrc, 17, nil)
	bst.InsertAsRChild(lrc, 22, nil)
	bst.InsertAsLChild(rrc, 27, nil)
	rrrc := bst.InsertAsRChild(rrc, 37, nil)
	bst.InsertAsRChild(lllc, 4, nil)
	bst.InsertAsLChild(rrlc, 13, nil)
	bst.InsertAsLChild(rrrc, 33, nil)
	bst.Print()
	oldsize := bst.Size
	n := bst.Insert(14, nil)
	bst.Print()
	if bst.Size != oldsize+1 {
		t.Errorf("Insert size got %v expected %v", bst.Size, oldsize+1)
	}
	if bst.Hot != n.Parent {
		t.Errorf("Insert Hot got %v expected %v", bst.Hot, n.Parent)
	}
	bst.TravIn1(bintree.WithPrintNodeKey(os.Stdout))
	fmt.Println()
}

func TestRemove(t *testing.T) {
	bst := NewBst()
	root := bst.InsertAsRoot(16, nil)
	lc := bst.InsertAsLChild(root, 10, nil)
	rc := bst.InsertAsRChild(root, 25, nil)
	llc := bst.InsertAsLChild(lc, 5, nil)
	rlc := bst.InsertAsRChild(lc, 11, nil)
	lrc := bst.InsertAsLChild(rc, 19, nil)
	rrc := bst.InsertAsRChild(rc, 28, nil)
	lllc := bst.InsertAsLChild(llc, 2, nil)
	bst.InsertAsRChild(llc, 8, nil)
	//lrlc := bst.InsertAsLChild(rlc, 2, nil)
	rrlc := bst.InsertAsRChild(rlc, 15, nil)
	bst.InsertAsLChild(lrc, 17, nil)
	bst.InsertAsRChild(lrc, 22, nil)
	bst.InsertAsLChild(rrc, 27, nil)
	rrrc := bst.InsertAsRChild(rrc, 37, nil)
	bst.InsertAsRChild(lllc, 4, nil)
	bst.InsertAsLChild(rrlc, 13, nil)
	bst.InsertAsLChild(rrrc, 33, nil)
	bst.Print()
	oldsize := bst.Size
	ok := bst.Remove(16)
	bst.Print()
	if ok != true {
		t.Errorf("Remove got %v expected %v", ok, true)
	}
	if bst.Size != oldsize-1 {
		t.Errorf("Remove size got %v expected %v", bst.Size, oldsize-1)
	}
	bst.TravIn1(bintree.WithPrintNodeKey(os.Stdout))
	fmt.Println()
	fmt.Println(bst.Root)
}

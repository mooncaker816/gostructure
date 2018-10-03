package bst_test

import (
	"fmt"
	"testing"

	_ "github.com/mooncaker816/gostructure/bst/avl"
	_ "github.com/mooncaker816/gostructure/bst/btree"
	_ "github.com/mooncaker816/gostructure/bst/redblack"
	_ "github.com/mooncaker816/gostructure/bst/splay"

	"github.com/mooncaker816/gostructure/bst"
)

func TestAVL(t *testing.T) {
	fmt.Println("AVL START")
	at := bst.New(bst.AVL)
	keys := []int{16, 10, 25, 5, 11, 19, 28, 2, 8, 15, 17, 22, 27, 37, 4, 33}
	for _, v := range keys {
		at.Insert(v, nil)
	}
	at.Print()
	at.Insert(13, nil)
	at.Print()
	at.Insert(30, nil)
	at.Print()
	at2 := bst.New(bst.AVL)
	at2.Insert(100, nil)
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
	at2.Walk(bst.LevelOrder, withprintheight())
	fmt.Println()

	keys2 := []int{16, 10, 25, 5, 11, 19, 28, 2, 8, 15, 17, 22, 27, 37, 4, 33}
	at = bst.New(bst.AVL)
	for _, v := range keys2 {
		at.Insert(v, nil)
	}
	at.Print()
	at.Remove(8)
	at.Print()
	at.Remove(16)
	at.Print()
	at.Remove(19)
	at.Print()
	at.Walk(bst.LevelOrder, withprintheight())
	fmt.Println()
	at2 = bst.New(bst.AVL)
	at2.Insert(5, nil)
	at2.Insert(1, nil)
	at2.Insert(10, nil)
	at2.Insert(11, nil)
	at2.Print()
	at2.Remove(1)
	at2.Print()
	at2.Walk(bst.LevelOrder, withprintheight())
	fmt.Println()
	fmt.Println("AVL END")
}

func TestSplay(t *testing.T) {
	fmt.Println("SPLAY START")
	st := bst.New(bst.Splay)
	for i := 0; i < 10; i++ {
		st.Insert(i, nil)
	}
	st.Print()
	fmt.Println(st.Search(0))
	st.Print()
	fmt.Println(st.Search(2))
	st.Print()
	fmt.Println(st.Search(8))
	st.Print()
	fmt.Println(st.Remove(5))
	st.Print()
	fmt.Println(st.Search(5))
	st.Print()
	fmt.Println("SPLAY END")
}

func TestBTree(t *testing.T) {
	fmt.Println("BTREE START")
	bt := bst.New(bst.BTree, 5)
	for i := 0; i < 20; i++ {
		bt.Insert(i, nil)
	}
	bt.Print()
	bt.Remove(5)
	bt.Print()
	fmt.Println("BTREE END")
}

func TestRBTree(t *testing.T) {
	fmt.Println("RBTREE START")
	rb := bst.New(bst.RBTree)
	for i := 0; i < 20; i++ {
		rb.Insert(i, nil)
	}
	rb.Print()
	rb.Walk(bst.LevelOrder, withprintheight())
	rb.Remove(8)
	rb.Print()
	rb.Walk(bst.LevelOrder, withprintheight())
	rb.Remove(11)
	rb.Print()
	rb.Walk(bst.LevelOrder, withprintheight())
	rb.Remove(13)
	rb.Print()
	rb.Walk(bst.LevelOrder, withprintheight())
	rb.Remove(7)
	rb.Print()
	rb.Walk(bst.LevelOrder, withprintheight())
	fmt.Println()
	fmt.Println("RBTREE END")
}

func withprintheight() func(n bst.Node) {
	return func(n bst.Node) {
		fmt.Printf("%d ", n.Height())
	}
}

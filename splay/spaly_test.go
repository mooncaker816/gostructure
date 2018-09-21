package splay

import (
	"fmt"
	"testing"
)

func TestSplay(t *testing.T) {
	st := NewSplayTree()
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
}

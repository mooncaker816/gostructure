package bst

import (
	"strconv"
)

// BST identifies a Binary Search Tree interface.
type BST interface {
	Search(key interface{}) (Node, bool)
	Insert(key, data interface{}) (Node, error)
	Remove(key interface{}) (Node, error)
	Root() Node
	Print()
	Walk(order Order, opts ...Option)
}
type Order uint8

const (
	PreOrder = iota
	InOrder
	PostOrder
	LevelOrder
)

type Option func(n Node)

// Class stands for the specific type of binary search tree, such as AVL,Red-Black Tree etc.

type Class uint8

const (
	AVL Class = iota + 1
	RBTree
	Splay
	BTree
	maxClass
)

var classes = make([]func(parms ...interface{}) BST, maxClass)

// New returns a new BST as per the provided class.
func New(c Class, parms ...interface{}) BST {
	// fmt.Println(parms...)
	if c > 0 && c < maxClass {
		f := classes[c]
		if f != nil {
			return f(parms...)
		}
	}
	panic("bst: requested BST function #" + strconv.Itoa(int(c)) + " is unavailable")
}

// // Available reports whether the given BST class is linked into the binary.
// func (c Class) Available() bool {
// 	return c < maxClass && classes[c] != nil
// }

// RegisterBST registers a function that returns a new instance of the given
// BST class. This is intended to be called from the init function in
// packages that implement BST.
func RegisterBST(c Class, f func(parms ...interface{}) BST) {
	if c >= maxClass {
		panic("bst: RegisterBST of unknown BST class")
	}
	classes[c] = f
}

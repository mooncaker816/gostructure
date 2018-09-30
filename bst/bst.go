package bst

import "strconv"

// BST identifies a Binary Search Tree interface.
type BST interface {
	Search(key interface{}) (*Node, bool)
	Insert(key, data interface{}) (*Node, error)
	Remove(key interface{}) (*Node, error)
	Root() *Node
	Print()
	TravLevel(opts ...Option)
	TravPre(opts ...Option)
	TravIn(opts ...Option)
	TravPost(opts ...Option)
}

// Class stands for the specific type of binary search tree, such as AVL,Red-Black Tree etc.
type Class uint

const (
	AVL Class = iota + 1
	BTree
	RBTree
	Splay
	maxClass
)

var classes = make([]func() BST, maxClass)

// New returns a new BST as per the provided class.
func New(c Class) BST {
	if c > 0 && c < maxClass {
		f := classes[c]
		if f != nil {
			return f()
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
func RegisterBST(c Class, f func() BST) {
	if c >= maxClass {
		panic("bst: RegisterBST of unknown BST class")
	}
	classes[c] = f
}

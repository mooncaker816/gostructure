package avl

import (
	"errors"

	"github.com/mooncaker816/gostructure/bst"
)

func init() {
	bst.RegisterBST(bst.AVL, New)
}

// AVL tree
type avl struct {
	root *node
	comp bst.Comparator
}

func avlOK(n *node) bool {
	lh, rh := -1, -1
	if bst.HasLChild(n) {
		lh = n.lchild.height
	}
	if bst.HasRChild(n) {
		rh = n.rchild.height
	}
	diff := lh - rh
	if diff > 1 || diff < -1 {
		return false
	}
	return true
}

// New returns a new empty avl tree with basic comparator
func New(parms ...interface{}) bst.BST {
	t := new(avl)
	t.comp = bst.BasicCompare
	for _, p := range parms {
		switch v := p.(type) {
		case bst.Comparator:
			t.comp = v
		}
	}
	return t
}

func (avl *avl) Root() bst.Node {
	return avl.root
}

func (avl *avl) Print() {
	bst.PrintWithUnitSize(avl.root, 2)
}

func (avl *avl) Search(key interface{}) (bst.Node, bool) {
	n, result := avl.searchIn(avl.root, key)
	if result == 0 {
		return n, true
	}
	return n, false
}

func (avl *avl) searchIn(n *node, key interface{}) (*node, int) {
	switch avl.comp(key, n.key) {
	case 0:
		return n, 0
	case -1:
		if bst.HasLChild(n) {
			return avl.searchIn(n.lchild, key)
		}
		return n, -1
	case 1:
		if bst.HasRChild(n) {
			return avl.searchIn(n.rchild, key)
		}
		return n, 1
	default:
		panic("compare result illegal")
	}
}

func (avl *avl) Insert(key, data interface{}) (bst.Node, error) {
	if avl.root == nil {
		avl.root = newNode(key, data)
		return avl.root, nil
	}
	n, result := avl.searchIn(avl.root, key)
	switch result {
	case 0:
		return nil, errors.New("insert node with duplicate key")
	case -1:
		new := newNode(key, data)
		bst.AttachLChild(n, new)
		if !bst.HasRChild(n) {
			n.updateHeightAbove()
		}
		avl.reBalance(n, true)
		return new, nil
	case 1:
		new := newNode(key, data)
		bst.AttachRChild(n, new)
		if !bst.HasLChild(n) {
			n.updateHeightAbove()
		}
		avl.reBalance(n, true)
		return new, nil
	}
	return nil, nil
}

func (avl *avl) reBalance(hot *node, insert bool) {
	for g := hot; g != nil; g = g.parent {
		if !avlOK(g) {
			x := g.parent
			p := g.tallerChild()
			v := p.tallerChild()
			var tmp *node
			if bst.IsLChild(g) {
				x.lchild = rotateAndUpdateHeight(v)
				tmp = x.lchild
			} else if bst.IsRChild(g) {
				x.rchild = rotateAndUpdateHeight(v)
				tmp = x.rchild
			} else {
				avl.root = rotateAndUpdateHeight(v)
				tmp = avl.root
			}
			if insert {
				break
			} else {
				g = tmp
			}
		} else {
			updateHeight(g)
		}
	}
}

func rotateAndUpdateHeight(n *node) *node {
	_, b, _ := bst.RotateAt(n, updateHeight)
	return b.(*node)
}

func (avl *avl) Remove(key interface{}) (bst.Node, error) {
	n, result := avl.searchIn(avl.root, key)
	if result != 0 {
		return nil, nil
	}
	hot, _ := bst.RemoveAt(n, avl.root)
	hot0, ok := hot.(*node)
	if ok {
		hot0.updateHeightAbove()
		avl.reBalance(hot0, false)
	}
	return hot0, nil
}

func (avl *avl) Walk(o bst.Order, opts ...bst.Option) {
	switch o {
	case bst.PreOrder:
		bst.TravPre(avl.root, opts...)
	case bst.InOrder:
		bst.TravIn(avl.root, opts...)
	case bst.PostOrder:
		bst.TravPost(avl.root, opts...)
	case bst.LevelOrder:
		bst.TravLevel(avl.root, opts...)
	default:
		panic("unsupported walk order")
	}
}

package redblack

import (
	"errors"

	"github.com/mooncaker816/gostructure/bst"
)

func init() {
	bst.RegisterBST(bst.RBTree, New)
}

type rbTree struct {
	root *node
	comp bst.Comparator
}

// New returns an empty redblack tree
func New(parms ...interface{}) bst.BST {
	t := new(rbTree)
	t.comp = bst.BasicCompare
	for _, p := range parms {
		switch v := p.(type) {
		case bst.Comparator:
			t.comp = v
		}
	}
	return t
}

func (rb *rbTree) Root() bst.Node { return rb.root }

func rbOK(n *node) bool {
	lh, rh := -1, -1
	if bst.HasLChild(n) {
		lh = n.lchild.height
	}
	if bst.HasRChild(n) {
		rh = n.rchild.height
	}
	if lh != rh {
		return false
	}
	if n.isRed() {
		return n.height == lh
	}
	return n.height == lh+1
}

func (rb *rbTree) Search(key interface{}) (bst.Node, bool) {
	n, result := rb.searchIn(rb.root, key)
	if result == 0 {
		return n, true
	}
	return n, false
}

func (rb *rbTree) searchIn(n *node, key interface{}) (*node, int) {
	switch rb.comp(key, n.key) {
	case 0:
		return n, 0
	case -1:
		if bst.HasLChild(n) {
			return rb.searchIn(n.lchild, key)
		}
		return n, -1
	case 1:
		if bst.HasRChild(n) {
			return rb.searchIn(n.rchild, key)
		}
		return n, 1
	default:
		panic("compare result illegal")
	}
}

func (rb *rbTree) Insert(key, data interface{}) (bst.Node, error) {
	if rb.root == nil {
		rb.root = newNode(key, data)
		rb.root.setBlack()
		return rb.root, nil
	}
	n, result := rb.searchIn(rb.root, key)
	switch result {
	case 0:
		return nil, errors.New("insert node with duplicate key")
	case -1:
		new := newNode(key, data)
		bst.AttachLChild(n, new)
		rb.solveDoubleRed(new)
		return new, nil
	case 1:
		new := newNode(key, data)
		bst.AttachRChild(n, new)
		rb.solveDoubleRed(new)
		return new, nil
	}
	return nil, nil
}

func (rb *rbTree) solveDoubleRed(n *node) {
	if bst.IsRoot(n) {
		n.setBlack()
		n.height++
		return
	}
	p := n.parent
	if p.isBlack() {
		return
	}
	g := p.parent
	u := bst.Sibling(p)
	u0 := u.(*node)
	if u0.isBlack() { // RR-1
		x := g.parent
		if bst.IsLChild(g) {
			x.lchild = rr1(n)
		} else if bst.IsRChild(g) {
			x.rchild = rr1(n)
		} else {
			rb.root = rr1(n)
		}
	} else { // RR-2
		p.setBlack()
		p.height++
		u0.setBlack()
		u0.height++
		if !bst.IsRoot(g) {
			g.setRed()
		}
		rb.solveDoubleRed(g)
	}
}

// roate + change color + update height for RR-1
func rr1(n *node) *node {
	a, b, c := bst.RotateAt(n)
	a.(*node).setRed()
	c.(*node).setRed()
	b.(*node).setBlack()
	a.(*node).updateHeight()
	b.(*node).updateHeight()
	c.(*node).updateHeight()
	return b.(*node)
}

func (rb *rbTree) Remove(key interface{}) (bst.Node, error) {
	n, result := rb.searchIn(rb.root, key)
	if result != 0 {
		return nil, nil
	}
	hot, r := bst.RemoveAt(n, rb.root)
	// fmt.Println(hot, r)
	hot0 := hot.(*node)
	if rb.root == nil {
		return nil, nil
	}
	if hot0 == nil {
		rb.root.setBlack()
		rb.root.updateHeight()
		return hot0, nil
	}
	if rbOK(hot0) {
		return hot0, nil
	}
	r0, ok := r.(*node)
	// if ok {
	if ok && r0.isRed() {
		r0.setBlack()
		r0.height++
		return hot0, nil
	}
	// }
	rb.solveDoubleBlack(r0, hot0)
	return hot0, nil
}

func (rb *rbTree) solveDoubleBlack(r, hot *node) {
	var p, s *node
	if r != nil {
		p = r.parent
	} else {
		p = hot
	}
	if p == nil {
		return
	}

	if r == p.lchild {
		s = p.rchild
	} else {
		s = p.lchild
	}
	if s.isBlack() {
		var t *node
		if s.rchild.isRed() {
			t = s.rchild
		}
		if s.lchild.isRed() {
			t = s.lchild
		}
		if t != nil { // BB-1 s至少有一个红孩子
			oldattr := p.attr
			x := p.parent
			if bst.IsLChild(p) {
				x.lchild = bb1(t, oldattr)
			} else if bst.IsRChild(p) {
				x.rchild = bb1(t, oldattr)
			} else {
				rb.root = bb1(t, oldattr)
			}
		} else {
			s.setRed()
			s.height--
			if p.isRed() { // BB-2R
				p.setBlack()
			} else { // BB-2B
				p.height--
				rb.solveDoubleBlack(p, hot)
			}
		}
	} else { // BB-3
		s.setBlack()
		p.setRed()
		hot = p
		var t *node
		if bst.IsLChild(s) {
			t = s.lchild
		} else {
			t = s.rchild
		}
		x := p.parent
		if bst.IsLChild(p) {
			x.lchild = bb3(t)
		} else if bst.IsRChild(p) {
			x.rchild = bb3(t)
		} else {
			rb.root = bb3(t)
		}
		rb.solveDoubleBlack(r, hot)
	}
}

func bb1(n *node, oldattr uint8) *node {
	a, b, c := bst.RotateAt(n)
	a0, b0, c0 := a.(*node), b.(*node), c.(*node)
	a0.updateHeight()
	b0.updateHeight()
	c0.updateHeight()

	b0.attr = b0.attr&0xfe | oldattr&1
	b0.updateHeight()

	if bst.HasLChild(b0) {
		b0.lchild.setBlack()
		b0.lchild.updateHeight()
	}
	if bst.HasRChild(b0) {
		b0.rchild.setBlack()
		b0.rchild.updateHeight()
	}
	return b0
}

func bb3(n *node) *node {
	a, b, c := bst.RotateAt(n)
	a0, b0, c0 := a.(*node), b.(*node), c.(*node)
	a0.updateHeight()
	b0.updateHeight()
	c0.updateHeight()
	return b0
}

func (rb *rbTree) Print() {
	bst.PrintWithUnitSize(rb.root, 2)
}

func (rb *rbTree) Walk(o bst.Order, opts ...bst.Option) {
	switch o {
	case bst.PreOrder:
		bst.TravPre(rb.root, opts...)
	case bst.InOrder:
		bst.TravIn(rb.root, opts...)
	case bst.PostOrder:
		bst.TravPost(rb.root, opts...)
	case bst.LevelOrder:
		bst.TravLevel(rb.root, opts...)
	default:
		panic("unsupported walk order")
	}
}

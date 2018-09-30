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
	root *bst.Node
	comp bst.Comparator
}

func avlOK(n *bst.Node) bool {
	lh, rh := -1, -1
	if n.HasLChild() {
		lh = n.LChild.Height
	}
	if n.HasRChild() {
		rh = n.RChild.Height
	}
	diff := lh - rh
	if diff > 1 || diff < -1 {
		return false
	}
	return true
}

// New returns a new empty avl tree with basic comparator
func New() bst.BST {
	return &avl{comp: bst.BasicCompare}
}

// NewWithComparator returns a new empty avl tree with provided comparator
func NewWithComparator(comp bst.Comparator) bst.BST {
	return &avl{comp: comp}
}

func (avl *avl) Root() *bst.Node {
	return avl.root
}

func (avl *avl) Print() {
	avl.root.PrintWithUnitSize(2)
}

func (avl *avl) Search(key interface{}) (*bst.Node, bool) {
	n, result := avl.searchIn(avl.root, key)
	if result == 0 {
		return n, true
	}
	return n, false
}

func (avl *avl) searchIn(n *bst.Node, key interface{}) (*bst.Node, int) {
	switch avl.comp(key, n.Key) {
	case 0:
		return n, 0
	case -1:
		if n.HasLChild() {
			return avl.searchIn(n.LChild, key)
		}
		return n, -1
	case 1:
		if n.HasRChild() {
			return avl.searchIn(n.RChild, key)
		}
		return n, 1
	default:
		panic("compare result illegal")
	}
}

func (avl *avl) Insert(key, data interface{}) (*bst.Node, error) {
	if avl.root == nil {
		avl.root = &bst.Node{Key: key, Data: data}
		return avl.root, nil
	}
	n, result := avl.searchIn(avl.root, key)
	switch result {
	case 0:
		return nil, errors.New("insert node with duplicate key")
	case -1:
		new := &bst.Node{Key: key, Data: data}
		n.AttachLChild(new)
		if !n.HasRChild() {
			n.UpdateHeightAbove()
		}
		avl.reBalance(n, true)
		return new, nil
	case 1:
		new := &bst.Node{Key: key, Data: data}
		n.AttachRChild(new)
		if !n.HasLChild() {
			n.UpdateHeightAbove()
		}
		avl.reBalance(n, true)
		return new, nil
	}
	return nil, nil
}

func (avl *avl) reBalance(hot *bst.Node, insert bool) {
	for g := hot; g != nil; g = g.Parent {
		if !avlOK(g) {
			x := g.Parent
			p := g.TallerChild()
			v := p.TallerChild()
			tmp := new(bst.Node)
			if g.IsLChild() {
				x.LChild = rotateAt(v)
				tmp = x.LChild
			} else if g.IsRChild() {
				x.RChild = rotateAt(v)
				tmp = x.RChild
			} else {
				avl.root = rotateAt(v)
				tmp = avl.root
			}
			if insert {
				break
			} else {
				g = tmp
			}
		} else {
			g.UpdateHeight()
		}
	}
}

// connect34 connect bst.Nodes as below
//	 	   b
//		a	  c
//	  T0 T1 T2 T3
func connect34(a, b, c, t1, t2, t3, t4 *bst.Node) *bst.Node {
	a.LChild = t1
	if t1 != nil {
		t1.Parent = a
	}
	a.RChild = t2
	if t2 != nil {
		t2.Parent = a
	}
	a.UpdateHeight()
	c.LChild = t3
	if t3 != nil {
		t3.Parent = c
	}
	c.RChild = t4
	if t4 != nil {
		t4.Parent = c
	}
	c.UpdateHeight()
	b.LChild = a
	b.RChild = c
	a.Parent = b
	c.Parent = b
	b.UpdateHeight()
	return b
}

func rotateAt(v *bst.Node) *bst.Node {
	p := v.Parent
	g := p.Parent
	if v.IsLChild() {
		if p.IsLChild() {
			p.Parent = g.Parent
			return connect34(v, p, g, v.LChild, v.RChild, p.RChild, g.RChild)
		}
		if p.IsRChild() {
			v.Parent = g.Parent
			return connect34(g, v, p, g.LChild, v.LChild, v.RChild, p.RChild)
		}
	}
	if v.IsRChild() {
		if p.IsLChild() {
			v.Parent = g.Parent
			return connect34(p, v, g, p.LChild, v.LChild, v.RChild, g.RChild)
		}
		if p.IsRChild() {
			p.Parent = g.Parent
			return connect34(g, p, v, g.LChild, p.LChild, v.LChild, v.RChild)
		}
	}
	return nil
}

func (avl *avl) Remove(key interface{}) (hot *bst.Node, err error) {
	n, result := avl.searchIn(avl.root, key)
	if result != 0 {
		return nil, nil
	}
	hot, _ = avl.removeAt(n)
	avl.reBalance(hot, false)
	return hot, nil
}

func (avl *avl) removeAt(n *bst.Node) (hot *bst.Node, err error) {
	// n has both left and right subtree
	if n.HasLChild() && n.HasRChild() {
		succ := n.Successor()
		bst.SwapKeyData(n, succ)
		hot = succ.Parent
		if succ == n.RChild {
			hot = n
			if succ.HasRChild() {
				succ.RChild.Parent = n
			}
			n.RChild = succ.RChild
		} else if succ.HasRChild() {
			// hot = succ.Parent
			succ.RChild.Parent = hot
			hot.LChild = succ.RChild
		} else {
			// hot = succ.Parent
			hot.LChild = nil
		}
		succ.Parent, succ.LChild, succ.RChild = nil, nil, nil
		hot.UpdateHeightAbove()
		return hot, nil
	}
	// n only has left subtree
	if n.HasLChild() && !n.HasRChild() {
		if n.IsLChild() {
			n.Parent.LChild = n.LChild
		} else if n.IsRChild() {
			n.Parent.RChild = n.LChild
		} else {
			avl.root = n.LChild
		}
		n.LChild.Parent = n.Parent
		hot = n.Parent
		n.Parent, n.LChild, n.RChild = nil, nil, nil
		hot.UpdateHeightAbove()
		return hot, nil
	}
	// n only has right subtree
	if !n.HasLChild() && n.HasRChild() {
		if n.IsLChild() {
			n.Parent.LChild = n.RChild
		} else if n.IsRChild() {
			n.Parent.RChild = n.RChild
		} else {
			avl.root = n.RChild
		}
		n.RChild.Parent = n.Parent
		hot = n.Parent
		n.Parent, n.LChild, n.RChild = nil, nil, nil
		hot.UpdateHeightAbove()
		return hot, nil
	}
	// n is leaf(or single root)
	if n.IsLChild() {
		n.Parent.LChild = nil
	} else if n.IsRChild() {
		n.Parent.RChild = nil
	} else {
		avl.root = nil
	}
	hot = n.Parent
	n.Parent, n.LChild, n.RChild = nil, nil, nil
	hot.UpdateHeightAbove()
	return hot, nil
}

func (avl *avl) TravLevel(opts ...bst.Option) {
	avl.root.TravLevel(opts...)
}

func (avl *avl) TravPre(opts ...bst.Option) {
	avl.root.TravPre(opts...)
}

func (avl *avl) TravIn(opts ...bst.Option) {
	avl.root.TravIn(opts...)
}

func (avl *avl) TravPost(opts ...bst.Option) {
	avl.root.TravPost(opts...)
}

package splay

import (
	"errors"

	"github.com/mooncaker816/gostructure/bst"
)

func init() {
	bst.RegisterBST(bst.Splay, New)
}

type splayTree struct {
	root *node
	comp bst.Comparator
}

// New returns a new empty splay tree with basic comparator
func New(parms ...interface{}) bst.BST {
	t := new(splayTree)
	t.comp = bst.BasicCompare
	for _, p := range parms {
		switch v := p.(type) {
		case bst.Comparator:
			t.comp = v
		}
	}
	return t
}

func (s *splayTree) Search(key interface{}) (bst.Node, bool) {
	n, result := s.searchIn(s.root, key)
	if result == 0 {
		return n, true
	}
	return n, false
}

func (s *splayTree) searchIn(n *node, key interface{}) (*node, int) {
	switch s.comp(key, n.key) {
	case 0:
		s.root = splay(n)
		return s.root, 0
	case -1:
		if n.lchild != nil {
			return s.searchIn(n.lchild, key)
		}
		s.root = splay(n)
		return s.root, -1
	case 1:
		if n.rchild != nil {
			return s.searchIn(n.rchild, key)
		}
		s.root = splay(n)
		return s.root, 1
	}
	return nil, 0
}

func splay(n *node) *node {
	if n == nil {
		return nil
	}
	for n.parent != nil && n.parent.parent != nil {
		p := n.parent
		g := n.parent.parent
		gp := g.parent
		if bst.IsLChild(n) {
			if bst.IsLChild(p) {
				bst.AttachLChild(g, p.rchild)
				bst.AttachLChild(p, n.rchild)
				bst.AttachRChild(p, g)
				bst.AttachRChild(n, p)
			} else {
				bst.AttachLChild(p, n.rchild)
				bst.AttachRChild(g, n.lchild)
				bst.AttachLChild(n, g)
				bst.AttachRChild(n, p)
			}
		} else {
			if bst.IsLChild(p) {
				bst.AttachRChild(p, n.lchild)
				bst.AttachLChild(g, n.rchild)
				bst.AttachRChild(n, g)
				bst.AttachLChild(n, p)
			} else {
				bst.AttachRChild(g, p.lchild)
				bst.AttachRChild(p, n.lchild)
				bst.AttachLChild(p, g)
				bst.AttachLChild(n, p)
			}
		}
		if gp != nil {
			if gp.lchild == g {
				bst.AttachLChild(gp, n)
			} else {
				bst.AttachRChild(gp, n)
			}
		} else {
			n.parent = nil
		}
	}
	if n.parent != nil {
		if bst.IsLChild(n) {
			bst.AttachLChild(n.parent, n.rchild)
			bst.AttachRChild(n, n.parent)
		} else {
			bst.AttachRChild(n.parent, n.lchild)
			bst.AttachLChild(n, n.parent)
		}
	}
	n.parent = nil
	return n
}

func (s *splayTree) Insert(key, data interface{}) (bst.Node, error) {
	if s.root == nil {
		s.root = newNode(key, data)
		return s.root, nil
	}
	n, result := s.searchIn(s.root, key)
	switch result {
	case 0:
		return nil, errors.New("insert node with duplicate key")
	case -1:
		new := newNode(key, data)
		bst.AttachRChild(new, n)
		bst.AttachLChild(new, n.lchild)
		n.lchild = nil
		s.root = new
		return new, nil
	case 1:
		new := newNode(key, data)
		bst.AttachLChild(new, n)
		bst.AttachRChild(new, n.rchild)
		n.rchild = nil
		s.root = new
		return new, nil
	}
	return nil, nil
}

func (s *splayTree) Remove(key interface{}) (bst.Node, error) {
	n, result := s.searchIn(s.root, key)
	if result != 0 {
		return nil, nil
	}
	// 此时待删节点位于 root
	if !bst.HasLChild(s.root) {
		s.root = s.root.rchild
		if s.root != nil {
			s.root.parent = nil
		}
		n.rchild = nil
	} else if !bst.HasRChild(s.root) {
		s.root = s.root.lchild
		s.root.parent = nil
		n.lchild = nil
	} else {
		// 临时切除左子树
		lc := s.root.lchild
		s.root.lchild = nil
		// 删除 root，右子树成为新的树
		s.root = s.root.rchild
		s.root.parent = nil
		n.rchild = nil
		// 在新树中再次查找原来的值，必然不存在，但会把最小值提升至顶端，且没有左子树
		// 该最小值一定比之前切除的左子树大，以此值为 root 重新连接原左子树即可
		s.searchIn(s.root, key)
		s.root.lchild = lc
		lc.parent = s.root
	}
	return n, nil
}

func (s *splayTree) Root() bst.Node {
	return s.root
}

func (s *splayTree) Print() {
	bst.PrintWithUnitSize(s.root, 2)
}

func (s *splayTree) Walk(o bst.Order, opts ...bst.Option) {
	switch o {
	case bst.PreOrder:
		bst.TravPre(s.root, opts...)
	case bst.InOrder:
		bst.TravIn(s.root, opts...)
	case bst.PostOrder:
		bst.TravPost(s.root, opts...)
	case bst.LevelOrder:
		bst.TravLevel(s.root, opts...)
	default:
		panic("unsupported walk order")
	}
}

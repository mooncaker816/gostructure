package splay

import (
	"errors"

	"github.com/mooncaker816/gostructure/bst"
)

func init() {
	bst.RegisterBST(bst.Splay, New)
}

type splayTree struct {
	root *bst.Node
	comp bst.Comparator
}

// New returns a new empty splay tree with basic comparator
func New() bst.BST {
	return &splayTree{comp: bst.BasicCompare}
}

func (s *splayTree) Search(key interface{}) (*bst.Node, bool) {
	n, result := s.searchIn(s.root, key)
	if result == 0 {
		return n, true
	}
	return n, false
}

func (s *splayTree) searchIn(n *bst.Node, key interface{}) (*bst.Node, int) {
	switch s.comp(key, n.Key) {
	case 0:
		s.root = splay(n)
		return s.root, 0
	case -1:
		if n.LChild != nil {
			return s.searchIn(n.LChild, key)
		}
		s.root = splay(n)
		return s.root, -1
	case 1:
		if n.RChild != nil {
			return s.searchIn(n.RChild, key)
		}
		s.root = splay(n)
		return s.root, 1
	}
	return nil, 0
}

func splay(n *bst.Node) *bst.Node {
	if n == nil {
		return nil
	}
	for n.Parent != nil && n.Parent.Parent != nil {
		p := n.Parent
		g := n.Parent.Parent
		gp := g.Parent
		if n.IsLChild() {
			if p.IsLChild() {
				g.AttachLChild(p.RChild)
				p.AttachLChild(n.RChild)
				p.AttachRChild(g)
				n.AttachRChild(p)
			} else {
				p.AttachLChild(n.RChild)
				g.AttachRChild(n.LChild)
				n.AttachLChild(g)
				n.AttachRChild(p)
			}
		} else {
			if p.IsLChild() {
				p.AttachRChild(n.LChild)
				g.AttachLChild(n.RChild)
				n.AttachRChild(g)
				n.AttachLChild(p)
			} else {
				g.AttachRChild(p.LChild)
				p.AttachRChild(n.LChild)
				p.AttachLChild(g)
				n.AttachLChild(p)
			}
		}
		if gp != nil {
			if gp.LChild == g {
				gp.AttachLChild(n)
			} else {
				gp.AttachRChild(n)
			}
		} else {
			n.Parent = nil
		}
	}
	if n.Parent != nil {
		if n.IsLChild() {
			n.Parent.AttachLChild(n.RChild)
			n.AttachRChild(n.Parent)
		} else {
			n.Parent.AttachRChild(n.LChild)
			n.AttachLChild(n.Parent)
		}
	}
	n.Parent = nil
	return n
}

func (s *splayTree) Insert(key, data interface{}) (*bst.Node, error) {
	if s.root == nil {
		s.root = &bst.Node{Key: key, Data: data}
		return s.root, nil
	}
	n, result := s.searchIn(s.root, key)
	switch result {
	case 0:
		return nil, errors.New("insert node with duplicate key")
	case -1:
		new := &bst.Node{Key: key, Data: data}
		new.AttachRChild(n)
		new.AttachLChild(n.LChild)
		n.LChild = nil
		s.root = new
		return new, nil
	case 1:
		new := &bst.Node{Key: key, Data: data}
		new.AttachLChild(n)
		new.AttachRChild(n.RChild)
		n.RChild = nil
		s.root = new
		return new, nil
	}
	return nil, nil
}

func (s *splayTree) Remove(key interface{}) (*bst.Node, error) {
	n, result := s.searchIn(s.root, key)
	if result != 0 {
		return nil, nil
	}
	// 此时待删节点位于 root
	if !s.root.HasLChild() {
		s.root = s.root.RChild
		if s.root != nil {
			s.root.Parent = nil
		}
		n.RChild = nil
	} else if !s.root.HasRChild() {
		s.root = s.root.LChild
		s.root.Parent = nil
		n.LChild = nil
	} else {
		// 临时切除左子树
		lc := s.root.LChild
		s.root.LChild = nil
		// 删除 root，右子树成为新的树
		s.root = s.root.RChild
		s.root.Parent = nil
		n.RChild = nil
		// 在新树中再次查找原来的值，必然不存在，但会把最小值提升至顶端，且没有左子树
		// 该最小值一定比之前切除的左子树大，以此值为 root 重新连接原左子树即可
		s.searchIn(s.root, key)
		s.root.LChild = lc
		lc.Parent = s.root
	}
	return n, nil
}

func (s *splayTree) Root() *bst.Node {
	return s.root
}

func (s *splayTree) Print() {
	s.root.PrintWithUnitSize(2)
}

func (s *splayTree) TravLevel(opts ...bst.Option) {
	s.root.TravLevel(opts...)
}

func (s *splayTree) TravPre(opts ...bst.Option) {
	s.root.TravPre(opts...)
}

func (s *splayTree) TravIn(opts ...bst.Option) {
	s.root.TravIn(opts...)
}

func (s *splayTree) TravPost(opts ...bst.Option) {
	s.root.TravPost(opts...)
}

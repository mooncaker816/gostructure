package splay

import (
	"github.com/mooncaker816/gostructure/bst"
)

type node struct {
	lchild *node
	rchild *node
	parent *node
	key    interface{}
	data   interface{}
}

func newNode(key, data interface{}) *node {
	return &node{key: key, data: data}
}

func (n *node) Key() interface{}  { return n.key }
func (n *node) Data() interface{} { return n.data }
func (n *node) Height() int       { return 0 }
func (n *node) LChild() bst.Node  { return n.lchild }
func (n *node) RChild() bst.Node  { return n.rchild }
func (n *node) Parent() bst.Node  { return n.parent }
func (n *node) Color() string     { return "" }

func (n *node) SetKey(key interface{})   { n.key = key }
func (n *node) SetData(data interface{}) { n.data = data }
func (n *node) SetLChild(lc bst.Node) {
	if lc == nil {
		n.lchild = nil
		return
	}
	lc0, ok := lc.(*node)
	if !ok {
		panic("inconsistent node type")
	}
	if n != nil {
		n.lchild = lc0
	}
}

func (n *node) SetRChild(rc bst.Node) {
	if rc == nil {
		n.rchild = nil
		return
	}
	rc0, ok := rc.(*node)
	if !ok {
		panic("inconsistent node type")
	}
	if n != nil {
		n.rchild = rc0
	}
}

func (n *node) SetParent(p bst.Node) {
	if p == nil {
		n.parent = nil
		return
	}
	p0, ok := p.(*node)
	if !ok {
		panic("inconsistent node type")
	}
	if n != nil {
		n.parent = p0
	}
}

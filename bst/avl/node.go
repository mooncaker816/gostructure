package avl

import (
	"github.com/mooncaker816/gostructure/bst"
)

type node struct {
	lchild *node
	rchild *node
	parent *node
	key    interface{}
	data   interface{}
	height int
}

func newNode(key, data interface{}) *node {
	return &node{key: key, data: data}
}

func (n *node) Key() interface{}         { return n.key }
func (n *node) SetKey(key interface{})   { n.key = key }
func (n *node) Data() interface{}        { return n.data }
func (n *node) SetData(data interface{}) { n.data = data }
func (n *node) Height() int              { return n.height }
func (n *node) LChild() bst.Node         { return n.lchild }
func (n *node) RChild() bst.Node         { return n.rchild }
func (n *node) Parent() bst.Node         { return n.parent }
func (n *node) Color() string            { return "" }

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

func updateHeight(n bst.Node) {
	n0 := n.(*node)
	n0.height = n0.maxHeightOfChildren() + 1
}

func (n *node) updateHeightAbove() {
	for max := n.maxHeightOfChildren(); n != nil && n.height != max+1; {
		n.height = max + 1
		n = n.parent
	}
}

func (n *node) maxHeightOfChildren() int {
	lH, rH := -1, -1
	if bst.HasLChild(n) {
		lH = n.lchild.height
	}
	if bst.HasRChild(n) {
		rH = n.rchild.height
	}
	return max(lH, rH)
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// Tallerchild 返回高度较高的那个孩子节点，若同高，返回和n同侧的节点
func (n *node) tallerChild() *node {
	if bst.HasLChild(n) && n.rchild == nil {
		return n.lchild
	}
	if bst.HasRChild(n) && n.lchild == nil {
		return n.rchild
	}

	if bst.HasLChild(n) && bst.HasRChild(n) {
		if n.lchild.height < n.rchild.height {
			return n.rchild
		}
		if n.lchild.height > n.rchild.height {
			return n.lchild
		}
		if bst.IsLChild(n) {
			return n.lchild
		}
		return n.rchild
	}
	return nil
}

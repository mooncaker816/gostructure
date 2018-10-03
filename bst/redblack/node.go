package redblack

import (
	"github.com/mooncaker816/gostructure/bst"
)

type node struct {
	lchild *node
	rchild *node
	parent *node
	key    interface{}
	data   interface{}
	height int // exact black height -1
	attr   uint8
}

func newNode(key, data interface{}) *node {
	return &node{key: key, data: data, height: -1} //默认红色，初始高度为实际黑高度-1
}

func (n *node) Key() interface{}  { return n.key }
func (n *node) Data() interface{} { return n.data }
func (n *node) Height() int       { return n.height + 1 }
func (n *node) LChild() bst.Node  { return n.lchild }
func (n *node) RChild() bst.Node  { return n.rchild }
func (n *node) Parent() bst.Node  { return n.parent }

func (n *node) Color() string {
	if n.isBlack() {
		return "B"
	}
	return "R"
}

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

// customized
func (n *node) updateHeight() {
	n.height = n.maxHeightOfChildren()
	if n.isBlack() {
		n.height++
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

const (
	Red = iota
	Black
)

func (n *node) isBlack() bool {
	return n == nil || n.attr&1 == 1 // 外部节点视为黑
}

func (n *node) isRed() bool {
	return !n.isBlack()
}

func (n *node) setRed() {
	n.attr &= 0xfe
}

func (n *node) setBlack() {
	n.attr |= 1
}

func (n *node) flipColor() {
	n.attr ^= 1
}

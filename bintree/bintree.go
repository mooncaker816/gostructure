package bintree

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mooncaker816/gostructure/queue"
	"github.com/mooncaker816/gostructure/stack"
)

// Node is an element in a binary tree
type Node struct {
	Parent *Node
	LChild *Node
	RChild *Node
	Data   interface{}
	Key    interface{}
	Height int
	Tree   *BinTree
}

// BinTree is common binary tree
type BinTree struct {
	Root *Node
	Size int
}

// IsRoot checks n is root or not
func (n *Node) IsRoot() bool { return n.Parent == nil }

// IsLChild checks n is left child or not
func (n *Node) IsLChild() bool { return !n.IsRoot() && n.Parent.LChild == n }

// IsRChild checks n is right child or not
func (n *Node) IsRChild() bool { return !n.IsRoot() && n.Parent.RChild == n }

// HasParent checks whether n has parent or not
func (n *Node) HasParent() bool { return !n.IsRoot() }

// HasLChild checks whether n has a left child
func (n *Node) HasLChild() bool { return n.LChild != nil }

// HasRChild checks whether n has a right child
func (n *Node) HasRChild() bool { return n.RChild != nil }

// HasChild checks whether n has a child
func (n *Node) HasChild() bool { return n.HasLChild() || n.HasRChild() }

// HasBothChild checks whether n has a child
func (n *Node) HasBothChild() bool { return n.HasLChild() && n.HasRChild() }

// IsLeaf checks n is leaf or not
func (n *Node) IsLeaf() bool { return !n.HasChild() }

// Sibling returns n's brother if any
func (n *Node) Sibling() *Node {
	if n.IsRoot() {
		return nil
	}
	if n.IsLChild() {
		return n.Parent.RChild
	}
	return n.Parent.LChild
}

// Uncle returns n's uncle if any
func (n *Node) Uncle() *Node {
	if n.IsRoot() || n.Parent.IsRoot() {
		return nil
	}
	if n.Parent.IsLChild() {
		return n.Parent.Parent.RChild
	}
	return n.Parent.Parent.LChild
}

// Size 统计当前节点后代总数，即以其为根的子树规模
func (n *Node) Size() int {
	count := 1
	if n.HasLChild() {
		count += n.LChild.Size()
	}
	if n.HasRChild() {
		count += n.RChild.Size()
	}
	return count
}

// InsertAsLChild inserts a new node(key,data) as left child of n
func (n *Node) InsertAsLChild(key, data interface{}) *Node {
	if n.LChild != nil {
		panic("current node already has a left child")
	}
	n.LChild = &Node{Parent: n, Data: data, Key: key, Tree: n.Tree}
	return n.LChild
}

// InsertAsRChild inserts a new node(key,data) as right child of n
func (n *Node) InsertAsRChild(key, data interface{}) *Node {
	if n.RChild != nil {
		panic("current node already has a right child")
	}
	n.RChild = &Node{Parent: n, Data: data, Key: key, Tree: n.Tree}
	return n.RChild
}

func (n *Node) updateHeight() {
	var lh, rh int
	if n.LChild == nil {
		lh = -1
	} else {
		lh = n.LChild.Height
	}
	if n.RChild == nil {
		rh = -1
	} else {
		rh = n.RChild.Height
	}
	if lh >= rh {
		n.Height = lh + 1
	} else {
		n.Height = rh + 1
	}
}

func (n *Node) updateHeightAbove() {
	tmp := n
	for tmp != nil {
		tmp.updateHeight()
		tmp = tmp.Parent
	}
}

// TravPreR 二叉树先序遍历算法（递归版）
func (n *Node) TravPreR(opts ...Option) []*Node {
	nodes := make([]*Node, 0)
	if n == nil {
		return nil
	}
	for _, opt := range opts {
		opt(n)
	}
	nodes = append(nodes, n)
	nodes = append(nodes, n.LChild.TravPreR()...)
	nodes = append(nodes, n.RChild.TravPreR()...)
	return nodes
}

// TravPre1 二叉树先序遍历算法（迭代版#1）
func (n *Node) TravPre1(opts ...Option) []*Node {
	s := stack.NewStack()
	nodes := make([]*Node, 0)
	if n != nil {
		s.Push(n)
	}
	for !s.Empty() {
		ni, _ := s.Pop()
		node := ni.(*Node)
		for _, opt := range opts {
			opt(node)
		}
		nodes = append(nodes, node)
		if node.HasRChild() {
			s.Push(node.RChild)
		}
		if node.HasLChild() {
			s.Push(node.LChild)
		}
	}
	return nodes
}

// TravPre2 二叉树先序遍历算法（迭代版#2）
func (n *Node) TravPre2(opts ...Option) []*Node {
	node := n
	s := stack.NewStack()
	nodes := make([]*Node, 0)
	for {
		node.visitAlongLeftBranch(s, &nodes, opts...)
		if s.Empty() {
			break
		}
		ni, _ := s.Pop()
		node = ni.(*Node)
	}
	return nodes
}

//从当前节点出发，沿左分支不断深入，直至没有左分支的节点；沿途节点遇到后立即访问
func (n *Node) visitAlongLeftBranch(s *stack.Stack, nodes *[]*Node, opts ...Option) {
	for n != nil {
		for _, opt := range opts {
			opt(n)
		}
		*nodes = append(*nodes, n)
		if n.HasRChild() {
			s.Push(n.RChild)
		}
		n = n.LChild
	}
}

// TravInR 二叉树中序遍历算法（递归版）
func (n *Node) TravInR(opts ...Option) []*Node {
	nodes := make([]*Node, 0)
	if n == nil {
		return nil
	}
	nodes = append(nodes, n.LChild.TravInR()...)
	for _, opt := range opts {
		opt(n)
	}
	nodes = append(nodes, n)
	nodes = append(nodes, n.RChild.TravInR()...)
	return nodes
}

// TravIn1 二叉树中序遍历算法（迭代版#1）
func (n *Node) TravIn1(opts ...Option) []*Node {
	nodes := make([]*Node, 0)
	s := stack.NewStack()
	for {
		n.goAlongLeftBranch(s)
		if s.Empty() {
			break
		}
		ni, _ := s.Pop()
		n = ni.(*Node)
		for _, opt := range opts {
			opt(n)
		}
		nodes = append(nodes, n)
		n = n.RChild
	}
	return nodes
}

//从当前节点出发，沿左分支不断深入，直至没有左分支的节点
func (n *Node) goAlongLeftBranch(s *stack.Stack) {
	for n != nil {
		s.Push(n)
		n = n.LChild
	}
}

// TravIn2 二叉树中序遍历算法（迭代版#2）
func (n *Node) TravIn2(opts ...Option) []*Node {
	nodes := make([]*Node, 0)
	s := stack.NewStack()
	for {
		if n != nil {
			s.Push(n)
			n = n.LChild
		} else if !s.Empty() {
			ni, _ := s.Pop()
			node := ni.(*Node)
			for _, opt := range opts {
				opt(node)
			}
			nodes = append(nodes, node)
			n = n.RChild
		} else {
			break
		}
	}
	return nodes
}

// TravIn3 二叉树中序遍历算法（迭代版#3）
func (n *Node) TravIn3(opts ...Option) []*Node {
	nodes := make([]*Node, 0)
	backtrack := false
	node := n
	for {
		if !backtrack && node.HasLChild() {
			node = node.LChild
		} else {
			for _, opt := range opts {
				opt(n)
			}
			nodes = append(nodes, n)
			if n.HasRChild() {
				n = n.RChild
				backtrack = false
			} else {
				n = n.succ()
				if n == nil {
					break
				}
				backtrack = true
			}
		}
	}
	return nodes
}

func (n *Node) succ() *Node {
	s := n
	if n.HasRChild() {
		s = n.RChild
		for s.HasLChild() {
			s = s.LChild
		}
		return s
	}
	for s.IsRChild() {
		s = s.Parent
	}
	s = s.Parent
	return s
}

// TravPostR 二叉树后序遍历算法（递归版）
func (n *Node) TravPostR(opts ...Option) []*Node {
	nodes := make([]*Node, 0)
	if n == nil {
		return nil
	}
	node := n
	nodes = append(nodes, node.LChild.TravInR()...)
	nodes = append(nodes, node.RChild.TravInR()...)
	for _, opt := range opts {
		opt(node)
	}
	nodes = append(nodes, node)
	return nodes
}

// TravPost1 二叉树的后序遍历（迭代版）
func (n *Node) TravPost1(opts ...Option) []*Node {
	s := stack.NewStack()
	nodes := make([]*Node, 0)
	if n != nil {
		s.Push(n)
	}
	for !s.Empty() {
		ni, _ := s.Peek()
		node := ni.(*Node)
		if node != n.Parent {
			gotoHLVFL(s)
		}
		ni, _ = s.Pop()
		node = ni.(*Node)
		for _, opt := range opts {
			opt(node)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

//在以S栈顶节点为根的子树中，找到最高左侧可见叶节点
func gotoHLVFL(s *stack.Stack) {
	for ni, ok := s.Peek(); ok; {
		node := ni.(*Node)
		if node.HasLChild() {
			if node.HasRChild() {
				s.Push(node.RChild)
			}
			s.Push(node.LChild)
		} else {
			s.Push(node.RChild)
		}
	}
	s.Pop()
}

// TravLevel 二叉树层次遍历
func (n *Node) TravLevel(opts ...Option) []*Node {
	q := queue.NewQueue()
	nodes := make([]*Node, 0)
	if n != nil {
		q.Enqueue(n)
	}
	for !q.Empty() {
		ni, _ := q.Dequeue()
		node := ni.(*Node)
		for _, opt := range opts {
			opt(node)
		}
		nodes = append(nodes, node)
		if node.HasLChild() {
			q.Enqueue(node.LChild)
		}
		if node.HasRChild() {
			q.Enqueue(node.RChild)
		}
	}
	return nodes
}

// Remove deletes the node n and the sub tree belongs to it, return count of nodes deleted and successfully deleted or not
func (t *BinTree) Remove(n *Node) (int, bool) {
	if n.Tree != t {
		return 0, false
	}
	if n.IsLChild() {
		n.Parent.LChild = nil
	} else {
		n.Parent.RChild = nil
	}
	n.Parent.updateHeightAbove()
	count := removeAt(n)
	t.Size -= count
	return count, true
}

func removeAt(n *Node) (count int) {
	if n == nil {
		return 0
	}
	count = 1
	count += removeAt(n.LChild) + removeAt(n.RChild)
	return
}

// Secede 二叉树子树分离算法：将子树x从当前树中摘除���将其封装为一棵独立子树返回
func (t *BinTree) Secede(n *Node) *BinTree {
	if n.Tree != t {
		return nil
	}
	if n.IsLChild() {
		n.Parent.LChild = nil
	} else {
		n.Parent.RChild = nil
	}
	n.Parent.updateHeightAbove()
	count := n.Size()
	t.Size -= count
	n.Parent = nil
	return &BinTree{n, count}
}

// InsertAsRoot creates a Node as root of a dummy tree
func (t *BinTree) InsertAsRoot(key, data interface{}) *Node {
	t.Root = &Node{Data: data, Key: key, Tree: t}
	t.Size = 1
	return t.Root
}

// InsertAsLChild insert as left child as provided node in the tree
func (t *BinTree) InsertAsLChild(n *Node, key, data interface{}) *Node {
	if n.Tree != t || n.LChild != nil {
		panic("can not insert as left child of provided node in the tree")
	}
	lc := n.InsertAsLChild(key, data)
	t.Size++
	n.updateHeightAbove()
	return lc
}

// InsertAsRChild insert as Right child as provided node in the tree
func (t *BinTree) InsertAsRChild(n *Node, key, data interface{}) *Node {
	if n.Tree != t || n.RChild != nil {
		panic("can not insert as right child of provided node in the tree")
	}
	rc := n.InsertAsRChild(key, data)
	t.Size++
	n.updateHeightAbove()
	return rc
}

// AttachAsLSubTree attaches t as n's left sub tree
func (t *BinTree) AttachAsLSubTree(n *Node, st *BinTree) {
	if n.LChild != nil || n.Tree != t {
		panic("can not attach this sub tree to current node")
	}
	n.LChild = st.Root
	n.LChild.Parent = n
	t.Size += st.Size
	n.updateHeightAbove()
}

// AttachAsRSubTree attaches t as n's right sub tree
func (t *BinTree) AttachAsRSubTree(n *Node, st *BinTree) {
	if n.RChild != nil || n.Tree != t {
		panic("can not attach this sub tree to current node")
	}
	n.RChild = st.Root
	n.RChild.Parent = n
	t.Size += st.Size
	n.updateHeightAbove()
}

// Option is a func to operate on a node
type Option func(n *Node)

// WithPrintNodeKey returns a func which prints node's key
func WithPrintNodeKey() Option {
	return func(n *Node) {
		fmt.Printf("%v ", n.Key)
	}
}

// TravPreR 二叉树前序遍历递归版
func (t *BinTree) TravPreR(opts ...Option) []*Node {
	return t.Root.TravPreR(opts...)
}

// TravPre1 二叉树前序遍历迭代版1
func (t *BinTree) TravPre1(opts ...Option) []*Node {
	return t.Root.TravPre1(opts...)
}

// TravPre2 二叉树前序遍历迭代版2
func (t *BinTree) TravPre2(opts ...Option) []*Node {
	return t.Root.TravPre2(opts...)
}

// TravInR 二叉树中序遍历递归版
func (t *BinTree) TravInR(opts ...Option) []*Node {
	return t.Root.TravInR(opts...)
}

// TravIn1 二叉树中序遍历迭代版1
func (t *BinTree) TravIn1(opts ...Option) []*Node {
	return t.Root.TravIn1(opts...)
}

// TravIn2 二叉树中序遍历迭代版2
func (t *BinTree) TravIn2(opts ...Option) []*Node {
	return t.Root.TravIn2(opts...)
}

// TravIn3 二叉树中序遍历迭代版3
func (t *BinTree) TravIn3(opts ...Option) []*Node {
	return t.Root.TravIn3(opts...)
}

// TravPostR 二叉树后序遍历递归版
func (t *BinTree) TravPostR(opts ...Option) []*Node {
	return t.Root.TravPostR(opts...)
}

// TravPost1 二叉树后序遍历迭代版1
func (t *BinTree) TravPost1(opts ...Option) []*Node {
	return t.Root.TravPost1(opts...)
}

// TravLevel 二叉树层次遍历
func (t *BinTree) TravLevel(opts ...Option) []*Node {
	return t.Root.TravLevel(opts...)
}

// 对n进行右旋, 返回旋转后的局部根节点
func (n *Node) zig() *Node {
	p := n.LChild // 指向n的左孩子
	// 将该左孩子替代n的位置
	p.Parent = n.Parent
	if p.Parent != nil {
		if n.IsLChild() {
			p.Parent.LChild = p
		}
		if n.IsRChild() {
			p.Parent.RChild = p
		}
	}
	// 将p的右孩子改为n的左孩子
	n.LChild = p.RChild
	if n.LChild != nil {
		n.LChild.Parent = n
	}
	// 将n改为p的右孩子
	p.RChild = n
	n.Parent = p
	return p
}

// 对n进行左旋, 返回旋转后的局部根节点
func (n *Node) zag() *Node {
	p := n.RChild // 指向n的右孩子
	// 将该右孩子替代n的位置
	p.Parent = n.Parent
	if p.Parent != nil {
		if n.IsLChild() {
			p.Parent.LChild = p
		}
		if n.IsRChild() {
			p.Parent.RChild = p
		}
	}
	// 将p的左孩子改为n的右孩子
	n.RChild = p.LChild
	if n.RChild != nil {
		n.RChild.Parent = n
	}
	// 将n改为p的左孩子
	p.LChild = n
	n.Parent = p
	return p
}

// Print 以树节点个数的位数为一个基本单元的长度，打印BinTree的拓扑结构到标准输出
func (t *BinTree) Print() {
	t.PrintWithUnitSize(len(strconv.Itoa(t.Size)))
}

// Fprint 以树节点个数的位数为一个基本单元的长度，打印BinTree的拓扑结构到io.Writer
func (t *BinTree) Fprint(w io.Writer) {
	t.FprintWithUnitSize(w, len(strconv.Itoa(t.Size)))
}

// PrintWithUnitSize 以指定的长度为一个基本单元，打印BinTree的拓扑结构到标准输出
func (t *BinTree) PrintWithUnitSize(size int) {
	t.FprintWithUnitSize(os.Stdout, size)
}

// FprintWithUnitSize 以指定的长度为一个基本单元，打印BinTree的拓扑结构到io.Writer
func (t *BinTree) FprintWithUnitSize(w io.Writer, size int) {
	if size <= 0 {
		panic("unit size can not be less than 1")
	}
	buf := bufio.NewWriter(w)
	unitSpace := strings.Repeat(" ", size)
	unitHen := strings.Repeat("─", size)
	q := queue.NewQueue()
	q.Enqueue(t.Root)
	nodewidth := 1<<uint((t.Root.Height+1)) - 1
	for l := 0; l <= t.Root.Height; l++ {
		for i := 0; i < 1<<uint(l); i++ {
			ni, _ := q.Dequeue()
			if ni == nil {
				q.Enqueue(nil)
				q.Enqueue(nil)
				buf.WriteString(fmt.Sprintf("%s", strings.Repeat(unitSpace, nodewidth)))
				buf.WriteString(fmt.Sprintf("%s", unitSpace))
				continue
			}
			node := ni.(*Node)
			if l < t.Root.Height {
				nodeLeftStr := strings.Repeat(unitSpace, (nodewidth-3)/4) +
					strings.Repeat(" ", size-1) + "┌" +
					//"┌" + strings.Repeat("─", size-1) +
					strings.Repeat(unitHen, (nodewidth-3)/4)
				nodeRightStr := strings.Repeat(unitHen, (nodewidth-3)/4) +
					//strings.Repeat(" ", size-1) + "┐" +
					strings.Repeat("─", size-1) + "┐" +
					strings.Repeat(unitSpace, (nodewidth-3)/4)
				if node.HasLChild() {
					q.Enqueue(node.LChild)
					buf.WriteString(fmt.Sprintf("%s", nodeLeftStr))
				} else {
					q.Enqueue(nil)
					buf.WriteString(fmt.Sprintf("%s", strings.Repeat(unitSpace, (nodewidth-1)/2)))
				}
				buf.WriteString(fmt.Sprintf("%*v", size, node.Key))
				if node.HasRChild() {
					q.Enqueue(node.RChild)
					buf.WriteString(fmt.Sprintf("%s", nodeRightStr))
				} else {
					q.Enqueue(nil)
					buf.WriteString(fmt.Sprintf("%s", strings.Repeat(unitSpace, (nodewidth-1)/2)))
				}
				buf.WriteString(fmt.Sprintf("%s", unitSpace))
			} else {
				buf.WriteString(fmt.Sprintf("%*v", size, node.Key))
				buf.WriteString(fmt.Sprintf("%s", unitSpace))
			}
		}
		nodewidth = (nodewidth - 1) / 2
		buf.WriteString(fmt.Sprintf("\n"))
	}
	buf.Flush()
}

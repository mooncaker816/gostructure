package bst

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Node struct {
	LChild *Node
	RChild *Node
	Parent *Node
	Key    interface{}
	Data   interface{}
	Height int
}

// type BinTree struct {
// 	Root *Node
// }

func newSingleNode(key, data interface{}) *Node {
	return &Node{Key: key, Data: data}
}

// // NewBinTree creates a bintree with a single root node
// func NewBinTree(root *Node) *BinTree {
// 	return &BinTree{Root: root}
// }

// HasLChild checks if n has left child
func (n *Node) HasLChild() bool { return n.LChild != nil }

// HasRChild checks if n has right child
func (n *Node) HasRChild() bool { return n.RChild != nil }

// IsLChild checks if n is a left child
func (n *Node) IsLChild() bool {
	return !n.IsRoot() && n == n.Parent.LChild
}

// IsRChild checks if n is a right child
func (n *Node) IsRChild() bool {
	return !n.IsRoot() && n == n.Parent.RChild
}

// IsRoot checks if n is root
func (n *Node) IsRoot() bool { return n.Parent == nil }

// Sibling gets n's sibling if exists
func (n *Node) Sibling() *Node {
	if n.IsRoot() {
		return nil
	}
	if n.IsLChild() {
		return n.Parent.RChild
	}
	return n.Parent.LChild
}

func (n *Node) AttachRChild(v *Node) {
	n.RChild = v
	if v != nil {
		v.Parent = n
	}
}

func (n *Node) AttachLChild(v *Node) {
	n.LChild = v
	if v != nil {
		v.Parent = n
	}
}

// AddLChild generates a new Node with the provided key and data, then adds to Node n as left child
func (n *Node) AddLChild(key, data interface{}) *Node {
	if n.HasLChild() {
		panic("node already has left child")
	}
	newNode := newSingleNode(key, data)
	n.LChild = newNode
	newNode.Parent = n
	if !n.HasRChild() {
		n.UpdateHeightAbove()
	}
	return newNode
}

// AddRChild generates a new Node with the provided key and data, then adds to Node n as right child
func (n *Node) AddRChild(key, data interface{}) *Node {
	if n.HasRChild() {
		panic("node already has right child")
	}
	newNode := newSingleNode(key, data)
	n.RChild = newNode
	newNode.Parent = n
	if !n.HasLChild() {
		n.UpdateHeightAbove()
	}
	return newNode
}

// AttachLSubTree will add left sub tree to n
func (n *Node) AttachLSubTree(left *Node) {
	if n.HasLChild() {
		panic("node already has left sub tree")
	}
	n.LChild = left
	left.Parent = n
	n.UpdateHeightAbove()
}

// AttachRSubTree will add right sub tree to n
func (n *Node) AttachRSubTree(right *Node) {
	if n.HasRChild() {
		panic("node already has right sub tree")
	}
	n.RChild = right
	right.Parent = n
	n.UpdateHeightAbove()
}

func (n *Node) UpdateHeight() {
	n.Height = n.maxHeightOfChildren() + 1
}

// UpdateHeightAbove updates height info for all the related nodes
func (n *Node) UpdateHeightAbove() {
	for max := n.maxHeightOfChildren(); n != nil && n.Height != max+1; {
		n.Height = max + 1
		n = n.Parent
	}
}

func (n *Node) maxHeightOfChildren() int {
	lH, rH := -1, -1
	if n.HasLChild() {
		lH = n.LChild.Height
	}
	if n.HasRChild() {
		rH = n.RChild.Height
	}
	return max(lH, rH)
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// Detach removes the whole sub tree
func (n *Node) Detach() {
	if n.IsRoot() {
		return
	}
	if n.IsLChild() {
		n.Parent.LChild = nil
		n.Parent.UpdateHeightAbove()
		n.Parent = nil
		return
	}
	n.Parent.RChild = nil
	n.Parent.UpdateHeightAbove()
	n.Parent = nil
}

// Size returns the count of sub tree's nodes
func (n *Node) Size() int {
	if n == nil {
		return 0
	}
	count := 1
	if n.HasLChild() {
		count += n.LChild.Size()
	}
	if n.HasRChild() {
		count += n.RChild.Size()
	}
	return count
}

// Level returns the level of the node in the tree, root is level 0
func (n *Node) Level() int {
	l := 0
	for !n.IsRoot() {
		l++
		n = n.Parent
	}
	return l
}

// Option is a useful function to dealing with the node during the traversal of the tree
type Option func(n *Node)

// TravPre provides the pre-order traversal
func (n *Node) TravPre(opts ...Option) {
	for _, opt := range opts {
		opt(n)
	}
	if n.HasLChild() {
		n.LChild.TravPre(opts...)
	}
	if n.HasRChild() {
		n.RChild.TravPre(opts...)
	}
}

// TravIn provides the in-order traversal
func (n *Node) TravIn(opts ...Option) {
	if n.HasLChild() {
		n.LChild.TravIn(opts...)
	}
	for _, opt := range opts {
		opt(n)
	}
	if n.HasRChild() {
		n.RChild.TravIn(opts...)
	}
}

// TravPost provides the post-order traversal
func (n *Node) TravPost(opts ...Option) {
	if n.HasLChild() {
		n.LChild.TravIn(opts...)
	}
	if n.HasRChild() {
		n.RChild.TravIn(opts...)
	}
	for _, opt := range opts {
		opt(n)
	}
}

// TravLevel provides the level-order traversal
func (n *Node) TravLevel(opts ...Option) {
	queue := make([]*Node, 0, n.Size())
	queue = append(queue, n)
	for len(queue) > 0 {
		visitNode := queue[0]
		queue = queue[1:]
		if visitNode.HasLChild() {
			queue = append(queue, visitNode.LChild)
		}
		if visitNode.HasRChild() {
			queue = append(queue, visitNode.RChild)
		}
		for _, opt := range opts {
			opt(visitNode)
		}
	}
}

// Print 以子树节点个数的位数为一个基本单元的长度，打印子树的拓扑结构到标准输出
func (n *Node) Print() {
	n.PrintWithUnitSize(len(strconv.Itoa(n.Size())))
}

// Fprint 以树节点个数的位数为一个基本单元的长度，打印子树的拓扑结构到io.Writer
func (n *Node) Fprint(w io.Writer) {
	n.FprintWithUnitSize(w, len(strconv.Itoa(n.Size())))
}

// PrintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到标准输出
func (n *Node) PrintWithUnitSize(size int) {
	n.FprintWithUnitSize(os.Stdout, size)
}

// FprintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到io.Writer，树宽为节点数
func (n *Node) FprintWithUnitSize(w io.Writer, size int) {
	buf := bufio.NewWriter(w)
	if n == nil {
		buf.WriteString("Empty tree!")
		buf.Flush()
		return
	}
	if size <= 0 {
		panic("unit size can not be less than 1")
	}

	total := n.Size()
	q := make([]nodePos, 0, total)
	prevlevel := 0
	line := make([]rune, total*size+1)
	for i := range line {
		line[i] = ' '
	}
	mid := n.LChild.Size()
	left, right := mid, mid
	if n.LChild != nil {
		left = mid - n.LChild.RChild.Size() - 1
	}
	if n.RChild != nil {
		right = mid + n.RChild.LChild.Size() + 1
	}
	q = append(q, nodePos{n, left, mid, right})
	for len(q) > 0 {
		np := q[0]
		q = q[1:]
		n := np.node
		if l := n.Level(); prevlevel != l {
			prevlevel = l
			buf.WriteString("\n")
			for i, r := range line {
				buf.WriteRune(r)
				line[i] = ' '
			}
		}

		np.fillNode(line, size)

		if n.HasLChild() {
			left, mid, right = np.computeLChildPos()
			q = append(q, nodePos{n.LChild, left, mid, right})
		}
		if n.HasRChild() {
			left, mid, right = np.computeRChildPos()
			q = append(q, nodePos{n.RChild, left, mid, right})
		}
	}
	buf.WriteString("\n")
	for _, r := range line {
		buf.WriteRune(r)
	}
	buf.WriteString("\n")
	buf.Flush()
}

func (np nodePos) fillNode(line []rune, size int) {
	if np.node.HasLChild() {
		i := np.left * size
		// for ; i < np.left*size; i++ {
		// 	line[i] = ' '
		// }
		for ; i < np.left*size+size-1; i++ {
			line[i] = ' '
		}
		line[i] = '┌'
		i++
		for ; i < np.mid*size; i++ {
			line[i] = '─'
		}
	}
	i := np.mid * size
	for _, r := range fmt.Sprintf("%*v", size, np.node.Key) {
		line[i] = r
		i++
	}
	if np.node.HasRChild() {
		// for ; i < np.right*size-size; i++ {
		// 	line[i] = '─'
		// }
		for ; i < np.right*size+size-1; i++ {
			line[i] = '─'
		}
		line[i] = '┐'
	}
}

type nodePos struct {
	node             *Node
	left, mid, right int
}

func (np nodePos) computeLChildPos() (left, mid, right int) {
	if np.node == nil {
		return
	}
	if np.node.LChild == nil {
		return mid, mid, mid
	}
	mid = np.mid - np.node.LChild.RChild.Size() - 1

	if np.node.LChild.LChild == nil {
		left = mid
	} else {
		left = mid - np.node.LChild.LChild.RChild.Size() - 1
	}
	if np.node.LChild.RChild == nil {
		right = mid
	} else {
		right = mid + np.node.LChild.RChild.LChild.Size() + 1
	}
	return
}

func (np nodePos) computeRChildPos() (left, mid, right int) {
	if np.node == nil {
		return
	}
	if np.node.RChild == nil {
		return mid, mid, mid
	}
	mid = np.mid + np.node.RChild.LChild.Size() + 1

	if np.node.RChild.LChild == nil {
		left = mid
	} else {
		left = mid - np.node.RChild.LChild.RChild.Size() - 1
	}
	if np.node.RChild.RChild == nil {
		right = mid
	} else {
		right = mid + np.node.RChild.RChild.LChild.Size() + 1
	}
	return
}

// TallerChild 返回高度较高的那个孩子节点，若同高，返回和n同侧的节点
func (n *Node) TallerChild() *Node {
	if n.HasLChild() && n.RChild == nil {
		return n.LChild
	}
	if n.HasRChild() && n.LChild == nil {
		return n.RChild
	}

	if n.HasLChild() && n.HasRChild() {
		if n.LChild.Height < n.RChild.Height {
			return n.RChild
		}
		if n.LChild.Height > n.RChild.Height {
			return n.LChild
		}
		if n.IsLChild() {
			return n.LChild
		}
		return n.RChild
	}
	return nil
}

func subTreeMin(n *Node) *Node {
	for n.HasLChild() {
		n = n.LChild
	}
	return n
}

func subTreeMax(n *Node) *Node {
	for n.HasRChild() {
		n = n.RChild
	}
	return n
}

// Successor returns the next larger node of n if it exists.
func (n *Node) Successor() *Node {
	if n.HasRChild() {
		return subTreeMin(n.RChild)
	}
	for n.IsRChild() {
		n = n.Parent
	}
	return n.Parent
}

// Predecessor returns the prev smaller node of n if it exists.
func (n *Node) Predecessor() *Node {
	if n.HasLChild() {
		return subTreeMax(n.LChild)
	}
	for n.IsLChild() {
		n = n.Parent
	}
	return n.Parent
}

func SwapKeyData(n1, n2 *Node) {
	n1.Key, n1.Data, n2.Key, n2.Data = n2.Key, n2.Data, n1.Key, n1.Data
}

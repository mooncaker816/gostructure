package bst

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type BNode struct {
	LChild *BNode
	RChild *BNode
	Parent *BNode
	Key    interface{}
	Data   interface{}
}

func newSingleNode(key, data interface{}) *BNode {
	return &BNode{Key: key, Data: data}
}

// HasLChild checks if n has left child
func (n *BNode) HasLChild() bool { return n.LChild != nil }

// HasRChild checks if n has right child
func (n *BNode) HasRChild() bool { return n.RChild != nil }

// IsLChild checks if n is a left child
func (n *BNode) IsLChild() bool {
	return !n.IsRoot() && n == n.Parent.LChild
}

// IsRChild checks if n is a right child
func (n *BNode) IsRChild() bool {
	return !n.IsRoot() && n == n.Parent.RChild
}

// IsRoot checks if n is root
func (n *BNode) IsRoot() bool { return n.Parent == nil }

// Sibling gets n's sibling if exists
func (n *BNode) Sibling() *BNode {
	if n.IsRoot() {
		return nil
	}
	if n.IsLChild() {
		return n.Parent.RChild
	}
	return n.Parent.LChild
}

func (n *BNode) AttachRChild(v *BNode) {
	n.RChild = v
	if v != nil {
		v.Parent = n
	}
}

func (n *BNode) AttachLChild(v *BNode) {
	n.LChild = v
	if v != nil {
		v.Parent = n
	}
}

// // AddLChild generates a new Node with the provided key and data, then adds to Node n as left child
// func (n *BNode) AddLChild(key, data interface{}) *BNode {
// 	if n.HasLChild() {
// 		panic("node already has left child")
// 	}
// 	newNode := newSingleNode(key, data)
// 	n.LChild = newNode
// 	newNode.Parent = n
// 	if !n.HasRChild() {
// 		n.UpdateHeightAbove()
// 	}
// 	return newNode
// }

// // AddRChild generates a new Node with the provided key and data, then adds to Node n as right child
// func (n *BNode) AddRChild(key, data interface{}) *BNode {
// 	if n.HasRChild() {
// 		panic("node already has right child")
// 	}
// 	newNode := newSingleNode(key, data)
// 	n.RChild = newNode
// 	newNode.Parent = n
// 	if !n.HasLChild() {
// 		n.UpdateHeightAbove()
// 	}
// 	return newNode
// }

// // AttachLSubTree will add left sub tree to n
// func (n *BNode) AttachLSubTree(left *BNode) {
// 	if n.HasLChild() {
// 		panic("node already has left sub tree")
// 	}
// 	n.LChild = left
// 	left.Parent = n
// 	n.UpdateHeightAbove()
// }

// // AttachRSubTree will add right sub tree to n
// func (n *BNode) AttachRSubTree(right *BNode) {
// 	if n.HasRChild() {
// 		panic("node already has right sub tree")
// 	}
// 	n.RChild = right
// 	right.Parent = n
// 	n.UpdateHeightAbove()
// }

// // Detach removes the whole sub tree
// func (n *BNode) Detach() {
// 	if n.IsRoot() {
// 		return
// 	}
// 	if n.IsLChild() {
// 		n.Parent.LChild = nil
// 		n.Parent.UpdateHeightAbove()
// 		n.Parent = nil
// 		return
// 	}
// 	n.Parent.RChild = nil
// 	n.Parent.UpdateHeightAbove()
// 	n.Parent = nil
// }

// Size returns the count of sub tree's nodes
func (n *BNode) Size() int {
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
func (n *BNode) Level() int {
	l := 0
	for !n.IsRoot() {
		l++
		n = n.Parent
	}
	return l
}

// Option is a useful function to dealing with the node during the traversal of the tree
// type Option func(n *BNode)

// TravPre provides the pre-order traversal
func (n *BNode) TravPre(opts ...Option) {
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
func (n *BNode) TravIn(opts ...Option) {
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
func (n *BNode) TravPost(opts ...Option) {
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
func (n *BNode) TravLevel(opts ...Option) {
	queue := make([]*BNode, 0, n.Size())
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
func (n *BNode) Print() {
	n.PrintWithUnitSize(len(strconv.Itoa(n.Size())))
}

// Fprint 以树节点个数的位数为一个基本单元的长度，打印子树的拓扑结构到io.Writer
func (n *BNode) Fprint(w io.Writer) {
	n.FprintWithUnitSize(w, len(strconv.Itoa(n.Size())))
}

// PrintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到标准输出
func (n *BNode) PrintWithUnitSize(size int) {
	n.FprintWithUnitSize(os.Stdout, size)
}

// FprintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到io.Writer，树宽为节点数
func (n *BNode) FprintWithUnitSize(w io.Writer, size int) {
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
	node             *BNode
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

func subTreeMin(n *BNode) *BNode {
	for n.HasLChild() {
		n = n.LChild
	}
	return n
}

func subTreeMax(n *BNode) *BNode {
	for n.HasRChild() {
		n = n.RChild
	}
	return n
}

// Successor returns the next larger node of n if it exists.
func (n *BNode) Successor() *BNode {
	if n.HasRChild() {
		return subTreeMin(n.RChild)
	}
	for n.IsRChild() {
		n = n.Parent
	}
	return n.Parent
}

// Predecessor returns the prev smaller node of n if it exists.
func (n *BNode) Predecessor() *BNode {
	if n.HasLChild() {
		return subTreeMax(n.LChild)
	}
	for n.IsLChild() {
		n = n.Parent
	}
	return n.Parent
}

func SwapKeyData(n1, n2 *BNode) {
	n1.Key, n1.Data, n2.Key, n2.Data = n2.Key, n2.Data, n1.Key, n1.Data
}

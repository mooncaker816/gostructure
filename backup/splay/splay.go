package splay

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

type Node struct {
	Parent *Node
	LChild *Node
	RChild *Node
	Data   interface{}
	Key    interface{}
}

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

type SplayTree struct {
	Root *Node
}

func NewSplayTree() *SplayTree {
	return new(SplayTree)
}

func (s *SplayTree) Print() {
	s.Root.PrintWithUnitSize(2)
}

func (s *SplayTree) PrintWithUnitSize(size int) {
	s.Root.PrintWithUnitSize(size)
}

func (s *SplayTree) Insert(key, data interface{}) (*Node, error) {
	if s.Root == nil {
		s.Root = &Node{Key: key, Data: data}
		return s.Root, nil
	}
	n, result := s.searchIn(s.Root, key)
	switch result {
	case 0:
		return nil, errors.New("insert node with duplicate key")
	case -1:
		new := &Node{Key: key, Data: data}
		new.attachRChild(n)
		new.attachLChild(n.LChild)
		n.LChild = nil
		s.Root = new
		return new, nil
	case 1:
		new := &Node{Key: key, Data: data}
		new.attachLChild(n)
		new.attachRChild(n.RChild)
		n.RChild = nil
		s.Root = new
		return new, nil
	}
	return nil, errors.New("unknow compare result")
}

func (s *SplayTree) Search(key interface{}) (*Node, bool) {
	n, result := s.searchIn(s.Root, key)
	if result == 0 {
		return n, true
	}
	return n, false
}

func (s *SplayTree) searchIn(n *Node, key interface{}) (*Node, int) {
	switch comp(key, n.Key) {
	case 0:
		s.Root = n.splay()
		return s.Root, 0
	case -1:
		if n.LChild != nil {
			return s.searchIn(n.LChild, key)
		}
		s.Root = n.splay()
		return s.Root, -1
	case 1:
		if n.RChild != nil {
			return s.searchIn(n.RChild, key)
		}
		s.Root = n.splay()
		return s.Root, 1
	}
	return nil, 0
}

func (s *SplayTree) Remove(key interface{}) (*Node, error) {
	n, result := s.searchIn(s.Root, key)
	if result != 0 {
		return nil, nil
	}
	// 此时待删节点位于 root
	if !s.Root.HasLChild() {
		s.Root = s.Root.RChild
		if s.Root != nil {
			s.Root.Parent = nil
		}
		n.RChild = nil
	} else if !s.Root.HasRChild() {
		s.Root = s.Root.LChild
		s.Root.Parent = nil
		n.LChild = nil
	} else {
		// 临时切除左子树
		lc := s.Root.LChild
		s.Root.LChild = nil
		// 删除 root，右子树成为新的树
		s.Root = s.Root.RChild
		s.Root.Parent = nil
		n.RChild = nil
		// 在新树中再次查找原来的值，必然不存在，但会把最小值提升至顶端，且没有左子树
		// 该最小值一定比之前切除的左子树大，以此值为 root 重新连接原左子树即可
		s.searchIn(s.Root, key)
		s.Root.LChild = lc
		lc.Parent = s.Root
	}
	return n, nil
}

func comp(a, b interface{}) int {

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	if va.Kind() != vb.Kind() {
		panic("can not compare between different types")
	}
	switch va.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if va.Int() == vb.Int() {
			return 0
		}
		if va.Int() < vb.Int() {
			return -1
		}
		return 1
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if va.Uint() == vb.Uint() {
			return 0
		}
		if va.Uint() < vb.Uint() {
			return -1
		}
		return 1
	case
		reflect.Float32,
		reflect.Float64:
		if va.Float() == vb.Float() {
			return 0
		}
		if va.Float() < vb.Float() {
			return -1
		}
		return 1
	case
		reflect.String:
		if va.String() == vb.String() {
			return 0
		}
		if va.String() < vb.String() {
			return -1
		}
		return 1
	default:
		panic("type not support for comparing,should customize compare method")
	}
}

func (n *Node) splay() *Node {
	if n == nil {
		return nil
	}
	for n.Parent != nil && n.Parent.Parent != nil {
		p := n.Parent
		g := n.Parent.Parent
		gp := g.Parent
		if n.IsLChild() {
			if p.IsLChild() {
				g.attachLChild(p.RChild)
				p.attachLChild(n.RChild)
				p.attachRChild(g)
				n.attachRChild(p)
			} else {
				p.attachLChild(n.RChild)
				g.attachRChild(n.LChild)
				n.attachLChild(g)
				n.attachRChild(p)
			}
		} else {
			if p.IsLChild() {
				p.attachRChild(n.LChild)
				g.attachLChild(n.RChild)
				n.attachRChild(g)
				n.attachLChild(p)
			} else {
				g.attachRChild(p.LChild)
				p.attachRChild(n.LChild)
				p.attachLChild(g)
				n.attachLChild(p)
			}
		}
		if gp != nil {
			if gp.LChild == g {
				gp.attachLChild(n)
			} else {
				gp.attachRChild(n)
			}
		} else {
			n.Parent = nil
		}
	}
	if n.Parent != nil {
		if n.IsLChild() {
			n.Parent.attachLChild(n.RChild)
			n.attachRChild(n.Parent)
		} else {
			n.Parent.attachRChild(n.LChild)
			n.attachLChild(n.Parent)
		}
	}
	n.Parent = nil
	return n
}

func (n *Node) attachRChild(v *Node) {
	n.RChild = v
	if v != nil {
		v.Parent = n
	}
}

func (n *Node) attachLChild(v *Node) {
	n.LChild = v
	if v != nil {
		v.Parent = n
	}
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

// PrintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到标准输出
func (n *Node) PrintWithUnitSize(size int) {
	n.FprintWithUnitSize(os.Stdout, size)
}

// // FprintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到io.Writer，树宽为节点数
// func (n *Node) FprintWithUnitSize(w io.Writer, size int) {
// 	buf := bufio.NewWriter(w)
// 	if n == nil {
// 		buf.WriteString("Empty tree!")
// 		buf.Flush()
// 		return
// 	}
// 	if size <= 0 {
// 		panic("unit size can not be less than 1")
// 	}
// 	unitSpace := strings.Repeat(" ", size)
// 	unitHen := strings.Repeat("─", size)
// 	type nodeOffset struct {
// 		node   *Node
// 		offset int
// 	}
// 	q := make([]nodeOffset, 0, n.Size())
// 	q = append(q, nodeOffset{n, 0})
// 	prevlevel := 0

// 	delta := 0 // 由于没有孩子节点需要过继给后续有孩子节点的偏移量
// 	for len(q) > 0 {
// 		no := q[0]
// 		q = q[1:]
// 		n := no.node
// 		offset := no.offset
// 		if l := n.Level(); prevlevel != l {
// 			delta = 0 // 原二叉树与对应的完全二叉树中缺失节点造成的空格缺失数
// 			prevlevel = l
// 			buf.WriteString("\n")
// 		}
// 		var nodeLeftStr, nodeRightStr string
// 		// 由父亲节点遗留下来的前缀偏移量
// 		buf.WriteString(strings.Repeat(unitSpace, offset))
// 		if n.IsRChild() && n.Sibling() == nil {
// 			offset++                   // 缺失左兄弟导致原本兄弟与兄弟之间的一个单位空格缺失，需传给该节点的后代节点
// 			buf.WriteString(unitSpace) // 在打印右儿子之前补上该空格
// 		}
// 		if n.HasLChild() {
// 			offset += delta // 将偏移量加上由之前同层的叶子节点造成的空儿子节点引起的空格数，完成过继后置零
// 			delta = 0
// 			//如果该节点有左孩子，优先将偏移量转移至左孩子（因为从左往右打印）
// 			q = append(q, nodeOffset{n.LChild, offset})
// 			nodeLeftStr = strings.Repeat(unitSpace, n.LChild.LChild.Size()) +
// 				strings.Repeat(" ", size-1) + "┌" +
// 				strings.Repeat(unitHen, n.LChild.RChild.Size())
// 		}
// 		if n.HasRChild() {
// 			if n.HasLChild() {
// 				offset = 0 //若有左孩子，则偏移量已经转移至左孩子，无需再转移给右孩子
// 			} else { // 不得已转移给右孩子
// 				offset += delta
// 				delta = 0
// 			}
// 			q = append(q, nodeOffset{n.RChild, offset})
// 			nodeRightStr = strings.Repeat(unitHen, n.RChild.LChild.Size()) +
// 				strings.Repeat("─", size-1) + "┐" +
// 				strings.Repeat(unitSpace, n.RChild.RChild.Size())
// 		}
// 		if !n.HasLChild() && !n.HasRChild() { // 叶节点需保存当前的偏移量再加上二个单位的空格，+= 防止连续叶节点导致偏移量丢失
// 			delta += offset
// 			delta++
// 			delta++
// 		}

// 		buf.WriteString(nodeLeftStr)
// 		buf.WriteString(fmt.Sprintf("%*v", size, n.Key))
// 		buf.WriteString(nodeRightStr)
// 		buf.WriteString(unitSpace)
// 		if n.IsLChild() && n.Sibling() == nil {
// 			// 缺失左兄弟导致原本兄弟与兄弟之间的一个单位空格缺失，需传给该节点的后代节点
// 			buf.WriteString(unitSpace) // 在打印右儿子之前补上该空格
// 		}
// 	}
// 	buf.WriteString("\n")
// 	buf.Flush()
// }

// Level returns the level of the node in the tree, root is level 0
func (n *Node) Level() int {
	l := 0
	for !n.IsRoot() {
		l++
		n = n.Parent
	}
	return l
}

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

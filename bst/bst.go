package bst

import "github.com/mooncaker816/gostructure/bintree"

// Bst 二叉搜索树
type Bst struct {
	*bintree.BinTree
	Hot        *bintree.Node // "命中节点"的父节点
	Comparator Comparator
}

// NewBst 新建Bst
func NewBst() *Bst {
	return NewBstWithComparator(BasicCompare)
}

// NewBstWithComparator 新建自定义比较器的Bst
func NewBstWithComparator(c Comparator) *Bst {
	return &Bst{new(bintree.BinTree), nil, c}
}

// Search 查找key为e的节点，失败返回nil，该nil已经区分
func (bst *Bst) Search(key interface{}) *bintree.Node {
	bst.Hot = nil
	return bst.searchIn(bst.Root, key)
}

func (bst *Bst) searchIn(n *bintree.Node, key interface{}) *bintree.Node {
	if n == nil || bst.Comparator(key, n.Key) == 0 {
		return n
	}
	bst.Hot = n
	if bst.Comparator(key, n.Key) == -1 {
		return bst.searchIn(n.LChild, key)
	}

	return bst.searchIn(n.RChild, key)

}

// Insert 按数据插入新节点并返回该节点
func (bst *Bst) Insert(key, data interface{}) (node *bintree.Node) {
	if node := bst.Search(key); node != nil {
		return node
	}
	if bst.Comparator(key, bst.Hot.Key) == -1 {
		node = bst.Hot.InsertAsLChild(key, data)
	} else {
		node = bst.Hot.InsertAsRChild(key, data)
	}
	node.UpdateHeightAbove()
	bst.Size++
	return node
}

// Remove 按key删除节点，若成功，则bst.Hot为正真删除节点的父节点，若失败，则为查找返回的Hot
func (bst *Bst) Remove(key interface{}) bool {
	n := bst.Search(key)
	if n == nil {
		return false
	}
	bst.RemoveAt(n)
	bst.Size--
	bst.Hot.UpdateHeightAbove()
	return true
}

func (bst *Bst) RemoveAt(n *bintree.Node) (succ *bintree.Node) {
	w := n
	if !n.HasLChild() {
		succ = n.RChild
		if n.IsLChild() {
			n.Parent.LChild = succ
		}
		if n.IsRChild() {
			n.Parent.RChild = succ
		}
		if n.IsRoot() {
			bst.Root = succ
		}
	} else if !n.HasRChild() {
		succ = n.LChild
		if n.IsLChild() {
			n.Parent.LChild = succ
		}
		if n.IsRChild() {
			n.Parent.RChild = succ
		}
		if n.IsRoot() {
			bst.Root = succ
		}
	} else {
		w = n.Succ()                    // 找到n的直接后继节点w，即右子树中左边最高节点
		n.Data, w.Data = w.Data, n.Data //交换n和w的数据项
		n.Key, w.Key = w.Key, n.Key     //交换n和w的key
		succ = w.RChild                 //待删节点的右子树
		if w.Parent == n {
			w.Parent.RChild = succ
		} else {
			w.Parent.LChild = succ
		}
	}
	bst.Hot = w.Parent
	if succ != nil {
		succ.Parent = bst.Hot
	}
	w.Parent, w.LChild, w.RChild = nil, nil, nil
	return succ
}

// Connect34 3+4 重构:
//	 	   b
//		a	  c
//	  T0 T1 T2 T3
func Connect34(a, b, c, T0, T1, T2, T3 *bintree.Node) *bintree.Node {
	a.LChild = T0
	if T0 != nil {
		T0.Parent = a
	}
	a.RChild = T1
	if T1 != nil {
		T1.Parent = a
	}
	a.UpdateHeight()
	c.LChild = T2
	if T2 != nil {
		T2.Parent = c
	}
	c.RChild = T3
	if T3 != nil {
		T3.Parent = c
	}
	c.UpdateHeight()

	b.LChild = a
	a.Parent = b
	b.RChild = c
	c.Parent = b
	b.UpdateHeight()
	return b
}

func RotateAt(v *bintree.Node) *bintree.Node {
	if v == nil {
		panic("can not rotate on nil Node")
	}
	p := v.Parent
	g := p.Parent
	if p.IsLChild() { // zig
		if v.IsLChild() { // zig-zig
			p.Parent = g.Parent
			return Connect34(v, p, g, v.LChild, v.RChild, p.RChild, g.RChild)
		}
		// zig-zag
		v.Parent = g.Parent
		return Connect34(p, v, g, p.LChild, v.LChild, v.RChild, g.RChild)

	}
	// zag
	if v.IsRChild() { //zag-zag
		p.Parent = g.Parent
		return Connect34(g, p, v, g.LChild, p.LChild, v.LChild, v.RChild)
	}
	//zag-zig
	v.Parent = g.Parent
	return Connect34(g, v, p, g.LChild, v.LChild, v.RChild, p.RChild)
}

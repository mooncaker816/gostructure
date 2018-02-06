package bst

import "github.com/mooncaker816/gostructure/bintree"

// Bst 二叉搜索树
type Bst struct {
	*bintree.BinTree
	hot        *bintree.Node // "命中节点"的父节点
	comparator Comparator
}

// NewBst 新建Bst
func NewBst() *Bst {
	return NewBstWithComparator(BasicCompare)
}

// NewBstWithComparator 新建自定义比较器的Bst
func NewBstWithComparator(c Comparator) *Bst {
	return &Bst{new(bintree.BinTree), nil, c}
}

// Search 查找key为e的节点
func (bst *Bst) Search(key interface{}) *bintree.Node {
	bst.hot = nil
	return bst.searchIn(bst.Root, key)
}

func (bst *Bst) searchIn(n *bintree.Node, key interface{}) *bintree.Node {
	if n == nil || bst.comparator(key, n.Key) == 0 {
		return n
	}
	bst.hot = n
	if bst.comparator(key, n.Key) == -1 {
		return bst.searchIn(n.LChild, key)
	}

	return bst.searchIn(n.RChild, key)

}

// Insert 按数据插入新节点并返回该节点
func (bst *Bst) Insert(key, data interface{}) (node *bintree.Node) {
	if node := bst.Search(key); node != nil {
		return node
	}
	if bst.comparator(key, bst.hot.Key) == -1 {
		node = bst.hot.InsertAsLChild(key, data)
	} else {
		node = bst.hot.InsertAsRChild(key, data)
	}
	node.UpdateHeightAbove()
	bst.Size++
	return node
}

// Remove 按key删除节点，若成功，则bst.hot为正真删除节点的父节点，若失败，则为查找返回的hot
func (bst *Bst) Remove(key interface{}) bool {
	n := bst.Search(key)
	if n == nil {
		return false
	}
	bst.removeAt(n)
	bst.Size--
	bst.hot.UpdateHeightAbove()
	return true
}

func (bst *Bst) removeAt(n *bintree.Node) (succ *bintree.Node) {
	w := n
	if !n.HasLChild() {
		succ = n.RChild
	} else if !n.HasRChild() {
		succ = n.LChild
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
	bst.hot = w.Parent
	if succ != nil {
		succ.Parent = bst.hot
	}
	w.Parent, w.LChild, w.RChild = nil, nil, nil
	return succ
}

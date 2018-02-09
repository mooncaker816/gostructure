package avl

import (
	"github.com/mooncaker816/gostructure/bintree"
	"github.com/mooncaker816/gostructure/bst"
)

// Avl Avl 平衡二叉树
type Avl struct {
	*bst.Bst
}

// NewAvl 新建默认比较器的Avl树
func NewAvl() *Avl {
	return &Avl{bst.NewBst()}
}

// NewAvlWithComparator 新建带自定义比较器的Avl树
func NewAvlWithComparator(c bst.Comparator) *Avl {
	return &Avl{bst.NewBstWithComparator(c)}
}

// Insert 插入新节点
func (avl *Avl) Insert(key, data interface{}) (node *bintree.Node) {
	if node := avl.Search(key); node != nil {
		return node
	}
	if avl.Comparator(key, avl.Hot.Key) == -1 {
		node = avl.Hot.InsertAsLChild(key, data)
	} else {
		node = avl.Hot.InsertAsRChild(key, data)
	}
	avl.Size++
	for g := avl.Hot; g != nil; g = g.Parent {
		if !g.IsAvlBalanced() {
			//x := bst.RotateAt(g.TallerChild().TallerChild())
			x := g.Parent
			if g.IsLChild() {
				x.LChild = bst.RotateAt(g.TallerChild().TallerChild())
			} else if g.IsRChild() {
				x.RChild = bst.RotateAt(g.TallerChild().TallerChild())
			} else {
				avl.Root = bst.RotateAt(g.TallerChild().TallerChild())
			}
			break
		} else {
			g.UpdateHeight()
		}
	}
	return node
}

// Remove 按key删除节点
func (avl *Avl) Remove(key interface{}) bool {
	n := avl.Search(key)
	if n == nil {
		return false
	}
	avl.RemoveAt(n)
	avl.Size--
	for g := avl.Hot; g != nil; g = g.Parent {
		if !g.IsAvlBalanced() {
			x := g.Parent
			if g.IsLChild() {
				x.LChild = bst.RotateAt(g.TallerChild().TallerChild())
				g = x.LChild
			} else if g.IsRChild() {
				x.RChild = bst.RotateAt(g.TallerChild().TallerChild())
				g = x.RChild
			} else {
				// g is root
				avl.Root = bst.RotateAt(g.TallerChild().TallerChild())
				g = avl.Root
			}
		}
		g.UpdateHeight()
	}
	return true
}

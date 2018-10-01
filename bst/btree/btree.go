package btree

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/mooncaker816/gostructure/bst"
)

type BTree struct {
	m    int // 阶数
	root *Node
	comp bst.Comparator
	hot  *Node
}

func New(m int) *BTree {
	return &BTree{m: m, comp: bst.BasicCompare}
}

func (b *BTree) Search(key interface{}) (*Node, bool) {
	n, _, ok := b.searchIn(b.root, key)
	return n, ok
}

func (b *BTree) searchIn(n *Node, key interface{}) (hot *Node, i int, ok bool) {
	i = sort.Search(len(n.Key), func(i int) bool {
		result := b.comp(n.Key[i], key)
		return result == 1 || result == 0
	})

	// 在当前节点找到了key
	if i < len(n.Key) && b.comp(n.Key[i], key) == 0 {
		return n, i, true
	}
	// 当前节点没有 key，且为叶子节点，则查找失败
	if n.Children == nil {
		return n, i, false
	}
	// 当前节点没有 key，且当前节点有孩子节点，继续查找
	return b.searchIn(n.Children[i], key)
}

// Insert returns the exact node which stores the newly inserted key
func (b *BTree) Insert(key, data interface{}) (*Node, error) {
	if b.root == nil {
		b.root = newNode(key, data, b.m)
		return b.root, nil
	}
	n, i, ok := b.searchIn(b.root, key)
	if ok {
		return nil, errors.New("insert with duplicate key")
	}
	n.Key = insert(n.Key, key, i)
	n.Data = insert(n.Data, data, i)
	// n.Children = append(n.Children, nil)
	// copy(n.Children[i+1:], n.Children[i:])
	// n.Children[i] = nil
	b.solveOverflow(n, key)
	return b.hot, nil
}

func insert(a []interface{}, v interface{}, i int) []interface{} {
	if i < 0 {
		panic("insert index can not be negetive")
	}
	a = append(a, v)
	if i < len(a)-1 {
		copy(a[i+1:], a[i:])
		a[i] = v
	}
	return a
}

func (b *BTree) solveOverflow(n *Node, origKey interface{}) {
	b.hot = n
	if b.m >= len(n.Key)+1 {
		return
	}
	mid := b.m / 2
	upKey, upData := n.Key[mid], n.Data[mid]
	sp := n.split(mid)
	p := n.Parent
	if p == nil {
		p = new(Node)
		b.root = p
		n.Parent = p
		p.Children = append(p.Children, n)
	}
	switch {
	case b.comp(origKey, upKey) == 0:
		b.hot = p
	case b.comp(origKey, upKey) == 1:
		b.hot = sp
	}

	i := sort.Search(len(p.Key), func(i int) bool {
		result := b.comp(p.Key[i], upKey)
		return result == 1 || result == 0
	})

	p.Key = insert(p.Key, upKey, i)
	p.Data = insert(p.Data, upData, i)

	p.Children = append(p.Children, sp)
	if i < len(p.Children)-2 {
		copy(p.Children[i+2:], p.Children[i+1:])
		p.Children[i+1] = sp
	}
	sp.Parent = p
	b.solveOverflow(p, origKey)
}

func (b *BTree) Remove(key interface{}) (*Node, error) {
	n, i, ok := b.searchIn(b.root, key)
	if !ok {
		return nil, nil
	}
	if n.Children != nil {
		succ := n.Children[i+1]
		for len(succ.Children) > 0 {
			succ = succ.Children[0]
		}
		n.Key[i] = succ.Key[0]
		n.Data[i] = succ.Data[0]
		succ.Key = succ.Key[1:]
		succ.Data = succ.Data[1:]
		b.solveUnderflow(succ)
		return b.hot, nil
	}
	n.Key = append(n.Key[:i], n.Key[i+1:]...)
	n.Data = append(n.Data[:i], n.Data[i+1:]...)
	b.solveUnderflow(n)
	return b.hot, nil
}

func (b *BTree) solveUnderflow(n *Node) {
	b.hot = n
	bottom := int(math.Ceil(float64(b.m)/2)) - 1
	// fmt.Println(bottom)
	if len(n.Key) >= bottom {
		return
	}
	p := n.Parent
	if p == nil {
		if len(n.Key) == 0 && len(n.Children) > 0 {
			b.root = n.Children[0]
			b.root.Parent = nil
			n.Children = nil
		}
		return
	}
	i := 0
	for ; i < len(p.Children); i++ {
		if p.Children[i] == n {
			break
		}
	}
	// 有左兄弟，且不处在下溢临界点
	if i > 0 {
		ls := p.Children[i-1]
		if len(ls.Key) > bottom {
			// 向父节点借关键码
			n.Key = insert(n.Key, p.Key[i-1], 0)
			n.Data = insert(n.Data, p.Data[i-1], 0)
			// 用左兄弟中最大的关键码填充父节点中被借出的关键码
			p.Key[i-1] = ls.Key[len(ls.Key)-1]
			p.Data[i-1] = ls.Data[len(ls.Data)-1]
			// 左兄弟删除最大关键码
			ls.Key = ls.Key[:len(ls.Key)-1]
			ls.Data = ls.Data[:len(ls.Data)-1]
			// 过继原来左兄弟最大关键码的右孩子给 n，作为最左面的孩子
			if len(ls.Children) > 0 {
				n.Children = append([]*Node{ls.Children[len(ls.Children)-1]}, n.Children...)
				if n.Children[0] != nil {
					n.Children[0].Parent = n
				}
				// 删除左兄弟的最右孩子
				ls.Children = ls.Children[:len(ls.Children)-1]
			}
			return
		}
	}
	// 有右兄弟，且右兄弟不处于下溢临界点
	if i < len(p.Children)-1 {
		rs := p.Children[i+1]
		if len(rs.Key) > bottom {
			// 向父节点借关键码
			n.Key = insert(n.Key, p.Key[i], len(n.Key))
			n.Data = insert(n.Data, p.Data[i], len(n.Data))
			// 用右兄弟中最小的关键码填充父节点中被借出的关键码
			p.Key[i] = rs.Key[0]
			p.Data[i] = rs.Data[0]
			// 右兄弟删除最小关键码
			rs.Key = rs.Key[1:]
			rs.Data = rs.Data[1:]
			// 过继原来右兄弟最小关键码的左孩子给 n，作为最右面的孩子
			if len(rs.Children) > 0 {
				n.Children = append(n.Children, rs.Children[0])
				if n.Children[len(n.Children)-1] != nil {
					n.Children[len(n.Children)-1].Parent = n
				}
				// 删除右兄弟的最左孩子
				rs.Children = rs.Children[1:]
			}
			return
		}
	}
	// 左右兄弟要么不存在，要么都处于自身难保的情况
	// 与左兄弟合并
	if i > 0 {
		ls := p.Children[i-1]
		// 首先和父节点关键码合并
		ls.Key = append(ls.Key, p.Key[i-1])
		ls.Data = append(ls.Data, p.Data[i-1])
		// 删除父节点关键码
		if i < len(p.Key) {
			copy(p.Key[i-1:], p.Key[i:])
			p.Key[len(p.Key)-1] = nil
			copy(p.Data[i-1:], p.Data[i:])
			p.Data[len(p.Data)-1] = nil
		}
		p.Key = p.Key[:len(p.Key)-1]
		p.Data = p.Data[:len(p.Data)-1]
		// 删除父节点关键码指向 n 的连接
		if i < len(p.Children)-1 {
			copy(p.Children[i:], p.Children[i+1:])
			p.Children[len(p.Children)-1] = nil
		}
		p.Children = p.Children[:len(p.Children)-1]
		// 与左兄弟合并
		ls.Key = append(ls.Key, n.Key...)
		ls.Data = append(ls.Data, n.Data...)
		for _, child := range n.Children {
			if child != nil {
				child.Parent = ls
			}
		}
		ls.Children = append(ls.Children, n.Children...)
	} else {
		// 与右兄弟合并
		rs := p.Children[i+1]
		rs.Key = append([]interface{}{p.Key[i]}, rs.Key...)
		rs.Data = append([]interface{}{p.Data[i]}, rs.Data...)
		// 删除父节点关键码
		if len(p.Key) > i+1 {
			copy(p.Key[i:], p.Key[i+1:])
			p.Key[len(p.Key)-1] = nil
			copy(p.Data[i:], p.Data[i+1:])
			p.Data[len(p.Data)-1] = nil
		}
		p.Key = p.Key[:len(p.Key)-1]
		p.Data = p.Data[:len(p.Data)-1]
		// 删除父节点关键码指向 n 的连接
		copy(p.Children[i:], p.Children[i+1:])
		p.Children[len(p.Children)-1] = nil
		p.Children = p.Children[:len(p.Children)-1]
		// 与右兄弟合并
		rs.Key = append(n.Key, rs.Key...)
		rs.Data = append(n.Data, rs.Data...)
		for _, child := range n.Children {
			if child != nil {
				child.Parent = rs
			}
		}
		rs.Children = append(n.Children, rs.Children...)
	}
	b.solveUnderflow(p)
}

func (b *BTree) Print() {
	q := make([]*Node, 1)
	q[0] = b.root
	levelFirst := b.root
	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		if n.Parent == levelFirst {
			fmt.Println()
			levelFirst = n
		}
		// fmt.Printf("%p ", n)
		fmt.Printf("$")
		for i, v := range n.Key {
			fmt.Printf("%v", v)
			if i < len(n.Key)-1 {
				fmt.Printf(",")
			}
		}
		fmt.Printf("$ ")
		// fmt.Printf("->")
		for _, c := range n.Children {
			if c != nil {
				q = append(q, c)
				// fmt.Printf("%p ", c)
			}
		}
	}
	fmt.Println()
}

func (b *BTree) TravLevel(opts ...Option) {
	q := make([]*Node, 1)
	q[0] = b.root
	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		for _, c := range n.Children {
			if c != nil {
				q = append(q, c)
			}
		}
		for _, opt := range opts {
			opt(n)
		}
	}
}

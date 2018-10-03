package btree

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/mooncaker816/gostructure/bst"
)

func init() {
	bst.RegisterBST(bst.BTree, New)
}

type bTree struct {
	m    int // 阶数
	root *node
	comp bst.Comparator
	hot  *node
}

func New(parms ...interface{}) bst.BST {
	bt := new(bTree)
	bt.comp = bst.BasicCompare
	for _, p := range parms {
		switch v := p.(type) {
		case bst.Comparator:
			bt.comp = v
		case int:
			bt.m = v
		}
	}
	return bt
}

func (b *bTree) Root() bst.Node {
	return b.root
}

func (b *bTree) Search(key interface{}) (bst.Node, bool) {
	n, _, ok := b.searchIn(b.root, key)
	return n, ok
}

func (b *bTree) searchIn(n *node, key interface{}) (hot *node, i int, ok bool) {
	i = sort.Search(len(n.key), func(i int) bool {
		result := b.comp(n.key[i], key)
		return result == 1 || result == 0
	})

	// 在当前节点找到了key
	if i < len(n.key) && b.comp(n.key[i], key) == 0 {
		return n, i, true
	}
	// 当前节点没有 key，且为叶子节点，则查找失败
	if n.children == nil {
		return n, i, false
	}
	// 当前节点没有 key，且当前节点有孩子节点，继续查找
	return b.searchIn(n.children[i], key)
}

// Insert returns the exact node which stores the newly inserted key
func (b *bTree) Insert(key, data interface{}) (bst.Node, error) {
	if b.root == nil {
		b.root = newNode(key, data, b.m)
		return b.root, nil
	}
	n, i, ok := b.searchIn(b.root, key)
	if ok {
		return nil, errors.New("insert with duplicate key")
	}
	n.key = insert(n.key, key, i)
	n.data = insert(n.data, data, i)
	// n.children = append(n.children, nil)
	// copy(n.children[i+1:], n.children[i:])
	// n.children[i] = nil
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

func (b *bTree) solveOverflow(n *node, origKey interface{}) {
	b.hot = n
	if b.m >= len(n.key)+1 {
		return
	}
	mid := b.m / 2
	upKey, upData := n.key[mid], n.data[mid]
	sp := n.split(mid)
	p := n.parent
	if p == nil {
		p = new(node)
		b.root = p
		n.parent = p
		p.children = append(p.children, n)
	}
	switch {
	case b.comp(origKey, upKey) == 0:
		b.hot = p
	case b.comp(origKey, upKey) == 1:
		b.hot = sp
	}

	i := sort.Search(len(p.key), func(i int) bool {
		result := b.comp(p.key[i], upKey)
		return result == 1 || result == 0
	})

	p.key = insert(p.key, upKey, i)
	p.data = insert(p.data, upData, i)

	p.children = append(p.children, sp)
	if i < len(p.children)-2 {
		copy(p.children[i+2:], p.children[i+1:])
		p.children[i+1] = sp
	}
	sp.parent = p
	b.solveOverflow(p, origKey)
}

func (b *bTree) Remove(key interface{}) (bst.Node, error) {
	n, i, ok := b.searchIn(b.root, key)
	if !ok {
		return nil, nil
	}
	if n.children != nil {
		succ := n.children[i+1]
		for len(succ.children) > 0 {
			succ = succ.children[0]
		}
		n.key[i] = succ.key[0]
		n.data[i] = succ.data[0]
		succ.key = succ.key[1:]
		succ.data = succ.data[1:]
		b.solveUnderflow(succ)
		return b.hot, nil
	}
	n.key = append(n.key[:i], n.key[i+1:]...)
	n.data = append(n.data[:i], n.data[i+1:]...)
	b.solveUnderflow(n)
	return b.hot, nil
}

func (b *bTree) solveUnderflow(n *node) {
	b.hot = n
	bottom := int(math.Ceil(float64(b.m)/2)) - 1
	// fmt.Println(bottom)
	if len(n.key) >= bottom {
		return
	}
	p := n.parent
	if p == nil {
		if len(n.key) == 0 && len(n.children) > 0 {
			b.root = n.children[0]
			b.root.parent = nil
			n.children = nil
		}
		return
	}
	i := 0
	for ; i < len(p.children); i++ {
		if p.children[i] == n {
			break
		}
	}
	// 有左兄弟，且不处在下溢临界点
	if i > 0 {
		ls := p.children[i-1]
		if len(ls.key) > bottom {
			// 向父节点借关键码
			n.key = insert(n.key, p.key[i-1], 0)
			n.data = insert(n.data, p.data[i-1], 0)
			// 用左兄弟中最大的关键码填充父节点中被借出的关键码
			p.key[i-1] = ls.key[len(ls.key)-1]
			p.data[i-1] = ls.data[len(ls.data)-1]
			// 左兄弟删除最大关键码
			ls.key = ls.key[:len(ls.key)-1]
			ls.data = ls.data[:len(ls.data)-1]
			// 过继原来左兄弟最大关键码的右孩子给 n，作为最左面的孩子
			if len(ls.children) > 0 {
				n.children = append([]*node{ls.children[len(ls.children)-1]}, n.children...)
				if n.children[0] != nil {
					n.children[0].parent = n
				}
				// 删除左兄弟的最右孩子
				ls.children = ls.children[:len(ls.children)-1]
			}
			return
		}
	}
	// 有右兄弟，且右兄弟不处于下溢临界点
	if i < len(p.children)-1 {
		rs := p.children[i+1]
		if len(rs.key) > bottom {
			// 向父节点借关键码
			n.key = insert(n.key, p.key[i], len(n.key))
			n.data = insert(n.data, p.data[i], len(n.data))
			// 用右兄弟中最小的关键码填充父节点中被借出的关键码
			p.key[i] = rs.key[0]
			p.data[i] = rs.data[0]
			// 右兄弟删除最小关键码
			rs.key = rs.key[1:]
			rs.data = rs.data[1:]
			// 过继原来右兄弟最小关键码的左孩子给 n，作为最右面的孩子
			if len(rs.children) > 0 {
				n.children = append(n.children, rs.children[0])
				if n.children[len(n.children)-1] != nil {
					n.children[len(n.children)-1].parent = n
				}
				// 删除右兄弟的最左孩子
				rs.children = rs.children[1:]
			}
			return
		}
	}
	// 左右兄弟要么不存在，要么都处于自身难保的情况
	// 与左兄弟合并
	if i > 0 {
		ls := p.children[i-1]
		// 首先和父节点关键码合并
		ls.key = append(ls.key, p.key[i-1])
		ls.data = append(ls.data, p.data[i-1])
		// 删除父节点关键码
		if i < len(p.key) {
			copy(p.key[i-1:], p.key[i:])
			p.key[len(p.key)-1] = nil
			copy(p.data[i-1:], p.data[i:])
			p.data[len(p.data)-1] = nil
		}
		p.key = p.key[:len(p.key)-1]
		p.data = p.data[:len(p.data)-1]
		// 删除父节点关键码指向 n 的连接
		if i < len(p.children)-1 {
			copy(p.children[i:], p.children[i+1:])
			p.children[len(p.children)-1] = nil
		}
		p.children = p.children[:len(p.children)-1]
		// 与左兄弟合并
		ls.key = append(ls.key, n.key...)
		ls.data = append(ls.data, n.data...)
		for _, child := range n.children {
			if child != nil {
				child.parent = ls
			}
		}
		ls.children = append(ls.children, n.children...)
	} else {
		// 与右兄弟合并
		rs := p.children[i+1]
		rs.key = append([]interface{}{p.key[i]}, rs.key...)
		rs.data = append([]interface{}{p.data[i]}, rs.data...)
		// 删除父节点关键码
		if len(p.key) > i+1 {
			copy(p.key[i:], p.key[i+1:])
			p.key[len(p.key)-1] = nil
			copy(p.data[i:], p.data[i+1:])
			p.data[len(p.data)-1] = nil
		}
		p.key = p.key[:len(p.key)-1]
		p.data = p.data[:len(p.data)-1]
		// 删除父节点关键码指向 n 的连接
		copy(p.children[i:], p.children[i+1:])
		p.children[len(p.children)-1] = nil
		p.children = p.children[:len(p.children)-1]
		// 与右兄弟合并
		rs.key = append(n.key, rs.key...)
		rs.data = append(n.data, rs.data...)
		for _, child := range n.children {
			if child != nil {
				child.parent = rs
			}
		}
		rs.children = append(n.children, rs.children...)
	}
	b.solveUnderflow(p)
}

func (b *bTree) Print() {
	q := make([]*node, 1)
	q[0] = b.root
	levelFirst := b.root
	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		if n.parent == levelFirst {
			fmt.Println()
			levelFirst = n
		}
		// fmt.Printf("%p ", n)
		fmt.Printf("$")
		for i, v := range n.key {
			fmt.Printf("%v", v)
			if i < len(n.key)-1 {
				fmt.Printf(",")
			}
		}
		fmt.Printf("$ ")
		// fmt.Printf("->")
		for _, c := range n.children {
			if c != nil {
				q = append(q, c)
				// fmt.Printf("%p ", c)
			}
		}
	}
	fmt.Println()
}

// LevelOrder only
func (b *bTree) Walk(o bst.Order, opts ...bst.Option) {
	if o != bst.LevelOrder {
		panic("unsupported walk order for B-Tree")
	}
	q := make([]*node, 1)
	q[0] = b.root
	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		for _, c := range n.children {
			if c != nil {
				q = append(q, c)
			}
		}
		for _, opt := range opts {
			opt(n)
		}
	}
}

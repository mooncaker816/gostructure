package bst

import "reflect"

// Node identifies a BST Node
type Node interface {
	LChild() Node
	RChild() Node
	Parent() Node
	SetLChild(Node)
	SetRChild(Node)
	SetParent(Node)
	Key() interface{}
	Data() interface{}
	SetKey(interface{})
	SetData(interface{})
	Height() int
	Color() string
}

func IsNil(n Node) bool {
	if n == nil {
		return true
	}
	return reflect.ValueOf(n).IsNil()
}

// IsRoot returns whether n is root node
func IsRoot(n Node) bool { return IsNil(n.Parent()) }

// IsLChild returns whether n is a left child node
func IsLChild(n Node) bool { return !IsRoot(n) && n == n.Parent().LChild() }

// IsRChild returns whether n is a right child node
func IsRChild(n Node) bool { return !IsRoot(n) && n == n.Parent().RChild() }

// HasLChild returns whether n has left child
func HasLChild(n Node) bool { return !IsNil(n.LChild()) }

// HasRChild returns whether n has right child
func HasRChild(n Node) bool { return !IsNil(n.RChild()) }

// IsLeaf returns whether n is leaf node
func IsLeaf(n Node) bool {
	if HasLChild(n) || HasRChild(n) {
		return false
	}
	return true
}

// Sibling returns sibling of n
func Sibling(n Node) Node {
	if IsRoot(n) {
		return nil
	}
	if IsLChild(n) {
		return n.Parent().RChild()
	}
	return n.Parent().LChild()
}

// AttachRChild connects rc as right child of n
func AttachRChild(n, rc Node) {
	n.SetRChild(rc)
	rc.SetParent(n)
}

// AttachLChild connects lc as left child of n
func AttachLChild(n, lc Node) {
	n.SetLChild(lc)
	lc.SetParent(n)
}

// Size returns the total node counts of the subtree rooted on n
func Size(n Node) int {
	if IsNil(n) {
		return 0
	}
	count := 1
	if HasLChild(n) {
		count += Size(n.LChild())
	}
	if HasRChild(n) {
		count += Size(n.RChild())
	}
	return count
}

// Level returns the level where the node lays on
func Level(n Node) int {
	l := 0
	for !IsNil(n) && !IsRoot(n) {
		l++
		n = n.Parent()
	}
	return l
}

// Successor returns the successor of n
func Successor(n Node) Node {
	if HasRChild(n) {
		return subTreeMin(n.RChild())
	}
	for IsRChild(n) {
		n = n.Parent()
	}
	return n.Parent()
}

// Predecessor returns the predecessor of n
func Predecessor(n Node) Node {
	if HasLChild(n) {
		return subTreeMax(n.LChild())
	}
	for IsLChild(n) {
		n = n.Parent()
	}
	return n.Parent()
}

func subTreeMin(n Node) Node {
	for HasLChild(n) {
		n = n.LChild()
	}
	return n
}

func subTreeMax(n Node) Node {
	for HasRChild(n) {
		n = n.RChild()
	}
	return n
}

// TravPre walks the subtree rooted on n by pre-order
func TravPre(n Node, opts ...Option) {
	for _, opt := range opts {
		opt(n)
	}
	if HasLChild(n) {
		TravPre(n.LChild(), opts...)
	}
	if HasRChild(n) {
		TravPre(n.RChild(), opts...)
	}
}

// TravIn walks the subtree rooted on n by in-order
func TravIn(n Node, opts ...Option) {
	if HasLChild(n) {
		TravIn(n.LChild(), opts...)
	}
	for _, opt := range opts {
		opt(n)
	}
	if HasRChild(n) {
		TravIn(n.RChild(), opts...)
	}
}

// TravPost walks the subtree rooted on n by post-order
func TravPost(n Node, opts ...Option) {
	if HasLChild(n) {
		TravIn(n.LChild(), opts...)
	}
	if HasRChild(n) {
		TravIn(n.RChild(), opts...)
	}
	for _, opt := range opts {
		opt(n)
	}
}

// TravLevel walks the subtree rooted on n by level-order
func TravLevel(n Node, opts ...Option) {
	queue := make([]Node, 0, Size(n))
	queue = append(queue, n)
	for len(queue) > 0 {
		visitNode := queue[0]
		queue = queue[1:]
		if HasLChild(visitNode) {
			queue = append(queue, visitNode.LChild())
		}
		if HasRChild(visitNode) {
			queue = append(queue, visitNode.RChild())
		}
		for _, opt := range opts {
			opt(visitNode)
		}
	}
}

// RotateAt use connect 3+4 strategy to reconstruct v,p,g which are all existing
func RotateAt(v Node, opts ...Option) (a, b, c Node) {
	p := v.Parent()
	g := p.Parent()
	if IsLChild(v) {
		if IsLChild(p) {
			p.SetParent(g.Parent())
			connect34(v, p, g, v.LChild(), v.RChild(), p.RChild(), g.RChild(), opts...)
			return v, p, g
		}
		if IsRChild(p) {
			v.SetParent(g.Parent())
			connect34(g, v, p, g.LChild(), v.LChild(), v.RChild(), p.RChild(), opts...)
			return g, v, p
		}
	}
	if IsRChild(v) {
		if IsLChild(p) {
			v.SetParent(g.Parent())
			connect34(p, v, g, p.LChild(), v.LChild(), v.RChild(), g.RChild(), opts...)
			return p, v, g
		}
		if IsRChild(p) {
			p.SetParent(g.Parent())
			connect34(g, p, v, g.LChild(), p.LChild(), v.LChild(), v.RChild(), opts...)
			return g, p, v
		}
	}
	return nil, nil, nil
}

// connect34 connect bst.Nodes as below
//	 	   b
//		a	  c
//	  T1 T2 T3 T4
func connect34(a, b, c, t1, t2, t3, t4 Node, opts ...Option) {
	AttachLChild(a, t1)
	AttachRChild(a, t2)
	for _, o := range opts {
		o(a)
	}

	AttachLChild(c, t3)
	AttachRChild(c, t4)
	for _, o := range opts {
		o(c)
	}

	AttachLChild(b, a)
	AttachRChild(b, c)
	for _, o := range opts {
		o(b)
	}
}

// RemoveAt removes "node n" and returns the exact removed node's parent as hot and the replacement node as r
func RemoveAt(n, root Node) (hot, r Node) {
	// n has both left and right subtree
	if HasLChild(n) && HasRChild(n) {
		succ := Successor(n)
		swapKeyData(n, succ)
		hot = succ.Parent()
		r = succ.RChild()
		if succ == n.RChild() {
			hot = n
			if HasRChild(succ) {
				succ.RChild().SetParent(n)
			}
			n.SetRChild(succ.RChild())
		} else if HasRChild(succ) {
			// hot = succ.parent
			succ.RChild().SetParent(hot)
			hot.SetLChild(succ.RChild())
		} else {
			// hot = succ.parent
			hot.SetLChild(nil)
		}
		release(succ)
		return hot, r
	}
	// n only has left subtree
	if HasLChild(n) && !HasRChild(n) {
		if IsLChild(n) {
			n.Parent().SetLChild(n.LChild())
		} else if IsRChild(n) {
			n.Parent().SetRChild(n.LChild())
		} else {
			root = n.LChild()
		}
		n.LChild().SetParent(n.Parent())
		hot = n.Parent()
		r = n.LChild()
		release(n)
		return hot, r
	}
	// n only has right subtree
	if !HasLChild(n) && HasRChild(n) {
		if IsLChild(n) {
			n.Parent().SetLChild(n.RChild())
		} else if IsRChild(n) {
			n.Parent().SetRChild(n.RChild())
		} else {
			root = n.RChild()
		}
		n.RChild().SetParent(n.Parent())
		hot = n.Parent()
		r = n.RChild()
		release(n)
		return hot, r
	}
	// n is leaf(or single root)
	if IsLChild(n) {
		n.Parent().SetLChild(nil)
	} else if IsRChild(n) {
		n.Parent().SetRChild(nil)
	} else {
		root = nil
	}
	hot = n.Parent()
	release(n)
	return hot, r
}

func swapKeyData(n1, n2 Node) {
	tmpKey, tmpData := n1.Key(), n1.Data()
	n1.SetKey(n2.Key())
	n1.SetData(n2.Data())
	n2.SetKey(tmpKey)
	n2.SetData(tmpData)
}

func release(n Node) {
	n.SetParent(nil)
	n.SetLChild(nil)
	n.SetRChild(nil)
}

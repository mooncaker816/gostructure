package suffixtree

import (
	"fmt"
	"sort"
)

// !!!!!every substring of s is a preﬁx of some sufﬁx of s

var text string

type Node struct {
	children   [256]*Node
	start      int
	end        *end
	index      int // suffix index of the original text, only for leaf node, internal node will be -1
	depth      int
	suffixLink *Node
	morePos    [][2]int
}

func (n *Node) String() string {
	if n.start == -1 {
		return ""
	}
	return text[n.start : n.end.end+1]
}

func (n *Node) edgeLen() int {
	return n.end.end - n.start + 1
}

type active struct {
	node   *Node
	edge   int
	length int
}

func (a *active) String() string {
	return fmt.Sprintf("Active [activeNode=%v, activeIndex=%d, activeLength=%d]", a.node, a.edge, a.length)
}

type SufTree struct {
	Root      *Node
	active    *active
	end       *end
	remaining int
}

func NewSufTree(s string) *SufTree {
	text = s + "$"
	st := new(SufTree)
	st.init()
	st.build(text)
	return st
}

func (st *SufTree) init() {
	st.newRoot()
	st.end = newEnd(-1) // 公共 end
	st.active = &active{node: st.Root, edge: -1, length: 0}
}

func (st *SufTree) newRoot() {
	root := st.newNode(-1, newEnd(-1))
	st.Root = root
	st.Root.suffixLink = root
}

func (st *SufTree) newNode(start int, end *end) *Node {
	n := new(Node)
	n.start = start
	n.end = end
	n.index = -1
	n.suffixLink = st.Root
	return n
}

func (st *SufTree) build(s string) {
	for i := 0; i < len(s); i++ {
		st.extend(i)
	}
	st.updateSuffixIndex()
}

// return skipped internal node count
func (st *SufTree) walkDown(n *Node) int {
	if l := n.edgeLen(); st.active.length >= l {
		st.active.edge += l
		st.active.length -= l
		st.active.node = n
		return 1
	}
	return 0
}

func (st *SufTree) extend(i int) {
	var prevNewNode *Node
	st.end.set(i)
	st.remaining++
	for st.remaining > 0 {
		// active.length 为 0 ，说明我们需要添加的字符是一个新的分支的首字符，结合 walkDown 看
		if st.active.length == 0 {
			st.active.edge = i
		}
		if n := st.tryEdge(); n == nil {
			// activeNode 的分支中没有以新增字符为起始的分支
			st.newEdge(i)
			if prevNewNode != nil {
				prevNewNode.suffixLink = st.active.node
				prevNewNode = nil
			}
		} else {
			// 当前 activeNode 的分支中有以新增字符为起始的分支
			// 此时需要判断 activeLength 是否能被该分支覆盖，
			// 如果超出了该分支的范围，我们就需要更新该分支对应的 node 为新的 activeNode，
			// 并更新相对的 acitveLength 和 activeEdge，再继续循环直到找到对应的 activeLength 位置的字符
			if st.walkDown(n) > 0 {
				n.morePos = append(n.morePos, [2]int{i - n.edgeLen(), i - 1})
				continue
			}
			// 如果已经存在，active.length 加一,继续下一 phase
			if text[n.start+st.active.length] == text[i] {
				if prevNewNode != nil && st.active.node != st.Root {
					prevNewNode.suffixLink = st.active.node
					prevNewNode = nil
				}
				st.active.length++
				break
			}
			// 不存在，在原有分支的基础上新开分支
			split := st.splitEdgeOn(n, i)
			if prevNewNode != nil {
				prevNewNode.suffixLink = split
			}
			prevNewNode = split
		}
		// 完成了一个 suffix 字符的添加
		st.remaining--
		if st.active.node == st.Root {
			if st.active.length > 0 {
				st.active.length--
				st.active.edge = i - st.remaining + 1
			}
		} else {
			st.active.node = st.active.node.suffixLink
		}
	}
}

func (st *SufTree) tryEdge() *Node {
	return st.active.node.children[text[st.active.edge]]
}

func (st *SufTree) newEdge(i int) {
	st.active.node.children[text[st.active.edge]] = st.newNode(i, st.end)
}

func (st *SufTree) setEdge(n *Node) {
	st.active.node.children[text[st.active.edge]] = n
}

func (st *SufTree) splitEdgeOn(n *Node, i int) *Node {
	// 缩短原来的 node 为 split，并将它替代 activeNode 中对应的原来更长的 child
	split := st.newNode(n.start, newEnd(n.start+st.active.length-1))
	split.morePos = append(split.morePos, [2]int{i - split.edgeLen(), i - 1})
	st.setEdge(split)

	// 后半部分作为 split 的分支接入
	n.start += st.active.length
	split.children[text[n.start]] = n

	// 新建分支
	split.children[text[i]] = st.newNode(i, st.end)

	return split
}

func (st *SufTree) updateSuffixIndex() {
	depth := 0
	st.Root.dfsUpdateIndex(depth)
}

func (n *Node) dfsUpdateIndex(depth int) {
	leaf := true
	n.depth = depth
	for i := 0; i < 256; i++ {
		if child := n.children[i]; child != nil {
			leaf = false
			child.dfsUpdateIndex(depth + child.edgeLen())
		}
	}
	if leaf {
		n.index = len(text) - depth
	}
}

// Print all the suffixes with index
func (st *SufTree) Print() {
	prefix := ""
	st.Root.dfsPrint(prefix)
}

func (n *Node) dfsPrint(prefix string) {
	leaf := true
	for i := 0; i < 256; i++ {
		if child := n.children[i]; child != nil {
			leaf = false
			child.dfsPrint(fmt.Sprintf("%s%s", prefix, child.String()))
		}
	}
	if leaf {
		fmt.Printf("[%d]:%s\n", n.index, prefix)
	}
}

func testPrint(n *Node) {
	leaf := true
	if n.start != -1 {
		fmt.Printf("%s", text[n.start:n.end.end+1])
	}
	for i := 0; i < 256; i++ {
		if n.children[i] != nil {
			if leaf && n.start != -1 {
				fmt.Printf("[%d]\n", n.index)
				fmt.Printf("start:%d,end:%d,depth:%d\n", n.start, n.end.end, n.depth)
			}
			leaf = false
			testPrint(n.children[i])
		}
	}
	if leaf {
		fmt.Printf("[%d]\n", n.index)
		fmt.Printf("start:%d,end:%d,depth:%d\n", n.start, n.end.end, n.depth)
	}
}

type end struct {
	end int
}

func (e *end) set(i int) {
	e.end = i
}

func newEnd(x int) *end {
	e := new(end)
	e.end = x
	return e
}

// Index returns the index of the substr if exists, otherwise returns -1
func (st *SufTree) Index(sub string) (idx int) {
	if sub == "" {
		return 0
	}
	if n, _ := st.walkPath(sub); n != nil {
		leaves := make([]*Node, 0)
		n.getLeaves(&leaves)
		sort.Slice(leaves, func(i, j int) bool {
			return leaves[i].index < leaves[j].index
		})
		return leaves[0].index
	}
	return -1
}

func (n *Node) getLeaves(leaves *[]*Node) {
	leaf := true
	for i := 0; i < 256; i++ {
		if child := n.children[i]; child != nil {
			leaf = false
			child.getLeaves(leaves)
		}
	}
	if leaf {
		*leaves = append(*leaves, n)
	}
}

func (n *Node) getLeavesCnt() int {
	cnt := 0
	leaf := true
	for i := 0; i < 256; i++ {
		if child := n.children[i]; child != nil {
			leaf = false
			cnt += child.getLeavesCnt()
		}
	}
	if leaf {
		return 1
	}
	return cnt
}

// IndexAny returns any index of sub if exists, otherwise returns -1.
func (st *SufTree) IndexAny(sub string) (idx int) {
	if sub == "" {
		return 0
	}
	if n, _ := st.walkPath(sub); n != nil {
		return n.getLeastLeaf().index
	}
	return -1
}

func (n *Node) getLeastLeaf() *Node {
	var l *Node
	for i := 0; i < 256; i++ {
		if child := n.children[i]; child != nil {
			if l = child.getLeastLeaf(); l != nil {
				break
			}
		}
	}
	return l
}

// walk the tree as per the path of p. if p exists, returns the last node of the path, else return nil
func (st *SufTree) walkPath(p string) (curr *Node, charCnt int) {
	if p == "" {
		return st.Root, 0
	}
	curr = st.Root.children[p[0]]
	for curr != nil {
		i, j := 0, curr.start
		for ; i < len(p) && j <= curr.end.end; i, j = i+1, j+1 {
			if p[i] != text[j] {
				return nil, charCnt
			}
			charCnt++
		}
		if i == len(p) {
			return curr, charCnt
		}
		if j == curr.end.end+1 {
			p = p[curr.edgeLen():]
			curr = curr.children[p[0]]
		}
	}
	return nil, -1
}

// Occurs returns the total count of p appears in text
func (st *SufTree) Occurs(p string) int {
	if n, _ := st.walkPath(p); n != nil {
		return n.getLeavesCnt()
	}
	return -1
}

// func (n *Node) getCharCntBelow() int {
// 	cnt := n.edgeLen()
// 	for i := 0; i < 256; i++ {
// 		if child := n.children[i]; child != nil {
// 			cnt += child.getCharCntBelow()
// 		}
// 	}
// 	return cnt
// }

func (n *Node) isLeaf() bool {
	return n.index != -1
}

// LongestRepeat returns the longest repeated substr in text
func (st *SufTree) LongestRepeat() string {
	deepest := new(Node)
	st.Root.getDeepestNonLeaf(deepest)
	if deepest != nil {
		for _, child := range deepest.children {
			if child != nil && child.isLeaf() {
				return text[child.index:][:deepest.depth]
			}
		}
	}
	return ""
}

func (n *Node) getDeepestNonLeaf(deepest *Node) {
	for _, child := range n.children {
		if child != nil {
			if child.isLeaf() {
				continue
			}
			if child.depth > deepest.depth {
				*deepest = *child
			}
			child.getDeepestNonLeaf(deepest)
		}
	}
	return
}

// LCS returns the longest common substr of text and pattern
func (st *SufTree) LCS(p string) string {
	if p == "" {
		return ""
	}
	lcs := ""
	currStr := ""
	curr := st.Root.children[p[0]]
	for len(p) > 0 {
		if curr.children[p[0]] == nil {
			if len(lcs) < len(currStr) {
				lcs = currStr
			}
			p = p[1:]
			curr = curr.suffixLink
		}
		i, j := 0, curr.start
		for ; i < len(p) && j <= curr.end.end; i, j = i+1, j+1 {
			if p[i] != text[j] {
				if len(lcs) < len(currStr) {
					lcs = currStr
				}
				for skip := 0; skip < i; skip++ {
					curr = curr.suffixLink
				}
				p = p[i:]
				currStr = ""
				continue
			}
		}
		if i == len(p) {
			currStr += p
		}
		if j == curr.end.end+1 {
			currStr += p[curr.edgeLen():]
			p = p[curr.edgeLen():]
			curr = curr.children[p[0]]
		}
	}
	return lcs
}

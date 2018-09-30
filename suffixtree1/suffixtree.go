package suftree

import "fmt"

type SufTree struct {
	Root    *Node
	Longest *Node
	Size    int
}

type Node struct {
	children   [26]*Node
	parent     *Node
	suffixLink *Node
}

func (n *Node) addChild(c byte, ln *Node) {
	n.children[indexOf(c)] = ln
}

func indexOf(c byte) int {
	return int(c - 'a')
}

func BuildSufTree(s string) *SufTree {
	// explicitly build the two-node suffix tree
	root := newRoot()
	longest := new(Node)
	longest.suffixLink = root
	longest.parent = root
	root.addChild(s[0], longest)
	// fmt.Printf("root:%p first:%p\n", root, longest)

	for i := 1; i < len(s); i++ {
		c := s[i]
		// fmt.Println(c)
		var curr, prev *Node
		curr = longest
		for curr.children[indexOf(c)] == nil {
			// create new node r1 with transition Current -c->r1
			r1 := new(Node)
			r1.parent = curr
			curr.addChild(c, r1)
			// fmt.Printf("%p %+v\n", r1, r1)
			// if we came from some previous node, make that node's suffix link point here
			if prev != nil {
				prev.suffixLink = r1
				// fmt.Printf("prev: %p %v\n", prev, prev)
			}
			// walk down the suffix links
			prev = r1
			curr = curr.suffixLink
		}
		// make the last suffix link
		if curr == root && curr.children[indexOf(c)] == prev {
			prev.suffixLink = root
		} else {
			prev.suffixLink = curr.children[indexOf(c)]
		}
		// fmt.Printf("prev: %p %v\n", prev, prev)

		longest = longest.children[indexOf(c)]
		// fmt.Printf("longest: %p\n", longest)
	}
	return &SufTree{Root: root, Longest: longest, Size: len(s)}
}

func newRoot() *Node {
	root := new(Node)
	root.suffixLink = root
	return root
}

func (s *SufTree) Print() {
	curr := s.Longest
	currSize := s.Size
	for currSize > 0 {
		b := make([]byte, currSize)
		n := curr
		for i := currSize - 1; i >= 0; i-- {
			b[i] = n.toChar()
			n = n.parent
		}
		fmt.Printf("%s\n", b)
		currSize--
		curr = curr.suffixLink
	}
}

func (n *Node) toChar() byte {
	if n.parent == nil {
		return 0
	}
	for i, child := range n.parent.children {
		if child == n {
			return 'a' + byte(i)
		}
	}
	return 0
}

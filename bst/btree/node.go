package btree

import "github.com/mooncaker816/gostructure/bst"

type node struct {
	parent   *node
	children []*node       // 分支
	key      []interface{} // 关键码
	data     []interface{} // 数据
}

func newNode(key, data interface{}, m int) *node {
	n := new(node)
	n.key = make([]interface{}, 0, m-1)
	n.data = make([]interface{}, 0, m-1)
	return n
}

func (n *node) Key() interface{}  { return n.key }
func (n *node) Data() interface{} { return n.data }
func (n *node) Height() int       { return 0 }
func (n *node) LChild() bst.Node  { return nil }
func (n *node) RChild() bst.Node  { return nil }
func (n *node) Parent() bst.Node  { return nil }
func (n *node) Color() string     { return "" }

func (n *node) SetKey(key interface{})   {}
func (n *node) SetData(data interface{}) {}
func (n *node) SetLChild(lc bst.Node)    {}
func (n *node) SetRChild(rc bst.Node)    {}
func (n *node) SetParent(p bst.Node)     {}

func (n *node) split(i int) *node {
	sp := new(node)
	sp.key = append(sp.key, n.key[i+1:]...)
	sp.data = append(sp.data, n.data[i+1:]...)
	n.key = n.key[:i]
	n.data = n.data[:i]
	if n.children != nil {
		sp.children = append(sp.children, n.children[i+1:]...)
		for _, child := range sp.children {
			if child != nil {
				child.parent = sp
			}
		}
		n.children = n.children[:i+1]
	}
	return sp
}

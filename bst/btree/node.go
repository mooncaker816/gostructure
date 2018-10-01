package btree

type Node struct {
	Parent   *Node
	Children []*Node       // 分支
	Key      []interface{} // 关键码
	Data     []interface{} // 数据
	// Height   int
}

func newNode(key, data interface{}, m int) *Node {
	n := new(Node)
	n.Key = make([]interface{}, 0, m-1)
	n.Data = make([]interface{}, 0, m-1)
	// n.Children = make([]*Node, 0, m)
	return n
}

func (n *Node) split(i int) *Node {
	sp := new(Node)
	sp.Key = append(sp.Key, n.Key[i+1:]...)
	sp.Data = append(sp.Data, n.Data[i+1:]...)
	n.Key = n.Key[:i]
	n.Data = n.Data[:i]
	if n.Children != nil {
		sp.Children = append(sp.Children, n.Children[i+1:]...)
		for _, child := range sp.Children {
			if child != nil {
				child.Parent = sp
			}
		}
		n.Children = n.Children[:i+1]
	}
	return sp
}

type Option func(*Node)

package bintree

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/mooncaker816/gostructure/queue"
)

func TestInsertAsRoot(t *testing.T) {
	q := queue.NewQueue()
	bt := new(BinTree)
	root := bt.InsertAsRoot(0, nil)
	q.Enqueue(root)
	key := 1
	for l := 0; l < 4; l++ {
		for i := 0; i < 1<<uint(l); i++ {
			ni, _ := q.Dequeue()
			n := ni.(*Node)
			q.Enqueue(bt.InsertAsLChild(n, key, nil))
			key++
			q.Enqueue(bt.InsertAsRChild(n, key, nil))

			key++
		}
	}
	if bt.Root.Height != 4 {
		t.Errorf("Got %v expected %v", bt.Root.Height, 4)
	}
	if bt.Size != 1<<uint(bt.Root.Height+1)-1 {
		t.Errorf("Got %v expected %v", bt.Size, 1<<uint(bt.Root.Height+1)-1)
	}
	bt.Print()
	bt2 := new(BinTree)
	root2 := bt2.InsertAsRoot("k", nil)
	i := bt2.InsertAsLChild(root2, "i", nil)
	bt2.InsertAsRChild(root2, "j", nil)
	h := bt2.InsertAsRChild(i, "h", nil)
	b := bt2.InsertAsLChild(h, "b", nil)
	g := bt2.InsertAsRChild(h, "g", nil)
	bt2.InsertAsRChild(b, "a", nil)
	e := bt2.InsertAsLChild(g, "e", nil)
	bt2.InsertAsRChild(g, "f", nil)
	bt2.InsertAsLChild(e, "c", nil)
	bt2.InsertAsRChild(e, "d", nil)
	if bt2.Root.Height != 5 {
		t.Errorf("Got %v expected %v", bt2.Root.Height, 5)
	}
	if bt2.Size != 11 {
		t.Errorf("Got %v expected %v", bt2.Size, 11)
	}
	bt2.Print()
}

func TestIsRoot(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if !a.IsRoot() {
		t.Errorf("Got %v expected %v", a.IsRoot(), true)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if b.IsRoot() {
		t.Errorf("Got %v expected %v", b.IsRoot(), false)
	}
}

func TestIsLChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.IsLChild() {
		t.Errorf("Got %v expected %v", a.IsLChild(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if !b.IsLChild() {
		t.Errorf("Got %v expected %v", b.IsLChild(), true)
	}
	c := bt.InsertAsRChild(a, 3, nil)
	if c.IsLChild() {
		t.Errorf("Got %v expected %v", c.IsLChild(), false)
	}
}

func TestIsRChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.IsRChild() {
		t.Errorf("Got %v expected %v", a.IsRChild(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if b.IsRChild() {
		t.Errorf("Got %v expected %v", b.IsRChild(), false)
	}
	c := bt.InsertAsRChild(a, 3, nil)
	if !c.IsRChild() {
		t.Errorf("Got %v expected %v", c.IsRChild(), true)
	}
}

func TestHasParent(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasParent() {
		t.Errorf("Got %v expected %v", a.HasParent(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if !b.HasParent() {
		t.Errorf("Got %v expected %v", b.HasParent(), true)
	}
}

func TestHasLChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasLChild() {
		t.Errorf("Got %v expected %v", a.HasLChild(), false)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	if !a.HasLChild() {
		t.Errorf("Got %v expected %v", a.HasLChild(), true)
	}
	if b.HasLChild() {
		t.Errorf("Got %v expected %v", b.HasLChild(), false)
	}
}

func TestHasRChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasRChild() {
		t.Errorf("Got %v expected %v", a.HasRChild(), false)
	}
	b := bt.InsertAsRChild(a, 2, nil)
	if !a.HasRChild() {
		t.Errorf("Got %v expected %v", a.HasRChild(), true)
	}
	if b.HasRChild() {
		t.Errorf("Got %v expected %v", b.HasRChild(), false)
	}
}

func TestHasChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasChild() {
		t.Errorf("Got %v expected %v", a.HasChild(), false)
	}
	b := bt.InsertAsRChild(a, 2, nil)
	if !a.HasChild() {
		t.Errorf("Got %v expected %v", a.HasChild(), true)
	}
	if b.HasChild() {
		t.Errorf("Got %v expected %v", b.HasChild(), false)
	}
}

func TestHasBothChild(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.HasBothChild() {
		t.Errorf("Got %v expected %v", a.HasBothChild(), false)
	}
	bt.InsertAsRChild(a, 2, nil)
	if a.HasBothChild() {
		t.Errorf("Got %v expected %v", a.HasBothChild(), false)
	}
	bt.InsertAsLChild(a, 3, nil)
	if !a.HasBothChild() {
		t.Errorf("Got %v expected %v", a.HasBothChild(), true)
	}
}

func TestIsLeaf(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if !a.IsLeaf() {
		t.Errorf("Got %v expected %v", a.IsLeaf(), true)
	}
	b := bt.InsertAsRChild(a, 2, nil)
	if !b.IsLeaf() {
		t.Errorf("Got %v expected %v", b.IsLeaf(), true)
	}
	if a.IsLeaf() {
		t.Errorf("Got %v expected %v", a.IsLeaf(), false)
	}
}

func TestSibling(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	if a.Sibling() != nil {
		t.Errorf("Got %v expected %v", a.Sibling(), nil)
	}
	b := bt.InsertAsLChild(a, 2, nil)
	c := bt.InsertAsRChild(a, 3, nil)
	if b.Sibling() == nil || b.Sibling() != c {
		t.Errorf("Got %v expected %v", b.Sibling(), c)
	}
	if c.Sibling() == nil || c.Sibling() != b {
		t.Errorf("Got %v expected %v", c.Sibling(), b)
	}
}

func TestUncle(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	b := bt.InsertAsLChild(a, 2, nil)
	c := bt.InsertAsRChild(a, 3, nil)
	d := bt.InsertAsLChild(b, 4, nil)
	e := bt.InsertAsRChild(b, 5, nil)
	f := bt.InsertAsLChild(c, 6, nil)
	g := bt.InsertAsRChild(c, 7, nil)
	if d.Uncle() == nil || d.Uncle() != c || d.Uncle() != e.Uncle() {
		t.Errorf("Got %v expected %v", d.Uncle(), c)
	}
	if f.Uncle() == nil || f.Uncle() != b || f.Uncle() != g.Uncle() {
		t.Errorf("Got %v expected %v", f.Uncle(), b)
	}
}

func TestSize(t *testing.T) {
	bt := new(BinTree)
	a := bt.InsertAsRoot(1, nil)
	b := bt.InsertAsLChild(a, 2, nil)
	c := bt.InsertAsRChild(a, 3, nil)
	d := bt.InsertAsLChild(b, 4, nil)
	f := bt.InsertAsLChild(c, 6, nil)
	g := bt.InsertAsRChild(c, 7, nil)
	if a.Size() != 6 {
		t.Errorf("Got %v expected %v", a.Size(), 6)
	}
	if b.Size() != 2 {
		t.Errorf("Got %v expected %v", b.Size(), 2)
	}
	if c.Size() != 3 {
		t.Errorf("Got %v expected %v", c.Size(), 3)
	}
	if d.Size() != 1 || d.Size() != f.Size() || d.Size() != g.Size() {
		t.Errorf("Got %v expected %v", d.Size(), 1)
	}
}

func TestTravPre(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot('k', nil)
	i := bt.InsertAsLChild(root, 'i', nil)
	j := bt.InsertAsRChild(root, 'j', nil)
	h := bt.InsertAsRChild(i, 'h', nil)
	b := bt.InsertAsLChild(h, 'b', nil)
	g := bt.InsertAsRChild(h, 'g', nil)
	a := bt.InsertAsRChild(b, 'a', nil)
	e := bt.InsertAsLChild(g, 'e', nil)
	f := bt.InsertAsRChild(g, 'f', nil)
	c := bt.InsertAsLChild(e, 'c', nil)
	d := bt.InsertAsRChild(e, 'd', nil)
	nodes := make([]*Node, 0)
	nodes = append(nodes, root, i, h, b, a, g, e, c, d, f, j)
	travrnodes := bt.TravPreR()
	if len(travrnodes) != len(nodes) {
		t.Errorf("TravPreR got %v expected %v", len(travrnodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != travrnodes[i].Key {
			t.Errorf("TravPreR got %v expected %v", travrnodes[i].Key, n.Key)
		}
	}
	trav1nodes := bt.TravPre1()
	if len(trav1nodes) != len(nodes) {
		t.Errorf("TravPre1 got %v expected %v", len(trav1nodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != trav1nodes[i].Key {
			t.Errorf("TravPre1 got %v expected %v", trav1nodes[i].Key, n.Key)
		}
	}
	trav2nodes := bt.TravPre2()
	if len(trav2nodes) != len(nodes) {
		t.Errorf("TravPre2 got %v expected %v", len(trav2nodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != trav2nodes[i].Key {
			t.Errorf("TravPre2 got %v expected %v", trav2nodes[i].Key, n.Key)
		}
	}
}

func TestTravIn(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot('k', nil)
	i := bt.InsertAsLChild(root, 'i', nil)
	j := bt.InsertAsRChild(root, 'j', nil)
	h := bt.InsertAsRChild(i, 'h', nil)
	b := bt.InsertAsLChild(h, 'b', nil)
	g := bt.InsertAsRChild(h, 'g', nil)
	a := bt.InsertAsRChild(b, 'a', nil)
	e := bt.InsertAsLChild(g, 'e', nil)
	f := bt.InsertAsRChild(g, 'f', nil)
	c := bt.InsertAsLChild(e, 'c', nil)
	d := bt.InsertAsRChild(e, 'd', nil)
	nodes := make([]*Node, 0)
	nodes = append(nodes, i, b, a, h, c, e, d, g, f, root, j)
	travrnodes := bt.TravInR()
	if len(travrnodes) != len(nodes) {
		t.Errorf("TravInR len got %v expected %v", len(travrnodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != travrnodes[i].Key {
			t.Errorf("TravInR got %v expected %v", travrnodes[i].Key, n.Key)
		}
	}
	trav1nodes := bt.TravIn1()
	if len(trav1nodes) != len(nodes) {
		t.Errorf("TravIn1 len got %v expected %v", len(trav1nodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != trav1nodes[i].Key {
			t.Errorf("TravIn1 got %v expected %v", trav1nodes[i].Key, n.Key)
		}
	}
	//bt.Print()
	trav2nodes := bt.TravIn2()
	if len(trav2nodes) != len(nodes) {
		t.Errorf("TravIn2 len got %v expected %v", len(trav2nodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != trav2nodes[i].Key {
			t.Errorf("TravIn2 got %v expected %v", trav2nodes[i].Key, n.Key)
		}
	}
	trav3nodes := bt.TravIn3()
	if len(trav3nodes) != len(nodes) {
		t.Errorf("TravIn3 len got %v expected %v", len(trav3nodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != trav3nodes[i].Key {
			t.Errorf("TravIn3 got %v expected %v", trav3nodes[i].Key, n.Key)
		}
	}
}

func TestTravPost(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot('k', nil)
	i := bt.InsertAsLChild(root, 'i', nil)
	j := bt.InsertAsRChild(root, 'j', nil)
	h := bt.InsertAsRChild(i, 'h', nil)
	b := bt.InsertAsLChild(h, 'b', nil)
	g := bt.InsertAsRChild(h, 'g', nil)
	a := bt.InsertAsRChild(b, 'a', nil)
	e := bt.InsertAsLChild(g, 'e', nil)
	f := bt.InsertAsRChild(g, 'f', nil)
	c := bt.InsertAsLChild(e, 'c', nil)
	d := bt.InsertAsRChild(e, 'd', nil)
	nodes := make([]*Node, 0)
	nodes = append(nodes, a, b, c, d, e, f, g, h, i, j, root)
	travrnodes := bt.TravPostR()
	if len(travrnodes) != len(nodes) {
		t.Errorf("TravPostR len got %v expected %v", len(travrnodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != travrnodes[i].Key {
			t.Errorf("TravPostR got %v expected %v", travrnodes[i].Key, n.Key)
		}
	}
	trav1nodes := bt.TravPost1()
	if len(trav1nodes) != len(nodes) {
		t.Errorf("TravPost1 len got %v expected %v", len(trav1nodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != trav1nodes[i].Key {
			t.Errorf("TravPost1 got %v expected %v", trav1nodes[i].Key, n.Key)
		}
	}
}

func TestTravLevel(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot('k', nil)
	i := bt.InsertAsLChild(root, 'i', nil)
	j := bt.InsertAsRChild(root, 'j', nil)
	h := bt.InsertAsRChild(i, 'h', nil)
	b := bt.InsertAsLChild(h, 'b', nil)
	g := bt.InsertAsRChild(h, 'g', nil)
	a := bt.InsertAsRChild(b, 'a', nil)
	e := bt.InsertAsLChild(g, 'e', nil)
	f := bt.InsertAsRChild(g, 'f', nil)
	c := bt.InsertAsLChild(e, 'c', nil)
	d := bt.InsertAsRChild(e, 'd', nil)
	nodes := make([]*Node, 0)
	nodes = append(nodes, root, i, j, h, b, g, a, e, f, c, d)
	travnodes := bt.TravLevel()
	if len(travnodes) != len(nodes) {
		t.Errorf("TravLevel len got %v expected %v", len(travnodes), len(nodes))
	}
	for i, n := range nodes {
		//fmt.Println(string(n.Key.(rune)))
		if n.Key != travnodes[i].Key {
			t.Errorf("TravLevel got %v expected %v", travnodes[i].Key, n.Key)
		}
	}
}

func TestRemove(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot("k", nil)
	i := bt.InsertAsLChild(root, "i", nil)
	j := bt.InsertAsRChild(root, "j", nil)
	h := bt.InsertAsRChild(i, "h", nil)
	b := bt.InsertAsLChild(h, "b", nil)
	g := bt.InsertAsRChild(h, "g", nil)
	a := bt.InsertAsRChild(b, "a", nil)
	e := bt.InsertAsLChild(g, "e", nil)
	f := bt.InsertAsRChild(g, "f", nil)
	c := bt.InsertAsLChild(e, "c", nil)
	d := bt.InsertAsRChild(e, "d", nil)
	bt.Print()
	fmt.Println("remove ", g.Key)
	bt.Remove(g)
	if bt.Size != 6 {
		t.Errorf("After removing got size %v expected %v", bt.Size, 6)
	}
	if root.Height != 4 {
		t.Errorf("After removing got height %v expected %v", root.Height, 4)
	}
	bt.Print()
	//w := os.Stdout
	wprer := new(bytes.Buffer)
	wpre1 := new(bytes.Buffer)
	wpre2 := new(bytes.Buffer)
	fmt.Println("pre order by recursive:")
	bt.TravPreR(WithPrintNodeKey(wprer))
	fmt.Println(wprer)
	fmt.Println("pre order by iteration1:")
	bt.TravPre1(WithPrintNodeKey(wpre1))
	fmt.Println(wpre1)
	fmt.Println("pre order by iteration2:")
	bt.TravPre2(WithPrintNodeKey(wpre2))
	fmt.Println(wpre2)
	if wprer.String() != wpre1.String() || wpre1.String() != wpre2.String() {
		t.Errorf("pre order is wrong")
	}
	winr := new(bytes.Buffer)
	win1 := new(bytes.Buffer)
	win2 := new(bytes.Buffer)
	win3 := new(bytes.Buffer)
	fmt.Println("in order by recursive:")
	bt.TravInR(WithPrintNodeKey(winr))
	fmt.Println(winr)
	fmt.Println("in order by iteration1:")
	bt.TravIn1(WithPrintNodeKey(win1))
	fmt.Println(win1)
	fmt.Println("in order by iteration2:")
	bt.TravIn2(WithPrintNodeKey(win2))
	fmt.Println(win2)
	fmt.Println("in order by iteration3:")
	bt.TravIn3(WithPrintNodeKey(win3))
	fmt.Println(win3)
	if winr.String() != win1.String() || win1.String() != win2.String() || win2.String() != win3.String() {
		t.Errorf("in order is wrong")
	}
	wpostr := new(bytes.Buffer)
	wpost1 := new(bytes.Buffer)
	fmt.Println("post order by recursive:")
	bt.TravPostR(WithPrintNodeKey(wpostr))
	fmt.Println(wpostr)
	fmt.Println("post order by iteration1:")
	bt.TravPost1(WithPrintNodeKey(wpost1))
	fmt.Println(wpost1)
	if wpostr.String() != wpost1.String() {
		t.Errorf("post order is wrong")
	}
	fmt.Println("level order by iteration1:")
	bt.TravLevel(WithPrintNodeKey(os.Stdout))
	fmt.Println()
	_ = j
	_ = a
	_ = f
	_ = c
	_ = d
}

func TestSecede(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot("k", nil)
	i := bt.InsertAsLChild(root, "i", nil)
	j := bt.InsertAsRChild(root, "j", nil)
	h := bt.InsertAsRChild(i, "h", nil)
	b := bt.InsertAsLChild(h, "b", nil)
	g := bt.InsertAsRChild(h, "g", nil)
	a := bt.InsertAsRChild(b, "a", nil)
	e := bt.InsertAsLChild(g, "e", nil)
	f := bt.InsertAsRChild(g, "f", nil)
	c := bt.InsertAsLChild(e, "c", nil)
	d := bt.InsertAsRChild(e, "d", nil)
	bt.Print()
	fmt.Println("Secede ", g.Key)
	subt := bt.Secede(g)
	if bt.Size != 6 {
		t.Errorf("After seceding got size %v expected %v", bt.Size, 6)
	}
	if root.Height != 4 {
		t.Errorf("After seceding got height %v expected %v", root.Height, 4)
	}
	bt.Print()
	if subt.Size != 5 {
		t.Errorf("After seceding sub tree got size %v expected %v", subt.Size, 5)
	}
	if subt.Root.Height != 2 {
		t.Errorf("After seceding sub tree got height %v expected %v", subt.Root.Height, 2)
	}
	subt.Print()
	//w := os.Stdout

	_ = j
	_ = a
	_ = f
	_ = c
	_ = d
}

func TestAttachAsLSubTree(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot("k", nil)
	i := bt.InsertAsLChild(root, "i", nil)
	j := bt.InsertAsRChild(root, "j", nil)
	h := bt.InsertAsRChild(i, "h", nil)
	b := bt.InsertAsLChild(h, "b", nil)
	g := bt.InsertAsRChild(h, "g", nil)
	a := bt.InsertAsRChild(b, "a", nil)
	e := bt.InsertAsLChild(g, "e", nil)
	f := bt.InsertAsRChild(g, "f", nil)
	c := bt.InsertAsLChild(e, "c", nil)
	d := bt.InsertAsRChild(e, "d", nil)
	bt.Print()
	bt2 := new(BinTree)
	root2 := bt2.InsertAsRoot("k", nil)
	i2 := bt2.InsertAsLChild(root2, "i", nil)
	j2 := bt2.InsertAsRChild(root2, "j", nil)
	h2 := bt2.InsertAsRChild(i2, "h", nil)
	b2 := bt2.InsertAsLChild(h2, "b", nil)
	g2 := bt2.InsertAsRChild(h2, "g", nil)
	a2 := bt2.InsertAsRChild(b2, "a", nil)
	e2 := bt2.InsertAsLChild(g2, "e", nil)
	f2 := bt2.InsertAsRChild(g2, "f", nil)
	c2 := bt2.InsertAsLChild(e2, "c", nil)
	d2 := bt2.InsertAsRChild(e2, "d", nil)
	bt2.Print()
	bt.AttachAsLSubTree(b, bt2)
	bt.PrintWithUnitSize(1)
	if bt.Size != 22 {
		t.Errorf("After attach got size %v expected %v", bt.Size, 22)
	}
	if root.Height != 9 {
		t.Errorf("After attach got height %v expected %v", root.Height, 9)
	}
	for _, n := range bt.TravPre1() {
		if n.Tree != bt {
			t.Errorf("After attach node %v got belongs to %v expected %v", n.Key, n.Tree, bt)
		}
	}
	_ = j
	_ = a
	_ = f
	_ = c
	_ = d

	_ = j2
	_ = a2
	_ = f2
	_ = c2
	_ = d2
}

func TestAttachAsRSubTree(t *testing.T) {
	bt := new(BinTree)
	root := bt.InsertAsRoot("k", nil)
	i := bt.InsertAsLChild(root, "i", nil)
	j := bt.InsertAsRChild(root, "j", nil)
	h := bt.InsertAsRChild(i, "h", nil)
	b := bt.InsertAsLChild(h, "b", nil)
	g := bt.InsertAsRChild(h, "g", nil)
	a := bt.InsertAsRChild(b, "a", nil)
	e := bt.InsertAsLChild(g, "e", nil)
	f := bt.InsertAsRChild(g, "f", nil)
	c := bt.InsertAsLChild(e, "c", nil)
	d := bt.InsertAsRChild(e, "d", nil)
	bt.Print()
	bt2 := new(BinTree)
	root2 := bt2.InsertAsRoot("k", nil)
	i2 := bt2.InsertAsLChild(root2, "i", nil)
	j2 := bt2.InsertAsRChild(root2, "j", nil)
	h2 := bt2.InsertAsRChild(i2, "h", nil)
	b2 := bt2.InsertAsLChild(h2, "b", nil)
	g2 := bt2.InsertAsRChild(h2, "g", nil)
	a2 := bt2.InsertAsRChild(b2, "a", nil)
	e2 := bt2.InsertAsLChild(g2, "e", nil)
	f2 := bt2.InsertAsRChild(g2, "f", nil)
	c2 := bt2.InsertAsLChild(e2, "c", nil)
	d2 := bt2.InsertAsRChild(e2, "d", nil)
	bt2.Print()
	// bt3 := bt2.Copy()
	// bt3.Print()
	bt.AttachAsRSubTree(j, bt2)
	bt.PrintWithUnitSize(1)
	// bt.AttachAsLSubTree(j2, bt3)
	// bt.PrintWithUnitSize(1)
	if bt.Size != 22 {
		t.Errorf("After attach got size %v expected %v", bt.Size, 22)
	}
	if root.Height != 7 {
		t.Errorf("After attach got height %v expected %v", root.Height, 7)
	}
	for _, n := range bt.TravPre1() {
		if n.Tree != bt {
			t.Errorf("After attach node %v got belongs to %v expected %v", n.Key, n.Tree, bt)
		}
	}
	_ = j
	_ = a
	_ = f
	_ = c
	_ = d

	_ = j2
	_ = a2
	_ = f2
	_ = c2
	_ = d2
}

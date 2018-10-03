package bst

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Print 以子树节点个数的位数为一个基本单元的长度，打印子树的拓扑结构到标准输出
func Print(n Node) {
	PrintWithUnitSize(n, len(strconv.Itoa(Size(n))))
}

// Fprint 以树节点个数的位数为一个基本单元的长度，打印子树的拓扑结构到io.Writer
func Fprint(n Node, w io.Writer) {
	FprintWithUnitSize(n, w, len(strconv.Itoa(Size(n))))
}

// PrintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到标准输出
func PrintWithUnitSize(n Node, size int) {
	FprintWithUnitSize(n, os.Stdout, size)
}

// FprintWithUnitSize 以指定的长度为一个基本单元，打印子树的拓扑结构到io.Writer，树宽为节点数
func FprintWithUnitSize(n Node, w io.Writer, size int) {
	buf := bufio.NewWriter(w)
	if IsNil(n) {
		buf.WriteString("Empty tree!")
		buf.Flush()
		return
	}
	if size <= 0 {
		panic("unit size can not be less than 1")
	}

	total := Size(n)
	q := make([]nodePos, 0, total)
	prevlevel := 0
	line := make([]rune, total*size+1)
	for i := range line {
		line[i] = ' '
	}
	mid := Size(n.LChild())
	left, right := mid, mid
	if !IsNil(n.LChild()) {
		left = mid - Size(n.LChild().RChild()) - 1
	}
	if !IsNil(n.RChild()) {
		right = mid + Size(n.RChild().LChild()) + 1
	}
	q = append(q, nodePos{n, left, mid, right})
	for len(q) > 0 {
		np := q[0]
		q = q[1:]
		n := np.node
		if l := Level(n); prevlevel != l {
			prevlevel = l
			buf.WriteString("\n")
			for i, r := range line {
				buf.WriteRune(r)
				line[i] = ' '
			}
		}

		np.fillNode(line, size)

		if HasLChild(n) {
			left, mid, right = np.computelchildPos()
			q = append(q, nodePos{n.LChild(), left, mid, right})
		}
		if HasRChild(n) {
			left, mid, right = np.computerchildPos()
			q = append(q, nodePos{n.RChild(), left, mid, right})
		}
	}
	buf.WriteString("\n")
	for _, r := range line {
		buf.WriteRune(r)
	}
	buf.WriteString("\n")
	buf.Flush()
}

func (np nodePos) fillNode(line []rune, size int) {
	if HasLChild(np.node) {
		i := np.left * size
		// for ; i < np.left*size; i++ {
		// 	line[i] = ' '
		// }
		for ; i < np.left*size+size-1; i++ {
			line[i] = ' '
		}
		line[i] = '┌'
		i++
		for ; i < np.mid*size; i++ {
			line[i] = '─'
		}
	}
	i := np.mid * size
	for _, r := range fmt.Sprintf("%*v%s", size, np.node.Key(), np.node.Color()) {
		line[i] = r
		i++
	}
	if HasRChild(np.node) {
		// for ; i < np.right*size-size; i++ {
		// 	line[i] = '─'
		// }
		for ; i < np.right*size+size-1; i++ {
			line[i] = '─'
		}
		line[i] = '┐'
	}
}

type nodePos struct {
	node             Node
	left, mid, right int
}

func (np nodePos) computelchildPos() (left, mid, right int) {
	if IsNil(np.node) {
		return
	}
	if IsNil(np.node.LChild()) {
		return mid, mid, mid
	}
	mid = np.mid - Size(np.node.LChild().RChild()) - 1

	if IsNil(np.node.LChild().LChild()) {
		left = mid
	} else {
		left = mid - Size(np.node.LChild().LChild().RChild()) - 1
	}
	if IsNil(np.node.LChild().RChild()) {
		right = mid
	} else {
		right = mid + Size(np.node.LChild().RChild().LChild()) + 1
	}
	return
}

func (np nodePos) computerchildPos() (left, mid, right int) {
	if IsNil(np.node) {
		return
	}
	if IsNil(np.node.RChild()) {
		return mid, mid, mid
	}
	mid = np.mid + Size(np.node.RChild().LChild()) + 1

	if IsNil(np.node.RChild().LChild()) {
		left = mid
	} else {
		left = mid - Size(np.node.RChild().LChild().RChild()) - 1
	}
	if IsNil(np.node.RChild().RChild()) {
		right = mid
	} else {
		right = mid + Size(np.node.RChild().RChild().LChild()) + 1
	}
	return
}

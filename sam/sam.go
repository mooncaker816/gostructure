package sam

import "fmt"

type state struct {
	len  int             // 该状态对应的子串等价类中最长的子串长度
	link *state          // suffix link
	next map[byte]*state // key 对应的是该状态接收的字符，value 对应的输出状态
}

func (sam *SAM) newState() *state {
	s := &state{len: 0, link: nil, next: make(map[byte]*state)}
	sam.st[sam.size] = s
	sam.size++
	return s
}

func (sam *SAM) cloneState(model *state) *state {
	clone := sam.newState()
	clone.next = copyMap(model.next)
	clone.link = model.link
	return clone
}

// SAM suffix automaton
type SAM struct {
	st   []*state // 状态集
	size int
	last *state
}

// NewSAM creates a SAM
func NewSAM(s string) *SAM {
	sam := new(SAM)
	sam.init(len(s))
	for i := 0; i < len(s); i++ {
		sam.extend(s[i])
	}
	return sam
}

func (sam *SAM) init(size int) {
	sam.st = make([]*state, size*2) // 状态最多不会超过 2*len(s) - 1
	// 初始状态
	sam.newState()

	sam.last = sam.st[0]
}

func (sam *SAM) extend(c byte) {
	// 新建一个状态
	curr := sam.newState()
	curr.len = sam.last.len + 1

	// 从上一状态开始依次追寻 suffix link，
	// 如果该状态没有对应 c 的出边，则添加对应的出边到 curr
	// 直到找到某一状态已经拥有 c 的出边，或者已经抵达初态的 link -1
	p := sam.last
	for p != nil {
		if _, ok := p.next[c]; ok {
			break
		}
		p.next[c] = curr
		p = p.link
	}

	if p == nil {
		curr.link = sam.st[0]
	} else {
		// 找到了状态 p，有对应 c 的出边，输出状态为 q
		q := p.next[c]
		// 如果状态 p 对应的子串等价类的最大长度比状态 q 等价类子串的最大长度大1
		// 说明 q 对应的子串都是 curr 的后缀，curr 可以直接添加 suffix link 到 q
		if p.len+1 == q.len {
			curr.link = q
		} else {
			// 否则需要从 q 中分离出一个新的状态用来接收 p 发出的 c
			// 新状态只负责从 q 中转移出“字符 c 对应的部分”，其 next 和 link 与 q 一致
			// 同时将 p 的 suffix link 中的 next[c] 改为新的状态
			clone := sam.cloneState(q)
			clone.len = p.len + 1
			for p != nil && p.next[c] == q {
				p.next[c] = clone
				p = p.link
			}
			// 将 q 和 curr 的 suffix link 指向分离出来的状态
			q.link = clone
			curr.link = clone
		}
	}
	sam.last = curr
}

// func (sam *SAM)
func copyMap(m map[byte]*state) map[byte]*state {
	nm := make(map[byte]*state, len(m))
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

// SubStrCount returns the number of all the distinct sub-strings in the sam
func (sam *SAM) SubStrCount() int {
	cnt := 0
	for i := 0; i < sam.size; i++ {
		fmt.Println(sam.st[i])
		if sam.st[i].link == nil {
			continue
		}
		cnt += sam.st[i].len - sam.st[i].link.len
	}
	return cnt
}

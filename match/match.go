package match

// BFMatch Brute Force match
func BFMatch(s, pattern string) int {
	if len(s) < len(pattern) {
		return -1
	}
	i, j := 0, 0
	for i < len(s) && j < len(pattern) {
		if s[i] == pattern[j] { // 当前字符匹配，继续向右匹配下一组
			i++
			j++
		} else { // 当前字符不匹配，字符串回溯到此时与模式串首字符对齐的字符的下一个字符，模式串回溯到首字符
			i = i - j + 1
			j = 0
		}
	}
	// find a match
	if j == len(pattern) {
		return i - j
	}
	// no match
	return -1
}

// KMPMatch KMP algorithm
func KMPMatch(s, pattern string) int {
	if len(pattern) == 0 {
		return 0
	}
	if len(s) < len(pattern) {
		return -1
	}
	next := genNextSmart(pattern)
	i, j := 0, 0
	for i < len(s) && j < len(pattern) {
		if j == -1 || s[i] == pattern[j] { // -1 为模式串的前缀通配哨兵，此时相当于将整个模式串相对于字符串右移一位
			i++
			j++
		} else {
			j = next[j]
		}
	}
	// find a match
	if j == len(pattern) {
		return i - j
	}
	// no match
	return -1
}

// next[j] 表示当模式串中的第 j 个字符跟文本串中的第 i 个字符匹配失配时，模式串下一步应该跳到的位置。
// 相当于下一次比较时，将模式串向右移动 j - next[j] 个字符，使得文本串的 i 和模式串的 next[j] 对齐
func genNext(pattern string) []int {
	next := make([]int, len(pattern))
	// j 表示 next[i], 已知 next[i] 递推求解 next[i+1]
	// 即对于 pattern[:i]来说，有 j 个相同的真前缀和真后缀，pattern[:j] = pattern[i-j:i]

	next[0] = -1
	i, j := 0, -1
	for i < len(pattern)-1 {
		if j == -1 || pattern[i] == pattern[j] {
			// 如果 pattern[j] 和 pattern[i] 相同，则说明 pattern[:i+1] 拥有 j+1 个相同的真前缀和真后缀，即next[i+1] = j+1
			next[i+1] = j + 1
			i++
			j++
		} else {
			// 如果 pattern[j] != pattern[i]，说明加上第 i 个字符后，pattern[i-j:i+1] != pattern[:j+1]，
			// 但是不排除存在一个较短的长度为 k 的子串，使得 pattern[:k] = pattern[i-k+1:i+1],
			// 如最初相同的真前缀和真后缀为 abcab，真后缀加上字符 c，得到 abcabc，真前缀加上字符 d，得到 abcabd
			// 显然他们仍然拥有较短的相同的真前缀 abc 和较短的真后缀 abc，所以 next[i+1] = 3
			j = next[j] //相同真前缀的真前缀
		}
	}
	return next
}

// 考虑到使用 next[i] 的情况是第 i 的字符失配，所以如果替换后比较的字符不变，相当于无用功，在构造 next 表的时候要避免这种情况
func genNextSmart(pattern string) []int {
	next := make([]int, len(pattern))
	next[0] = -1
	i, j := 0, -1
	for i < len(pattern)-1 {
		if j == -1 || pattern[i] == pattern[j] {
			// 此时在构造 next[i+1]，使用 next[i+1] 的前提是 pattern[i+1] 失配，
			// 为了保证重新对齐后的第 j+1 字符不会因为和 pattern[i+1] 相同而再次失配，直接替换为 next[j+1]
			// 相当于提前做了再次
			if pattern[i+1] != pattern[j+1] {
				next[i+1] = j + 1
			} else {
				next[i+1] = next[j+1]
			}
			i++
			j++
		} else {
			j = next[j]
		}
	}
	return next
}

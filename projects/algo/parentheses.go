package algo

// ValidParentheses reports whether s's brackets ((), [], {}) are balanced
// and correctly nested. O(n) time using the Stack built earlier — every
// opener gets pushed, every closer must match the most recent unmatched
// opener, which is exactly a LIFO's job.
func ValidParentheses(s string) bool {
	pairs := map[byte]byte{')': '(', ']': '[', '}': '{'}
	var stack Stack[byte]
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '(', '[', '{':
			stack.Push(ch)
		case ')', ']', '}':
			top, ok := stack.Pop()
			if !ok || top != pairs[ch] {
				return false
			}
		}
	}
	return stack.Len() == 0
}

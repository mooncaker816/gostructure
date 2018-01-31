package stack

import (
	"fmt"
	"strings"
)

// Stack is a slice with any type of data -- LIFO
type Stack struct {
	elements []interface{}
}

// NewStack creates a stack with default 0 capacity
func NewStack() *Stack { return NewStackCap(0) }

// NewStackCap creates a stack with provided capactity
func NewStackCap(cap int) *Stack {
	if cap < 0 {
		panic("stack's capacity can not be negtive")
	}
	return &Stack{make([]interface{}, 0, cap)}
}

// Len returns the count of the stack's elements
func (s *Stack) Len() int { return len(s.elements) }

// Empty checks if the stack is empty or not
func (s *Stack) Empty() bool { return s.Len() == 0 }

// Push adds an emlement to the stack's top
func (s *Stack) Push(value interface{}) {
	s.elements = append(s.elements, value)
}

// Pop returns the top element from the stack and removes it from stack, if the stack is empty,then returns nil,false
func (s *Stack) Pop() (interface{}, bool) {
	value, ok := s.Peek()
	if ok {
		s.elements = s.elements[:len(s.elements)-1]
	}
	return value, ok
}

// Peek returns the top element from the stack without remove it from stack, if it's empty, then return nil,false
func (s *Stack) Peek() (interface{}, bool) {
	if s.Empty() {
		return nil, false
	}
	return s.elements[len(s.elements)-1], true
}

// String formats the stack's elememts to a string from top to bottom a,b,c
func (s *Stack) String() string {
	vals := make([]string, 0)
	for i := len(s.elements) - 1; i >= 0; i-- {
		vals = append(vals, fmt.Sprintf("%v", s.elements[i]))
	}
	return strings.Join(vals, ",")
}

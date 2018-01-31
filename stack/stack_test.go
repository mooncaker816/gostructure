package stack

import (
	"fmt"
	"testing"
)

func TestLen(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if s.Len() != 3 {
		t.Errorf("Got %v expected %v", s.Len(), 3)
	}
	s.Peek()
	if s.Len() != 3 {
		t.Errorf("Got %v expected %v", s.Len(), 3)
	}
	s.Pop()
	if s.Len() != 2 {
		t.Errorf("Got %v expected %v", s.Len(), 2)
	}
	s.Pop()
	s.Pop()
	if s.Len() != 0 {
		t.Errorf("Got %v expected %v", s.Len(), 0)
	}
}

func TestEmpty(t *testing.T) {
	s := NewStackCap(10)
	if !s.Empty() {
		t.Errorf("Got %v expected %v", s.Empty(), true)
	}
	s.Push('a')
	if s.Empty() {
		t.Errorf("Got %v expected %v", s.Empty(), false)
	}
	s.Peek()
	if s.Empty() {
		t.Errorf("Got %v expected %v", s.Empty(), false)
	}
	s.Pop()
	if !s.Empty() {
		t.Errorf("Got %v expected %v", s.Empty(), true)
	}
}

func TestPush(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	str := "3,2,1"
	if fmt.Sprintf("%v", s) != str {
		t.Errorf("Got %v expected %v", fmt.Sprintf("%v", s), str)
	}
	v, ok := s.Pop()
	if !ok || v.(int) != 3 {
		t.Errorf("Got %v expected %v", v, 3)
	}
	if fmt.Sprintf("%v", s) != "2,1" {
		t.Errorf("Got %v expected %v", fmt.Sprintf("%v", s), "2,1")
	}
}

func TestPop(t *testing.T) {
	s := NewStackCap(5)
	s.Push("a")
	s.Push("b")
	s.Push("c")
	len1 := s.Len()
	if e, ok := s.Pop(); ok {
		if e.(string) != "c" {
			t.Errorf("Got %v expected %v", e, "c")
		}
	}
	if len1-1 != s.Len() {
		t.Errorf("Got %v expected %v", s.Len(), len1-1)
	}
	if e, ok := s.Pop(); ok {
		if e.(string) != "b" {
			t.Errorf("Got %v expected %v", e, "b")
		}
	}
	if len1-2 != s.Len() {
		t.Errorf("Got %v expected %v", s.Len(), len1-2)
	}
}

func TestPeek(t *testing.T) {
	s := NewStack()
	e, ok := s.Peek()
	if ok || e != nil {
		t.Errorf("Got %v expected %v", e, nil)
	}
	s.Push(1)
	s.Push(2)
	e, ok = s.Peek()
	if !ok || e != 2 {
		t.Errorf("Got %v expected %v", e, 2)
	}
	s.Push(3)
	e, ok = s.Peek()
	if !ok || e != 3 {
		t.Errorf("Got %v expected %v", e, 3)
	}
	s.Pop()
	s.Pop()
	e, ok = s.Peek()
	if !ok || e != 1 {
		t.Errorf("Got %v expected %v", e, 1)
	}
}

package queue

import (
	"fmt"
	"testing"
)

func TestLen(t *testing.T) {
	q := NewQueue()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if q.Len() != 3 {
		t.Errorf("Got %v expected %v", q.Len(), 3)
	}
	q.Peek()
	if q.Len() != 3 {
		t.Errorf("Got %v expected %v", q.Len(), 3)
	}
	q.Dequeue()
	if q.Len() != 2 {
		t.Errorf("Got %v expected %v", q.Len(), 2)
	}
	q.Dequeue()
	q.Dequeue()
	if q.Len() != 0 {
		t.Errorf("Got %v expected %v", q.Len(), 0)
	}
}

func TestEmpty(t *testing.T) {
	q := NewQueueCap(10)
	if !q.Empty() {
		t.Errorf("Got %v expected %v", q.Empty(), true)
	}
	q.Enqueue('a')
	if q.Empty() {
		t.Errorf("Got %v expected %v", q.Empty(), false)
	}
	q.Peek()
	if q.Empty() {
		t.Errorf("Got %v expected %v", q.Empty(), false)
	}
	q.Dequeue()
	if !q.Empty() {
		t.Errorf("Got %v expected %v", q.Empty(), true)
	}
}

func TestEnqueue(t *testing.T) {
	q := NewQueue()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	str := "1,2,3"
	if fmt.Sprintf("%v", q) != str {
		t.Errorf("Got %v expected %v", fmt.Sprintf("%v", q), str)
	}
	v, ok := q.Dequeue()
	if !ok || v.(int) != 1 {
		t.Errorf("Got %v expected %v", v, 1)
	}
	if fmt.Sprintf("%v", q) != "2,3" {
		t.Errorf("Got %v expected %v", fmt.Sprintf("%v", q), "2,3")
	}
}

func TestDequeue(t *testing.T) {
	q := NewQueueCap(5)
	q.Enqueue("a")
	q.Enqueue("b")
	q.Enqueue("c")
	len1 := q.Len()
	if e, ok := q.Dequeue(); ok {
		if e.(string) != "a" {
			t.Errorf("Got %v expected %v", e, "a")
		}
	}
	if len1-1 != q.Len() {
		t.Errorf("Got %v expected %v", q.Len(), len1-1)
	}
	if e, ok := q.Dequeue(); ok {
		if e.(string) != "b" {
			t.Errorf("Got %v expected %v", e, "b")
		}
	}
	if len1-2 != q.Len() {
		t.Errorf("Got %v expected %v", q.Len(), len1-2)
	}
}

func TestPeek(t *testing.T) {
	q := NewQueue()
	e, ok := q.Peek()
	if ok || e != nil {
		t.Errorf("Got %v expected %v", e, nil)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	e, ok = q.Peek()
	if !ok || e != 1 {
		t.Errorf("Got %v expected %v", e, 1)
	}
	q.Enqueue(3)
	e, ok = q.Peek()
	if !ok || e != 1 {
		t.Errorf("Got %v expected %v", e, 1)
	}
	q.Dequeue()
	q.Dequeue()
	e, ok = q.Peek()
	if !ok || e != 3 {
		t.Errorf("Got %v expected %v", e, 3)
	}
}

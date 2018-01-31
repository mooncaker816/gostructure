package queue

import (
	"fmt"
	"strings"
)

// Queue is a slice with any type of data -- FIFO
type Queue struct {
	elements []interface{}
}

// NewQueue creates a queue with default 0 capacity
func NewQueue() *Queue { return NewQueueCap(0) }

// NewQueueCap creates a queue with provided capactity
func NewQueueCap(cap int) *Queue {
	if cap < 0 {
		panic("queue's capacity can not be negtive")
	}
	return &Queue{make([]interface{}, 0, cap)}
}

// Len returns the count of the queue's elements
func (q *Queue) Len() int { return len(q.elements) }

// Empty checks if the queue is empty or not
func (q *Queue) Empty() bool { return q.Len() == 0 }

// Enqueue adds an emlement to the queue's tail
func (q *Queue) Enqueue(value interface{}) {
	q.elements = append(q.elements, value)
}

// Dequeue returns the front element from the queue and removes it from queue, if the queue is empty,then returns nil,false
func (q *Queue) Dequeue() (interface{}, bool) {
	value, ok := q.Peek()
	if ok {
		q.elements = q.elements[1:len(q.elements)]
	}
	return value, ok
}

// Peek returns the front element from the queue without remove it from queue, if it's empty, then return nil,false
func (q *Queue) Peek() (interface{}, bool) {
	if q.Empty() {
		return nil, false
	}
	return q.elements[0], true
}

// String formats the queue's elememts to a string from front to tail like a,b,c
func (q *Queue) String() string {
	vals := make([]string, 0)
	for _, v := range q.elements {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return strings.Join(vals, ",")
}

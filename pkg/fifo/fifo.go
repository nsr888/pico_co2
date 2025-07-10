// Package fifo provides a fixed-size, no-allocation FIFO queue for int16 values on the stack.
package fifo

const MaxCapacity = 128 // maximum supported capacity

// FIFO16 implements a circular buffer queue of int16 with runtime capacity up to MaxCapacity.
type FIFO16 struct {
	buf      [MaxCapacity]int16 // pre-allocated storage on the stack
	capacity int                // actual usable capacity (≤ MaxCapacity)
	head     int                // index of the oldest element
	tail     int                // index to write the next element
	count    int                // number of elements stored
}

// NewFIFO16 creates a FIFO with the given capacity (1..MaxCapacity). No heap allocation.
func NewFIFO16(cap int) *FIFO16 {
	if cap <= 0 {
		cap = 1
	} else if cap > MaxCapacity {
		cap = MaxCapacity
	}
	return &FIFO16{capacity: cap}
}

// Reset clears the queue back to empty state.
func (q *FIFO16) Reset() {
	q.head, q.tail, q.count = 0, 0, 0
}

// Len returns the number of elements in the queue.
func (q *FIFO16) Len() int {
	return q.count
}

// IsEmpty reports whether the queue has no elements.
func (q *FIFO16) IsEmpty() bool {
	return q.count == 0
}

// IsFull reports whether the queue has reached its capacity.
func (q *FIFO16) IsFull() bool {
	return q.count == q.capacity
}

// Enqueue adds v at the tail. If full, it drops the oldest to make space.
// No allocations or pointers.
func (q *FIFO16) Enqueue(v int16) {
	if q.count == q.capacity {
		// drop oldest
		q.head = (q.head + 1) % q.capacity
		q.count--
	}
	q.buf[q.tail] = v
	q.tail = (q.tail + 1) % q.capacity
	q.count++
}

// Dequeue removes and returns the oldest element; ok=false if empty.
func (q *FIFO16) Dequeue() (v int16, ok bool) {
	if q.count == 0 {
		return 0, false
	}
	v = q.buf[q.head]
	q.head = (q.head + 1) % q.capacity
	q.count--
	return v, true
}

// PeekAll calls fn(v) for each element from oldest to newest, without removing.
func (q *FIFO16) PeekAll(fn func(int16)) {
	idx := q.head
	for i := 0; i < q.count; i++ {
		fn(q.buf[idx])
		idx = (idx + 1) % q.capacity
	}
}

func (q *FIFO16) CopyTo() []int16 {
	out := make([]int16, q.count)
	idx := q.head
	for i := 0; i < q.count; i++ {
		out[i] = q.buf[idx]
		idx = (idx + 1) % q.capacity
	}
	return out
}

/*
Usage example (no heap, everything on stack):

import (
    "fmt"
    "time"
    "your/module/path/fifo"
)

func main() {
    // create a queue of up to 64 samples (stack-allocated array)
    q := fifo.NewFIFO16(64)

    for {
        q.Enqueue(readSensor())
        time.Sleep(100 * time.Millisecond)

        // every minute at second 0
        if time.Now().Second() == 0 {
            fmt.Println("Measurements:")
            q.PeekAll(func(v int16) {
                fmt.Println(v)
            })
        }
    }
}
*/

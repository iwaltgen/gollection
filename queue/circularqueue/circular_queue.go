package circularqueue

import (
	"fmt"
	"strings"
)

// Queue is circular queue
type Queue struct {
	elements       []interface{}
	size           int     // available elements count
	base           int     // available elements first index
	growthFactor   float32 // arrived max size then grow by `factor` percent (0 means never grow)
	shrinkFactor   float32 // shrink when size is `factor` percent of capacity (0 means never shrink)
	guaranteedSize int     // prevent shrink when less than guaranteed size[default: 8] (8 means size less than 8 never shrink)
}

// New instantiates a new empty queue
func New(opts ...Option) *Queue {
	queue := &Queue{}
	for _, opt := range opts {
		opt(queue)
	}

	if queue.Capacity() == 0 {
		queue.guaranteedSize = 8
		queue.elements = make([]interface{}, queue.guaranteedSize)
	}
	return queue
}

// Add appends a value at the end of the queue
func (queue *Queue) Add(values ...interface{}) {
	if 0.0 < queue.growthFactor {
		count := len(values)
		if queue.Capacity() < queue.size+count {
			queue.growBy(count)
		}

		for _, value := range values {
			queue.elements[queue.writerIndex()] = value
			queue.size++
		}
		return
	}

	capacity := queue.Capacity()
	for _, value := range values {
		queue.elements[queue.writerIndex()] = value

		if queue.size < capacity {
			queue.size++
		} else {
			queue.base = queue.elementIndex(1)
		}
	}
}

// Peek returns first element on the queue without removing it, or nil if queue is empty.
// Second return parameter is true, unless the queue was empty and there was nothing to peek.
func (queue *Queue) Peek() (interface{}, bool) {
	if queue.size == 0 {
		return nil, false
	}

	return queue.elements[queue.base], true
}

// Poll get one element from the queue front.
func (queue *Queue) Poll() (interface{}, bool) {
	if queue.size == 0 {
		return nil, false
	}

	value := queue.elements[queue.base]
	queue.elements[queue.base] = nil
	queue.size--
	queue.base = queue.elementIndex(1)

	queue.shrink()
	return value, true
}

// PollUntil get the elements at the given count from the queue front.
func (queue *Queue) PollUntil(count int) ([]interface{}, bool) {
	if queue.size < count {
		return nil, false
	}

	values := make([]interface{}, count)
	for i := 0; i < count; i++ {
		index := queue.elementIndex(i)
		values[i] = queue.elements[index]
		queue.elements[index] = nil
	}
	queue.size -= count
	queue.base = queue.elementIndex(count)

	queue.shrink()
	return values, true
}

// Remove remove one element from the queue front.
func (queue *Queue) Remove() bool {
	if queue.size == 0 {
		return false
	}

	queue.elements[queue.base] = nil
	queue.size--
	queue.base = queue.elementIndex(1)

	queue.shrink()
	return true
}

// RemoveUntil removes the elements at the given count from the queue front.
func (queue *Queue) RemoveUntil(count int) bool {
	if queue.size < count {
		return false
	}

	for i := 0; i < count; i++ {
		queue.elements[queue.elementIndex(i)] = nil
	}
	queue.size -= count
	queue.base = queue.elementIndex(count)

	queue.shrink()
	return true
}

// Element returns the element at index.
// Second return parameter is true if index is within bounds of the array and array is not empty, otherwise false.
func (queue *Queue) Element(index int) (interface{}, bool) {
	if !queue.withinRange(index) {
		return nil, false
	}

	return queue.elements[queue.elementIndex(index)], true
}

// IndexOf returns index of provided predicate
func (queue *Queue) IndexOf(predicate func(interface{}) bool) int {
	if queue.size == 0 {
		return -1
	}

	for i := 0; i < queue.size; i++ {
		element := queue.elements[queue.elementIndex(i)]
		if predicate(element) {
			return i
		}
	}
	return -1
}

// Empty returns true if queue does not contain any elements.
func (queue *Queue) Empty() bool {
	return queue.size == 0
}

// Size returns number of elements within the queue.
func (queue *Queue) Size() int {
	return queue.size
}

// Capacity returns capacity the queue.
func (queue *Queue) Capacity() int {
	return cap(queue.elements)
}

// Clear removes all elements from the queue.
func (queue *Queue) Clear() {
	queue.base, queue.size = 0, 0
	queue.elements = make([]interface{}, queue.guaranteedSize)
}

// Values returns all elements in the queue (FIFO order).
func (queue *Queue) Values() []interface{} {
	if queue.size == 0 {
		return nil
	}

	values := make([]interface{}, queue.size)
	queue.copyTo(values)
	return values
}

// Stqueue returns a stqueue representation of container
func (queue *Queue) String() string {
	str := "CircularQueue\n"
	values := make([]string, 0, queue.size)
	for i := 0; i < queue.size; i++ {
		values = append(values, fmt.Sprintf("%v", queue.elements[queue.elementIndex(i)]))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within bounds of the queue
func (queue *Queue) withinRange(index int) bool {
	return index >= 0 && index < queue.size
}

// append elements index
func (queue *Queue) writerIndex() int {
	return queue.elementIndex(queue.size)
}

// access index to elements index
func (queue *Queue) elementIndex(index int) int {
	return (queue.base + index) % queue.Capacity()
}

func (queue *Queue) copyTo(values []interface{}) {
	start := queue.base
	end := queue.writerIndex()
	if start < end {
		copy(values, queue.elements[start:end])
		return
	}

	capacity := queue.Capacity()
	copy(values, queue.elements[start:])
	copy(values[capacity-start:], queue.elements[:end])
}

func (queue *Queue) resize(capacity int) {
	items := make([]interface{}, capacity)
	queue.copyTo(items)

	queue.elements = items
	queue.base = 0
}

func (queue *Queue) growBy(n int) {
	newCapacity := int(queue.growthFactor * float32(queue.Capacity()+n))
	queue.resize(newCapacity)
}

func (queue *Queue) shrink() {
	if queue.shrinkFactor == 0.0 || queue.size < queue.guaranteedSize {
		return
	}

	capacity := queue.Capacity()
	if queue.size <= int(float32(capacity)*queue.shrinkFactor) {
		queue.resize(queue.size)
	}
}

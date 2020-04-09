package circularqueue

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueNew(t *testing.T) {
	queue := New()

	assert.Zero(t, queue.Size())
	assert.True(t, queue.Empty())
}

func TestQueueAdd(t *testing.T) {
	queue := New()
	queue.Add("a")
	queue.Add("b", "c")

	assert.False(t, queue.Empty())
	assert.Equal(t, 3, queue.Size())

	element, ok := queue.Element(2)
	assert.Equal(t, "c", element)
	assert.True(t, ok)
}

func TestQueueAddOverwrite(t *testing.T) {
	queue := New()
	queue.Add("a", "b", "c", "d")
	queue.Add("e", "f", "g", "h")

	assert.False(t, queue.Empty())

	capacity := 8
	assert.Equal(t, capacity, queue.Capacity())
	assert.Equal(t, capacity, queue.Size())

	element, ok := queue.Peek()
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	queue.Add("i", "j")

	assert.Equal(t, capacity, queue.Capacity())
	assert.Equal(t, capacity, queue.Size())

	element, ok = queue.Poll()
	assert.Equal(t, "c", element)
	assert.True(t, ok)

	element, ok = queue.Poll()
	assert.Equal(t, "d", element)
	assert.True(t, ok)

	assert.True(t, queue.Remove(4))

	element, ok = queue.Poll()
	assert.Equal(t, "i", element)
	assert.True(t, ok)

	element, ok = queue.Poll()
	assert.Equal(t, "j", element)
	assert.True(t, ok)

	assert.True(t, queue.Empty())
	assert.Equal(t, capacity, queue.Capacity())
}

func TestQueueAddGrow(t *testing.T) {
	queue := New(GrowthFactor(2.0, 0.0))
	queue.Add("a", "b", "c", "d")
	queue.Add("e", "f", "g", "h")

	assert.False(t, queue.Empty())

	capacity := 8
	assert.Equal(t, capacity, queue.Capacity())
	assert.Equal(t, capacity, queue.Size())

	element, ok := queue.Peek()
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	queue.Add("i", "j")

	doubleCapacity := int((capacity + 2) * 2.0)
	assert.Equal(t, doubleCapacity, queue.Capacity())
	assert.Equal(t, 10, queue.Size())

	element, ok = queue.Poll()
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	element, ok = queue.Poll()
	assert.Equal(t, "b", element)
	assert.True(t, ok)

	assert.True(t, queue.Remove(5))

	element, ok = queue.Poll()
	assert.Equal(t, "h", element)
	assert.True(t, ok)

	assert.True(t, queue.Remove(1))

	element, ok = queue.Poll()
	assert.Equal(t, "j", element)
	assert.True(t, ok)

	assert.True(t, queue.Empty())
	assert.Equal(t, doubleCapacity, queue.Capacity())
}

func TestQueueAddGrowShrink(t *testing.T) {
	queue := New(GrowthFactor(2.0, 0.5))
	queue.Add("a", "b", "c", "d")
	queue.Add("e", "f", "g", "h")

	assert.False(t, queue.Empty())

	capacity := 8
	assert.Equal(t, capacity, queue.Capacity())
	assert.Equal(t, capacity, queue.Size())

	element, ok := queue.Peek()
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	queue.Add("i", "j")

	doubleCapacity := int((capacity + 2) * 2.0)
	assert.Equal(t, doubleCapacity, queue.Capacity())
	assert.Equal(t, 10, queue.Size())

	element, ok = queue.Poll()
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	assert.Equal(t, 9, queue.Capacity())
	assert.Equal(t, 9, queue.Size())

	element, ok = queue.Poll()
	assert.Equal(t, "b", element)
	assert.True(t, ok)

	assert.True(t, queue.Remove(5))

	element, ok = queue.Poll()
	assert.Equal(t, "h", element)
	assert.True(t, ok)

	assert.True(t, queue.Remove(1))

	element, ok = queue.Poll()
	assert.Equal(t, "j", element)
	assert.True(t, ok)

	assert.True(t, queue.Empty())
	assert.Equal(t, 9, queue.Capacity())
}

func TestQueueShrinkGuaranteedSize(t *testing.T) {
	queue := New(GrowthFactor(2.0, 0.5), GuaranteedSize(10))
	queue.Add("a", "b", "c", "d")
	queue.Add("e", "f", "g", "h")

	assert.False(t, queue.Empty())

	capacity := 10
	assert.Equal(t, capacity, queue.Capacity())
	assert.Equal(t, 8, queue.Size())

	element, ok := queue.Peek()
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	queue.Add("i", "j")

	assert.Equal(t, capacity, queue.Capacity())
	assert.Equal(t, capacity, queue.Size())

	queue.Add("k")

	doubleCapacity := int((capacity + 1) * 2.0)
	assert.Equal(t, doubleCapacity, queue.Capacity())
	assert.Equal(t, 11, queue.Size())

	element, ok = queue.Poll()
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	assert.Equal(t, capacity, queue.Capacity())
	assert.Equal(t, capacity, queue.Size())

	element, ok = queue.Poll()
	assert.Equal(t, "b", element)
	assert.True(t, ok)

	assert.True(t, queue.Remove(6))

	element, ok = queue.Poll()
	assert.Equal(t, "i", element)
	assert.True(t, ok)

	assert.True(t, queue.Remove(1))

	element, ok = queue.Poll()
	assert.Equal(t, "k", element)
	assert.True(t, ok)

	assert.True(t, queue.Empty())
	assert.Equal(t, capacity, queue.Capacity())
}

func TestQueuePeek(t *testing.T) {
	queue := New()

	element, ok := queue.Peek()
	assert.Nil(t, element)
	assert.False(t, ok)

	queue.Add(1)
	queue.Add(2)
	queue.Add(3)

	element, ok = queue.Peek()
	assert.Equal(t, 1, element)
	assert.True(t, ok)
}

func TestQueuePoll(t *testing.T) {
	queue := New()

	queue.Add(1)
	queue.Add(2)
	queue.Add(3)

	element, ok := queue.Poll()
	assert.Equal(t, 1, element)
	assert.True(t, ok)

	element, ok = queue.Peek()
	assert.Equal(t, 2, element)
	assert.True(t, ok)

	element, ok = queue.Poll()
	assert.Equal(t, 2, element)
	assert.True(t, ok)

	element, ok = queue.Poll()
	assert.Equal(t, 3, element)
	assert.True(t, ok)

	element, ok = queue.Poll()
	assert.Nil(t, element)
	assert.False(t, ok)

	assert.True(t, queue.Empty())
	assert.Zero(t, queue.Size())
}

func TestQueueIndexOf(t *testing.T) {
	queue := New()

	assert.Equal(t, -1, queue.IndexOf(func(interface{}) bool {
		return true
	}))

	queue.Add("a")
	queue.Add("b", "c")

	dataset := []struct {
		name      string
		predicate func(interface{}) bool
		expected  int
	}{
		{"PredicateA", func(v interface{}) bool { return v.(string) == "a" }, 0},
		{"PredicateB", func(v interface{}) bool { return v.(string) == "b" }, 1},
		{"PredicateC", func(v interface{}) bool { return v.(string) == "c" }, 2},
		{"PredicateD", func(v interface{}) bool { return v.(string) == "d" }, -1},
	}

	for _, data := range dataset {
		t.Run(data.name, func(t *testing.T) {
			assert.Equal(t, data.expected, queue.IndexOf(data.predicate))
		})
	}
}

func TestQueueTake(t *testing.T) {
	queue := New()
	queue.Add("a", "b", "c")
	queue.Add("d", "e", "f")

	index := queue.IndexOf(func(v interface{}) bool { return v.(string) == "c" })
	assert.Equal(t, 2, index)

	count := index
	elements, ok := queue.Take(count)
	assert.Len(t, elements, count)
	assert.Equal(t, []interface{}{"a", "b"}, elements)
	assert.True(t, ok)

	count = queue.Size()
	ok = queue.Remove(count)
	assert.True(t, ok)
	assert.True(t, queue.Empty())

	_, ok = queue.Take(1)
	assert.False(t, ok)
	assert.True(t, queue.Empty())

	ok = queue.Remove(1)
	assert.False(t, ok)
	assert.True(t, queue.Empty())
}

func TestQueueRemove(t *testing.T) {
	queue := New()
	queue.Add("a")
	queue.Add("b", "c")

	assert.True(t, queue.Remove(2))
	assert.False(t, queue.Remove(2))
	assert.True(t, queue.Remove(1))
	assert.False(t, queue.Remove(1))

	queue.Add("d", "e")
	assert.True(t, queue.Remove(1))
	assert.True(t, queue.Remove(1))
	assert.False(t, queue.Remove(1))

	assert.True(t, queue.Empty())
	assert.Zero(t, queue.Size())
}

func TestQueueElement(t *testing.T) {
	queue := New()
	queue.Add("a")
	queue.Add("b", "c")

	element, ok := queue.Element(0)
	assert.Equal(t, "a", element)
	assert.True(t, ok)

	element, ok = queue.Element(1)
	assert.Equal(t, "b", element)
	assert.True(t, ok)

	element, ok = queue.Element(2)
	assert.Equal(t, "c", element)
	assert.True(t, ok)

	element, ok = queue.Element(3)
	assert.Nil(t, element)
	assert.False(t, ok)

	queue.Remove(1)
	element, ok = queue.Element(0)
	assert.Equal(t, "b", element)
	assert.True(t, ok)
}

func TestQueueClear(t *testing.T) {
	queue := New(FixedSize(16))
	queue.Add("a", "b", "c", "d", "e", "f", "g", "h")

	assert.Equal(t, 16, queue.Capacity())
	assert.Equal(t, 8, queue.Size())

	queue.Clear()
	assert.True(t, queue.Empty())
	assert.Zero(t, queue.Size())
	assert.Equal(t, 16, queue.Capacity())
}

func TestQueueValues(t *testing.T) {
	queue := New()

	list := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h"}
	queue.Add(list...)

	assert.False(t, queue.Empty())

	values := queue.Values()
	assert.Equal(t, list, values)

	// TODO: empty queue values test

	// TODO: separate string test
	var vstrings []string
	for _, v := range list {
		vstrings = append(vstrings, v.(string))
	}

	qstring := queue.String()
	assert.True(t, strings.Contains(qstring, "CircularQueue"))
	assert.True(t, strings.Contains(qstring, strings.Join(vstrings, ", ")))
}

func benchmarkElement(b *testing.B, queue *Queue, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			queue.Element(n)
		}
	}
}

func benchmarkAdd(b *testing.B, queue *Queue, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			queue.Add(n)
		}
	}
}

func benchmarkRemove(b *testing.B, queue *Queue, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			queue.Remove(1)
		}
	}
}

func BenchmarkArrayqueueGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkElement(b, queue, size)
}

func BenchmarkArrayqueueGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkElement(b, queue, size)
}

func BenchmarkArrayqueueGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkElement(b, queue, size)
}

func BenchmarkArrayqueueGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkElement(b, queue, size)
}

func BenchmarkArrayqueueAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := New()
	b.StartTimer()
	benchmarkAdd(b, queue, size)
}

func BenchmarkArrayqueueAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, queue, size)
}

func BenchmarkArrayqueueAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, queue, size)
}

func BenchmarkArrayqueueAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, queue, size)
}

func BenchmarkArrayqueueRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, queue, size)
}

func BenchmarkArrayqueueRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, queue, size)
}

func BenchmarkArrayqueueRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, queue, size)
}

func BenchmarkArrayqueueRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	queue := New()
	for n := 0; n < size; n++ {
		queue.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, queue, size)
}

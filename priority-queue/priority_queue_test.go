package pq_test

import (
	"testing"
	"time"

	"math/rand"

	"math"

	"github.com/kyroy/priority-queue"
	"github.com/stretchr/testify/assert"
)

const (
	randMultiplier = 3000.
)

func TestInsertAndLen(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
	}{
		{name: "1", input: generateRandomInput(1)},
		{name: "100", input: generateRandomInput(100)},
		{name: "1000", input: generateRandomInput(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := pq.NewPriorityQueue()
			length := 0
			assert.Equal(t, queue.Len(), length)
			for _, i := range test.input {
				queue.Insert(i, i)
				length++
				assert.Equal(t, queue.Len(), length)
				checkSorted(t, queue)
			}
		})
	}
}

func TestPopLowest(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		queue := pq.NewPriorityQueue()
		assert.Nil(t, queue.PopLowest())
	})

	tests := []struct {
		name  string
		input []float64
	}{
		{name: "10", input: generateRandomInput(10)},
		{name: "100", input: generateRandomInput(100)},
		{name: "1000", input: generateRandomInput(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := pq.NewPriorityQueue()
			for _, i := range test.input {
				queue.Insert(i, i)
			}
			lastVal := queue.PopLowest().(float64)
			for currentInt, _ := queue.Get(0); queue.Len() > 1; currentInt, _ = queue.Get(0) {
				currentVal := currentInt.(float64)
				assert.Equal(t, queue.PopLowest().(float64), currentVal)
				assert.True(t, lastVal <= currentVal, "expect %f <= %f", lastVal, currentVal)
				lastVal = currentVal
			}

		})

	}
}

func TestPopHighest(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		queue := pq.NewPriorityQueue()
		assert.Nil(t, queue.PopHighest())
	})

	tests := []struct {
		name  string
		input []float64
	}{
		{name: "10", input: generateRandomInput(10)},
		{name: "100", input: generateRandomInput(100)},
		{name: "1000", input: generateRandomInput(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := pq.NewPriorityQueue()
			for _, i := range test.input {
				queue.Insert(i, i)
			}
			lastVal := queue.PopHighest().(float64)
			i := 0
			for currentInt, _ := queue.Get(queue.Len() - 1); queue.Len() > 1; currentInt, _ = queue.Get(queue.Len() - 1) {
				currentVal := currentInt.(float64)
				assert.Equal(t, queue.PopHighest().(float64), currentVal)
				assert.True(t, lastVal >= currentVal, "expect %f >= %f", i, lastVal, currentVal)
				lastVal = currentVal
			}

		})

	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
	}{
		{name: "empty", input: []float64{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := pq.NewPriorityQueue()
			for _, itm := range test.input {
				queue.Insert(itm, itm)
			}
			lastPrio := -randMultiplier
			for i := 0; i < len(test.input); i++ {
				_, prio := queue.Get(i)
				assert.True(t, lastPrio <= prio, "[%d] expect %f <= %f", i, prio, lastPrio)
				lastPrio = prio
			}
		})
	}
}

func TestMinPrioSize(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		input []float64
	}{
		{name: "10", size: 5, input: generateRandomInput(10)},
		{name: "100", size: 10, input: generateRandomInput(100)},
		{name: "1000", size: 20, input: generateRandomInput(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lastPrio := randMultiplier
			queue := pq.NewPriorityQueue(pq.WithMinPrioSize(test.size))
			for _, itm := range test.input {
				queue.Insert(itm, itm)
				_, smallestElem := queue.Get(0)
				assert.True(t, smallestElem <= lastPrio, "expect %f <= %f", smallestElem, lastPrio)
				lastPrio = math.Min(lastPrio, smallestElem)
				assert.True(t, queue.Len() <= test.size)
			}
		})
	}
}

func TestMaxPrioSize(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		input []float64
	}{
		{name: "10", size: 5, input: generateRandomInput(10)},
		{name: "100", size: 10, input: generateRandomInput(100)},
		{name: "1000", size: 20, input: generateRandomInput(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lastPrio := -randMultiplier
			queue := pq.NewPriorityQueue(pq.WithMaxPrioSize(test.size))
			for _, itm := range test.input {
				queue.Insert(itm, itm)
				_, largestElem := queue.Get(queue.Len() - 1)
				assert.True(t, lastPrio <= largestElem, "expect %f <= %f", lastPrio, largestElem)
				lastPrio = math.Max(lastPrio, largestElem)
				assert.True(t, queue.Len() <= test.size)
			}
		})
	}
}

func generateRandomInput(size int) []float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	input := make([]float64, size)
	for i := 0; i < size; i++ {
		input[i] = r.Float64()*randMultiplier - (randMultiplier / 2.)
	}
	return input
}

func checkSorted(t *testing.T, queue *pq.PriorityQueue) {
	lastPrio := -randMultiplier
	for i := 0; i < queue.Len(); i++ {
		_, prio := queue.Get(i)
		assert.True(t, lastPrio <= prio, "[%d] expect %f <= %f", i, lastPrio, prio)
		lastPrio = prio
	}
}

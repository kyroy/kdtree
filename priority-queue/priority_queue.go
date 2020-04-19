// Package pq implements a priority queue data structure on top of golang's sort.Interface.
package pq

import (
	"sort"
)

// PriorityQueue represents the queue
type PriorityQueue struct {
	items      *items
	sizeOption Option
}

type optionEnum int

const (
	none optionEnum = iota
	minPrioSize
	maxPrioSize
)

// Option represents the option
type Option struct {
	option optionEnum
	value  int
}

// NewPriorityQueue returns an initialized PriorityQueue with the given options
func NewPriorityQueue(options ...Option) *PriorityQueue {
	pq := PriorityQueue{
		items: &items{},
	}
	for _, o := range options {
		switch o.option {
		case none, minPrioSize, maxPrioSize:
			pq.sizeOption = o
		}
	}
	return &pq
}

// WithMinPrioSize limits the size of the priority queue to 'size' by keeping the elements with the lowest priorities.
func WithMinPrioSize(size int) Option {
	return Option{
		option: minPrioSize,
		value:  size,
	}
}

// WithMaxPrioSize limits the size of the priority queue to 'size' by keeping the elements with the highest priorities.
func WithMaxPrioSize(size int) Option {
	return Option{
		option: maxPrioSize,
		value:  size,
	}
}

// Len returns the number of elements in the queue.
func (p *PriorityQueue) Len() int {
	return p.items.Len()
}

// Insert inserts a new element into the queue.
func (p *PriorityQueue) Insert(v interface{}, priority float64) {
	*p.items = append(*p.items, &item{value: v, priority: priority})
	sort.Sort(p.items)
	switch p.sizeOption.option {
	case minPrioSize:
		if p.sizeOption.value < len(*p.items) {
			*p.items = (*p.items)[:p.sizeOption.value]
		}
	case maxPrioSize:
		diff := len(*p.items) - p.sizeOption.value
		if diff > 0 {
			*p.items = (*p.items)[diff:]
		}
	}
}

// PopLowest removes the element with the lowest priority from the queue and returns it.
// If the queue is empty, nil is returned.
func (p *PriorityQueue) PopLowest() interface{} {
	if len(*p.items) == 0 {
		return nil
	}
	x := (*p.items)[0]
	*p.items = (*p.items)[1:]
	return x.value
}

// PopHighest removes the element with the highest priority from the queue and returns it.
// If the queue is empty, nil is returned.
func (p *PriorityQueue) PopHighest() interface{} {
	l := len(*p.items) - 1
	if l < 0 {
		return nil
	}
	x := (*p.items)[l]
	*p.items = (*p.items)[:l]
	return x.value
}

// Get returns the element with the i-th priority (from low to high).
// So the element with the lowest priority has the index 0.
// And the element with the highest priority has the index queue.Len()-1.
func (p *PriorityQueue) Get(i int) (interface{}, float64) {
	x := (*p.items)[i]
	return x.value, x.priority
}

type items []*item

func (p items) Len() int           { return len(p) }
func (p items) Less(i, j int) bool { return p[i].priority < p[j].priority }
func (p items) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type item struct {
	value    interface{}
	priority float64
}

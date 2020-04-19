/*
 * Copyright 2020 Dennis Kuhnert
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package kdtree implements a k-d tree data structure.
package kdtree

import (
	"fmt"
	"github.com/kyroy/kdtree/kdrange"
	"github.com/kyroy/priority-queue"
	"math"
	"sort"
)

// Point specifies one element of the k-d tree.
type Point interface {
	// Dimensions returns the total number of dimensions.
	Dimensions() int
	// Dimension returns the value of the i-th dimension.
	Dimension(i int) float64
}

// KDTree represents the k-d tree.
type KDTree struct {
	root *node
}

// New returns a balanced k-d tree.
func New(points []Point) *KDTree {
	return &KDTree{
		root: newKDTree(points, 0),
	}
}

func newKDTree(points []Point, axis int) *node {
	if len(points) == 0 {
		return nil
	}
	if len(points) == 1 {
		return &node{Point: points[0]}
	}

	sort.Sort(&byDimension{dimension: axis, points: points})
	mid := len(points) / 2
	root := points[mid]
	nextDim := (axis + 1) % root.Dimensions()
	return &node{
		Point: root,
		Left:  newKDTree(points[:mid], nextDim),
		Right: newKDTree(points[mid+1:], nextDim),
	}
}

// String returns a string representation of the k-d tree.
func (t *KDTree) String() string {
	return fmt.Sprintf("[%s]", printTreeNode(t.root))
}

func printTreeNode(n *node) string {
	if n != nil && (n.Left != nil || n.Right != nil) {
		return fmt.Sprintf("[%s %s %s]", printTreeNode(n.Left), n.String(), printTreeNode(n.Right))
	}
	return fmt.Sprintf("%s", n)
}

// Insert adds a point to the k-d tree.
func (t *KDTree) Insert(p Point) {
	if t.root == nil {
		t.root = &node{Point: p}
	} else {
		t.root.Insert(p, 0)
	}
}

// Remove removes and returns the first point from the tree that equals the given point p in all dimensions.
// Returns nil if not found.
func (t *KDTree) Remove(p Point) Point {
	if t.root == nil || p == nil {
		return nil
	}
	n, sub := t.root.Remove(p, 0)
	if n == t.root {
		t.root = sub
	}
	if n == nil {
		return nil
	}
	return n.Point
}

// Balance rebalances the k-d tree by recreating it.
func (t *KDTree) Balance() {
	t.root = newKDTree(t.Points(), 0)
}

// Points returns all points in the k-d tree.
// The tree is traversed in-order.
func (t *KDTree) Points() []Point {
	if t.root == nil {
		return []Point{}
	}
	return t.root.Points()
}

// KNN returns the k-nearest neighbours of the given point.
// The points are sorted by the distance to the given points. Starting with the nearest.
func (t *KDTree) KNN(p Point, k int) []Point {
	if t.root == nil || p == nil || k == 0 {
		return []Point{}
	}

	nearestPQ := pq.NewPriorityQueue(pq.WithMinPrioSize(k))
	knn(p, k, t.root, 0, nearestPQ)

	points := make([]Point, 0, k)
	for i := 0; i < k && 0 < nearestPQ.Len(); i++ {
		o := nearestPQ.PopLowest().(*node).Point
		points = append(points, o)
	}

	return points
}

// RangeSearch returns all points in the given range r.
//
// Returns an empty slice when input is nil or len(r) does not equal Point.Dimensions().
func (t *KDTree) RangeSearch(r kdrange.Range) []Point {
	if t.root == nil || r == nil || len(r) != t.root.Dimensions() {
		return []Point{}
	}

	return t.root.RangeSearch(r, 0)
}

func knn(p Point, k int, start *node, currentAxis int, nearestPQ *pq.PriorityQueue) {
	if p == nil || k == 0 || start == nil {
		return
	}

	var path []*node
	currentNode := start

	// 1. move down
	for currentNode != nil {
		path = append(path, currentNode)
		if p.Dimension(currentAxis) < currentNode.Dimension(currentAxis) {
			currentNode = currentNode.Left
		} else {
			currentNode = currentNode.Right
		}
		currentAxis = (currentAxis + 1) % p.Dimensions()
	}

	// 2. move up
	currentAxis = (currentAxis - 1 + p.Dimensions()) % p.Dimensions()
	for path, currentNode = popLast(path); currentNode != nil; path, currentNode = popLast(path) {
		currentDistance := distance(p, currentNode)
		checkedDistance := getKthOrLastDistance(nearestPQ, k-1)
		if currentDistance < checkedDistance {
			nearestPQ.Insert(currentNode, currentDistance)
			checkedDistance = getKthOrLastDistance(nearestPQ, k-1)
		}

		// check other side of plane
		if planeDistance(p, currentNode.Dimension(currentAxis), currentAxis) < checkedDistance {
			var next *node
			if p.Dimension(currentAxis) < currentNode.Dimension(currentAxis) {
				next = currentNode.Right
			} else {
				next = currentNode.Left
			}
			knn(p, k, next, (currentAxis+1)%p.Dimensions(), nearestPQ)
		}
		currentAxis = (currentAxis - 1 + p.Dimensions()) % p.Dimensions()
	}
}

func distance(p1, p2 Point) float64 {
	sum := 0.
	for i := 0; i < p1.Dimensions(); i++ {
		sum += math.Pow(p1.Dimension(i)-p2.Dimension(i), 2.0)
	}
	return math.Sqrt(sum)
}

func planeDistance(p Point, planePosition float64, dim int) float64 {
	return math.Abs(planePosition - p.Dimension(dim))
}

func popLast(arr []*node) ([]*node, *node) {
	l := len(arr) - 1
	if l < 0 {
		return arr, nil
	}
	return arr[:l], arr[l]
}

func getKthOrLastDistance(nearestPQ *pq.PriorityQueue, i int) float64 {
	if nearestPQ.Len() <= i {
		return math.MaxFloat64
	}
	_, prio := nearestPQ.Get(i)
	return prio
}

//
//
// byDimension
//

type byDimension struct {
	dimension int
	points    []Point
}

// Len is the number of elements in the collection.
func (b *byDimension) Len() int {
	return len(b.points)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (b *byDimension) Less(i, j int) bool {
	return b.points[i].Dimension(b.dimension) < b.points[j].Dimension(b.dimension)
}

// Swap swaps the elements with indexes i and j.
func (b *byDimension) Swap(i, j int) {
	b.points[i], b.points[j] = b.points[j], b.points[i]
}

//
//
// node
//

type node struct {
	Point
	Left  *node
	Right *node
}

func (n *node) String() string {
	return fmt.Sprintf("%v", n.Point)
}

func (n *node) Points() []Point {
	var points []Point
	if n.Left != nil {
		points = n.Left.Points()
	}
	points = append(points, n.Point)
	if n.Right != nil {
		points = append(points, n.Right.Points()...)
	}
	return points
}

func (n *node) Insert(p Point, axis int) {
	if p.Dimension(axis) < n.Point.Dimension(axis) {
		if n.Left == nil {
			n.Left = &node{Point: p}
		} else {
			n.Left.Insert(p, (axis+1)%n.Point.Dimensions())
		}
	} else {
		if n.Right == nil {
			n.Right = &node{Point: p}
		} else {
			n.Right.Insert(p, (axis+1)%n.Point.Dimensions())
		}
	}
}

// Remove returns (returned node, substitute node)
func (n *node) Remove(p Point, axis int) (*node, *node) {
	for i := 0; i < n.Dimensions(); i++ {
		if n.Dimension(i) != p.Dimension(i) {
			if n.Left != nil {
				returnedNode, substitutedNode := n.Left.Remove(p, (axis+1)%n.Dimensions())
				if returnedNode != nil {
					if returnedNode == n.Left {
						n.Left = substitutedNode
					}
					return returnedNode, nil
				}
			}
			if n.Right != nil {
				returnedNode, substitutedNode := n.Right.Remove(p, (axis+1)%n.Dimensions())
				if returnedNode != nil {
					if returnedNode == n.Right {
						n.Right = substitutedNode
					}
					return returnedNode, nil
				}
			}
			return nil, nil
		}
	}

	// equals, remove n

	if n.Left != nil {
		largest := n.Left.FindLargest(axis, nil)
		removed, sub := n.Left.Remove(largest, (axis+1)%n.Dimensions())

		removed.Left = n.Left
		removed.Right = n.Right
		if n.Left == removed {
			removed.Left = sub
		}
		return n, removed
	}

	if n.Right != nil {
		smallest := n.Right.FindSmallest(axis, nil)
		removed, sub := n.Right.Remove(smallest, (axis+1)%n.Dimensions())

		removed.Left = n.Left
		removed.Right = n.Right
		if n.Right == removed {
			removed.Right = sub
		}
		return n, removed
	}

	// n.Left == nil && n.Right == nil
	return n, nil
}

func (n *node) FindSmallest(axis int, smallest *node) *node {
	if smallest == nil || n.Dimension(axis) < smallest.Dimension(axis) {
		smallest = n
	}
	if n.Left != nil {
		smallest = n.Left.FindSmallest(axis, smallest)
	}
	if n.Right != nil {
		smallest = n.Right.FindSmallest(axis, smallest)
	}
	return smallest
}

func (n *node) FindLargest(axis int, largest *node) *node {
	if largest == nil || n.Dimension(axis) > largest.Dimension(axis) {
		largest = n
	}
	if n.Left != nil {
		largest = n.Left.FindLargest(axis, largest)
	}
	if n.Right != nil {
		largest = n.Right.FindLargest(axis, largest)
	}
	return largest
}

func (n *node) RangeSearch(r kdrange.Range, axis int) []Point {
	points := []Point{}

	for dim, limit := range r {
		if limit[0] > n.Dimension(dim) || limit[1] < n.Dimension(dim) {
			goto checkChildren
		}
	}
	points = append(points, n.Point)

checkChildren:
	if n.Left != nil && n.Dimension(axis) >= r[axis][0] {
		points = append(points, n.Left.RangeSearch(r, (axis+1)%n.Dimensions())...)
	}
	if n.Right != nil && n.Dimension(axis) <= r[axis][1] {
		points = append(points, n.Right.RangeSearch(r, (axis+1)%n.Dimensions())...)
	}

	return points
}

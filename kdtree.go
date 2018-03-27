// Package kdtree implements a k-d tree data structure.
package kdtree

import (
	"fmt"
	"math"
	"sort"

	"github.com/kyroy/priority-queue"
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

// NewKDTree returns an initialized k-d tree.
func NewKDTree(points []Point) *KDTree {
	return &KDTree{
		root: newKDTree(points, 0),
	}
}

func newKDTree(points []Point, axis int) *node {
	if points == nil || len(points) == 0 {
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

// Remove // TODO planned
//func (t *KDTree) Remove(p Point) {
//	// requires equals method? or based in Dim()?
//}

// Points returns all points in the k-d tree.
// The tree is traversed in-order.
func (t *KDTree) Points() []Point {
	if t.root == nil {
		return []Point{}
	}
	return t.root.Points()
}

// KNN returns the k-nearest neighbours of the given point.
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

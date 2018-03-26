package kdtree_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/jupp0r/go-priority-queue"
	"github.com/kyroy/kdtree"
	. "github.com/kyroy/kdtree/points"
	"github.com/stretchr/testify/assert"
)

func TestNewKDTree(t *testing.T) {
	tests := []struct {
		name   string
		input  []kdtree.Point
		output []kdtree.Point
	}{
		{
			name:   "nil",
			input:  nil,
			output: []kdtree.Point{},
		},
		{
			name:   "empty",
			input:  []kdtree.Point{},
			output: []kdtree.Point{},
		},
		{
			name:   "1",
			input:  []kdtree.Point{&Point2D{X: 1., Y: 2.}},
			output: []kdtree.Point{&Point2D{X: 1., Y: 2.}},
		},
		{
			name:   "2 equal",
			input:  []kdtree.Point{&Point2D{X: 1., Y: 2.}, &Point2D{X: 1., Y: 2.}},
			output: []kdtree.Point{&Point2D{X: 1., Y: 2.}, &Point2D{X: 1., Y: 2.}},
		},
		{
			name:   "sort 1 dim",
			input:  []kdtree.Point{&Point2D{X: 1.1, Y: 1.2}, &Point2D{X: 1.3, Y: 1.0}, &Point2D{X: 0.9, Y: 1.3}},
			output: []kdtree.Point{&Point2D{X: 0.9, Y: 1.3}, &Point2D{X: 1.1, Y: 1.2}, &Point2D{X: 1.3, Y: 1.0}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := kdtree.NewKDTree(test.input)
			assert.Equal(t, test.output, tree.Points())
		})

	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		tree     *kdtree.KDTree
		expected string
	}{
		{name: "empty", tree: &kdtree.KDTree{}, expected: "[<nil>]"},
		{name: "1 elem", tree: kdtree.NewKDTree([]kdtree.Point{&Point2D{X: 2, Y: 3}}), expected: "[{2.00 3.00}]"},
		{name: "2 elem", tree: kdtree.NewKDTree([]kdtree.Point{&Point2D{X: 2, Y: 3}, &Point2D{X: 3.4, Y: 1}}), expected: "[[{2.00 3.00} {3.40 1.00} <nil>]]"},
		{name: "3 elem", tree: kdtree.NewKDTree([]kdtree.Point{&Point2D{X: 2, Y: 3}, &Point2D{X: 1.4, Y: 7.1}, &Point2D{X: 3.4, Y: 1}}), expected: "[[{1.40 7.10} {2.00 3.00} {3.40 1.00}]]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.tree.String())
		})
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name      string
		treeInput *kdtree.KDTree
		input     []kdtree.Point
		output    []kdtree.Point
	}{
		{
			name:   "empty tree",
			input:  []kdtree.Point{&Point2D{X: 1., Y: 2.}},
			output: []kdtree.Point{&Point2D{X: 1., Y: 2.}},
		},
		{
			name:      "1 dim",
			treeInput: kdtree.NewKDTree([]kdtree.Point{&Point2D{X: 1., Y: 2.}}),
			input:     []kdtree.Point{&Point2D{X: 0.9, Y: 2.1}},
			output:    []kdtree.Point{&Point2D{X: 0.9, Y: 2.1}, &Point2D{X: 1., Y: 2.}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.treeInput == nil {
				test.treeInput = kdtree.NewKDTree(nil)
			}
			for _, elem := range test.input {
				test.treeInput.Insert(elem)
			}
			assert.Equal(t, test.output, test.treeInput.Points())
		})
	}
}

func TestInsertWithGenerator(t *testing.T) {
	tests := []struct {
		name  string
		input []kdtree.Point
	}{
		{name: "p:10,k:5", input: generateLargeCaseData(10)},
		{name: "p:100,k:5", input: generateLargeCaseData(100)},
		{name: "p:1000,k:5", input: generateLargeCaseData(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := kdtree.NewKDTree(nil)
			for _, elem := range test.input {
				tree.Insert(elem)
				_ = tree.Points()
				// TODO assert
				//assert.Equal(t, , treePoints)
			}
		})
	}
}

// TestKNN ...
func TestKNN(t *testing.T) {
	tests := []struct {
		name   string
		target kdtree.Point
		k      int
		input  []kdtree.Point
		output []kdtree.Point
	}{
		{
			name:   "nil",
			target: nil,
			k:      3,
			input:  []kdtree.Point{&Point2D{X: 1., Y: 2.}},
			output: []kdtree.Point{},
		},
		{
			name:   "empty",
			target: &Point2D{X: 1., Y: 2.},
			k:      3,
			input:  []kdtree.Point{},
			output: []kdtree.Point{},
		},
		{
			name:   "k >> points",
			target: &Point2D{X: 1., Y: 2.},
			k:      10,
			input:  []kdtree.Point{&Point2D{X: 1., Y: 2.}, &Point2D{X: 0.9, Y: 2.1}, &Point2D{X: 1.1, Y: 1.9}},
			output: []kdtree.Point{&Point2D{X: 1., Y: 2.}, &Point2D{X: 0.9, Y: 2.1}, &Point2D{X: 1.1, Y: 1.9}},
		},
		{
			name:   "small 2D example",
			target: &Point2D{X: 9, Y: 4},
			k:      3,
			input:  []kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}},
			output: []kdtree.Point{&Point2D{X: 8, Y: 5}, &Point2D{X: 7, Y: 4}, &Point2D{X: 8, Y: 2}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := kdtree.NewKDTree(test.input)
			assert.Equal(t, test.output, tree.KNN(test.target, test.k))
		})
	}
}

func TestKNNWithGenerator(t *testing.T) {
	tests := []struct {
		name   string
		target kdtree.Point
		k      int
		input  []kdtree.Point
	}{
		{name: "p:100,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(100)},
		{name: "p:1000,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(1000)},
		{name: "p:10000,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(10000)},
		{name: "p:100000,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(100000)},
		{name: "p:1000000,k:10", target: &Point2D{}, k: 10, input: generateLargeCaseData(1000000)},
		{name: "p:1000000,k:20", target: &Point2D{}, k: 20, input: generateLargeCaseData(1000000)},
		{name: "p:1000000,k:30", target: &Point2D{}, k: 30, input: generateLargeCaseData(1000000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := kdtree.NewKDTree(test.input)
			assert.Equal(t, prioQueueKNN(test.input, test.target, test.k), tree.KNN(test.target, test.k))
		})
	}
}

// benchmarks

var resultTree *kdtree.KDTree
var resultPoints []kdtree.Point

// BenchmarkNewKDTree ...
func BenchmarkNewKDTree(b *testing.B) {
	benchmarks := []struct {
		name  string
		input []kdtree.Point
	}{
		{name: "100", input: generateLargeCaseData(100)},
		{name: "1000", input: generateLargeCaseData(1000)},
		{name: "10000", input: generateLargeCaseData(10000)},
		{name: "100000", input: generateLargeCaseData(100000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var t *kdtree.KDTree
			for i := 0; i < b.N; i++ {
				t = kdtree.NewKDTree(bm.input)
			}
			resultTree = t
		})
	}
}

// BenchmarkKNN ...
func BenchmarkKNN(b *testing.B) {
	benchmarks := []struct {
		name   string
		target kdtree.Point
		k      int
		input  []kdtree.Point
	}{
		{name: "p:100,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(100)},
		{name: "p:1000,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(1000)},
		{name: "p:10000,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(10000)},
		{name: "p:100000,k:5", target: &Point2D{}, k: 5, input: generateLargeCaseData(100000)},
	}
	for _, bm := range benchmarks {
		var res []kdtree.Point
		tree := kdtree.NewKDTree(bm.input)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				res = tree.KNN(bm.target, bm.k)
			}
			resultPoints = res
		})
	}
}

// helpers

func generateLargeCaseData(size int) []kdtree.Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var points []kdtree.Point
	for i := 0; i < size; i++ {
		points = append(points, &Point2D{X: r.Float64()*3000 - 1500, Y: r.Float64()*3000 - 1500})
	}

	return points
}

func prioQueueKNN(points []kdtree.Point, p kdtree.Point, k int) []kdtree.Point {
	knn := make([]kdtree.Point, 0, k)
	if p == nil {
		return knn
	}

	nnPQ := pq.New()
	for _, point := range points {
		nnPQ.Insert(point, distance(p, point))
	}

	for i := 0; i < k; i++ {
		point, err := nnPQ.Pop()
		if err != nil {
			break
		}
		knn = append(knn, point.(kdtree.Point))
	}
	return knn
}

func distance(p1, p2 kdtree.Point) float64 {
	sum := 0.
	for i := 0; i < p1.Dimensions(); i++ {
		sum += math.Pow(p1.Dimension(i)-p2.Dimension(i), 2.0)
	}
	return math.Sqrt(sum)
}

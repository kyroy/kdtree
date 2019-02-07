/*
 * Copyright 2018 Dennis Kuhnert
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

package kdtree

import (
	"math/rand"
	"testing"
	"time"

	"github.com/jupp0r/go-priority-queue"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		input  []Point
		output []Point
	}{
		{
			name:   "nil",
			input:  nil,
			output: []Point{},
		},
		{
			name:   "empty",
			input:  []Point{},
			output: []Point{},
		},
		{
			name:   "1",
			input:  []Point{&SamplePoint2D{X: 1., Y: 2.}},
			output: []Point{&SamplePoint2D{X: 1., Y: 2.}},
		},
		{
			name:   "2 equal",
			input:  []Point{&SamplePoint2D{X: 1., Y: 2.}, &SamplePoint2D{X: 1., Y: 2.}},
			output: []Point{&SamplePoint2D{X: 1., Y: 2.}, &SamplePoint2D{X: 1., Y: 2.}},
		},
		{
			name: "sort 1 dim",
			input: []Point{&SamplePoint2D{X: 1.1, Y: 1.2},
				&SamplePoint2D{X: 1.3, Y: 1.0},
				&SamplePoint2D{X: 0.9, Y: 1.3}},
			output: []Point{&SamplePoint2D{X: 0.9, Y: 1.3},
				&SamplePoint2D{X: 1.1, Y: 1.2},
				&SamplePoint2D{X: 1.3, Y: 1.0}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := New(test.input)
			assert.Equal(t, test.output, tree.Points())
		})

	}
}

func TestKDTree_String(t *testing.T) {
	tests := []struct {
		name     string
		tree     *KDTree
		expected string
	}{
		{name: "empty", tree: &KDTree{}, expected: "[<nil>]"},
		{name: "1 elem", tree: New([]Point{&SamplePoint2D{X: 2, Y: 3}}), expected: "[{2.00 3.00}]"},
		{name: "2 elem", tree: New([]Point{&SamplePoint2D{X: 2, Y: 3},
			&SamplePoint2D{X: 3.4, Y: 1}}),
			expected: "[[{2.00 3.00} {3.40 1.00} <nil>]]"},
		{name: "3 elem", tree: New([]Point{&SamplePoint2D{X: 2, Y: 3},
			&SamplePoint2D{X: 1.4, Y: 7.1},
			&SamplePoint2D{X: 3.4, Y: 1}}),
			expected: "[[{1.40 7.10} {2.00 3.00} {3.40 1.00}]]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.tree.String())
		})
	}
}

func TestKDTree_Insert(t *testing.T) {
	tests := []struct {
		name      string
		treeInput *KDTree
		input     []Point
		output    []Point
	}{
		{
			name:   "empty tree",
			input:  []Point{&SamplePoint2D{X: 1., Y: 2.}},
			output: []Point{&SamplePoint2D{X: 1., Y: 2.}},
		},
		{
			name:      "1 dim",
			treeInput: New([]Point{&SamplePoint2D{X: 1., Y: 2.}}),
			input:     []Point{&SamplePoint2D{X: 0.9, Y: 2.1}},
			output:    []Point{&SamplePoint2D{X: 0.9, Y: 2.1}, &SamplePoint2D{X: 1., Y: 2.}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.treeInput == nil {
				test.treeInput = New(nil)
			}
			for _, elem := range test.input {
				test.treeInput.Insert(elem)
			}
			assert.Equal(t, test.output, test.treeInput.Points())
		})
	}
}

func TestKDTree_InsertWithGenerator(t *testing.T) {
	tests := []struct {
		name  string
		input []Point
	}{
		{name: "p:10,k:5", input: generateTestCaseData(10)},
		{name: "p:100,k:5", input: generateTestCaseData(100)},
		{name: "p:1000,k:5", input: generateTestCaseData(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := New(nil)
			for _, elem := range test.input {
				tree.Insert(elem)
				_ = tree.Points()
				// TODO assert
				//assert.Equal(t, , treePoints)
			}
		})
	}
}

func TestKDTree_Remove(t *testing.T) {
	tests := []struct {
		name       string
		treeInput  *KDTree
		preRemove  []Point
		input      Point
		treeOutput string
		output     Point
	}{
		{
			name:       "empty tree",
			treeInput:  New([]Point{}),
			input:      &SamplePoint2D{},
			treeOutput: "[<nil>]",
			output:     nil,
		},
		{
			name:       "nil input",
			treeInput:  New([]Point{&SamplePoint2D{X: 1., Y: 2.}}),
			input:      nil,
			treeOutput: "[{1.00 2.00}]",
			output:     nil,
		},
		{
			name:       "remove root",
			treeInput:  New([]Point{&SamplePoint2D{X: 1., Y: 2.}}),
			input:      &SamplePoint2D{X: 1., Y: 2.},
			treeOutput: "[<nil>]",
			output:     &SamplePoint2D{X: 1., Y: 2.},
		},
		{
			name: "remove root with children",
			treeInput: New([]Point{&SamplePoint2D{X: 1., Y: 2.},
				&SamplePoint2D{X: 1.2, Y: 2.2},
				&SamplePoint2D{X: 1.3, Y: 2.3},
				&SamplePoint2D{X: 1.1, Y: 2.1},
				&SamplePoint2D{X: -1.3, Y: -2.2}}),
			input:      &SamplePoint2D{X: 1.1, Y: 2.1},
			treeOutput: "[[{-1.30 -2.20} {1.00 2.00} [{1.20 2.20} {1.30 2.30} <nil>]]]",
			output:     &SamplePoint2D{X: 1.1, Y: 2.1},
		},
		// x(5,4)
		// y(3,6)                       (7, 7)
		// x(2, 2)          (2, 10)     (8, 2)          (9, 9)
		// y(1, 3) (4, 1)   (1,8) nil   (7, 4) (8, 5)   (6, 8) nil
		// [[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]
		{
			name: "not existing",
			treeInput: New(
				[]Point{
					&SamplePoint2D{X: 1, Y: 3},
					&SamplePoint2D{X: 1, Y: 8},
					&SamplePoint2D{X: 2, Y: 2},
					&SamplePoint2D{X: 2, Y: 10},
					&SamplePoint2D{X: 3, Y: 6},
					&SamplePoint2D{X: 4, Y: 1},
					&SamplePoint2D{X: 5, Y: 4},
					&SamplePoint2D{X: 6, Y: 8},
					&SamplePoint2D{X: 7, Y: 4},
					&SamplePoint2D{X: 7, Y: 7},
					&SamplePoint2D{X: 8, Y: 2},
					&SamplePoint2D{X: 8, Y: 5},
					&SamplePoint2D{X: 9, Y: 9}}),
			input:      &SamplePoint2D{X: 1., Y: 1.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     nil,
		},
		{
			name: "remove leaf",
			treeInput: New(
				[]Point{
					&SamplePoint2D{X: 1, Y: 3},
					&SamplePoint2D{X: 1, Y: 8},
					&SamplePoint2D{X: 2, Y: 2},
					&SamplePoint2D{X: 2, Y: 10},
					&SamplePoint2D{X: 3, Y: 6},
					&SamplePoint2D{X: 4, Y: 1},
					&SamplePoint2D{X: 5, Y: 4},
					&SamplePoint2D{X: 6, Y: 8},
					&SamplePoint2D{X: 7, Y: 4},
					&SamplePoint2D{X: 7, Y: 7},
					&SamplePoint2D{X: 8, Y: 2},
					&SamplePoint2D{X: 8, Y: 5},
					&SamplePoint2D{X: 9, Y: 9}}),
			input:      &SamplePoint2D{X: 8., Y: 5.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} <nil>] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &SamplePoint2D{X: 8., Y: 5.},
		},
		{
			name: "remove leaf",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			input:      &SamplePoint2D{X: 6., Y: 8.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} {9.00 9.00}]]]",
			output:     &SamplePoint2D{X: 6., Y: 8.},
		},
		{
			name: "remove with 1 replace, right child nil",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			input:      &SamplePoint2D{X: 9., Y: 9.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} {6.00 8.00}]]]",
			output:     &SamplePoint2D{X: 9., Y: 9.},
		},
		{
			name: "remove with 1 replace, left child nil",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			preRemove:  []Point{&SamplePoint2D{X: 1, Y: 3}},
			input:      &SamplePoint2D{X: 2., Y: 2.},
			treeOutput: "[[[{4.00 1.00} {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &SamplePoint2D{X: 2., Y: 2.},
		},
		{
			name: "remove with 1 replace",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			input:      &SamplePoint2D{X: 8., Y: 2.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[<nil> {7.00 4.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &SamplePoint2D{X: 8., Y: 2.},
		},
		{
			name: "remove with 1 replace, deep",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			input:      &SamplePoint2D{X: 3., Y: 6.},
			treeOutput: "[[[[<nil> {2.00 2.00} {4.00 1.00}] {1.00 3.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &SamplePoint2D{X: 3., Y: 6.},
		},
		{
			name: "remove with 1 replace, deep",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			input:      &SamplePoint2D{X: 7., Y: 7.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} <nil>] {8.00 5.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &SamplePoint2D{X: 7., Y: 7.},
		},
		{
			name: "remove with left nil",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			preRemove: []Point{&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6}},
			input:      &SamplePoint2D{X: 5., Y: 4.},
			treeOutput: "[[<nil> {6.00 8.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} {9.00 9.00}]]]",
			output:     &SamplePoint2D{X: 5., Y: 4.},
		},
		{
			name: "remove with sub left nil",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}}),
			preRemove: []Point{&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 2, Y: 2}},
			input:      &SamplePoint2D{X: 5., Y: 4.},
			treeOutput: "[[[<nil> {1.00 8.00} {2.00 10.00}] {3.00 6.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &SamplePoint2D{X: 5., Y: 4.},
		},
		// x (4,1)
		// y (1,3)                                                       (5,4)
		// x (3,1)                       (3,6)                           (7,3)                       (8,5)
		// y (2,2)         (4,2)         (2,8)           (3,8)           (5,3)         (9,2)         (7,7)         (9,8)
		// x (2,1) (1,3)   (3,1) (3,3)   (1,8) (2,10)    (4,4) (3,9)     (6,2) (6,4)   (8,2) (7,4)   (6,5) (6,8)   (9,6) (9,9)
		// [[[[[{2.00 1.00} {2.00 2.00} {1.00 3.00}] {3.00 1.00} [{3.00 1.00} {4.00 2.00} {3.00 3.00}]] {1.00 3.00} [[{1.00 8.00} {2.00 8.00} {2.00 10.00}] {3.00 6.00} [{4.00 4.00} {3.00 8.00} {3.00 9.00}]]] {4.00 1.00} [[[{6.00 2.00} {5.00 3.00} {6.00 4.00}] {7.00 3.00} [{8.00 2.00} {9.00 2.00} {7.00 4.00}]] {5.00 4.00} [[{6.00 5.00} {7.00 7.00} {6.00 8.00}] {8.00 5.00} [{9.00 6.00} {9.00 8.00} {9.00 9.00}]]]]]
		{
			name: "remove (3,1) with 2 replace",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 4, Y: 2},
				&SamplePoint2D{X: 9, Y: 2},
				&SamplePoint2D{X: 6, Y: 5},
				&SamplePoint2D{X: 3, Y: 8},
				&SamplePoint2D{X: 6, Y: 2},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 3, Y: 3},
				&SamplePoint2D{X: 6, Y: 4},
				&SamplePoint2D{X: 9, Y: 8},
				&SamplePoint2D{X: 2, Y: 1},
				&SamplePoint2D{X: 2, Y: 8},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 7, Y: 3},
				&SamplePoint2D{X: 3, Y: 9},
				&SamplePoint2D{X: 4, Y: 4},
				&SamplePoint2D{X: 5, Y: 3},
				&SamplePoint2D{X: 9, Y: 6}}),
			input:      &SamplePoint2D{X: 3., Y: 1.},
			treeOutput: "[[[[[<nil> {2.00 1.00} {1.00 3.00}] {2.00 2.00} [{3.00 1.00} {4.00 2.00} {3.00 3.00}]] {1.00 3.00} [[{1.00 8.00} {2.00 8.00} {2.00 10.00}] {3.00 6.00} [{4.00 4.00} {3.00 8.00} {3.00 9.00}]]] {4.00 1.00} [[[{6.00 2.00} {5.00 3.00} {6.00 4.00}] {7.00 3.00} [{8.00 2.00} {9.00 2.00} {7.00 4.00}]] {5.00 4.00} [[{6.00 5.00} {7.00 7.00} {6.00 8.00}] {8.00 5.00} [{9.00 6.00} {9.00 8.00} {9.00 9.00}]]]]]",
			output:     &SamplePoint2D{X: 3., Y: 1.},
		},
		{
			name: "remove (5,4) with 1 replace, deep 3",
			treeInput: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 4, Y: 2},
				&SamplePoint2D{X: 9, Y: 2},
				&SamplePoint2D{X: 6, Y: 5},
				&SamplePoint2D{X: 3, Y: 8},
				&SamplePoint2D{X: 6, Y: 2},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 3, Y: 3},
				&SamplePoint2D{X: 6, Y: 4},
				&SamplePoint2D{X: 9, Y: 8},
				&SamplePoint2D{X: 2, Y: 1},
				&SamplePoint2D{X: 2, Y: 8},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 7, Y: 3},
				&SamplePoint2D{X: 3, Y: 9},
				&SamplePoint2D{X: 4, Y: 4},
				&SamplePoint2D{X: 5, Y: 3},
				&SamplePoint2D{X: 9, Y: 6}}),
			input:      &SamplePoint2D{X: 5., Y: 4.},
			treeOutput: "[[[[[{2.00 1.00} {2.00 2.00} {1.00 3.00}] {3.00 1.00} [{3.00 1.00} {4.00 2.00} {3.00 3.00}]] {1.00 3.00} [[{1.00 8.00} {2.00 8.00} {2.00 10.00}] {3.00 6.00} [{4.00 4.00} {3.00 8.00} {3.00 9.00}]]] {4.00 1.00} [[[{6.00 2.00} {5.00 3.00} <nil>] {7.00 3.00} [{8.00 2.00} {9.00 2.00} {7.00 4.00}]] {6.00 4.00} [[{6.00 5.00} {7.00 7.00} {6.00 8.00}] {8.00 5.00} [{9.00 6.00} {9.00 8.00} {9.00 9.00}]]]]]",
			output:     &SamplePoint2D{X: 5., Y: 4.},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preRemove != nil {
				for _, r := range test.preRemove {
					test.treeInput.Remove(r)
				}
			}
			o := test.treeInput.Remove(test.input)
			if c, ok := o.(Point); ok {
				assertPointsEqual(t, test.output, c)
			} else {
				assert.Equal(t, test.output, o)
			}
			assert.Equal(t, test.treeOutput, test.treeInput.String())
		})
	}
}

func TestKDTree_Balance(t *testing.T) {
	tests := []struct {
		name       string
		treeInput  *KDTree
		preRemove  []Point
		treeOutput string
	}{
		{
			name:       "empty tree",
			treeInput:  New([]Point{}),
			treeOutput: "[<nil>]",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preRemove != nil {
				for _, r := range test.preRemove {
					test.treeInput.Remove(r)
				}
			}
			test.treeInput.Balance()
			assert.Equal(t, test.treeOutput, test.treeInput.String())
		})
	}
}

func TestKDTree_BalanceNoNilNode(t *testing.T) {
	tests := []struct {
		name   string
		input  []Point
		add    int
		remove int
	}{
		// add
		{name: "0->1", input: generateTestCaseData(0), add: 1},
		{name: "0->3", input: generateTestCaseData(0), add: 3},
		{name: "0->7", input: generateTestCaseData(0), add: 7},
		{name: "0->15", input: generateTestCaseData(0), add: 15},
		// remove
		{name: "5->1", input: generateTestCaseData(5), remove: 4},
		{name: "32->3", input: generateTestCaseData(32), remove: 29},
		{name: "17->7", input: generateTestCaseData(17), remove: 10},
		{name: "17->15", input: generateTestCaseData(17), remove: 2},
		// remove & add
		{name: "50->8->15", input: generateTestCaseData(50), remove: 42, add: 7},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := New(test.input)
			for i := 0; i < test.remove; i++ {
				tree.Remove(test.input[i])
			}
			for _, p := range generateTestCaseData(test.add) {
				tree.Insert(p)
			}
			tree.Balance()
			assert.NotContains(t, tree.String(), "<nil>")
		})
	}
}

// TestKNN ...
func TestKDTree_KNN(t *testing.T) {
	tests := []struct {
		name   string
		target Point
		k      int
		input  []Point
		output []Point
	}{
		{
			name:   "nil",
			target: nil,
			k:      3,
			input:  []Point{&SamplePoint2D{X: 1., Y: 2.}},
			output: []Point{},
		},
		{
			name:   "empty",
			target: &SamplePoint2D{X: 1., Y: 2.},
			k:      3,
			input:  []Point{},
			output: []Point{},
		},
		{
			name:   "k >> points",
			target: &SamplePoint2D{X: 1., Y: 2.},
			k:      10,
			input: []Point{&SamplePoint2D{X: 1., Y: 2.},
				&SamplePoint2D{X: 0.9, Y: 2.1},
				&SamplePoint2D{X: 1.1, Y: 1.9}},
			output: []Point{&SamplePoint2D{X: 1., Y: 2.},
				&SamplePoint2D{X: 0.9, Y: 2.1},
				&SamplePoint2D{X: 1.1, Y: 1.9}},
		},
		{
			name:   "small 2D example",
			target: &SamplePoint2D{X: 9, Y: 4},
			k:      3,
			input: []Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9}},
			output: []Point{&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 8, Y: 2}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := New(test.input)
			assert.Equal(t, test.output, tree.KNN(test.target, test.k))
		})
	}
}

func TestKDTree_KNNWithGenerator(t *testing.T) {
	tests := []struct {
		name   string
		target Point
		k      int
		input  []Point
	}{
		{name: "p:100,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(100)},
		{name: "p:1000,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(1000)},
		{name: "p:10000,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(10000)},
		{name: "p:100000,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(100000)},
		{name: "p:1000000,k:10", target: &SamplePoint2D{}, k: 10, input: generateTestCaseData(1000000)},
		{name: "p:1000000,k:20", target: &SamplePoint2D{}, k: 20, input: generateTestCaseData(1000000)},
		{name: "p:1000000,k:30", target: &SamplePoint2D{}, k: 30, input: generateTestCaseData(1000000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := New(test.input)
			assert.Equal(t, prioQueueKNN(test.input, test.target, test.k), tree.KNN(test.target, test.k))
		})
	}
}

func TestKDTree_RangeSearch(t *testing.T) {
	tests := []struct {
		name     string
		tree     *KDTree
		input    kdrangeRange
		expected []Point
	}{
		{
			name:     "nil",
			tree:     New(generateTestCaseData(5)),
			input:    nil,
			expected: []Point{},
		},
		{
			name:     "wrong dim",
			tree:     New(generateTestCaseData(5)),
			input:    NewKDRange(),
			expected: []Point{},
		},
		{
			name: "out of range x (lower)",
			tree: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 4, Y: 2},
				&SamplePoint2D{X: 9, Y: 2},
				&SamplePoint2D{X: 6, Y: 5},
				&SamplePoint2D{X: 3, Y: 8},
				&SamplePoint2D{X: 6, Y: 2},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 3, Y: 3},
				&SamplePoint2D{X: 6, Y: 4},
				&SamplePoint2D{X: 9, Y: 8},
				&SamplePoint2D{X: 2, Y: 1},
				&SamplePoint2D{X: 2, Y: 8},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 7, Y: 3},
				&SamplePoint2D{X: 3, Y: 9},
				&SamplePoint2D{X: 4, Y: 4},
				&SamplePoint2D{X: 5, Y: 3},
				&SamplePoint2D{X: 9, Y: 6}}),
			input:    NewKDRange(-2, -1, 2, 10),
			expected: []Point{},
		},
		{
			name: "out of range y (lower)",
			tree: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 4, Y: 2},
				&SamplePoint2D{X: 9, Y: 2},
				&SamplePoint2D{X: 6, Y: 5},
				&SamplePoint2D{X: 3, Y: 8},
				&SamplePoint2D{X: 6, Y: 2},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 3, Y: 3},
				&SamplePoint2D{X: 6, Y: 4},
				&SamplePoint2D{X: 9, Y: 8},
				&SamplePoint2D{X: 2, Y: 1},
				&SamplePoint2D{X: 2, Y: 8},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 7, Y: 3},
				&SamplePoint2D{X: 3, Y: 9},
				&SamplePoint2D{X: 4, Y: 4},
				&SamplePoint2D{X: 5, Y: 3},
				&SamplePoint2D{X: 9, Y: 6}}),
			input:    NewKDRange(2, 10, -2, -1),
			expected: []Point{},
		},
		{
			name: "out of range x (higher)",
			tree: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 4, Y: 2},
				&SamplePoint2D{X: 9, Y: 2},
				&SamplePoint2D{X: 6, Y: 5},
				&SamplePoint2D{X: 3, Y: 8},
				&SamplePoint2D{X: 6, Y: 2},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 3, Y: 3},
				&SamplePoint2D{X: 6, Y: 4},
				&SamplePoint2D{X: 9, Y: 8},
				&SamplePoint2D{X: 2, Y: 1},
				&SamplePoint2D{X: 2, Y: 8},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 7, Y: 3},
				&SamplePoint2D{X: 3, Y: 9},
				&SamplePoint2D{X: 4, Y: 4},
				&SamplePoint2D{X: 5, Y: 3},
				&SamplePoint2D{X: 9, Y: 6}}),
			input:    NewKDRange(20, 30, 2, 10),
			expected: []Point{},
		},
		{
			name: "out of range y (higher)",
			tree: New([]Point{&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 1, Y: 8},
				&SamplePoint2D{X: 2, Y: 2},
				&SamplePoint2D{X: 2, Y: 10},
				&SamplePoint2D{X: 3, Y: 6},
				&SamplePoint2D{X: 4, Y: 1},
				&SamplePoint2D{X: 5, Y: 4},
				&SamplePoint2D{X: 6, Y: 8},
				&SamplePoint2D{X: 7, Y: 4},
				&SamplePoint2D{X: 7, Y: 7},
				&SamplePoint2D{X: 8, Y: 2},
				&SamplePoint2D{X: 8, Y: 5},
				&SamplePoint2D{X: 9, Y: 9},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 4, Y: 2},
				&SamplePoint2D{X: 9, Y: 2},
				&SamplePoint2D{X: 6, Y: 5},
				&SamplePoint2D{X: 3, Y: 8},
				&SamplePoint2D{X: 6, Y: 2},
				&SamplePoint2D{X: 1, Y: 3},
				&SamplePoint2D{X: 3, Y: 3},
				&SamplePoint2D{X: 6, Y: 4},
				&SamplePoint2D{X: 9, Y: 8},
				&SamplePoint2D{X: 2, Y: 1},
				&SamplePoint2D{X: 2, Y: 8},
				&SamplePoint2D{X: 3, Y: 1},
				&SamplePoint2D{X: 7, Y: 3},
				&SamplePoint2D{X: 3, Y: 9},
				&SamplePoint2D{X: 4, Y: 4},
				&SamplePoint2D{X: 5, Y: 3},
				&SamplePoint2D{X: 9, Y: 6}}),
			input:    NewKDRange(2, 10, 20, 30),
			expected: []Point{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.tree.RangeSearch(test.input))
		})
	}
}

func TestKDTree_RangeSearchWithGenerator(t *testing.T) {
	tests := []struct {
		name  string
		input []Point
		r     kdrangeRange
	}{
		{name: "nodes: 100 range: -100 50 -50 100", input: generateTestCaseData(100), r: NewKDRange(-100, 50, -50, 100)},
		{name: "nodes: 1000 range: -100 50 -50 100", input: generateTestCaseData(1000), r: NewKDRange(-100, 50, -50, 100)},
		{name: "nodes: 10000 range: -100 50 -50 100", input: generateTestCaseData(10000), r: NewKDRange(-100, 50, -50, 100)},
		{name: "nodes: 100000 range: -500 250 -250 500", input: generateTestCaseData(100000), r: NewKDRange(-500, 250, -250, 500)},
		{name: "nodes: 1000000 range: -500 250 -250 500", input: generateTestCaseData(1000000), r: NewKDRange(-500, 250, -250, 500)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := New(test.input)
			assert.ElementsMatch(t, filterRangeSearch(test.input, test.r), tree.RangeSearch(test.r))
		})
	}
}

// benchmarks

var resultTree *KDTree
var resultPoints []Point

func BenchmarkNew(b *testing.B) {
	benchmarks := []struct {
		name  string
		input []Point
	}{
		{name: "100", input: generateTestCaseData(100)},
		{name: "1000", input: generateTestCaseData(1000)},
		{name: "10000", input: generateTestCaseData(10000)},
		{name: "100000", input: generateTestCaseData(100000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var t *KDTree
			for i := 0; i < b.N; i++ {
				t = New(bm.input)
			}
			resultTree = t
		})
	}
}

func BenchmarkKNN(b *testing.B) {
	benchmarks := []struct {
		name   string
		target Point
		k      int
		input  []Point
	}{
		{name: "p:100,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(100)},
		{name: "p:1000,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(1000)},
		{name: "p:10000,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(10000)},
		{name: "p:100000,k:5", target: &SamplePoint2D{}, k: 5, input: generateTestCaseData(100000)},
	}
	for _, bm := range benchmarks {
		var res []Point
		tree := New(bm.input)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				res = tree.KNN(bm.target, bm.k)
			}
			resultPoints = res
		})
	}
}

// helpers

func generateTestCaseData(size int) []Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var pointsSlice []Point
	for i := 0; i < size; i++ {
		pointsSlice = append(pointsSlice,
			&SamplePoint2D{X: r.Float64()*3000 - 1500, Y: r.Float64()*3000 - 1500})
	}

	return pointsSlice
}

func prioQueueKNN(points []Point, p Point, k int) []Point {
	knn := make([]Point, 0, k)
	if p == nil {
		return knn
	}

	nnPQ := pq.New()
	for _, point := range points {
		nnPQ.Insert(point, p.Distance(point))
	}

	for i := 0; i < k; i++ {
		point, err := nnPQ.Pop()
		if err != nil {
			break
		}
		knn = append(knn, point.(Point))
	}
	return knn
}

func filterRangeSearch(points []Point, r kdrangeRange) []Point {
	result := make([]Point, 0)

pointLoop:
	for _, point := range points {
		for i, d := range r {
			if d[0] > point.Dimension(i) || d[1] < point.Dimension(i) {
				continue pointLoop
			}
		}
		result = append(result, point)
	}

	return result
}

func assertPointsEqual(t *testing.T, p1 Point, p2 Point) {
	assert.Equal(t, p1.Dimensions(), p2.Dimensions())
	for i := 0; i < p1.Dimensions(); i++ {
		assert.Equal(t, p1.Dimension(i), p2.Dimension(i))
	}
}

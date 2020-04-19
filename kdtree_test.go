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

package kdtree_test

import (
	"github.com/jupp0r/go-priority-queue"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/kdrange"
	. "github.com/kyroy/kdtree/points"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
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
			tree := kdtree.New(test.input)
			assert.Equal(t, test.output, tree.Points())
		})

	}
}

func TestKDTree_String(t *testing.T) {
	tests := []struct {
		name     string
		tree     *kdtree.KDTree
		expected string
	}{
		{name: "empty", tree: &kdtree.KDTree{}, expected: "[<nil>]"},
		{name: "1 elem", tree: kdtree.New([]kdtree.Point{&Point2D{X: 2, Y: 3}}), expected: "[{2.00 3.00}]"},
		{name: "2 elem", tree: kdtree.New([]kdtree.Point{&Point2D{X: 2, Y: 3}, &Point2D{X: 3.4, Y: 1}}), expected: "[[{2.00 3.00} {3.40 1.00} <nil>]]"},
		{name: "3 elem", tree: kdtree.New([]kdtree.Point{&Point2D{X: 2, Y: 3}, &Point2D{X: 1.4, Y: 7.1}, &Point2D{X: 3.4, Y: 1}}), expected: "[[{1.40 7.10} {2.00 3.00} {3.40 1.00}]]"},
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
			treeInput: kdtree.New([]kdtree.Point{&Point2D{X: 1., Y: 2.}}),
			input:     []kdtree.Point{&Point2D{X: 0.9, Y: 2.1}},
			output:    []kdtree.Point{&Point2D{X: 0.9, Y: 2.1}, &Point2D{X: 1., Y: 2.}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.treeInput == nil {
				test.treeInput = kdtree.New(nil)
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
		input []kdtree.Point
	}{
		{name: "p:10,k:5", input: generateTestCaseData(10)},
		{name: "p:100,k:5", input: generateTestCaseData(100)},
		{name: "p:1000,k:5", input: generateTestCaseData(1000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := kdtree.New(nil)
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
		treeInput  *kdtree.KDTree
		preRemove  []kdtree.Point
		input      kdtree.Point
		treeOutput string
		output     kdtree.Point
	}{
		{
			name:       "empty tree",
			treeInput:  kdtree.New([]kdtree.Point{}),
			input:      &Point2D{},
			treeOutput: "[<nil>]",
			output:     nil,
		},
		{
			name:       "nil input",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1., Y: 2.}}),
			input:      nil,
			treeOutput: "[{1.00 2.00}]",
			output:     nil,
		},
		{
			name:       "remove root",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1., Y: 2.}}),
			input:      &Point2D{X: 1., Y: 2.},
			treeOutput: "[<nil>]",
			output:     &Point2D{X: 1., Y: 2.},
		},
		{
			name:       "remove root with children",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1., Y: 2.}, &Point2D{X: 1.2, Y: 2.2}, &Point2D{X: 1.3, Y: 2.3}, &Point2D{X: 1.1, Y: 2.1}, &Point2D{X: -1.3, Y: -2.2}}),
			input:      &Point2D{X: 1.1, Y: 2.1},
			treeOutput: "[[{-1.30 -2.20} {1.00 2.00} [{1.20 2.20} {1.30 2.30} <nil>]]]",
			output:     &Point2D{X: 1.1, Y: 2.1},
		},
		// x(5,4)
		// y(3,6)                       (7, 7)
		// x(2, 2)          (2, 10)     (8, 2)          (9, 9)
		// y(1, 3) (4, 1)   (1,8) nil   (7, 4) (8, 5)   (6, 8) nil
		// [[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]
		{
			name:       "not existing",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			input:      &Point2D{X: 1., Y: 1.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     nil,
		},
		{
			name:       "remove leaf",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			input:      &Point2D{X: 8., Y: 5.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} <nil>] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &Point2D{X: 8., Y: 5.},
		},
		{
			name:       "remove leaf",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			input:      &Point2D{X: 6., Y: 8.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} {9.00 9.00}]]]",
			output:     &Point2D{X: 6., Y: 8.},
		},
		{
			name:       "remove with 1 replace, right child nil",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			input:      &Point2D{X: 9., Y: 9.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} {6.00 8.00}]]]",
			output:     &Point2D{X: 9., Y: 9.},
		},
		{
			name:       "remove with 1 replace, left child nil",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			preRemove:  []kdtree.Point{&Point2D{X: 1, Y: 3}},
			input:      &Point2D{X: 2., Y: 2.},
			treeOutput: "[[[{4.00 1.00} {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &Point2D{X: 2., Y: 2.},
		},
		{
			name:       "remove with 1 replace",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			input:      &Point2D{X: 8., Y: 2.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[<nil> {7.00 4.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &Point2D{X: 8., Y: 2.},
		},
		{
			name:       "remove with 1 replace, deep",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			input:      &Point2D{X: 3., Y: 6.},
			treeOutput: "[[[[<nil> {2.00 2.00} {4.00 1.00}] {1.00 3.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &Point2D{X: 3., Y: 6.},
		},
		{
			name:       "remove with 1 replace, deep",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			input:      &Point2D{X: 7., Y: 7.},
			treeOutput: "[[[[{1.00 3.00} {2.00 2.00} {4.00 1.00}] {3.00 6.00} [{1.00 8.00} {2.00 10.00} <nil>]] {5.00 4.00} [[{7.00 4.00} {8.00 2.00} <nil>] {8.00 5.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &Point2D{X: 7., Y: 7.},
		},
		{
			name:       "remove with left nil",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			preRemove:  []kdtree.Point{&Point2D{X: 4, Y: 1}, &Point2D{X: 1, Y: 3}, &Point2D{X: 2, Y: 2}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}},
			input:      &Point2D{X: 5., Y: 4.},
			treeOutput: "[[<nil> {6.00 8.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} {9.00 9.00}]]]",
			output:     &Point2D{X: 5., Y: 4.},
		},
		{
			name:       "remove with sub left nil",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}}),
			preRemove:  []kdtree.Point{&Point2D{X: 4, Y: 1}, &Point2D{X: 1, Y: 3}, &Point2D{X: 2, Y: 2}},
			input:      &Point2D{X: 5., Y: 4.},
			treeOutput: "[[[<nil> {1.00 8.00} {2.00 10.00}] {3.00 6.00} [[{7.00 4.00} {8.00 2.00} {8.00 5.00}] {7.00 7.00} [{6.00 8.00} {9.00 9.00} <nil>]]]]",
			output:     &Point2D{X: 5., Y: 4.},
		},
		// x (4,1)
		// y (1,3)                                                       (5,4)
		// x (3,1)                       (3,6)                           (7,3)                       (8,5)
		// y (2,2)         (4,2)         (2,8)           (3,8)           (5,3)         (9,2)         (7,7)         (9,8)
		// x (2,1) (1,3)   (3,1) (3,3)   (1,8) (2,10)    (4,4) (3,9)     (6,2) (6,4)   (8,2) (7,4)   (6,5) (6,8)   (9,6) (9,9)
		// [[[[[{2.00 1.00} {2.00 2.00} {1.00 3.00}] {3.00 1.00} [{3.00 1.00} {4.00 2.00} {3.00 3.00}]] {1.00 3.00} [[{1.00 8.00} {2.00 8.00} {2.00 10.00}] {3.00 6.00} [{4.00 4.00} {3.00 8.00} {3.00 9.00}]]] {4.00 1.00} [[[{6.00 2.00} {5.00 3.00} {6.00 4.00}] {7.00 3.00} [{8.00 2.00} {9.00 2.00} {7.00 4.00}]] {5.00 4.00} [[{6.00 5.00} {7.00 7.00} {6.00 8.00}] {8.00 5.00} [{9.00 6.00} {9.00 8.00} {9.00 9.00}]]]]]
		{
			name:       "remove (3,1) with 2 replace",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}, &Point2D{X: 3, Y: 1}, &Point2D{X: 4, Y: 2}, &Point2D{X: 9, Y: 2}, &Point2D{X: 6, Y: 5}, &Point2D{X: 3, Y: 8}, &Point2D{X: 6, Y: 2}, &Point2D{X: 1, Y: 3}, &Point2D{X: 3, Y: 3}, &Point2D{X: 6, Y: 4}, &Point2D{X: 9, Y: 8}, &Point2D{X: 2, Y: 1}, &Point2D{X: 2, Y: 8}, &Point2D{X: 3, Y: 1}, &Point2D{X: 7, Y: 3}, &Point2D{X: 3, Y: 9}, &Point2D{X: 4, Y: 4}, &Point2D{X: 5, Y: 3}, &Point2D{X: 9, Y: 6}}),
			input:      &Point2D{X: 3., Y: 1.},
			treeOutput: "[[[[[<nil> {2.00 1.00} {1.00 3.00}] {2.00 2.00} [{3.00 1.00} {4.00 2.00} {3.00 3.00}]] {1.00 3.00} [[{1.00 8.00} {2.00 8.00} {2.00 10.00}] {3.00 6.00} [{4.00 4.00} {3.00 8.00} {3.00 9.00}]]] {4.00 1.00} [[[{6.00 2.00} {5.00 3.00} {6.00 4.00}] {7.00 3.00} [{8.00 2.00} {9.00 2.00} {7.00 4.00}]] {5.00 4.00} [[{6.00 5.00} {7.00 7.00} {6.00 8.00}] {8.00 5.00} [{9.00 6.00} {9.00 8.00} {9.00 9.00}]]]]]",
			output:     &Point2D{X: 3., Y: 1.},
		},
		{
			name:       "remove (5,4) with 1 replace, deep 3",
			treeInput:  kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}, &Point2D{X: 3, Y: 1}, &Point2D{X: 4, Y: 2}, &Point2D{X: 9, Y: 2}, &Point2D{X: 6, Y: 5}, &Point2D{X: 3, Y: 8}, &Point2D{X: 6, Y: 2}, &Point2D{X: 1, Y: 3}, &Point2D{X: 3, Y: 3}, &Point2D{X: 6, Y: 4}, &Point2D{X: 9, Y: 8}, &Point2D{X: 2, Y: 1}, &Point2D{X: 2, Y: 8}, &Point2D{X: 3, Y: 1}, &Point2D{X: 7, Y: 3}, &Point2D{X: 3, Y: 9}, &Point2D{X: 4, Y: 4}, &Point2D{X: 5, Y: 3}, &Point2D{X: 9, Y: 6}}),
			input:      &Point2D{X: 5., Y: 4.},
			treeOutput: "[[[[[{2.00 1.00} {2.00 2.00} {1.00 3.00}] {3.00 1.00} [{3.00 1.00} {4.00 2.00} {3.00 3.00}]] {1.00 3.00} [[{1.00 8.00} {2.00 8.00} {2.00 10.00}] {3.00 6.00} [{4.00 4.00} {3.00 8.00} {3.00 9.00}]]] {4.00 1.00} [[[{6.00 2.00} {5.00 3.00} <nil>] {7.00 3.00} [{8.00 2.00} {9.00 2.00} {7.00 4.00}]] {6.00 4.00} [[{6.00 5.00} {7.00 7.00} {6.00 8.00}] {8.00 5.00} [{9.00 6.00} {9.00 8.00} {9.00 9.00}]]]]]",
			output:     &Point2D{X: 5., Y: 4.},
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
			if c, ok := o.(kdtree.Point); ok {
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
		treeInput  *kdtree.KDTree
		preRemove  []kdtree.Point
		treeOutput string
	}{
		{
			name:       "empty tree",
			treeInput:  kdtree.New([]kdtree.Point{}),
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
		input  []kdtree.Point
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
			tree := kdtree.New(test.input)
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
			tree := kdtree.New(test.input)
			assert.Equal(t, test.output, tree.KNN(test.target, test.k))
		})
	}
}

func TestKDTree_KNNWithGenerator(t *testing.T) {
	tests := []struct {
		name   string
		target kdtree.Point
		k      int
		input  []kdtree.Point
	}{
		{name: "p:100,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(100)},
		{name: "p:1000,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(1000)},
		{name: "p:10000,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(10000)},
		{name: "p:100000,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(100000)},
		{name: "p:1000000,k:10", target: &Point2D{}, k: 10, input: generateTestCaseData(1000000)},
		{name: "p:1000000,k:20", target: &Point2D{}, k: 20, input: generateTestCaseData(1000000)},
		{name: "p:1000000,k:30", target: &Point2D{}, k: 30, input: generateTestCaseData(1000000)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := kdtree.New(test.input)
			assert.Equal(t, prioQueueKNN(test.input, test.target, test.k), tree.KNN(test.target, test.k))
		})
	}
}

func TestKDTree_RangeSearch(t *testing.T) {
	tests := []struct {
		name     string
		tree     *kdtree.KDTree
		input    kdrange.Range
		expected []kdtree.Point
	}{
		{
			name:     "nil",
			tree:     kdtree.New(generateTestCaseData(5)),
			input:    nil,
			expected: []kdtree.Point{},
		},
		{
			name:     "wrong dim",
			tree:     kdtree.New(generateTestCaseData(5)),
			input:    kdrange.New(),
			expected: []kdtree.Point{},
		},
		{
			name:     "out of range x (lower)",
			tree:     kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}, &Point2D{X: 3, Y: 1}, &Point2D{X: 4, Y: 2}, &Point2D{X: 9, Y: 2}, &Point2D{X: 6, Y: 5}, &Point2D{X: 3, Y: 8}, &Point2D{X: 6, Y: 2}, &Point2D{X: 1, Y: 3}, &Point2D{X: 3, Y: 3}, &Point2D{X: 6, Y: 4}, &Point2D{X: 9, Y: 8}, &Point2D{X: 2, Y: 1}, &Point2D{X: 2, Y: 8}, &Point2D{X: 3, Y: 1}, &Point2D{X: 7, Y: 3}, &Point2D{X: 3, Y: 9}, &Point2D{X: 4, Y: 4}, &Point2D{X: 5, Y: 3}, &Point2D{X: 9, Y: 6}}),
			input:    kdrange.New(-2, -1, 2, 10),
			expected: []kdtree.Point{},
		},
		{
			name:     "out of range y (lower)",
			tree:     kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}, &Point2D{X: 3, Y: 1}, &Point2D{X: 4, Y: 2}, &Point2D{X: 9, Y: 2}, &Point2D{X: 6, Y: 5}, &Point2D{X: 3, Y: 8}, &Point2D{X: 6, Y: 2}, &Point2D{X: 1, Y: 3}, &Point2D{X: 3, Y: 3}, &Point2D{X: 6, Y: 4}, &Point2D{X: 9, Y: 8}, &Point2D{X: 2, Y: 1}, &Point2D{X: 2, Y: 8}, &Point2D{X: 3, Y: 1}, &Point2D{X: 7, Y: 3}, &Point2D{X: 3, Y: 9}, &Point2D{X: 4, Y: 4}, &Point2D{X: 5, Y: 3}, &Point2D{X: 9, Y: 6}}),
			input:    kdrange.New(2, 10, -2, -1),
			expected: []kdtree.Point{},
		},
		{
			name:     "out of range x (higher)",
			tree:     kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}, &Point2D{X: 3, Y: 1}, &Point2D{X: 4, Y: 2}, &Point2D{X: 9, Y: 2}, &Point2D{X: 6, Y: 5}, &Point2D{X: 3, Y: 8}, &Point2D{X: 6, Y: 2}, &Point2D{X: 1, Y: 3}, &Point2D{X: 3, Y: 3}, &Point2D{X: 6, Y: 4}, &Point2D{X: 9, Y: 8}, &Point2D{X: 2, Y: 1}, &Point2D{X: 2, Y: 8}, &Point2D{X: 3, Y: 1}, &Point2D{X: 7, Y: 3}, &Point2D{X: 3, Y: 9}, &Point2D{X: 4, Y: 4}, &Point2D{X: 5, Y: 3}, &Point2D{X: 9, Y: 6}}),
			input:    kdrange.New(20, 30, 2, 10),
			expected: []kdtree.Point{},
		},
		{
			name:     "out of range y (higher)",
			tree:     kdtree.New([]kdtree.Point{&Point2D{X: 1, Y: 3}, &Point2D{X: 1, Y: 8}, &Point2D{X: 2, Y: 2}, &Point2D{X: 2, Y: 10}, &Point2D{X: 3, Y: 6}, &Point2D{X: 4, Y: 1}, &Point2D{X: 5, Y: 4}, &Point2D{X: 6, Y: 8}, &Point2D{X: 7, Y: 4}, &Point2D{X: 7, Y: 7}, &Point2D{X: 8, Y: 2}, &Point2D{X: 8, Y: 5}, &Point2D{X: 9, Y: 9}, &Point2D{X: 3, Y: 1}, &Point2D{X: 4, Y: 2}, &Point2D{X: 9, Y: 2}, &Point2D{X: 6, Y: 5}, &Point2D{X: 3, Y: 8}, &Point2D{X: 6, Y: 2}, &Point2D{X: 1, Y: 3}, &Point2D{X: 3, Y: 3}, &Point2D{X: 6, Y: 4}, &Point2D{X: 9, Y: 8}, &Point2D{X: 2, Y: 1}, &Point2D{X: 2, Y: 8}, &Point2D{X: 3, Y: 1}, &Point2D{X: 7, Y: 3}, &Point2D{X: 3, Y: 9}, &Point2D{X: 4, Y: 4}, &Point2D{X: 5, Y: 3}, &Point2D{X: 9, Y: 6}}),
			input:    kdrange.New(2, 10, 20, 30),
			expected: []kdtree.Point{},
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
		input []kdtree.Point
		r     kdrange.Range
	}{
		{name: "nodes: 100 range: -100 50 -50 100", input: generateTestCaseData(100), r: kdrange.New(-100, 50, -50, 100)},
		{name: "nodes: 1000 range: -100 50 -50 100", input: generateTestCaseData(1000), r: kdrange.New(-100, 50, -50, 100)},
		{name: "nodes: 10000 range: -100 50 -50 100", input: generateTestCaseData(10000), r: kdrange.New(-100, 50, -50, 100)},
		{name: "nodes: 100000 range: -500 250 -250 500", input: generateTestCaseData(100000), r: kdrange.New(-500, 250, -250, 500)},
		{name: "nodes: 1000000 range: -500 250 -250 500", input: generateTestCaseData(1000000), r: kdrange.New(-500, 250, -250, 500)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := kdtree.New(test.input)
			assert.ElementsMatch(t, filterRangeSearch(test.input, test.r), tree.RangeSearch(test.r))
		})
	}
}

// TestKDTree_RemoveAxisInversion is a targeted test for issue #6.
//
// https://github.com/kyroy/kdtree/issues/6
//
// Remove wasn't correctly taking into account the axis when searching for
// replacements/substitutes. This caused an incorrect result when removing the
// root node from this tree.
//
// This is because the {171, 176} node starts on the 'left' branch of the
// {238, 155} node, which is correct if indexed by the X axis. When the root
// node is removed, {238, 155} instead becomes indexed on the Y axis, but
// {171, 176} was being left on the 'left' branch.
//
// This test verifies the fix and should help prevent regressions
func TestKDTree_RemoveAxisInversion(t *testing.T) {
	tree := kdtree.New([]kdtree.Point{
		&Point2D{X: 171, Y: 176},
		&Point2D{X: 238, Y: 155},
		&Point2D{X: 257, Y: 246},
		&Point2D{X: 181, Y: 265},
		&Point2D{X: 206, Y: 282},
		&Point2D{X: 265, Y: 176},
		&Point2D{X: 284, Y: 209},
		&Point2D{X: 296, Y: 168},
		&Point2D{X: 280, Y: 225},
		&Point2D{X: 288, Y: 283},
		&Point2D{X: 289, Y: 292},
	})
	search := &Point2D{X: 150, Y: 218}
	remove := &Point2D{X: 265, Y: 176}

	tree.Remove(remove)

	fewNN := tree.KNN(search, 1)
	manyNN := tree.KNN(search, 10)

	assertPointsEqual(t, fewNN[0], manyNN[0])
}

func TestKDTree_RemoveAxisInversionGenerator(t *testing.T) {
	for dims := 2; dims <= 4; dims++ {
		maxSize := int(math.Pow(float64(dims), 4))

		tree := kdtree.New(nil)
		arr := make([]kdtree.Point, 0, maxSize+1)
		for i := 0; i < 1000; i++ {
			p := generateTestPoint(dims)

			// Two KNN queries
			fewNN := tree.KNN(p, 1)
			manyNN := tree.KNN(p, maxSize)

			if len(arr) > 0 {
				assertPointsEqual(t, fewNN[0], manyNN[0])
			}

			// Add in the new point
			arr = append(arr, p)
			tree.Insert(p)

			// Limit the max number of elements - which will also
			// introduce some churn in the tree
			if len(arr) > maxSize {
				idx := rand.Intn(len(arr))
				tree.Remove(arr[idx])
				arr[idx] = arr[len(arr)-1]
				arr = arr[:len(arr)-1]
			}
		}
	}
}

// benchmarks

var resultTree *kdtree.KDTree
var resultPoints []kdtree.Point

func BenchmarkNew(b *testing.B) {
	benchmarks := []struct {
		name  string
		input []kdtree.Point
	}{
		{name: "100", input: generateTestCaseData(100)},
		{name: "1000", input: generateTestCaseData(1000)},
		{name: "10000", input: generateTestCaseData(10000)},
		{name: "100000", input: generateTestCaseData(100000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var t *kdtree.KDTree
			for i := 0; i < b.N; i++ {
				t = kdtree.New(bm.input)
			}
			resultTree = t
		})
	}
}

func BenchmarkKNN(b *testing.B) {
	benchmarks := []struct {
		name   string
		target kdtree.Point
		k      int
		input  []kdtree.Point
	}{
		{name: "p:100,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(100)},
		{name: "p:1000,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(1000)},
		{name: "p:10000,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(10000)},
		{name: "p:100000,k:5", target: &Point2D{}, k: 5, input: generateTestCaseData(100000)},
	}
	for _, bm := range benchmarks {
		var res []kdtree.Point
		tree := kdtree.New(bm.input)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				res = tree.KNN(bm.target, bm.k)
			}
			resultPoints = res
		})
	}
}

// helpers

func generateTestCaseData(size int) []kdtree.Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var points []kdtree.Point
	for i := 0; i < size; i++ {
		points = append(points, &Point2D{X: r.Float64()*3000 - 1500, Y: r.Float64()*3000 - 1500})
	}

	return points
}

func generateTestPoint(dimensions int) kdtree.Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	values := make([]float64, dimensions)
	for j := range values {
		values[j] = r.Float64()*3000 - 1500
	}
	return NewPoint(values, nil)
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

func filterRangeSearch(points []kdtree.Point, r kdrange.Range) []kdtree.Point {
	result := make([]kdtree.Point, 0)

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

func distance(p1, p2 kdtree.Point) float64 {
	sum := 0.
	for i := 0; i < p1.Dimensions(); i++ {
		sum += math.Pow(p1.Dimension(i)-p2.Dimension(i), 2.0)
	}
	return math.Sqrt(sum)
}

func assertPointsEqual(t *testing.T, p1 kdtree.Point, p2 kdtree.Point) {
	assert.Equal(t, p1.Dimensions(), p2.Dimensions())
	for i := 0; i < p1.Dimensions(); i++ {
		assert.Equal(t, p1.Dimension(i), p2.Dimension(i), "assert equal dimension %d", i)
	}
}

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

package points_test

import (
	"github.com/kyroy/kdtree/points"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPoint(t *testing.T) {
	tests := []struct {
		name        string
		coordinates []float64
		data        interface{}
	}{
		{name: "nil nil", coordinates: nil, data: nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := points.NewPoint(test.coordinates, test.data)
			assert.Equal(t, test.coordinates, p.Coordinates)
			assert.Equal(t, test.data, p.Data)
			assert.Equal(t, len(test.coordinates), p.Dimensions())
			for i, v := range test.coordinates {
				assert.Equal(t, v, p.Dimension(i))
			}
		})
	}
}

func TestPoint_Dimensions(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
	}{
		{name: "empty", input: []float64{}},
		{name: "1", input: []float64{1}},
		{name: "2", input: []float64{2.34, 42.}},
		{name: "6", input: []float64{2.34, 42., 2.7, -1.2, 4.3, -0.2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := points.Point{Coordinates: test.input}
			assert.Equal(t, p.Dimensions(), len(test.input))
		})
	}
}

func TestPoint_Dimension(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
	}{
		{name: "empty", input: []float64{}},
		{name: "1", input: []float64{1}},
		{name: "2", input: []float64{2.34, 42.}},
		{name: "6", input: []float64{2.34, 42., 2.7, -1.2, 4.3, -0.2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := points.Point{Coordinates: test.input}
			for i, itm := range test.input {
				assert.Equal(t, p.Dimension(i), itm)
			}
		})
	}
}

func TestPoint_String(t *testing.T) {
	tests := []struct {
		name     string
		point    points.Point
		expected string
	}{
		{name: "empty", point: points.Point{}, expected: "{[] <nil>}"},
		{name: "string data", point: points.Point{Coordinates: []float64{1, 2}, Data: "testme"}, expected: "{[1 2] testme}"},
		{name: "float data", point: points.Point{Coordinates: []float64{1, 2}, Data: 4.3}, expected: "{[1 2] 4.3}"},
		{name: "int data", point: points.Point{Coordinates: []float64{1, 2}, Data: 42}, expected: "{[1 2] 42}"},
		{
			name: "struct data",
			point: points.Point{
				Coordinates: []float64{1, 2},
				Data: struct {
					A int
					B string
				}{
					A: 4,
					B: "teststruct",
				},
			},
			expected: "{[1 2] {4 teststruct}}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.point.String(), test.expected)
		})
	}
}

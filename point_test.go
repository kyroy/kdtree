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
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SamplePoint represents a n-dimensional point of the k-d tree.
type SamplePoint struct {
	Coordinates []float64
	Data        interface{}
}

// NewSamplePoint creates a new point at the given coordinates and contains the given data.
func NewSamplePoint(coordinates []float64, data interface{}) *SamplePoint {
	return &SamplePoint{
		Coordinates: coordinates,
		Data:        data,
	}
}

// Dimensions returns the total number of dimensions.
func (p *SamplePoint) Dimensions() int {
	return len(p.Coordinates)
}

// Dimension returns the value of the i-th dimension.
func (p *SamplePoint) Dimension(i int) float64 {
	return p.Coordinates[i]
}

// String returns the string representation of the point.
func (p *SamplePoint) String() string {
	return fmt.Sprintf("{%v %v}", p.Coordinates, p.Data)
}

// Distance implements Distance for points.Point
func (p1 *SamplePoint) Distance(p2 Point) float64 {
	sum := 0.
	for i := 0; i < p1.Dimensions(); i++ {
		sum += math.Pow(p1.Dimension(i)-p2.Dimension(i), 2.0)
	}
	return math.Sqrt(sum)
}

// PlaneDistance implements PlaneDistance for points.Point
func (p *SamplePoint) PlaneDistance(planePosition float64, dim int) float64 {
	return math.Abs(planePosition - p.Dimension(dim))
}

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
			p := NewSamplePoint(test.coordinates, test.data)
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
			p := SamplePoint{Coordinates: test.input}
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
			p := SamplePoint{Coordinates: test.input}
			for i, itm := range test.input {
				assert.Equal(t, p.Dimension(i), itm)
			}
		})
	}
}

func TestPoint_String(t *testing.T) {
	tests := []struct {
		name     string
		point    SamplePoint
		expected string
	}{
		{name: "empty", point: SamplePoint{}, expected: "{[] <nil>}"},
		{name: "string data", point: SamplePoint{Coordinates: []float64{1, 2}, Data: "testme"}, expected: "{[1 2] testme}"},
		{name: "float data", point: SamplePoint{Coordinates: []float64{1, 2}, Data: 4.3}, expected: "{[1 2] 4.3}"},
		{name: "int data", point: SamplePoint{Coordinates: []float64{1, 2}, Data: 42}, expected: "{[1 2] 42}"},
		{
			name: "struct data",
			point: SamplePoint{
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

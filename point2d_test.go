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

// Point2D ...
type SamplePoint2D struct {
	X float64
	Y float64
}

// Dimensions ...
func (p *SamplePoint2D) Dimensions() int {
	return 2
}

// Dimension ...
func (p *SamplePoint2D) Dimension(i int) float64 {
	if i == 0 {
		return p.X
	}
	return p.Y
}

// String ...
func (p *SamplePoint2D) String() string {
	return fmt.Sprintf("{%.2f %.2f}", p.X, p.Y)
}

// Distance implements *Point Distance for Point2D
func (p *SamplePoint2D) Distance(p2 Point) float64 {
	sum := 0.
	for i := 0; i < p.Dimensions(); i++ {
		sum += math.Pow(p.Dimension(i)-p2.Dimension(i), 2.0)
	}
	return math.Sqrt(sum)
}

// PlaneDistance implements PlaneDistance for SamplePoint2D
func (p *SamplePoint2D) PlaneDistance(planePosition float64, dim int) float64 {
	return math.Abs(planePosition - p.Dimension(dim))
}

func TestPoint2D_Dimensions(t *testing.T) {
	assert.Equal(t, (&SamplePoint2D{}).Dimensions(), 2)
}

func TestPoint2D_TestDimension(t *testing.T) {
	tests := []struct {
		name  string
		point SamplePoint2D
	}{
		{name: "empty", point: SamplePoint2D{}},
		{name: "X", point: SamplePoint2D{X: 2}},
		{name: "Y", point: SamplePoint2D{Y: 3}},
		{name: "XY", point: SamplePoint2D{X: 2, Y: 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.point.Dimension(0), test.point.X)
			assert.Equal(t, test.point.Dimension(1), test.point.Y)
		})
	}
}

func TestPoint2D_String(t *testing.T) {
	tests := []struct {
		name     string
		point    SamplePoint2D
		expected string
	}{
		{name: "empty", point: SamplePoint2D{}, expected: "{0.00 0.00}"},
		{name: "X", point: SamplePoint2D{X: 2}, expected: "{2.00 0.00}"},
		{name: "Y", point: SamplePoint2D{Y: 3}, expected: "{0.00 3.00}"},
		{name: "XY", point: SamplePoint2D{X: 2, Y: 3}, expected: "{2.00 3.00}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.point.String(), test.expected)
		})
	}
}

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

// Package points contains multiple example implementations of the kdtree.Point interface.
package points

import "fmt"

// Point represents a n-dimensional point of the k-d tree.
type Point struct {
	Coordinates []float64
	Data        interface{}
}

// NewPoint creates a new point at the given coordinates and contains the given data.
func NewPoint(coordinates []float64, data interface{}) *Point {
	return &Point{
		Coordinates: coordinates,
		Data:        data,
	}
}

// Dimensions returns the total number of dimensions.
func (p *Point) Dimensions() int {
	return len(p.Coordinates)
}

// Dimension returns the value of the i-th dimension.
func (p *Point) Dimension(i int) float64 {
	return p.Coordinates[i]
}

// String returns the string representation of the point.
func (p *Point) String() string {
	return fmt.Sprintf("{%v %v}", p.Coordinates, p.Data)
}

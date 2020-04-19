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

// Package kdrange contains k-dimensional range struct and helpers.
package kdrange

// Range represents a range in k-dimensional space.
type Range [][2]float64

// New creates a new Range.
//
// It accepts a sequence of min/max pairs that define the Range.
// For example a 2-dimensional rectangle with the with 2 and height 3, starting at (1,2):
//
//     r := NewRange(1, 3, 2, 5)
//
// I.e.:
//     x (dim 0): 1 <= x <= 3
//     y (dim 1): 2 <= y <= 5
func New(limits ...float64) Range {
	if limits == nil || len(limits)%2 != 0 {
		return nil
	}
	r := make([][2]float64, len(limits)/2)
	for i := range r {
		r[i] = [2]float64{limits[2*i], limits[2*i+1]}
	}
	return r
}

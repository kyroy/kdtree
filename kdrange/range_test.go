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

package kdrange_test

import (
	"github.com/kyroy/kdtree/kdrange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRange(t *testing.T) {
	tests := []struct {
		name   string
		input  []float64
		output kdrange.Range
	}{
		{
			name:   "nil",
			input:  nil,
			output: nil,
		},
		{
			name:   "uneven",
			input:  []float64{0},
			output: nil,
		},
		{
			name:   "empty",
			input:  []float64{},
			output: [][2]float64{},
		},
		{
			name:   "2d",
			input:  []float64{1, 2, 3, 4},
			output: [][2]float64{{1, 2}, {3, 4}},
		},
		{
			name:   "3d",
			input:  []float64{4, 1, 6, 19, 3, 1},
			output: [][2]float64{{4, 1}, {6, 19}, {3, 1}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.output, kdrange.New(test.input...))
		})
	}
}

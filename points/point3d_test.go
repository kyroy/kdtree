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

func TestPoint3D_Dimensions(t *testing.T) {
	assert.Equal(t, (&points.Point3D{}).Dimensions(), 3)
}

func TestPoint3D_TestDimension(t *testing.T) {
	tests := []struct {
		name  string
		point points.Point3D
	}{
		{name: "empty", point: points.Point3D{}},
		{name: "X", point: points.Point3D{X: 2}},
		{name: "Y", point: points.Point3D{Y: 3}},
		{name: "Z", point: points.Point3D{Y: 4}},
		{name: "XYZ", point: points.Point3D{X: 2, Y: 3, Z: 4}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.point.Dimension(0), test.point.X)
			assert.Equal(t, test.point.Dimension(1), test.point.Y)
			assert.Equal(t, test.point.Dimension(2), test.point.Z)
		})
	}
}

func TestPoint3D_String(t *testing.T) {
	tests := []struct {
		name     string
		point    points.Point3D
		expected string
	}{
		{name: "empty", point: points.Point3D{}, expected: "{0.00 0.00 0.00}"},
		{name: "X", point: points.Point3D{X: 2}, expected: "{2.00 0.00 0.00}"},
		{name: "Y", point: points.Point3D{Y: 3}, expected: "{0.00 3.00 0.00}"},
		{name: "Z", point: points.Point3D{Z: 4}, expected: "{0.00 0.00 4.00}"},
		{name: "XY", point: points.Point3D{X: 2, Y: 3, Z: 4}, expected: "{2.00 3.00 4.00}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.point.String(), test.expected)
		})
	}
}

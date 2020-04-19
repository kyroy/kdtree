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

package points

import "fmt"

// Point3D ...
type Point3D struct {
	X float64
	Y float64
	Z float64
}

// Dimensions ...
func (p *Point3D) Dimensions() int {
	return 3
}

// Dimension ...
func (p *Point3D) Dimension(i int) float64 {
	switch i {
	case 0:
		return p.X
	case 1:
		return p.Y
	default:
		return p.Z
	}
}

// String ...
func (p *Point3D) String() string {
	return fmt.Sprintf("{%.2f %.2f %.2f}", p.X, p.Y, p.Z)
}

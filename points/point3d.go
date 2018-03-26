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

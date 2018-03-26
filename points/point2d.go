package points

import "fmt"

// Point2D ...
type Point2D struct {
	X float64
	Y float64
}

// Dimensions ...
func (p *Point2D) Dimensions() int {
	return 2
}

// Dimension ...
func (p *Point2D) Dimension(i int) float64 {
	if i == 0 {
		return p.X
	}
	return p.Y
}

// String ...
func (p *Point2D) String() string {
	return fmt.Sprintf("{%.2f %.2f}", p.X, p.Y)
}

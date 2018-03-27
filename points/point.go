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

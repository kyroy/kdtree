package points

import "fmt"

// Point ...
type Point struct {
	Coordinates []float64
	Data        interface{}
}

// Dimensions ...
func (p *Point) Dimensions() int {
	return len(p.Coordinates)
}

// Dimension ...
func (p *Point) Dimension(i int) float64 {
	return p.Coordinates[i]
}

// String ...
func (p *Point) String() string {
	return fmt.Sprintf("{%v %v}", p.Coordinates, p.Data)
}

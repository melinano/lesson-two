package visitor

import (
	"fmt"
	"math"
)

type AreaCalculator struct {
	area float64
}

func (a *AreaCalculator) visitForSquare(s *Square) {
	// Calculate area for square.
	// Then assign it to the area instance variable.
	fmt.Println("Calculating area for square")
	a.area = float64(s.Side * s.Side)
	fmt.Printf("Area of square %f", a.area)
}

func (a *AreaCalculator) visitForCircle(s *Circle) {
	fmt.Println("Calculating area for circle")
	a.area = math.Pow(float64(s.Radius), 2) * math.Pi
	fmt.Printf("Area of circle %f\n", a.area)
}
func (a *AreaCalculator) visitForRectangle(s *Rectangle) {
	fmt.Println("Calculating area for rectangle\n")
	a.area = float64(s.Width * s.Height)
	fmt.Printf("Area of rectangle %f\n", a.area)
}

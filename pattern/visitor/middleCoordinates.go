package visitor

import "fmt"

type MiddleCoordinates struct {
	x float64
	y float64
}

func (a *MiddleCoordinates) visitForSquare(s *Square) {
	// Calculate middle point coordinates for square.
	// Then assign in to the x and y instance variable.
	fmt.Println("Calculating middle point coordinates for square")
	a.x = s.Side / 2
	a.y = s.Side / 2
	fmt.Printf("Center of square %f;%f\n", a.x, a.y)
}

func (a *MiddleCoordinates) visitForCircle(c *Circle) {
	fmt.Println("Calculating middle point coordinates for circle")
	a.x = c.Radius
	a.y = c.Radius
	fmt.Printf("Center of circle %f;%f\n", a.x, a.y)
}
func (a *MiddleCoordinates) visitForRectangle(t *Rectangle) {
	fmt.Println("Calculating middle point coordinates for rectangle")
	a.x = t.Width / 2
	a.y = t.Height / 2
	fmt.Printf("Center of rectangle %f;%f\n", a.x, a.y)
}

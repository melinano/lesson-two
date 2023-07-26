package visitor

type Circle struct {
	Radius float64
}

func (c *Circle) Accept(v Visitor) {
	v.visitForCircle(c)
}

func (c *Circle) getType() string {
	return "Circle"
}

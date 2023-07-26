package visitor

type Rectangle struct {
	Height float64
	Width  float64
}

func (t *Rectangle) Accept(v Visitor) {
	v.visitForRectangle(t)
}

func (t *Rectangle) getType() string {
	return "rectangle"
}

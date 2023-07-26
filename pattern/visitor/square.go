package visitor

type Square struct {
	Side float64
}

func (s *Square) Accept(v Visitor) {
	v.visitForSquare(s)
}

func (s *Square) getType() string {
	return "Square"
}

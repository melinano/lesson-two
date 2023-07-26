package chain_of_responsability

import "fmt"

type Medical struct {
	next Department
}

func (m *Medical) SetNext(next Department) {
	m.next = next
}

func (m *Medical) execute(p *Patient) {
	if p.medicineDone {
		fmt.Println("Medicine already given to patient")
		m.next.execute(p)
		return
	}
	fmt.Println("Medical giving medicine to patient")
	p.medicineDone = true
	m.next.execute(p)
}

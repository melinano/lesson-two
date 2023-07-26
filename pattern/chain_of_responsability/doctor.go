package chain_of_responsability

import "fmt"

type Doctor struct {
	next Department
}

func (d *Doctor) SetNext(next Department) {
	d.next = next
}

func (d *Doctor) execute(p *Patient) {
	if p.doctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		d.next.execute(p)
		return
	}
	fmt.Println("Doctor checking patient")
	p.doctorCheckUpDone = true
	d.next.execute(p)
}

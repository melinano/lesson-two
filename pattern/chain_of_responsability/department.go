package chain_of_responsability

type Department interface {
	execute(*Patient)
	SetNext(Department)
}

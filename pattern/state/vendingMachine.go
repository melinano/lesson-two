package state

import "fmt"

type VendingMachine struct {
	hasItem       State
	itemRequested State
	hasMoney      State
	noItem        State

	currentState State

	itemCount int
	itemPrice int
}

func NewVendingMachine(itemCount, itemPrice int) *VendingMachine {
	// setting itemCount and Price
	v := &VendingMachine{
		itemCount: itemCount,
		itemPrice: itemPrice,
	}
	// assigning the created VendingMachine to each state
	hasItemState := &HasItemState{
		vendingMachine: v,
	}
	itemRequestedState := &ItemRequestedState{
		vendingMachine: v,
	}
	hasMoneyState := &HasMoneyState{
		vendingMachine: v,
	}
	noItemState := &NoItemState{
		vendingMachine: v,
	}

	// setting current state
	v.setState(hasItemState)
	v.hasItem = hasItemState
	v.itemRequested = itemRequestedState
	v.hasMoney = hasMoneyState
	v.noItem = noItemState
	return v
}

// functions that execute a specific state-function dependent on the current state
func (v *VendingMachine) RequestItem() error {
	return v.currentState.requestItem()
}

func (v *VendingMachine) AddItem(count int) error {
	return v.currentState.addItem(count)
}

func (v *VendingMachine) InsertMoney(money int) error {
	return v.currentState.insertMoney(money)
}

func (v *VendingMachine) DispenseItem() error {
	return v.currentState.dispenseItem()
}

func (v *VendingMachine) setState(s State) {
	v.currentState = s
}

// incrementing the ItemCount in the VendingMachine
func (v *VendingMachine) incrementItemCount(count int) {
	fmt.Printf("Adding %d items\n", count)
	v.itemCount = v.itemCount + count
}

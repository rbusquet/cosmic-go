package model

import "time"

type OrderLine struct {
	OrderID  string
	SKU      string
	Quantity int
}

type Batch struct {
	Reference         string
	SKU               string
	PurchasedQuantity int
	ETA               time.Time

	allocations map[OrderLine]bool
}

func NewBatch(ref string, sku string, qty int, eta time.Time) Batch {
	allocations := make(map[OrderLine]bool)
	return Batch{Reference: ref, SKU: sku, ETA: eta, PurchasedQuantity: qty, allocations: allocations}
}

func (b *Batch) AvailableQuantity() int {
	return b.PurchasedQuantity - b.AllocatedQuantity()
}

func (b *Batch) AllocatedQuantity() int {
	allocated := 0
	for line := range b.allocations {
		allocated += line.Quantity
	}
	return allocated
}

func (b *Batch) Allocate(line OrderLine) {
	if b.allocations == nil {
		b.allocations = make(map[OrderLine]bool)
	}
	if b.CanAllocate(line) {
		b.allocations[line] = true
	}
}

func (b *Batch) CanAllocate(line OrderLine) bool {
	return b.SKU == line.SKU && b.AvailableQuantity() >= line.Quantity
}

func (b *Batch) Deallocate(line OrderLine) {
	if b.allocations[line] {
		delete(b.allocations, line)
	}
}

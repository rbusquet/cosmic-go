package model

import (
	"time"
)

type OrderLine struct {
	OrderID  string
	SKU      string
	Quantity int
}

type Batch struct {
	Reference string
	SKU       string
	ETA       time.Time

	purchasedQuantity int
	allocations       map[OrderLine]bool
}

func NewBatch(ref string, sku string, qty int, eta time.Time) Batch {
	allocations := make(map[OrderLine]bool)
	return Batch{ref, sku, eta, qty, allocations}
}

func (b *Batch) AvailableQuantity() int {
	return b.purchasedQuantity - b.AllocatedQuantity()
}

func (b *Batch) AllocatedQuantity() int {
	allocated := 0
	for line := range b.allocations {
		allocated += line.Quantity
	}
	return allocated
}

func (b *Batch) Allocate(line OrderLine) {
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

package model_test

import (
	"testing"
	"time"

	"github.com/rbusquet/cosmic-go/model"
	"github.com/stretchr/testify/assert"
)

type input struct {
	batch model.Batch
	line  model.OrderLine
}

func makeBatchAndLine(sku string, batchQty int, lineQty int) input {
	batch := model.NewBatch("batch-001", sku, batchQty, time.Now())
	orderLine := model.OrderLine{OrderID: "order-123", SKU: sku, Quantity: lineQty}
	return input{batch, orderLine}
}

func TestBatchAllocate(t *testing.T) {
	t.Run("allocating to a batch reduces the available quantity", func(t *testing.T) {
		batch := model.NewBatch("batch-001", "SMALL-TABLE", 20, time.Now())
		line := model.OrderLine{OrderID: "order-ref", SKU: "SMALL-TABLE", Quantity: 2}
		batch.Allocate(line)

		assert.Equal(t, 18, batch.AvailableQuantity())
	})
	t.Run("allocation is idempotent", func(t *testing.T) {
		batch := model.NewBatch("batch-001", "SMALL-TABLE", 20, time.Now())
		line := model.OrderLine{OrderID: "order-ref", SKU: "SMALL-TABLE", Quantity: 2}

		batch.Allocate(line)
		batch.Allocate(line)

		assert.Equal(t, 18, batch.AvailableQuantity())
	})

}

func TestBatchCanAllocate(t *testing.T) {
	tests := map[string]struct {
		input input
		want  bool
	}{
		"can allocate if available greater than required": {
			makeBatchAndLine("ELEGANT-LAMP", 20, 2),
			true,
		},
		"cannot allocate if available smaller than required": {
			makeBatchAndLine("ELEGANT-LAMP", 2, 20),
			false,
		},
		"can allocate if available equal to required": {
			makeBatchAndLine("ELEGANT-LAMP", 2, 2),
			true,
		},
		"cannot allocate if skus do not match": {
			input{
				model.NewBatch("batch-001", "UNCOMFORTABLE-CHAIR", 100, time.Time{}),
				model.OrderLine{OrderID: "order-123", SKU: "EXPENSIVE-TOASTER", Quantity: 10},
			},
			false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.input.batch.CanAllocate(test.input.line)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestBatchDeallocate(t *testing.T) {
	input := makeBatchAndLine("DECORATIVE-TRINKET", 20, 2)
	input.batch.Deallocate(input.line)
	if input.batch.AvailableQuantity() != 20 {
		t.Errorf("Got %d; want %d", input.batch.AvailableQuantity(), 20)
	}
}

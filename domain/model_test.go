package domain_test

import (
	"testing"
	"time"

	"github.com/rbusquet/cosmic-go/domain"
)

type input struct {
	batch domain.Batch
	line  domain.OrderLine
}

func makeBatchAndLine(sku string, batchQty int, lineQty int) input {
	batch := domain.NewBatch("batch-001", sku, batchQty, time.Now())
	orderLine := domain.OrderLine{"order-123", sku, lineQty}
	return input{batch, orderLine}
}

func TestBatchAllocate(t *testing.T) {
	t.Run("allocating to a batch reduces the available quantity", func(t *testing.T) {
		batch := domain.NewBatch("batch-001", "SMALL-TABLE", 20, time.Now())
		line := domain.OrderLine{"order-ref", "SMALL-TABLE", 2}
		batch.Allocate(line)
		if batch.AvailableQuantity() != 18 {
			t.Errorf("Got %d; want %d", batch.AvailableQuantity(), 18)
		}
	})
	t.Run("allocation is idempotent", func(t *testing.T) {
		batch := domain.NewBatch("batch-001", "SMALL-TABLE", 20, time.Now())
		line := domain.OrderLine{"order-ref", "SMALL-TABLE", 2}

		batch.Allocate(line)
		batch.Allocate(line)

		if batch.AvailableQuantity() != 18 {
			t.Errorf("Got %d; want %d", batch.AvailableQuantity(), 18)
		}
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
				domain.NewBatch("batch-001", "UNCOMFORTABLE-CHAIR", 100, time.Time{}),
				domain.OrderLine{"order-123", "EXPENSIVE-TOASTER", 10},
			},
			false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.input.batch.CanAllocate(test.input.line)
			if got != test.want {
				t.Errorf("expected: %v, got: %v", test.want, got)
			}
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

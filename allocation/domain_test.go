package allocation

import (
	"reflect"
	"testing"
	"time"
)

type input struct {
	batch Batch
	line  OrderLine
}

func makeBatchAndLine(sku string, batchQty int, lineQty int) input {
	batch := NewBatch("batch-001", sku, batchQty, time.Now())
	orderLine := OrderLine{"order-123", sku, lineQty}
	return input{batch, orderLine}
}

func TestBatchAllocate(t *testing.T) {
	t.Run("allocating to a batch reduces the available quantity", func(t *testing.T) {
		batch := NewBatch("batch-001", "SMALL-TABLE", 20, time.Now())
		line := OrderLine{"order-ref", "SMALL-TABLE", 2}
		batch.Allocate(line)
		if batch.AvailableQuantity() != 18 {
			t.Errorf("Got %d; want %d", batch.AvailableQuantity(), 18)
		}
	})
	t.Run("allocation is idempotent", func(t *testing.T) {
		batch := NewBatch("batch-001", "SMALL-TABLE", 20, time.Now())
		line := OrderLine{"order-ref", "SMALL-TABLE", 2}

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
				NewBatch("batch-001", "UNCOMFORTABLE-CHAIR", 100, time.Time{}),
				OrderLine{"order-123", "EXPENSIVE-TOASTER", 10},
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

func TestAllocate(t *testing.T) {
	tests := map[string]struct {
		batches   []Batch
		line      OrderLine
		want      []int
		allocated string
	}{
		"prefers current stock batches to shipments": {
			[]Batch{
				NewBatch("in-stock-batch", "RETRO-CLOCK", 100, time.Time{}),
				NewBatch("shipment-batch", "RETRO-CLOCK", 100, time.Now().Add(24*time.Hour)),
			},
			OrderLine{"oref", "RETRO-CLOCK", 10},
			[]int{90, 100},
			"in-stock-batch",
		},
		"prefers prefers earlier batches": {
			[]Batch{
				NewBatch("speedy-batch", "MINIMALIST-SPOON", 100, time.Now()),
				NewBatch("normal-batch", "MINIMALIST-SPOON", 100, time.Now().Add(24*time.Hour)),
				NewBatch("slow-batch", "MINIMALIST-SPOON", 100, time.Now().Add(48*time.Hour)),
			},
			OrderLine{"oref", "MINIMALIST-SPOON", 10},
			[]int{90, 100, 100},
			"speedy-batch",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ref, err := Allocate(test.line, test.batches...)
			if err != nil {
				t.Fatalf("Error: %+v", err)
			}
			var got []int
			for _, batch := range test.batches {
				got = append(got, batch.AvailableQuantity())
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Got %v; want %v", got, test.want)
			}
			if ref != test.allocated {
				t.Errorf("For return, got %s; want %s", ref, test.allocated)
			}
		})
	}

	t.Run("raises out of stock exception if cannot allocate", func(t *testing.T) {
		batch := NewBatch("batch1", "SMALL-FORK", 10, time.Now())

		_, err := Allocate(OrderLine{"order1", "SMALL-FORK", 10}, batch)
		if err != nil {
			t.Fatalf("Error: %+v", err)
		}
		_, err = Allocate(OrderLine{"order2", "SMALL-FORK", 10}, batch)
		if err == nil {
			t.Errorf("Expected error.")
		}
	})
}

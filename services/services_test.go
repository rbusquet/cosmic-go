package services_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rbusquet/cosmic-go/domain"
	"github.com/rbusquet/cosmic-go/services"
)

func TestAllocate(t *testing.T) {
	tests := map[string]struct {
		batches   []domain.Batch
		line      domain.OrderLine
		want      []int
		allocated string
	}{
		"prefers current stock batches to shipments": {
			[]domain.Batch{
				domain.NewBatch("in-stock-batch", "RETRO-CLOCK", 100, time.Time{}),
				domain.NewBatch("shipment-batch", "RETRO-CLOCK", 100, time.Now().Add(24*time.Hour)),
			},
			domain.OrderLine{OrderID: "oref", SKU: "RETRO-CLOCK", Quantity: 10},
			[]int{90, 100},
			"in-stock-batch",
		},
		"prefers prefers earlier batches": {
			[]domain.Batch{
				domain.NewBatch("speedy-batch", "MINIMALIST-SPOON", 100, time.Now()),
				domain.NewBatch("normal-batch", "MINIMALIST-SPOON", 100, time.Now().Add(24*time.Hour)),
				domain.NewBatch("slow-batch", "MINIMALIST-SPOON", 100, time.Now().Add(48*time.Hour)),
			},
			domain.OrderLine{
				OrderID:  "oref",
				SKU:      "MINIMALIST-SPOON",
				Quantity: 10,
			},
			[]int{90, 100, 100},
			"speedy-batch",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ref, err := services.Allocate(test.line, test.batches...)
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
		batch := domain.NewBatch("batch1", "SMALL-FORK", 10, time.Now())

		_, err := services.Allocate(domain.OrderLine{
			OrderID:  "order1",
			SKU:      "SMALL-FORK",
			Quantity: 10,
		}, batch)
		if err != nil {
			t.Fatalf("Error: %+v", err)
		}
		_, err = services.Allocate(domain.OrderLine{
			OrderID:  "order2",
			SKU:      "SMALL-FORK",
			Quantity: 10,
		}, batch)
		if err == nil {
			t.Errorf("Expected error.")
		}
	})
}

package services_test

import (
	"testing"
	"time"

	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/services"
	"github.com/stretchr/testify/assert"
)

func TestAllocate(t *testing.T) {
	tests := map[string]struct {
		batches   []model.Batch
		line      model.OrderLine
		want      []int
		allocated string
	}{
		"prefers current stock batches to shipments": {
			[]model.Batch{
				model.NewBatch("in-stock-batch", "RETRO-CLOCK", 100, time.Time{}),
				model.NewBatch("shipment-batch", "RETRO-CLOCK", 100, time.Now().Add(24*time.Hour)),
			},
			model.OrderLine{OrderID: "oref", SKU: "RETRO-CLOCK", Quantity: 10},
			[]int{90, 100},
			"in-stock-batch",
		},
		"prefers prefers earlier batches": {
			[]model.Batch{
				model.NewBatch("speedy-batch", "MINIMALIST-SPOON", 100, time.Now()),
				model.NewBatch("normal-batch", "MINIMALIST-SPOON", 100, time.Now().Add(24*time.Hour)),
				model.NewBatch("slow-batch", "MINIMALIST-SPOON", 100, time.Now().Add(48*time.Hour)),
			},
			model.OrderLine{
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
			assert.NoError(t, err)
			var got []int
			for _, batch := range test.batches {
				got = append(got, batch.AvailableQuantity())
			}
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.allocated, ref)
		})
	}

	t.Run("raises out of stock exception if cannot allocate", func(t *testing.T) {
		batch := model.NewBatch("batch1", "SMALL-FORK", 10, time.Now())

		_, err := services.Allocate(model.OrderLine{
			OrderID:  "order1",
			SKU:      "SMALL-FORK",
			Quantity: 10,
		}, batch)
		assert.NoError(t, err)
		_, err = services.Allocate(model.OrderLine{
			OrderID:  "order2",
			SKU:      "SMALL-FORK",
			Quantity: 10,
		}, batch)
		assert.EqualError(t, err, "OutOfStock: no stock for sku SMALL-FORK")
	})
}

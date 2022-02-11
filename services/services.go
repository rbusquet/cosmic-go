package services

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/rbusquet/cosmic-go/domain"
)

func Allocate(line domain.OrderLine, batches ...domain.Batch) (string, error) {
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].ETA.Before(batches[j].ETA)
	})
	for _, batch := range batches {
		if batch.CanAllocate(line) {
			batch.Allocate(line)
			return batch.Reference, nil
		}
	}
	return "", errors.Errorf("Out of stock: no stock for sku %s", line.SKU)
}

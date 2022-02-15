package services

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/repository"
)

func isValidSku(sku string, batches ...*model.Batch) bool {
	for _, b := range batches {
		if b.SKU == sku {
			return true
		}
	}
	return false
}

func Allocate(orderId, sku string, quantity int, repo repository.Repository) (batchref string, err error) {
	batches := repo.List()

	line := model.OrderLine{OrderID: orderId, SKU: sku, Quantity: quantity}

	if !isValidSku(line.SKU, batches...) {
		return "", errors.Errorf("Invalid SKU %s", line.SKU)
	}
	result, err := model.Allocate(line, batches...)
	repo.Save(batches...)
	return result, err
}

func AddBatch(batchref, sku string, quantity int, eta time.Time, repo repository.Repository) {
	batch := model.NewBatch(batchref, sku, quantity, eta)
	repo.Add(batch)
}

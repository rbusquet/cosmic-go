package services

import (
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

func Allocate(line model.OrderLine, repo repository.Repository) (batchref string, err error) {
	batches := repo.List()

	if !isValidSku(line.SKU, batches...) {
		return "", errors.Errorf("Invalid SKU %s", line.SKU)
	}
	result, err := model.Allocate(line, batches...)
	repo.Save(batches...)
	return result, err
}

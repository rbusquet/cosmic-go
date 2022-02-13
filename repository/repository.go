package repository

import (
	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/orm"
	"gorm.io/gorm"
)

type IRepository interface {
	Add(batch *model.Batch)
	Get(reference string) *model.Batch
	List() []model.Batch
}

type GormRepository struct {
	DB *gorm.DB
}

func (r *GormRepository) Add(batch *model.Batch) {
	r.DB.Table("batches").Create(batch)
}

func (r *GormRepository) Get(reference string) *model.Batch {
	var batch orm.Batches
	r.DB.Preload("Allocations").Find(&batch, "reference = ?", reference)
	for _, line := range batch.Allocations {
		batch.Batch.Allocate(line.OrderLine)
	}
	return &batch.Batch
}

func (r *GormRepository) List() []model.Batch {
	var batches []orm.Batches
	r.DB.Model(&orm.Batches{}).Preload("Allocations").Find(&batches)

	var parsedBatches []model.Batch
	for _, batch := range batches {
		for _, line := range batch.Allocations {
			batch.Batch.Allocate(line.OrderLine)
		}
		parsedBatches = append(parsedBatches, batch.Batch)
	}
	return parsedBatches
}

type FakeRepository struct {
	batches map[string]*model.Batch
}

func (r *FakeRepository) Add(batch *model.Batch) {
	r.batches[batch.Reference] = batch
}

func (r *FakeRepository) Get(reference string) *model.Batch {
	return r.batches[reference]
}

func (r *FakeRepository) List() []model.Batch {
	var batches []model.Batch
	for _, batch := range r.batches {
		batches = append(batches, *batch)
	}
	return batches
}

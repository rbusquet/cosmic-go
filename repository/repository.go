package repository

import (
	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/orm"
	"gorm.io/gorm"
)

type Repository interface {
	Add(batch *model.Batch)
	Get(reference string) *model.Batch
	List() []*model.Batch
	Save(batches ...*model.Batch)
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

func (r *GormRepository) List() []*model.Batch {
	var batches []orm.Batches
	r.DB.Model(&orm.Batches{}).Preload("Allocations").Find(&batches)

	var parsedBatches []*model.Batch
	for _, batch := range batches {
		bBatch := batch.Batch
		for _, line := range batch.Allocations {
			bBatch.Allocate(line.OrderLine)
		}
		parsedBatches = append(parsedBatches, &bBatch)
	}
	return parsedBatches
}

func (r *GormRepository) Save(batches ...*model.Batch) {
	for _, batch := range batches {
		var batchID uint
		r.DB.Table("batches").Where("reference = ?", batch.Reference).Select("id").Scan(&batchID)
		var oallocations []orm.OrderLines
		for _, allocation := range batch.Allocations() {
			oallocations = append(oallocations, orm.OrderLines{OrderLine: allocation})
		}
		obatch := orm.Batches{Model: gorm.Model{ID: batchID}, Batch: *batch, Allocations: oallocations}
		r.DB.Where("reference = ?", batch.Reference).Updates(obatch)
	}
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

func (r *FakeRepository) List() []*model.Batch {
	var batches []*model.Batch
	for _, batch := range r.batches {
		batches = append(batches, batch)
	}
	return batches
}

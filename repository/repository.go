package repository

import (
	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/orm"
	"gorm.io/gorm"
)

type Repository interface {
	Add(batch *model.Batch) uint
	Get(reference string) *model.Batch
	List() []*model.Batch
	Save(batches ...*model.Batch)
}

type GormRepository struct {
	DB *gorm.DB
}

func (r *GormRepository) Add(batch *model.Batch) uint {
	toSave := orm.Batches{Batch: *batch}
	r.DB.Create(&toSave)
	return toSave.ID
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

type fakeRepository struct {
	batches map[uint]*model.Batch
	Saved   bool
	LastID  uint
}

func NewFakeRepository(batches ...*model.Batch) *fakeRepository {
	r := new(fakeRepository)
	r.batches = make(map[uint]*model.Batch)
	for _, b := range batches {
		r.LastID += 1
		r.batches[r.LastID] = b
	}
	return r
}

func (r *fakeRepository) Add(batch *model.Batch) uint {
	r.LastID += 1
	r.batches[r.LastID] = batch
	return r.LastID
}

func (r *fakeRepository) Get(reference string) *model.Batch {
	for _, batch := range r.batches {
		if batch.Reference == reference {
			return batch
		}
	}
	return nil
}

func (r *fakeRepository) List() []*model.Batch {
	var batches []*model.Batch
	for _, batch := range r.batches {
		batches = append(batches, batch)
	}
	return batches
}

func (r *fakeRepository) Save(batches ...*model.Batch) {
	var newBatches []*model.Batch
	for _, batch := range batches {
		saved := false
		for id, currentBatch := range r.batches {
			if batch.Reference == currentBatch.Reference {
				r.batches[id] = batch
				saved = true
				break
			}
		}
		if !saved {
			newBatches = append(newBatches, batch)
		}
	}
	for _, batch := range newBatches {
		r.LastID += 1
		r.batches[r.LastID] = batch
	}
	r.Saved = true
}

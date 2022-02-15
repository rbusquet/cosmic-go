package services_test

import (
	"testing"
	"time"

	"github.com/rbusquet/cosmic-go/repository"
	"github.com/rbusquet/cosmic-go/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServicesSuite struct {
	suite.Suite
}

func (suite *ServicesSuite) TestReturnsAllocation() {
	repo := repository.NewFakeRepository()
	services.AddBatch("b1", "COMPLICATED-LAMP", 100, time.Time{}, repo)

	result, err := services.Allocate("o1", "COMPLICATED-LAMP", 10, repo)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "b1", result)
}

func (suite *ServicesSuite) TestErrorForInvalidSku() {
	repo := repository.NewFakeRepository()
	services.AddBatch("b1", "AREALSKU", 100, time.Time{}, repo)

	_, err := services.Allocate("o1", "NONEXISTENTSKU", 10, repo)

	assert.EqualError(suite.T(), err, "Invalid SKU NONEXISTENTSKU")
}

func (suite *ServicesSuite) TestSaves() {
	repo := repository.NewFakeRepository()
	services.AddBatch("b1", "SOMETHING-ELSE", 100, time.Time{}, repo)

	services.Allocate("o1", "SOMETHING-ELSE", 10, repo)

	assert.Equal(suite.T(), true, repo.Saved)
}

func (suite *ServicesSuite) TestPrefersWarehouseStockBatchesToShipments() {
	tomorrow := time.Now().AddDate(0, 0, 1)
	repo := repository.NewFakeRepository()
	services.AddBatch("in-stock-batch", "RETRO-CLOCK", 100, time.Time{}, repo)
	services.AddBatch("shipment-batch", "RETRO-CLOCK", 100, tomorrow, repo)

	batch, err := services.Allocate("oref", "RETRO-CLOCK", 10, repo)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "in-stock-batch", batch)
}

func (suite *ServicesSuite) TestPrefersEarlierBatches() {
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	later := today.AddDate(0, 1, 0)
	repo := repository.NewFakeRepository()
	services.AddBatch("normal-batch", "MINIMALIST-SPOON", 100, tomorrow, repo)
	services.AddBatch("speedy-batch", "MINIMALIST-SPOON", 100, today, repo)
	services.AddBatch("slow-batch", "MINIMALIST-SPOON", 100, later, repo)

	batch, err := services.Allocate("order1", "MINIMALIST-SPOON", 10, repo)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "speedy-batch", batch)
}

func (suite *ServicesSuite) TestReturnsAllocatedBatchRef() {
	tomorrow := time.Now().AddDate(0, 0, 1)
	repo := repository.NewFakeRepository()
	services.AddBatch("in-stock-batch", "HIGHBROW-POSTER", 100, time.Time{}, repo)
	services.AddBatch("shipment-batch", "HIGHBROW-POSTER", 100, tomorrow, repo)

	allocation, err := services.Allocate("oref", "HIGHBROW-POSTER", 10, repo)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "in-stock-batch", allocation)
}

func (suite *ServicesSuite) TestRaisesOutOfStockExceptionIfCannotAllocate() {
	repo := repository.NewFakeRepository()
	services.AddBatch("batch1", "SMALL-FORK", 10, time.Now(), repo)

	_, err := services.Allocate("order1", "SMALL-FORK", 10, repo)
	assert.NoError(suite.T(), err)

	_, err = services.Allocate("order2", "SMALL-FORK", 10, repo)
	assert.EqualError(suite.T(), err, "Out of stock for SKU SMALL-FORK")
}

func (suite *ServicesSuite) TestAddBatch() {
	repo := repository.NewFakeRepository()

	services.AddBatch("b1", "CRUNCHY-ARMCHAIR", 100, time.Now(), repo)

	assert.NotNil(suite.T(), repo.Get("b1"))
}

func TestServicesSuite(t *testing.T) {
	suite.Run(t, new(ServicesSuite))
}

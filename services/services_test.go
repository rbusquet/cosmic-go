package services_test

import (
	"testing"
	"time"

	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/repository"
	"github.com/rbusquet/cosmic-go/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServicesSuite struct {
	suite.Suite
}

func (suite *ServicesSuite) SetupTest() {
}

func (suite *ServicesSuite) TestReturnsAllocation() {
	line := model.OrderLine{OrderID: "o1", SKU: "COMPLICATED-LAMP", Quantity: 10}
	batch := model.NewBatch("b1", "COMPLICATED-LAMP", 100, time.Time{})
	repo := repository.NewFakeRepository(batch)

	result, err := services.Allocate(line, repo)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "b1", result)
}

func (suite *ServicesSuite) TestErrorForInvalidSku() {
	line := model.OrderLine{OrderID: "o1", SKU: "NONEXISTENTSKU", Quantity: 10}
	batch := model.NewBatch("b1", "AREALSKU", 100, time.Time{})
	repo := repository.NewFakeRepository(batch)

	_, err := services.Allocate(line, repo)
	assert.EqualError(suite.T(), err, "Invalid SKU NONEXISTENTSKU")
}

func (suite *ServicesSuite) TestSaves() {
	line := model.OrderLine{OrderID: "o1", SKU: "SOMETHING-ELSE", Quantity: 10}
	batch := model.NewBatch("b1", "SOMETHING-ELSE", 100, time.Time{})
	repo := repository.NewFakeRepository(batch)

	services.Allocate(line, repo)
	assert.Equal(suite.T(), true, repo.Saved)
}

func (suite *ServicesSuite) TestPrefersWarehouseStockBatchesToShipments() {
	tomorrow := time.Now().AddDate(0, 0, 1)
	inStockBatch := model.NewBatch("in-stock-batch", "RETRO-CLOCK", 100, time.Time{})
	shipmentBatch := model.NewBatch("shipment-batch", "RETRO-CLOCK", 100, tomorrow)
	repo := repository.NewFakeRepository(inStockBatch, shipmentBatch)

	line := model.OrderLine{OrderID: "oref", SKU: "RETRO-CLOCK", Quantity: 10}

	services.Allocate(line, repo)
	assert.Equal(suite.T(), 90, inStockBatch.AvailableQuantity())
	assert.Equal(suite.T(), 100, shipmentBatch.AvailableQuantity())
}

func (suite *ServicesSuite) TestPrefersEarlierBatches() {
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	later := today.AddDate(0, 1, 0)
	earliest := model.NewBatch("speedy-batch", "MINIMALIST-SPOON", 100, today)
	medium := model.NewBatch("normal-batch", "MINIMALIST-SPOON", 100, tomorrow)
	latest := model.NewBatch("slow-batch", "MINIMALIST-SPOON", 100, later)
	repo := repository.NewFakeRepository(earliest, medium, latest)
	line := model.OrderLine{OrderID: "order1", SKU: "MINIMALIST-SPOON", Quantity: 10}

	services.Allocate(line, repo)

	assert.Equal(suite.T(), 90, earliest.AvailableQuantity())
	assert.Equal(suite.T(), 100, medium.AvailableQuantity())
	assert.Equal(suite.T(), 100, latest.AvailableQuantity())
}

func (suite *ServicesSuite) TestReturnsAllocatedBatchRef() {
	tomorrow := time.Now().AddDate(0, 0, 1)
	inStockBatch := model.NewBatch("in-stock-batch", "HIGHBROW-POSTER", 100, time.Time{})
	shipmentBatch := model.NewBatch("shipment-batch", "HIGHBROW-POSTER", 100, tomorrow)
	repo := repository.NewFakeRepository(inStockBatch, shipmentBatch)

	line := model.OrderLine{OrderID: "oref", SKU: "HIGHBROW-POSTER", Quantity: 10}

	allocation, err := services.Allocate(line, repo)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), inStockBatch.Reference, allocation)
}

func (suite *ServicesSuite) TestRaisesOutOfStockExceptionIfCannotAllocate() {
	batch := model.NewBatch("batch1", "SMALL-FORK", 10, time.Now())
	line := model.OrderLine{OrderID: "order1", SKU: "SMALL-FORK", Quantity: 10}
	_, err := model.Allocate(line, batch)
	assert.NoError(suite.T(), err)
	line2 := model.OrderLine{OrderID: "order2", SKU: "SMALL-FORK", Quantity: 10}
	_, err = model.Allocate(line2, batch)
	assert.EqualError(suite.T(), err, "Out of stock for SKU SMALL-FORK")
}

func TestServicesSuite(t *testing.T) {
	suite.Run(t, new(ServicesSuite))
}

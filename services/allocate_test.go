package services_test

import (
	"testing"
	"time"

	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AllocateSuite struct {
	suite.Suite
}

func (suite *AllocateSuite) TestPrefersCurrentStockBatchesToShipments() {
	tomorrow := time.Now().AddDate(0, 0, 1)
	inStockBatch := model.NewBatch("in-stock-batch", "RETRO-CLOCK", 100, time.Time{})
	shipmentBatch := model.NewBatch("shipment-batch", "RETRO-CLOCK", 100, tomorrow)
	line := model.OrderLine{OrderID: "oref", SKU: "RETRO-CLOCK", Quantity: 10}

	services.Allocate(line, inStockBatch, shipmentBatch)
	assert.Equal(suite.T(), 90, inStockBatch.AvailableQuantity())
	assert.Equal(suite.T(), 100, shipmentBatch.AvailableQuantity())
}

func (suite *AllocateSuite) TestPrefersEarlierBatches() {
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	later := today.AddDate(0, 1, 0)
	earliest := model.NewBatch("speedy-batch", "MINIMALIST-SPOON", 100, today)
	medium := model.NewBatch("normal-batch", "MINIMALIST-SPOON", 100, tomorrow)
	latest := model.NewBatch("slow-batch", "MINIMALIST-SPOON", 100, later)
	line := model.OrderLine{OrderID: "order1", SKU: "MINIMALIST-SPOON", Quantity: 10}

	services.Allocate(line, earliest, medium, latest)
	assert.Equal(suite.T(), 90, earliest.AvailableQuantity())
	assert.Equal(suite.T(), 100, medium.AvailableQuantity())
	assert.Equal(suite.T(), 100, latest.AvailableQuantity())
}

func (suite *AllocateSuite) TestRetrusnAllocatedBatchRef() {
	tomorrow := time.Now().AddDate(0, 0, 1)
	inStockBatch := model.NewBatch("in-stock-batch", "HIGHBROW-POSTER", 100, time.Time{})
	shipmentBatch := model.NewBatch("shipment-batch", "HIGHBROW-POSTER", 100, tomorrow)
	line := model.OrderLine{OrderID: "oref", SKU: "HIGHBROW-POSTER", Quantity: 10}

	allocation, err := services.Allocate(line, inStockBatch, shipmentBatch)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), inStockBatch.Reference, allocation)
}

// def test_raises_out_of_stock_exception_if_cannot_allocate():
//     batch = Batch("batch1", "SMALL-FORK", 10, eta=today)
//     allocate(OrderLine("order1", "SMALL-FORK", 10), [batch])

//     with pytest.raises(OutOfStock, match="SMALL-FORK"):
//         allocate(OrderLine("order2", "SMALL-FORK", 1), [batch])

func (suite *AllocateSuite) TestRaisesOutOfStockExceptionIfCannotAllocate() {
	batch := model.NewBatch("batch1", "SMALL-FORK", 10, time.Now())
	line := model.OrderLine{OrderID: "order1", SKU: "SMALL-FORK", Quantity: 10}
	_, err := services.Allocate(line, batch)
	assert.NoError(suite.T(), err)
	line2 := model.OrderLine{OrderID: "order2", SKU: "SMALL-FORK", Quantity: 10}
	_, err = services.Allocate(line2, batch)
	assert.EqualError(suite.T(), err, "Out of stock for SKU SMALL-FORK")
}

func TestAllocateSuite(t *testing.T) {
	suite.Run(t, new(AllocateSuite))
}

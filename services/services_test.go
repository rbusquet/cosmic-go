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

func TestServicesSuite(t *testing.T) {
	suite.Run(t, new(ServicesSuite))
}

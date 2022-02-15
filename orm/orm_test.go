package orm_test

import (
	"testing"

	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ORMSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *ORMSuite) SetupTest() {
	suite.db = orm.InitDB(&orm.Config{Debug: true, AutoMigrate: true})
}

func (suite *ORMSuite) TestCanLoadLines() {
	suite.db.Exec(
		"INSERT INTO order_lines (order_id, sku, quantity) VALUES "+
			"(?, ?, ?),(?, ?, ?),(?, ?, ?)",
		"order1", "RED-CHAIR", 12,
		"order1", "RED-TABLE", 13,
		"order2", "BLUE-LIPSTICK", 14,
	)
	expected := []model.OrderLine{
		{OrderID: "order1", SKU: "RED-CHAIR", Quantity: 12},
		{OrderID: "order1", SKU: "RED-TABLE", Quantity: 13},
		{OrderID: "order2", SKU: "BLUE-LIPSTICK", Quantity: 14},
	}
	var actual []model.OrderLine
	suite.db.Table("order_lines").Find(&actual)

	assert.ElementsMatch(suite.T(), expected, actual)
}

func TestORMSuite(t *testing.T) {
	suite.Run(t, new(ORMSuite))
}

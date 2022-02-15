package repository_test

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/orm"
	"github.com/rbusquet/cosmic-go/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RepositorySuite struct {
	suite.Suite
	db  *gorm.DB
	env map[string]string
}

func (suite *RepositorySuite) SetupSuite() {
	godotenv.Load(".config")
}

func (suite *RepositorySuite) SetupTest() {
	db := orm.InitDB(&orm.Config{Debug: true, AutoMigrate: true})
	suite.db = db.Begin()
}

func (suite *RepositorySuite) TearDownTest() {
	suite.db.Rollback()
}

func (suite *RepositorySuite) TestRepositoryCanSaveBatch() {
	zeroTime := time.Time{}
	batch := model.NewBatch("batch1", "RUSTY-SOAPDISH", 100, zeroTime.UTC())
	repo := repository.GormRepository{DB: suite.db}
	repo.Add(batch)

	var reference, sku string
	var purchased_quantity int
	var eta time.Time
	rows, err := suite.db.Table("batches").Select("reference", "sku", "purchased_quantity", "eta").Rows()
	assert.NoError(suite.T(), err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&reference, &sku, &purchased_quantity, &eta)
		expected := []interface{}{"batch1", "RUSTY-SOAPDISH", 100, zeroTime.UTC()}
		actual := []interface{}{reference, sku, purchased_quantity, eta.UTC()}
		assert.Equal(suite.T(), expected, actual)
	}
}

func (suite *RepositorySuite) insertOrderLine() int {
	suite.db.Exec(
		"INSERT INTO order_lines (order_id, sku, quantity) "+
			"VALUES (?, ?, ?)", "order1", "GENERIC-SOFA", 12,
	)
	var id int
	rows, err := suite.db.Table("order_lines").Select("id").Where(
		"order_id=? AND sku=?", "order1", "GENERIC-SOFA",
	).Rows()
	assert.NoError(suite.T(), err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	return id
}

func (suite *RepositorySuite) insertBatch(batchID string) int {
	suite.db.Exec(
		"INSERT INTO batches (reference, sku, purchased_quantity) "+
			"VALUES (?, ?, ?)", batchID, "GENERIC-SOFA", 100,
	)
	var id int
	rows, err := suite.db.Table("batches").Select("id").Where(
		"reference=? AND sku=?", batchID, "GENERIC-SOFA",
	).Rows()
	assert.NoError(suite.T(), err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	return id
}

func (suite *RepositorySuite) insertAllocation(orderLineID int, batchID int) {
	suite.db.Exec(
		"INSERT INTO allocations (order_lines_id, batches_id) "+
			"VALUES (?, ?)", orderLineID, batchID,
	)
}

func (suite *RepositorySuite) TestRepositoryCanRetrieveABatchWithAllocations() {
	orderLineID := suite.insertOrderLine()
	batch1ID := suite.insertBatch("batch1")
	suite.insertBatch("batch2")
	suite.insertAllocation(orderLineID, batch1ID)
	repo := repository.GormRepository{DB: suite.db}

	result := repo.List()
	assert.Len(suite.T(), result, 2)

	retrieved := repo.Get("batch1")
	expected := model.NewBatch("batch1", "GENERIC-SOFA", 100, time.Time{})

	assert.Equal(suite.T(), expected.Reference, retrieved.Reference)
	assert.Equal(suite.T(), expected.SKU, retrieved.SKU)
	assert.Equal(suite.T(), expected.PurchasedQuantity, retrieved.PurchasedQuantity)
	assert.Equal(suite.T(), 12, retrieved.AllocatedQuantity())
}

func (suite *RepositorySuite) TestRepositoryCansaveBatchWithAllocations() {
	suite.insertBatch("batch1")
	suite.insertBatch("batch2")

	repo := repository.GormRepository{DB: suite.db}

	retrieved := repo.Get("batch1")
	orderLine := model.OrderLine{OrderID: "o1", SKU: retrieved.SKU, Quantity: 1}
	retrieved.Allocate(orderLine)

	repo.Save(retrieved)

	assert.Equal(suite.T(), retrieved.AllocatedQuantity(), repo.Get("batch1").AllocatedQuantity())
}
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

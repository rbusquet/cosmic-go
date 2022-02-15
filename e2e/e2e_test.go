package e2e_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/rbusquet/cosmic-go/e2e"
	"github.com/rbusquet/cosmic-go/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type E2ESuite struct {
	suite.Suite
	db           *gorm.DB
	batchesAdded map[uint]bool
	skusAdded    map[interface{}]bool
}

func (suite *E2ESuite) SetupSuite() {
	os.Setenv("DATABASE_HOST", "../allocate.db")
}

func (suite *E2ESuite) TearDownSuite() {
	os.Unsetenv("DATABASE_HOST")
}

func (suite *E2ESuite) SetupTest() {
	suite.db = orm.InitDB(&orm.Config{Debug: true, AutoMigrate: true})
	suite.batchesAdded = make(map[uint]bool)
	suite.skusAdded = make(map[interface{}]bool)
}

func (suite *E2ESuite) addStock(ref, sku, qty, eta interface{}) {
	suite.db.Exec(
		"INSERT INTO batches (reference, sku, purchased_quantity, eta) VALUES "+
			"(?, ?, ?, ?)",
		ref, sku, qty, eta,
	)

	row := suite.db.Raw("SELECT id FROM batches WHERE reference=? AND sku=?", ref, sku).Row()
	var batchId uint
	row.Scan(&batchId)

	suite.batchesAdded[batchId] = true
	suite.skusAdded[sku] = true
}

func (suite *E2ESuite) TearDownTest() {
	for batch := range suite.batchesAdded {
		suite.db.Exec(
			"DELETE FROM allocations WHERE batches_id=?",
			batch,
		)
		suite.db.Exec(
			"DELETE FROM batches WHERE id=?", batch,
		)
	}
	for sku := range suite.skusAdded {
		suite.db.Exec("DELETE FROM order_lines WHERE sku=?", sku)
	}
}

func (suite *E2ESuite) TestApiReturns201AndAllocatedBatch() {
	sku := e2e.RandomSku()
	othersku := e2e.RandomSku("other")
	earlybatch := e2e.RandomBatchref("1")
	laterbatch := e2e.RandomBatchref("2")
	otherbatch := e2e.RandomBatchref("3")

	suite.addStock(laterbatch, sku, "100", "2011-01-02")
	suite.addStock(earlybatch, sku, "100", "2011-01-01")
	suite.addStock(otherbatch, othersku, "100", nil)

	resp, err := http.PostForm("http://localhost:8080/allocate",
		url.Values{"orderid": {e2e.RandomOrderid()}, "sku": {sku}, "qty": {"3"}},
	)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(suite.T(), 201, resp.StatusCode)
	assert.Equal(suite.T(), earlybatch, response["batchref"])
}

func (suite *E2ESuite) TestUnhappy400AndErrorMessage() {
	unknownSku := e2e.RandomSku()
	orderid := e2e.RandomOrderid()
	data := url.Values{"orderid": {orderid}, "sku": {unknownSku}, "qty": {"20"}}

	resp, err := http.PostForm("http://localhost:8080/allocate", data)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 400, resp.StatusCode)
	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(suite.T(), fmt.Sprintf("Invalid SKU %s", unknownSku), response["message"])
}

func TestE2ESuite(t *testing.T) {
	if _, ok := os.LookupEnv("RUN_E2E_TESTS"); !ok {
		t.Skip()
	}
	suite.Run(t, new(E2ESuite))
}
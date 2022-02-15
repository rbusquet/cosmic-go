package e2e_test

import (
	"encoding/json"
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

func (suite *E2ESuite) SetupTest() {
	suite.db = orm.InitDB("../allocate.db", "sqlite", true)
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

func (suite *E2ESuite) TestApiReturnsAllocation() {
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
	assert.Equal(suite.T(), earlybatch, response["batchref"])
}

func (suite *E2ESuite) TestAllocationsArePersisted() {
	sku := e2e.RandomSku()
	batch1 := e2e.RandomBatchref("1")
	batch2 := e2e.RandomBatchref("2")
	order1 := e2e.RandomOrderid("1")
	order2 := e2e.RandomOrderid("2")

	suite.addStock(batch1, sku, 10, "2011-01-01")
	suite.addStock(batch2, sku, 10, "2011-01-02")

	line1 := url.Values{"orderid": {order1}, "sku": {sku}, "qty": {"10"}}
	line2 := url.Values{"orderid": {order2}, "sku": {sku}, "qty": {"10"}}

	resp, err := http.PostForm("http://localhost:8080/allocate", line1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)
	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(suite.T(), batch1, response["batchref"])

	resp, err = http.PostForm("http://localhost:8080/allocate", line2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(suite.T(), batch2, response["batchref"])
}

//     # first order uses up all stock in batch 1
//     r = requests.post(f"{url}/allocate", json=line1)
//     assert r.status_code == 201
//     assert r.json()["batchref"] == batch1

//     # second order should go to batch 2
//     r = requests.post(f"{url}/allocate", json=line2)
//     assert r.status_code == 201
//     assert r.json()["batchref"] == batch2

func TestE2ESuite(t *testing.T) {
	if _, ok := os.LookupEnv("RUN_E2E_TESTS"); !ok {
		t.Skip()
	}
	suite.Run(t, new(E2ESuite))
}

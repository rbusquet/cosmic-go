package allocate

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rbusquet/cosmic-go/model"
	"github.com/rbusquet/cosmic-go/repository"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

type Params struct {
	Orderid string `json:"orderid" form:"orderid"`
	Sku     string `json:"sku" form:"sku"`
	Qty     int    `json:"qty" form:"qty"`
}

type Result struct {
	Batchref string `json:"batchref" form:"batchref"`
}

func (h *Handler) AllocateEndpoint(c echo.Context) error {

	return h.DB.Transaction(func(tx *gorm.DB) error {
		repo := repository.GormRepository{DB: tx}
		batches := repo.List()

		req := new(Params)
		if err := c.Bind(req); err != nil {
			return err
		}

		line := model.OrderLine{OrderID: req.Orderid, SKU: req.Sku, Quantity: req.Qty}
		batchref, err := model.Allocate(line, batches...)
		if err != nil {
			return err
		}
		repo.Save(batches...)

		return c.JSON(http.StatusCreated, Result{Batchref: batchref})
	})
}

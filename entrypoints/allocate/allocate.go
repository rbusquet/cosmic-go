package allocate

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rbusquet/cosmic-go/repository"
	"github.com/rbusquet/cosmic-go/services"
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
	Message  string `json:"message" form:"message"`
	Batchref string `json:"batchref" form:"batchref"`
}

func (h *Handler) AllocateEndpoint(c echo.Context) error {
	return h.DB.Transaction(func(tx *gorm.DB) error {
		repo := repository.GormRepository{DB: tx}
		req := new(Params)
		if err := c.Bind(req); err != nil {
			return err
		}

		batchref, err := services.Allocate(req.Orderid, req.Sku, req.Qty, &repo)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, Result{Batchref: batchref})
	})
}

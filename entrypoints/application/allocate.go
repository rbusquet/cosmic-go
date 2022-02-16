package application

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

type LineSerializer struct {
	Orderid string `json:"orderid" form:"orderid"`
	Sku     string `json:"sku" form:"sku"`
	Qty     int    `json:"qty" form:"qty"`
}

func (h *Handler) AllocateEndpoint(c echo.Context) error {
	tx := c.Get("db")
	if tx == nil {
		tx = h.DB
	}
	repo := repository.GormRepository{DB: tx.(*gorm.DB)}
	req := new(LineSerializer)
	if err := c.Bind(req); err != nil {
		return err
	}

	batchref, err := services.Allocate(req.Orderid, req.Sku, req.Qty, &repo)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"batchref": batchref})
}

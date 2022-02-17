package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rbusquet/cosmic-go/services"
)

type LineSerializer struct {
	Orderid string `json:"orderid" form:"orderid"`
	Sku     string `json:"sku" form:"sku"`
	Qty     int    `json:"qty" form:"qty"`
}

func (h *Handler) AllocateEndpoint(c echo.Context) error {
	repo := h.InitRepo(c)
	data := new(LineSerializer)
	if err := c.Bind(data); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	batchref, err := services.Allocate(data.Orderid, data.Sku, data.Qty, repo)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"batchref": batchref})
}

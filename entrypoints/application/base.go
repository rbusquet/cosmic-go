package application

import (
	"github.com/labstack/echo/v4"
	"github.com/rbusquet/cosmic-go/repository"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) InitRepo(ctx echo.Context) repository.Repository {
	db := ctx.Get("db")
	if db == nil {
		db = h.DB
	}
	return &repository.GormRepository{DB: db.(*gorm.DB)}
}

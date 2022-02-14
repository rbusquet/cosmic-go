package orm

import (
	"github.com/glebarez/sqlite"
	"github.com/rbusquet/cosmic-go/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderLines struct {
	gorm.Model
	model.OrderLine
}

type Batches struct {
	gorm.Model
	model.Batch
	Allocations []OrderLines `gorm:"many2many:allocations;"`
}

var clients = map[string]func(dsn string) gorm.Dialector{
	"sqlite":   sqlite.Open,
	"postgres": postgres.Open,
}

func InitDB(dns string, driver string, debug bool) *gorm.DB {
	var db *gorm.DB
	if client, ok := clients[driver]; ok {
		var err error
		db, err = gorm.Open(client(dns), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}
	if debug {
		db = db.Debug()
	}
	db.AutoMigrate(&OrderLines{}, &Batches{})
	return db
}

package orm

import (
	"os"

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

type Config struct {
	Debug       bool
	AutoMigrate bool
}

func InitDB(config *Config) *gorm.DB {
	dns := ":memory:"
	driver := "sqlite"

	if envDns, ok := os.LookupEnv("DATABASE_HOST"); ok {
		dns = envDns
	}
	if envDriver, ok := os.LookupEnv("DATABASE_DRIVER"); ok {
		driver = envDriver
	}
	var db *gorm.DB
	if client, ok := clients[driver]; ok {
		var err error
		db, err = gorm.Open(client(dns), &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			panic("failed to connect database")
		}
	}
	if config.Debug {
		db = db.Debug()
	}
	if config.AutoMigrate {
		db.AutoMigrate(&OrderLines{}, &Batches{})
	}
	return db
}

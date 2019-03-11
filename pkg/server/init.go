package server

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/taufanmahaputra/forex/pkg/controller"
	"github.com/taufanmahaputra/forex/pkg/lib/config"
	"github.com/taufanmahaputra/forex/pkg/repository"
	"github.com/taufanmahaputra/forex/pkg/service"
	"log"
	"time"
)

var rateRepository repository.RateRepository

var rateService service.RateService

var rateController *controller.RateController

func Init() error {
	cfg := config.GetConfig()

	db := getSQLDB(cfg)

	rateRepository = repository.RateRepository{DB: db}

	rateService = service.InitRateService(rateRepository)

	rateController = controller.InitRateController(rateService)

	return nil
}

func getSQLDB(cfg config.Config) *gorm.DB {
	dbConfig := "host=" + cfg.Database.Host +
		" port=" + cfg.Database.Port +
		" user=" + cfg.Database.User +
		" dbname=" + cfg.Database.Database +
		" sslmode=" + cfg.Database.SSL

	db, err := gorm.Open("postgres", dbConfig)

	for err != nil { // Re-establish db connection every 3 secs if failed
		log.Println(err)
		time.Sleep(3 * time.Second)
		db, err = gorm.Open("postgres", dbConfig)
	}

	return db
}

package server

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/taufanmahaputra/forex/pkg/controller"
	"github.com/taufanmahaputra/forex/pkg/repository"
	"github.com/taufanmahaputra/forex/pkg/service"
	"log"
	"time"
)

var rateRepository repository.RateRepository

var rateController *controller.RateController

var rateService service.RateService

func Init() error {
	db := getSQLDB()

	rateRepository = repository.RateRepository{DB: db}

	rateService = service.InitRateService(rateRepository)

	rateController = controller.InitRateController(rateService)

	return nil
}

func getSQLDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=postgres port=5432 user=postgres dbname=forex sslmode=disable")

	for err != nil { // Re-establish db connection every 3 secs if failed
		log.Println(err)
		time.Sleep(3 * time.Second)
		db, err = gorm.Open("postgres", "host=postgres port=5432 user=postgres dbname=forex sslmode=disable")
	}
	return db
}
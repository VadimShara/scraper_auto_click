package main

import (
	"context"
	"database/sql"
	"fmt"
	"scraper-first/internal/config"
	"scraper-first/internal/service"
	"scraper-first/internal/service/db"
	"scraper-first/pkg/client"
	"scraper-first/pkg/logging"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/tebeka/selenium"
)

func main() {
	const port = 4444

	var wg sync.WaitGroup

	logger := logging.GetLogger()
	logger.Info("Start application")

	edgeDriverPath := "C:/Users/vadim/Downloads/edgedriver_win64/msedgedriver.exe"

	seleniumService, err := selenium.NewChromeDriverService(edgeDriverPath, 9515)
	if err != nil {
		logger.Fatalf("Error starting the Edge: %v", err)
	}
	defer seleniumService.Stop()

	caps := selenium.Capabilities{"browserName": "MicrosoftEdge"}

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		logger.Fatalf("Error connecting to Edge: %v", err)
	}
	defer wd.Quit()

	cfg := config.GetConfig()

	pg, err := store.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	client := postgresql.NewRepository(pg, logger)

	db, err := sql.Open("pgx", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database))

	if err != nil {
		logger.Fatal(err)
	}

	err = goose.Up(db, "C:/Users/vadim/scraper-first/migrations/sql")
	if err != nil {
		logger.Fatal(err)
	}

	wg.Add(1)

	go Parsing(&wg, wd, client, logger)
	
	router := gin.Default()

	handler := service.NewHandler(client, logger)

	handler.Register(router)
	
	router.Run(":8080")

	wg.Wait()
	logger.Info("parsing is complete")
}

func Parsing(wg *sync.WaitGroup, wd selenium.WebDriver, client service.Repository, logger *logging.Logger) {
	defer wg.Done()
	parser := service.NewParser(wd, client, logger)

	if err := parser.AllMarks(); err != nil {
		logger.Fatal(err)
	}

}
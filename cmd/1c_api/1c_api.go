package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/donskova1ex/1cServices/internal/processors"
	"github.com/donskova1ex/1cServices/internal/repositories"

	openapi "github.com/donskova1ex/1cServices/openapi"
)

func main() {

	logJSONHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logJSONHandler)
	slog.SetDefault(logger)

	sqlDSN := os.Getenv("SQL_DSN")
	if sqlDSN == "" {
		logger.Error("empty SQL_DSN")
		os.Exit(1)
	}
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		logger.Error("empty API_PORT")
		os.Exit(1)
	}

	db, err := repositories.NewSQLDB(sqlDSN)
	if err != nil {
		logger.Error(
			"cannot connect to database",
			slog.String("err", err.Error()),
		)
		return
	}

	defer db.Close()
	repository := repositories.NewRepository(db)
	pdnCalcProcessor := processors.NewPDNParametres(repository, logger)
	PDNcalculationAPIService := openapi.NewPDNcalculationAPIService(pdnCalcProcessor, logger)
	PDNcalculationAPIController := openapi.NewPDNcalculationAPIController(PDNcalculationAPIService)

	router := openapi.NewRouter(PDNcalculationAPIController)

	httpServer := http.Server{
		Addr:     ":" + apiPort,
		ErrorLog: slog.NewLogLogger(logJSONHandler, slog.LevelError),
		Handler:  router,
	}
	logger.Info("application started", slog.String("port", apiPort))
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("failed to start server", slog.String("err", err.Error()))
	}

	log.Fatal(http.ListenAndServe(":1616", router))
}

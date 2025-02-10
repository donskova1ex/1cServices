package main

import (
	"context"
	"errors"
	"github.com/donskova1ex/1cServices/internal"
	"github.com/donskova1ex/1cServices/internal/processors"
	"github.com/donskova1ex/1cServices/internal/repositories"
	"log/slog"
	"net/http"
	"os"
	"time"

	openapi "github.com/donskova1ex/1cServices/openapi"
)

func main() {

	logJSONHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logJSONHandler)
	slog.SetDefault(logger)
	sqlDSN := "server=192.168.3.23,1430;user id=1C_user;password=MhO52KbhaC;database=crm_real_data;"
	apiPort := "8080"
	//sqlDSN := os.Getenv("SQL_DSN")
	if sqlDSN == "" {
		logger.Error("empty SQL_DSN")
		os.Exit(1)
	}
	//apiPort := os.Getenv("API_PORT")
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

	repository := repositories.NewRepository(db)

	pdnCalcProcessor := processors.NewPDNParametres(repository, logger)
	PDNcalculationAPIService := openapi.NewPDNcalculationAPIService(pdnCalcProcessor, logger)
	PDNcalculationAPIController := openapi.NewPDNcalculationAPIController(PDNcalculationAPIService)

	rkoByDivisionProcessor := processors.NewDivisionRko(repository, logger)
	RkoByDivisionAPIService := openapi.NewRkoByDivisionAPIService(rkoByDivisionProcessor, logger)
	RkoByDivisionAPIController := openapi.NewRkoByDivisionAPIController(RkoByDivisionAPIService)

	router := openapi.NewRouter(PDNcalculationAPIController, RkoByDivisionAPIController)

	httpServer := http.Server{
		Addr:     ":" + apiPort,
		ErrorLog: slog.NewLogLogger(logJSONHandler, slog.LevelError),
		Handler:  router,
	}
	Closer := internal.NewCloser()

	Closer.Add(func() error {
		logger.Info("closing db connection")
		if err := db.Close(); err != nil {
			logger.Error("error closing db", slog.String("err", err.Error()))
			return err
		}
		logger.Info("db connection closed successfully")
		return nil
	})

	Closer.Add(func() error {
		logger.Info("shutting down HTTP server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error("error shutting down HTTP server", slog.String("err", err.Error()))
			return err
		}
		logger.Info("HTTP server shut down successfully")
		return nil
	})
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()
	go func() {
		ctx, cancel := context.WithCancel(shutdownCtx)
		defer cancel()
		Closer.Run(ctx, logger)
		os.Exit(0)
	}()
	logger.Info("application started", slog.String("port", apiPort))
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start server", slog.String("err", err.Error()))
		return
	}

	//log.Fatal(http.ListenAndServe(":1616", router))
}

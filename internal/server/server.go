package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/adapters"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/adapters/handler"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/adapters/ionCannonClient"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/common/logger"
	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/services"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// App represents the Endor service application.
type App struct {
	logger *zap.SugaredLogger
	svc    *services.EndorService
	srv    adapters.ServerHTTP
}

// Initialize initializes the Endor service application.
func (a *App) Initialize() error {
	err := checkConfig()
	if err != nil {
		return err
	}

	// Configure logger based on environment
	env := os.Getenv("ENV")
	var logLevel string
	switch env {
	case "dev", "test":
		logLevel = logger.DEBUG
	default:
		logLevel = logger.WARN
	}
	a.logger = logger.NewLogger(logLevel, false)

	// TODO: Create instances of IonCannonClient from configuration file
	ionCannon1 := ionCannonClient.NewIonCannonClient(os.Getenv("ION_CANNON_URL1"))
	ionCannon2 := ionCannonClient.NewIonCannonClient(os.Getenv("ION_CANNON_URL2"))
	ionCannon3 := ionCannonClient.NewIonCannonClient(os.Getenv("ION_CANNON_URL3"))
	ionCannons := []adapters.IonCannon{ionCannon1, ionCannon2, ionCannon3}

	a.svc = services.NewEndorService(a.logger, ionCannons)
	validate := validator.New()
	a.srv = handler.NewHTTPServer(a.svc, validate)

	return nil
}

// Run starts the Endor service.
func (a *App) Run() {
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	defer serverStopCtx()

	// trap sigterm or interupt and gracefully shutdown the server
	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)
	go func() {
		// Block until a signal is received.
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				a.logger.Errorln("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := a.srv.Shutdown(shutdownCtx)
		if err != nil {
			a.logger.Errorln(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := a.srv.ListenAndServe("")
	if err != nil && err != http.ErrServerClosed {
		a.logger.Errorln("Error starting server", "error", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

// checkConfig checks if the required configuration is set.
func checkConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file:", err)
	}

	requiredVariables := []struct {
		name string
		val  string
	}{
		{"ENV", os.Getenv("ENV")},
		{"ION_CANNON_URL1", os.Getenv("ION_CANNON_URL1")},
		{"ION_CANNON_URL2", os.Getenv("ION_CANNON_URL2")},
		{"ION_CANNON_URL3", os.Getenv("ION_CANNON_URL3")},
	}

	missingVariables := []string{}
	for _, v := range requiredVariables {
		if v.val == "" {
			missingVariables = append(missingVariables, v.name)
		}
	}

	if len(missingVariables) > 0 {
		return fmt.Errorf("One or more required environment variables are not set: %s", strings.Join(missingVariables, ", "))
	}

	return nil
}

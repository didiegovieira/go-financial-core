package main

import (
	"context"
	"log"
	"time"

	"github.com/didiegovieira/go-financial-core/apps/core/di"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/settings"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

// @title           Financial Core API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1/financial

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	settings.Load()

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	sdk, err := telemetry.Start(startupCtx)
	cancelStartup()
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := sdk.Shutdown(shutdownCtx); err != nil {
			log.Printf("Failed to shutdown telemetry: %v", err)
		}
	}()

	api, cleanup, err := di.InitializeApi()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	api.Start()
}

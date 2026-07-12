package main

import (
	"context"
	"log"

	"github.com/didiegovieira/go-financial-core/apps/core/di"
	"github.com/didiegovieira/go-financial-core/apps/core/internal/settings"
	"github.com/signalfx/splunk-otel-go/distro"
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

	sdk, err := distro.Run()
	if err != nil {
		log.Fatalf("Failed to run distro: %v", err)
	}
	defer func() {
		if err := sdk.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown SDK: %v", err)
		}
	}()

	api, cleanup, err := di.InitializeApi()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	api.Start()
}

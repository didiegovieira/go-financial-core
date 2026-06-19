package settings

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type (
	Specification struct {
		Environment string `envconfig:"ENVIRONMENT" default:"dev"`
		HttpServer  HttpServerSpecification
		Database    DatabaseSpecification
		Kafka       KafkaSpecification
		Metrics     MetricsSpecification
	}

	HttpServerSpecification struct {
		Port         string        `envconfig:"HTTP_SERVER_PORT" default:":8080"`
		ReadTimeout  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"15s"`
		WriteTimeout time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"15s"`
	}

	DatabaseSpecification struct {
		Host     string `envconfig:"DB_HOST" default:"localhost"`
		Port     int    `envconfig:"DB_PORT" default:"5432"`
		User     string `envconfig:"DB_USER" default:"postgres"`
		Password string `envconfig:"DB_PASSWORD" default:"postgres"`
		Name     string `envconfig:"DB_NAME" default:"payments_db"`
	}

	KafkaSpecification struct {
		Brokers []string `envconfig:"KAFKA_BROKERS" default:"kafka:9092"`
	}

	MetricsSpecification struct {
		Name            string `envconfig:"OTEL_SERVICE_NAME" default:"github.com/didiegovieira/go-payments-core"`
		Url             string `envconfig:"OTEL_EXPORTER_JAEGER_ENDPOINT" default:"http://localhost:4317"`
		Token           string `envconfig:"SPLUNK_ACCESS_TOKEN"`
		Resource        string `envconfig:"OTEL_RESOURCE_ATTRIBUTES" default:"service.name=github.com/didiegovieira/go-payments-core"`
		TracesExporter  string `envconfig:"OTEL_TRACES_EXPORTER" default:"jaeger"`
		MetricsExporter string `envconfig:"OTEL_METRICS_EXPORTER" default:""`
	}
)

var Settings Specification

func Load() {
	if os.Getenv("ENVIRONMENT") != "test" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning loading .env file: %v", err)
		}
	}

	Init()
}

func Init() {
	if err := envconfig.Process("", &Settings); err != nil {
		panic(err.Error())
	}
}

func (s *Specification) IsProduction() bool {
	return s.Environment == "prod"
}

func (s *Specification) IsLocal() bool {
	return s.Environment == "local"
}

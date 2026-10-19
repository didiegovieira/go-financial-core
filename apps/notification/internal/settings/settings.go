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
		Telemetry   TelemetrySpecification
	}

	HttpServerSpecification struct {
		Port         string        `envconfig:"HTTP_SERVER_PORT" default:":8083"`
		ReadTimeout  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"15s"`
		WriteTimeout time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"15s"`
	}

	DatabaseSpecification struct {
		Postgres PostgresSpecification
		MySQL    MySQLSpecification
		MongoDB  MongoDBSpecification
	}

	PostgresSpecification struct {
		Host         string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port         int    `envconfig:"POSTGRES_PORT" default:"5432"`
		User         string `envconfig:"POSTGRES_USER" default:"postgres"`
		Password     string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
		DatabaseName string `envconfig:"POSTGRES_DATABASE" default:"notification_db"`
	}

	MySQLSpecification struct {
		Host         string `envconfig:"MYSQL_HOST" default:"localhost"`
		Port         int    `envconfig:"MYSQL_PORT" default:"3306"`
		User         string `envconfig:"MYSQL_USER" default:"root"`
		Password     string `envconfig:"MYSQL_PASSWORD" default:"root"`
		DatabaseName string `envconfig:"MYSQL_DATABASE" default:"notification_db"`
	}

	MongoDBSpecification struct {
		URI          string `envconfig:"MONGODB_URI" default:"mongodb://localhost:27017"`
		DatabaseName string `envconfig:"MONGODB_DATABASE" default:"notification_db"`
	}

	KafkaSpecification struct {
		Brokers         []string `envconfig:"KAFKA_BROKERS" default:"kafka:9092"`
		ConsumerGroupID string   `envconfig:"KAFKA_CONSUMER_GROUP_ID" default:"notification-service"`
	}

	TelemetrySpecification struct {
		Name            string `envconfig:"OTEL_SERVICE_NAME" default:"github.com/didiegovieira/go-financial-core"`
		Url             string `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:"http://localhost:4317"`
		Resource        string `envconfig:"OTEL_RESOURCE_ATTRIBUTES" default:"service.name=github.com/didiegovieira/go-financial-core"`
		TracesExporter  string `envconfig:"OTEL_TRACES_EXPORTER" default:"otlp"`
		MetricsExporter string `envconfig:"OTEL_METRICS_EXPORTER" default:"none"`
	}
)

var Settings Specification

func Load() {
	if os.Getenv("ENVIRONMENT") != "test" {
		if _, statErr := os.Stat(".env"); statErr == nil {
			if err := godotenv.Load(); err != nil {
				log.Printf("Warning loading .env file: %v", err)
			}
		} else if !os.IsNotExist(statErr) {
			log.Printf("Warning loading .env file: %v", statErr)
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

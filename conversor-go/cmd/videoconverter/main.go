package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"imersaofctube/internal/converter"
	"imersaofctube/internal/rabbitmq"
	"github.com/streadway/amqp"
	_ "github.com/lib/pq"

)

func connectPostgres() (*sql.DB, error) {
	user := getEnvOrDefault("POSTGRES_USER", "postgres")
	pass := getEnvOrDefault("POSTGRES_PASSWORD", "root")
	dbname := getEnvOrDefault("POSTGRES_DBNAME", "mydb")
	host := getEnvOrDefault("POSTGRES_HOST", "localhost")
	sslmode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", user, pass, dbname, host, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err!= nil {
    slog.Error("Error opening database", slog.String("connStr", connStr))
		return nil, err
  }
	err = db.Ping()
	if err != nil {
		slog.Error("Error connecting to database", slog.String("connStr", connStr))
    return nil, err
	}
	slog.Info("Database connection established")
	return db, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	db, err := connectPostgres()
	if err!= nil {
    panic(err)
  }
	defer db.Close()

	rabbitMQURL := getEnvOrDefault("RABBITMQ_URL", "amqp://admin:admin@rabbitmq:5672/")
	rabbitClient, err := rabbitmq.NewRabbitClient(rabbitMQURL)
	if err!= nil {
    panic(err)
  }
	defer rabbitClient.Close()

	convertionExch := getEnvOrDefault("CONVERSION_EXCHANGE", "conversion_exchange")
	queueName := getEnvOrDefault("QUEUE_NAME","video_conversion_queue")
	conversionKey := getEnvOrDefault("CONVERSION_KEY","conversion")

	confirmationKey := getEnvOrDefault("CONFIRMATION_KEY","finish-conversion")
	confirmationQueue := getEnvOrDefault("CONFIRMATION_QUEUE", "finish_conversion_queue")
	

	vc := converter.NewVideoConverter(rabbitClient, db)

	msgs, err := rabbitClient.ConsumeMessages(convertionExch, conversionKey, queueName)
	if err!= nil {
    slog.Error("Failed to consume messages", slog.String("error", err.Error()))
  }
	
	for d := range msgs {
		go func(delivery amqp.Delivery) {
			vc.Handle(delivery, convertionExch, confirmationKey, confirmationQueue)
		}(d)
	}
}

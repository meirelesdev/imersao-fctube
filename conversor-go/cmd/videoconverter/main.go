package main

import (
	"context"
	"database/sql"
	"fmt"
	"imersaofctube/internal/converter"
	"imersaofctube/pkg/log"
	"imersaofctube/pkg/rabbitmq"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func connectPostgres() (*sql.DB, error) {
	user := getEnvOrDefault("POSTGRES_USER", "postgres")
	pass := getEnvOrDefault("POSTGRES_PASSWORD", "root")
	dbname := getEnvOrDefault("POSTGRES_DBNAME", "mydb")
	host := getEnvOrDefault("POSTGRES_HOST", "localhost")
	sslmode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", user, pass, dbname, host, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
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
	isDebug := getEnvOrDefault("DEBUG", "false") == "true"
	logger := log.NewLogger(isDebug)
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	db, err := connectPostgres()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQURL := getEnvOrDefault("RABBITMQ_URL", "amqp://admin:admin@rabbitmq:5672/")
	rabbitClient, err := rabbitmq.NewRabbitClient(ctx, rabbitMQURL)
	if err != nil {
		panic(err)
	}
	defer rabbitClient.Close()

	convertionExch := getEnvOrDefault("CONVERSION_EXCHANGE", "conversion_exchange")
	queueName := getEnvOrDefault("QUEUE_NAME", "video_conversion_queue")
	conversionKey := getEnvOrDefault("CONVERSION_KEY", "conversion")

	confirmationKey := getEnvOrDefault("CONFIRMATION_KEY", "finish-conversion")
	confirmationQueue := getEnvOrDefault("CONFIRMATION_QUEUE", "finish_conversion_queue")
	rootPath := getEnvOrDefault("VIDEO_ROOT_PATH", "/media/uploads")
	fmt.Println("Root path:", rootPath)

	vc := converter.NewVideoConverter(rabbitClient, db, rootPath)

	msgs, err := rabbitClient.ConsumeMessages(convertionExch, conversionKey, queueName)
	if err != nil {
		slog.Error("Failed to consume messages", slog.String("error", err.Error()))
	}

	var wg sync.WaitGroup

	go func() {
		for d := range msgs {
			wg.Add(1)
			go func(delivery amqp.Delivery) {
				defer wg.Done()
				vc.Handle(delivery, convertionExch, confirmationKey, confirmationQueue)
			}(d)
		}
	}()

	slog.Info("Waiting for messages from RabbitMQ...")
	<-signalChan
	slog.Info("Shoutdown signal received, finalizing process...")

	cancel()
	wg.Wait()

	slog.Info("Process finalized, exiting...")
}

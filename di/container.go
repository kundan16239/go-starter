package di

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"project-structure/internal/user"
	"project-structure/pkg/db"
	"project-structure/pkg/shared/logger"
	"project-structure/pkg/shared/messaging"
)

type App struct {
	Server   *http.Server
	Logger   logger.Logger
	Port     string
	Shutdown func()
}

func InitializeApp() (*App, func(), error) {
	_ = godotenv.Load()

	log := logger.NewLogger()

	// DB setup
	client, mongoDB, err := db.Connect()
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB", "error", err)
	}

	cleanup := func() {
		_ = client.Disconnect(context.Background())
	}

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())

	// Repositories
	userRepo := user.NewMongoRepository(mongoDB)

	// Services
	userService := user.NewService(userRepo)

	// Handlers
	userHandler := user.NewHandler(userService, log)

	// Router
	userHandler.RegisterRoutes(router)

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive", "time": time.Now().UTC()})
	})

	// Messaging
	rabbitmqURI := os.Getenv("RABBITMQ_URI")
	if rabbitmqURI == "" {
		rabbitmqURI = "amqp://guest:guest@localhost:5672/"
	}

	watermillLogger := watermill.NewStdLogger(false, false)
	amqpCfg := amqp.NewDurableQueueConfig(rabbitmqURI)

	publisher, err := amqp.NewPublisher(amqpCfg, watermillLogger)
	if err != nil {
		return nil, cleanup, err
	}
	subscriber, err := amqp.NewSubscriber(amqpCfg, watermillLogger)
	if err != nil {
		return nil, cleanup, err
	}

	// DLQ reprocessing route
	router.POST("/dlq/reprocess", func(c *gin.Context) {
		queue := c.Query("queue")
		if queue == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "queue param required"})
			return
		}
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0")) // 0 = unlimited
		dlq := "dlq_" + queue
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		msgs, err := subscriber.Subscribe(ctx, dlq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		reprocessed := 0
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					goto done
				}
				newMsg := message.NewMessage(watermill.NewUUID(), msg.Payload)
				newMsg.Metadata.Set("x-retry-count", "0")
				if err := publisher.Publish(queue, newMsg); err != nil {
					log.Error("republish failed", "queue", queue, "error", err)
					msg.Nack()
					continue
				}
				msg.Ack()
				reprocessed++
				if limit > 0 && reprocessed >= limit {
					goto done
				}
			case <-ctx.Done():
				goto done
			}
		}
	done:
		c.JSON(http.StatusOK, gin.H{"dlq": dlq, "reprocessed": reprocessed})
	})

	// Start all consumers
	queues := []string{"test_queue", "test1_queue", "test2_queue"}
	processor := messaging.NewEventProcessor(userService)
	for _, q := range queues {
		go messaging.ConsumeQueue(subscriber, publisher, q, processor)
	}

	// HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		Server:   server,
		Logger:   log,
		Port:     port,
		Shutdown: cleanup,
	}, cleanup, nil
}

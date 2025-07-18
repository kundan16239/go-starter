package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-structure/di"
)

func main() {
	// Initialize container
	app, cleanup, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	// Start server
	go func() {
		app.Logger.Info("Starting HTTP server", "port", app.Port)
		if err := app.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			app.Logger.Fatal("HTTP server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Server.Shutdown(ctx); err != nil {
		app.Logger.Error("Forced shutdown error", "error", err)
	}

	app.Logger.Info("Server exited")
}

// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"syscall"
// 	"time"

// 	"github.com/ThreeDotsLabs/watermill"
// 	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
// 	"github.com/ThreeDotsLabs/watermill/message"
// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"

// 	"project-structure/internal/user"
// 	"project-structure/pkg/db"
// 	"project-structure/pkg/shared/logger"
// 	"project-structure/pkg/shared/messaging"
// )

// var (
// 	rabbitmqURI = "amqp://guest:guest@152.53.230.202:5672/"
// 	queues      = []string{"test_queue", "test1_queue", "test2_queue"}
// )

// func main() {
// 	// Load environment variables from .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using system environment variables")
// 	}

// 	// Initialize logger
// 	logger := logger.NewLogger()

// 	// Queue Initialize
// 	watermillLogger := watermill.NewStdLogger(false, false)
// 	amqpCfg := amqp.NewDurableQueueConfig(rabbitmqURI)

// 	publisher, err := amqp.NewPublisher(amqpCfg, watermillLogger)
// 	if err != nil {
// 		log.Fatalf("❌ publisher: %v", err)
// 	}
// 	subscriber, err := amqp.NewSubscriber(amqpCfg, watermillLogger)
// 	if err != nil {
// 		log.Fatalf("❌ subscriber: %v", err)
// 	}

// 	// Connect to MongoDB
// 	client, db, err := db.Connect()
// 	if err != nil {
// 		logger.Fatal("Failed to connect to MongoDB", "error", err)
// 	}
// 	defer func() {
// 		if err := client.Disconnect(context.Background()); err != nil {
// 			logger.Error("Failed to disconnect from MongoDB", "error", err)
// 		}
// 	}()

// 	// Initialize repositories
// 	userRepo := user.NewMongoRepository(db)

// 	// Initialize services
// 	userService := user.NewService(userRepo)

// 	// Initialize handlers
// 	userHandler := user.NewHandler(userService, logger)

// 	// Setup Gin router
// 	router := gin.Default()

// 	// Add middleware
// 	router.Use(gin.Logger())
// 	router.Use(gin.Recovery())

// 	// Register routes
// 	userHandler.RegisterRoutes(router)

// 	router.GET("/healthz", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":    "alive",
// 			"timestamp": time.Now().UTC(),
// 		})
// 	})

// 	router.POST("/dlq/reprocess", func(c *gin.Context) {
// 		queue := c.Query("queue")
// 		if queue == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "queue param required"})
// 			return
// 		}
// 		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0")) // 0 = unlimited

// 		dlq := "dlq_" + queue
// 		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 		defer cancel()

// 		msgs, err := subscriber.Subscribe(ctx, dlq)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		reprocessed := 0
// 		for {
// 			select {
// 			case msg, ok := <-msgs:
// 				if !ok {
// 					goto done // channel closed
// 				}
// 				// republish to original queue with retry counter reset
// 				newMsg := message.NewMessage(watermill.NewUUID(), msg.Payload)
// 				newMsg.Metadata.Set("x-retry-count", "0")

// 				if err := publisher.Publish(queue, newMsg); err != nil {
// 					log.Printf("❌ republish to %s failed: %v", queue, err)
// 					msg.Nack()
// 					continue
// 				}
// 				msg.Ack()
// 				reprocessed++
// 				if limit > 0 && reprocessed >= limit {
// 					goto done
// 				}
// 			case <-ctx.Done():
// 				goto done
// 			}
// 		}
// 	done:
// 		c.JSON(http.StatusOK, gin.H{"dlq": dlq, "reprocessed": reprocessed})
// 	})

// 	// start a consumer goroutine for each queue
// 	eventProcessor := messaging.NewEventProcessor(userService)
// 	for _, q := range queues {
// 		go messaging.ConsumeQueue(subscriber, publisher, q, eventProcessor)
// 	}

// 	// Get port from environment variable
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	server := &http.Server{
// 		Addr:         ":" + port,
// 		Handler:      router,
// 		ReadTimeout:  15 * time.Second,
// 		WriteTimeout: 15 * time.Second,
// 		IdleTimeout:  60 * time.Second,
// 	}

// 	go func() {
// 		logger.Info("Starting server on", "port", port)
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Fatal("Server failed to start", "error", err)
// 		}
// 	}()

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	logger.Info("Shutting down server...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	if err := server.Shutdown(ctx); err != nil {
// 		logger.Error("Server forced to shutdown", "error", err)
// 	}

// 	logger.Info("Server exited")
// }

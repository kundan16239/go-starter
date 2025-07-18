package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/time/rate"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type QueueConfig struct {
	MaxRetries int           // Maximum number of retry attempts
	RetryDelay time.Duration // Base delay between retries
	RateLimit  time.Duration // Minimum interval between processing messages (0 for no limit)
	Burst      int           // Maximum burst capacity (for token bucket)
}

var queueConfigs = map[string]QueueConfig{
	"user_queue": {
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
		RateLimit:  0, // No rate limiting
		Burst:      1,
	},
	"order_queue": {
		MaxRetries: 5,
		RetryDelay: 2 * time.Second,
		RateLimit:  100 * time.Millisecond, // Max 10 messages/second
		Burst:      5,
	},
	"payment_queue": {
		MaxRetries: 7,
		RetryDelay: 5 * time.Second,
		RateLimit:  1 * time.Second, // Max 1 message/second
		Burst:      1,
	},
}

func getQueueConfig(queue string) QueueConfig {
	if config, exists := queueConfigs[queue]; exists {
		return config
	}
	// Default configuration
	return QueueConfig{
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
		RateLimit:  0,
		Burst:      1,
	}
}

type RateLimiter struct {
	limiter *rate.Limiter
	config  QueueConfig
}

func NewRateLimiter(config QueueConfig) *RateLimiter {
	if config.RateLimit <= 0 {
		return nil
	}
	messagesPerSecond := float64(time.Second) / float64(config.RateLimit)
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(messagesPerSecond), config.Burst),
		config:  config,
	}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	if rl == nil {
		return nil
	}
	return rl.limiter.Wait(ctx)
}

func ConsumeQueue(sub message.Subscriber, pub message.Publisher, queue string, eventProcessor *EventProcessor) {
	msgs, err := sub.Subscribe(context.Background(), queue)
	if err != nil {
		log.Fatalf("❌ subscribe %s: %v", queue, err)
	}

	queueConfig := getQueueConfig(queue)
	rateLimiter := NewRateLimiter(queueConfig)

	for msg := range msgs {
		if err := rateLimiter.Wait(context.Background()); err != nil {
			log.Printf("⚠️ [%s] rate limiter error: %v", queue, err)
			msg.Nack() // Couldn't process due to rate limiting
			continue
		}

		retry := currentRetry(msg)
		log.Printf("📥 [%s] %s | retry %d", queue, msg.UUID, retry)

		var evt Event
		if err := json.Unmarshal(msg.Payload, &evt); err != nil {
			log.Printf("❌ [%s] bad JSON: %v", queue, err)
			msg.Ack()
			continue
		}

		if err := dispatch(queue, evt, eventProcessor); err != nil {
			log.Printf("❌ [%s] dispatch error: %v", queue, err)

			if retry >= queueConfig.MaxRetries {
				dlq := "dlq_" + queue
				dlqMsg := message.NewMessage(watermill.NewUUID(), msg.Payload)
				dlqMsg.Metadata.Set("x-retry-count", fmt.Sprintf("%d", retry))
				dlqMsg.Metadata.Set("x-original-queue", queue)
				_ = pub.Publish(dlq, dlqMsg)
				msg.Ack()
				continue
			}

			time.Sleep(queueConfig.RetryDelay)

			retryMsg := message.NewMessage(watermill.NewUUID(), msg.Payload)
			retryMsg.Metadata.Set("x-retry-count", fmt.Sprintf("%d", retry+1))
			_ = pub.Publish(queue, retryMsg)

			msg.Ack()
			continue
		}

		log.Printf("✅ [%s] %s processed type=%s", queue, msg.UUID, evt.Type)
		msg.Ack()
	}
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────
func currentRetry(msg *message.Message) int {
	var r int
	fmt.Sscanf(msg.Metadata.Get("x-retry-count"), "%d", &r)
	return r
}

// dispatch routes by queue THEN by type
func dispatch(queue string, evt Event, eventProcessor *EventProcessor) error {
	switch queue {
	case "test_queue":
		return handleTestQueue(evt, eventProcessor)
	case "test1_queue":
		return handleTest1Queue(evt)
	case "test2_queue":
		return handleTest2Queue(evt)
	default:
		return fmt.Errorf("unsupported queue %s", queue)
	}
}

// ─────────────────────────────────────────────────────────────
// Handlers for each queue
// ─────────────────────────────────────────────────────────────
func handleTestQueue(evt Event, eventProcessor *EventProcessor) error {
	switch evt.Type {
	case "say_hello":
		var p struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(evt.Data, &p); err != nil {
			return err
		}
		log.Printf("👋 Hello %s (test_queue)", p.Name)
		eventProcessor.UserService.DeleteUser(context.Background(), "1")
		return nil
	default:
		return fmt.Errorf("test_queue: unsupported type %s", evt.Type)
	}
}

func handleTest1Queue(evt Event) error {
	switch evt.Type {
	case "print_number":
		var p struct {
			Number int `json:"number"`
		}
		if err := json.Unmarshal(evt.Data, &p); err != nil {
			return err
		}
		log.Printf("🔢 Number %d (test1_queue)", p.Number)
		return nil
	default:
		return fmt.Errorf("test1_queue: unsupported type %s", evt.Type)
	}
}

func handleTest2Queue(evt Event) error {
	switch evt.Type {
	case "log_message":
		var p struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(evt.Data, &p); err != nil {
			return err
		}
		log.Printf("📝 Message: %s (test2_queue)", p.Message)
		return nil
	default:
		return fmt.Errorf("test2_queue: unsupported type %s", evt.Type)
	}
}

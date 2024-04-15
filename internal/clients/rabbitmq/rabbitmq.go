package rabbitmq

import (
	"encoding/json"
	"github.com/fishmanDK/avito_test_task/internal/service"
	"github.com/wagslane/go-rabbitmq"
	"log"
	"time"
)

const (
	subject           = "deleteRequest"
	checkTime         = 1 * time.Second
	streamName        = "delete-request"
	streamDescription = "requests for delete banners"
)

type RabbitMQConsumer struct {
	deleteBannerService *service.DeleteService
	consumer            *rabbitmq.Consumer
}

type Message struct {
	BannerID  int64 `json:"bannerID,omitempty"`
	TagID     int64 `json:"tagID,omitempty"`
	FeatureID int64 `json:"featureID,omitempty"`
}

func NewRabbitMQConsumer(deleteBannerService *service.DeleteService) (*RabbitMQConsumer, error) {
	const op = "rabbitmq.NewRabbitMQPublisher"

	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := rabbitmq.NewConsumer(
		conn,
		"my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMQConsumer{
		consumer:            consumer,
		deleteBannerService: deleteBannerService,
	}, nil
}

func (rq *RabbitMQConsumer) SubscribeAndReadMessage() error {
	const op = "rabbitmq.SubscribeAndReadMessage"

	err := rq.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("consumed: %v", string(d.Body))
		var msg Message
		err := json.Unmarshal(d.Body, &msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return rabbitmq.NackRequeue
		}
		if msg.BannerID != 0 {
			err = rq.deleteBannerService.DeleteBanner(msg.BannerID)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return rabbitmq.NackRequeue
			}
		} else {
			err = rq.deleteBannerService.DeleteBannerByParams(msg.TagID, msg.FeatureID)
			if err != nil {
				log.Printf("Error partial delete: %v", err)
				return rabbitmq.NackRequeue
			}
		}

		return rabbitmq.Ack
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

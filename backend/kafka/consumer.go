package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"github.com/VK-Container-Pinger/backend/models"
	"github.com/VK-Container-Pinger/backend/repository"
	"github.com/VK-Container-Pinger/backend/cache"
)

type Consumer struct {
	Repo *repository.PingRepository
	RedisClient *redis.Client
}

type ConsumerGroupHandler struct {
	Consumer *Consumer
}

//реализация интерфейса sarama.ConsumerGroupHandler
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Получено сообщение: topic=%s, partition=%d, offset=%d", message.Topic, message.Partition, message.Offset)
		
		var pr models.PingResult
		if err := json.Unmarshal(message.Value, &pr); err != nil {
			log.Printf("Ошибка десериализации сообщения: %s", err)
			continue
		}
		
		if err := h.Consumer.Repo.InsertPingResult(pr); err != nil {
			log.Printf("Ошибка записи в базу данных: %s", err)
			continue
		}else{
			if err := cache.DeleteCache(h.Consumer.RedisClient, "pings_cache"); err != nil {
				log.Printf("Ошибка удаления кэша: %s", err)
			}
		}
		
		session.MarkMessage(message, "")
	}
	return nil
}

func StartKafkaConsumer(broker string, groupID string, topics []string, consumer *Consumer) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Return.Errors = true
	
	group, err := sarama.NewConsumerGroup([]string{broker}, groupID, config)
	if err != nil {
		log.Fatalf("Ошибка создания группы потребителей: %s", err)
	}
	defer group.Close()
	
	handler := &ConsumerGroupHandler{Consumer: consumer}
	
	ctx := context.Background()
	for {
		if err := group.Consume(ctx, topics, handler); err != nil {
			log.Fatalf("Ошибка потребления: %s", err)
			time.Sleep(3 * time.Second)
		}
	}
}
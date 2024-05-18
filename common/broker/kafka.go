package broker

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"log"
)

type KafkaClient struct {
	brokers []string
	config  *sarama.Config
}

func NewKafkaClient(brokers []string) *KafkaClient {
	// kafka configuration
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return &KafkaClient{
		brokers: brokers,
		config:  config,
	}
}

func (k *KafkaClient) Produce(value interface{}, key, topic EventType) (int32, int64, error) {
	producer, err := sarama.NewSyncProducer(k.brokers, k.config)

	if err != nil {
		zap.L().Error("Error creating the Sarama producer", zap.Error(err))
	}

	defer producer.Close()

	// data marshalling
	o, err := json.Marshal(value)

	if err != nil {
		zap.L().Error("Error marshaling value", zap.Error(err))
		return 0, 0, err
	}

	message := &sarama.ProducerMessage{
		Topic: string(topic),
		Key:   sarama.ByteEncoder([]byte(key)),
		Value: sarama.ByteEncoder(o),
	}

	return producer.SendMessage(message)
}

func (k *KafkaClient) Consume(topic EventType) {
	consumer, err := sarama.NewConsumer(k.brokers, k.config)

	if err != nil {
		zap.L().Error("Error creating the Sarama consumer", zap.Error(err))
	}

	// Consume messages
	partitionConsumer, err := consumer.ConsumePartition(string(topic), 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln("Failed to start Sarama partition consumer:", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		log.Printf("Message received: key=%s, value=%s, offset=%d\n", string(message.Key), string(message.Value), message.Offset)
	}
}

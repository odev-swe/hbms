package main

import (
	"context"
	"encoding/json"
	pb "github.com/odev-swe/hbms/common/api"
	"github.com/odev-swe/hbms/common/broker"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	services BookingServices
}

func NewKafkaConsumer(services BookingServices) *KafkaConsumer {
	return &KafkaConsumer{services: services}
}

func (c *KafkaConsumer) Listen(client *broker.KafkaClient) {
	consumerPartition := client.Consume(broker.BookingTopic)

	defer consumerPartition.Close()

	for msg := range consumerPartition.Messages() {

		if string(msg.Key) == string(broker.BookingPaid) {
			var booking *pb.Booking

			err := json.Unmarshal(msg.Value, &booking)

			if err != nil {
				break
			}

			b, err := c.services.UpdateBooking(context.Background(), booking)

			if err != nil {
				zap.L().Error("Error update booking", zap.Error(err), zap.Any("booking", booking))
			}

			zap.L().Info("Updated booking", zap.Any("booking", b))

		}
	}
}

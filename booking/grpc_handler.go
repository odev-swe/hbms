package main

import (
	"context"
	pb "github.com/odev-swe/hbms/common/api"
	"github.com/odev-swe/hbms/common/broker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCHandler struct {
	pb.UnimplementedBookingServicesServer

	*bookingService
	*broker.KafkaClient
}

func NewGRPCHandler(svr *grpc.Server, bookingService *bookingService, kafkaClient *broker.KafkaClient) {

	h := &GRPCHandler{
		bookingService: bookingService,
		KafkaClient:    kafkaClient,
	}

	pb.RegisterBookingServicesServer(svr, h)
}

func (g *GRPCHandler) CreateBooking(ctx context.Context, booking *pb.BookingRequest) (*pb.Booking, error) {

	data, err := g.bookingService.CreateBooking(ctx, booking)
	if err != nil {
		return nil, err
	}

	// send booking.created message
	partition, offset, err := g.KafkaClient.Produce(data, broker.BookingCreated, broker.BookingTopic)
	if err != nil {
		return nil, err
	}

	zap.L().Info("Message booking.created published successfully", zap.Int32("partition", partition), zap.Int64("offset", offset))

	return data, nil
}

func (g *GRPCHandler) UpdateBooking(ctx context.Context, booking *pb.Booking) (*pb.Booking, error) {

	return g.bookingService.UpdateBooking(ctx, booking)
}

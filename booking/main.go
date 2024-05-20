package main

import (
	"context"
	"flag"
	"github.com/odev-swe/hbms/common/broker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

func main() {
	// flag section
	uri := flag.String("uri", "mongodb://root:example@localhost:27017", "mongodb uri")
	brokerURI := flag.String("broker", "localhost:9092", "broker uri")

	// logger section
	logger, _ := zap.NewProduction()

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	// db section
	client, err := connectMongoDB(*uri)

	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	// kafka section
	var brokers = []string{*brokerURI}

	kc := broker.NewKafkaClient(brokers)

	// grpc section
	conn, err := net.Listen("tcp", ":4000")

	if err != nil {
		logger.Error("failed to listen", zap.Error(err))
	}

	grpc := grpc.NewServer()
	store := NewStore(client)
	services := NewBookingService(store)

	NewGRPCHandler(grpc, services, kc)

	// kafka consumer sesction
	consumer := NewKafkaConsumer(services)
	go consumer.Listen(kc)

	logger.Info("Starting server", zap.String("port", "4000"))
	err = grpc.Serve(conn)

	if err != nil {
		logger.Error("failed to serve", zap.Error(err))
	}

}

func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	return client, nil

}

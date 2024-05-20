package main

import (
	"context"
	pb "github.com/odev-swe/hbms/common/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	db *mongo.Client
}

const (
	ColName = "bookings"
	DbName  = "bookings"
)

func NewStore(db *mongo.Client) *store {
	return &store{
		db: db,
	}
}

func (s *store) Create(ctx context.Context, booking Booking) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(ColName)

	res, err := col.InsertOne(ctx, booking)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (s *store) Update(ctx context.Context, id string, booking *pb.Booking) error {
	col := s.db.Database(DbName).Collection(ColName)

	bID, _ := primitive.ObjectIDFromHex(id)

	_, err := col.UpdateOne(ctx, bson.M{"_id": bID}, bson.M{"$set": bson.M{"isPaid": booking.IsPaid, "paymentLink": booking.PaymentLink}})

	if err != nil {
		return err
	}
	return nil
}

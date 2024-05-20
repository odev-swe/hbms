package main

import (
	"context"
	pb "github.com/odev-swe/hbms/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type BookingServices interface {
	CreateBooking(ctx context.Context, pb *pb.BookingRequest) (*pb.Booking, error)
	UpdateBooking(ctx context.Context, pb *pb.Booking) (*pb.Booking, error)
}

type BookingStore interface {
	Create(ctx context.Context, b Booking) (primitive.ObjectID, error)
	Update(ctx context.Context, id string, b *pb.Booking) error
}

type Booking struct {
	Id           string    `bson:"id,omitempty"`
	GuestName    string    `bson:"guestName,omitempty"`
	GuestEmail   string    `bson:"guestEmail,omitempty"`
	CheckInDate  time.Time `bson:"checkInDate,omitempty"`
	CheckOutDate time.Time `bson:"checkOutDate,omitempty"`
	RoomId       string    `bson:"roomId,omitempty"`
	IsCheckIn    bool      `bson:"isCheckIn"`
	IsCheckOut   bool      `bson:"isCheckOut"`
	IsPaid       bool      `bson:"isPaid"`
	PaymentLink  string    `bson:"paymentLink"`
}

func (b Booking) toProto() *pb.Booking {
	return &pb.Booking{
		ID:           b.Id,
		GuestName:    b.GuestName,
		GuestEmail:   b.GuestEmail,
		CheckInDate:  timestamppb.New(b.CheckInDate),
		CheckOutDate: timestamppb.New(b.CheckOutDate),
		RoomID:       b.RoomId,
		IsCheckIn:    b.IsCheckIn,
		IsCheckOut:   b.IsCheckOut,
		PaymentLink:  b.PaymentLink,
	}
}

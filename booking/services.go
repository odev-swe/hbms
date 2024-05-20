package main

import (
	"context"
	pb "github.com/odev-swe/hbms/common/api"
)

type bookingService struct {
	db *store
}

func NewBookingService(store *store) *bookingService {
	return &bookingService{
		db: store,
	}
}

func (b *bookingService) CreateBooking(ctx context.Context, p *pb.BookingRequest) (*pb.Booking, error) {

	data := Booking{
		GuestName:    p.GuestName,
		GuestEmail:   p.GuestEmail,
		CheckInDate:  p.CheckInDate.AsTime(),
		CheckOutDate: p.CheckOutDate.AsTime(),
		RoomId:       p.RoomID,
		IsCheckOut:   false,
		IsCheckIn:    false,
		IsPaid:       false,
		PaymentLink:  "",
	}

	id, err := b.db.Create(ctx, data)

	if err != nil {
		return nil, err
	}

	data.Id = id.Hex()

	return data.toProto(), nil
}

func (b *bookingService) UpdateBooking(ctx context.Context, pb *pb.Booking) (*pb.Booking, error) {
	err := b.db.Update(ctx, pb.ID, pb)

	if err != nil {
		return nil, err
	}

	return pb, err
}

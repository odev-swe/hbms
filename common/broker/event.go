package broker

type EventType string

const (
	BookingTopic   EventType = "booking"
	BookingCreated EventType = "booking.created"
	BookingPaid    EventType = "booking.paid"
)

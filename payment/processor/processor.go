package processor

type PaymentProcessor interface {
	CreatePaymentLink() (string, error)
}

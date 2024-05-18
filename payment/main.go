package main

import (
	"flag"
	"github.com/odev-swe/hbms/common/broker"
	"github.com/odev-swe/hbms/payment/processor/stripe"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	// flag section
	APIKey := flag.String("stripe-api-key", "sk_test_51Owb5kBXNiy8ssD3neOYKRiQwfu6XPpB7nwk1lfau9g2Mocoxytw28xF0Is0o7yeRYUgtTH15VOxpuy2N4ebCzMB00tXNLJIs8", "Stripe-API key")
	domain := flag.String("domain", "http://localhost:2000", "Domain")
	brokers := flag.String("brokers", "localhost:9092", "Brokers")

	// logger section
	logger, _ := zap.NewProduction()

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	// stripe section
	s := stripe.NewProcessor(*APIKey, *domain)

	url, err := s.CreatePaymentLink()
	if err != nil {
		logger.Error("Failed to create payment link", zap.Error(err))
	}

	logger.Info("Payment link created", zap.String("url", url))

	// http section
	mux := http.NewServeMux()

	var bkrs []string

	bkrs = append(bkrs, *brokers)

	kc := broker.NewKafkaClient(bkrs)

	h := NewPaymentHTTPHandler(kc)

	h.registerRoutes(mux)

	err = http.ListenAndServe(":2000", mux)

	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

}

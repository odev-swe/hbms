package main

import (
	"encoding/json"
	"fmt"
	"github.com/odev-swe/hbms/common/broker"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

type PaymentHTTPHandler struct {
	kc *broker.KafkaClient
}

func NewPaymentHTTPHandler(kc *broker.KafkaClient) *PaymentHTTPHandler {
	return &PaymentHTTPHandler{
		kc: kc,
	}
}

func (h *PaymentHTTPHandler) registerRoutes(router *http.ServeMux) {
	fh := http.FileServer(http.Dir("./static"))

	router.Handle("/", fh)
	router.HandleFunc("GET /health-check", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
	})
	router.HandleFunc("POST /webhook", h.handlerWebhook)
}

func (h *PaymentHTTPHandler) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := "whsec_2194bb61daecbda5ebcc7e6856e6da76112f793ad6d0d0e5130456298c531375"
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"),
		endpointSecret, webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		partition, offset, err := h.kc.Produce(session, broker.BookingPaid, broker.BookingTopic)

		if err != nil {
			zap.L().Error("Error getting partition", zap.Error(err))
		}

		zap.L().Info("Message published successfully", zap.String("event", string(broker.BookingPaid)), zap.Int32("partition", partition), zap.Int64("offset", offset))
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

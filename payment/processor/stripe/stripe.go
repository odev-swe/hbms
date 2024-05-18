package stripe

import (
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type Stripe struct {
	Domain string
}

func NewProcessor(APIKey, Domain string) *Stripe {
	stripe.Key = APIKey

	return &Stripe{
		Domain: Domain,
	}
}

func (s Stripe) CreatePaymentLink() (string, error) {
	params := &stripe.CheckoutSessionParams{
		Metadata: map[string]string{
			"asd": "asd",
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Provide the exact Price ID (for example, pr_1234) of the product you want to sell
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyMYR)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Room A"),
					},
					UnitAmount: stripe.Int64(10000),
				},
				Quantity: stripe.Int64(1),
			},
		},
		PaymentMethodTypes: stripe.StringSlice([]string{"fpx", "card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(s.Domain + "/success.html"),
		CancelURL:          stripe.String(s.Domain + "/cancel.html"),
	}

	res, err := session.New(params)
	if err != nil {
		return "", err
	}

	return res.URL, nil
}

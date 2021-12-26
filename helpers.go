package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v72"
)

func errorResponseSetter(err error, c chargeResponse) chargeResponse {
	if stripeErr, ok := err.(*stripe.Error); ok {
		switch stripeErr.Code {
		case stripe.ErrorCodeCardDeclined:
			c.ErrorMessage = "Card declined"
		case stripe.ErrorCodeExpiredCard:
			c.ErrorMessage = "Card expired"
		case stripe.ErrorCodeIncorrectCVC:
			c.ErrorMessage = "Incorrect CVC"
		case stripe.ErrorCodeInvalidChargeAmount:
			c.ErrorMessage = "Invalid Charge amount"
		case stripe.ErrorCodeChargeAlreadyRefunded:
			c.ErrorMessage = "Charge already refunded"
		case stripe.ErrorCodeChargeAlreadyCaptured:
			c.ErrorMessage = "Charge already captured"
		default:
			c.ErrorMessage = "Failed to create charge"
		}
	} else {
		c.ErrorMessage = "Other error occurred"
	}
	return c
}

func invalidIdChecker(w http.ResponseWriter, Id string, formId string, c chargeResponse) {
	c.ErrorMessage = "Charge Id mismatch"
	if Id == "" {
		c.ErrorMessage = "No charge Id provided"
	}
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		log.Panicf("JSON Encoding error: %v", err)
	}
}

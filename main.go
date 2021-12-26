package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_KEY")
}

func main() {

	r := mux.NewRouter()
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logMiddleware := NewLogMiddleware(logger)
	r.Use(logMiddleware.Func())
	r.Use(jsonResponseMiddleware)

	r.HandleFunc("/api/v1/get_charges", ListChargesHandler).
		Methods("GET")

	r.HandleFunc("/api/v1/create_charge", CreateChargeHandler).
		Methods("POST")

	r.HandleFunc(`/api/v1/capture_charge/{chargeId:[\w]+}`, CaptureChargeHandler).
		Methods("POST")

	r.HandleFunc(`/api/v1/create_refund/{chargeId:[\w]+}`, CreateRefundHandler).
		Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

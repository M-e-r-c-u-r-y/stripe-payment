package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
	"github.com/stripe/stripe-go/v72/refund"
)

type chargeResponse struct {
	Id             string `json:"id,omitempty"`
	Success        bool   `json:"success"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	SuccessMessage string `json:"successMessage,omitempty"`
}

type listchargeResponse struct {
	Ids []string `json:"ids"`
}

func ListChargesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	listchargeresponse := listchargeResponse{}
	params := &stripe.ChargeListParams{}
	i := charge.List(params)
	for i.Next() {
		c := i.Charge()
		listchargeresponse.Ids = append(listchargeresponse.Ids, c.ID)
	}
	jsonResponse, _ := json.Marshal(listchargeresponse)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func CaptureChargeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	Id := r.FormValue("id")
	if Id != vars["chargeId"] {
		chargeresponse := chargeResponse{Success: false}
		chargeresponse.ErrorMessage = "Charge Id mismatch"
		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}
	if Id == "" {
		chargeresponse := chargeResponse{Success: false}
		chargeresponse.ErrorMessage = "No charge Id provided"
		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}
	c, err := charge.Capture(
		Id,
		nil,
	)
	if err != nil {
		chargeresponse := chargeResponse{Success: false}
		if stripeErr, ok := err.(*stripe.Error); ok {
			switch stripeErr.Code {
			case stripe.ErrorCodeCardDeclined:
				chargeresponse.ErrorMessage = "Card declined"
			case stripe.ErrorCodeExpiredCard:
				chargeresponse.ErrorMessage = "Card expired"
			case stripe.ErrorCodeIncorrectCVC:
				chargeresponse.ErrorMessage = "Incorrect CVC"
			case stripe.ErrorCodeInvalidChargeAmount:
				chargeresponse.ErrorMessage = "Invalid Charge amount"
			default:
				chargeresponse.ErrorMessage = "Failed to create charge"
			}
		} else {
			chargeresponse.ErrorMessage = "Other error occurred"
		}

		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}
	chargeresponse := chargeResponse{Success: true}
	chargeresponse.Id = c.ID
	chargeresponse.SuccessMessage = fmt.Sprintf("Charged amount %v successfully, transaction id is: %v", c.Amount, c.ID)
	jsonResponse, _ := json.Marshal(chargeresponse)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func CreateRefundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	Id := r.FormValue("id")
	if Id != vars["chargeId"] {
		chargeresponse := chargeResponse{Success: false}
		chargeresponse.ErrorMessage = "Charge Id mismatch"
		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}
	if Id == "" {
		chargeresponse := chargeResponse{Success: false}
		chargeresponse.ErrorMessage = "No charge Id provided"
		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}
	params := &stripe.RefundParams{
		Charge: stripe.String(Id),
	}
	refud_data, err := refund.New(params)
	if err != nil {
		chargeresponse := chargeResponse{Success: false}
		if _, ok := err.(*stripe.Error); ok {
			chargeresponse.ErrorMessage = "Charge already refunded or invalid transaction id"
		} else {
			chargeresponse.ErrorMessage = "Other error occurred"
		}

		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}
	chargeresponse := chargeResponse{Success: true}
	chargeresponse.Id = refud_data.ID
	chargeresponse.SuccessMessage = fmt.Sprintf("Charged amount %v successfully refunded, refund id is: %v", refud_data.Amount, refud_data.ID)
	jsonResponse, _ := json.Marshal(chargeresponse)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func CreateChargeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	amount_string := r.FormValue("amount")
	sourceToken := r.FormValue("sourceToken")
	amount, err := strconv.ParseInt(amount_string, 10, 64)
	if err != nil {
		chargeresponse := chargeResponse{Success: false, ErrorMessage: "Invalid amount"}

		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	} else {
		params := &stripe.ChargeParams{
			Amount:   stripe.Int64(amount),
			Currency: stripe.String(string(stripe.CurrencyINR)),
			Capture:  stripe.Bool(false),
			Source:   &stripe.SourceParams{Token: stripe.String(sourceToken)},
		}
		c, err := charge.New(params)
		if err != nil {
			chargeresponse := chargeResponse{Success: false}
			if stripeErr, ok := err.(*stripe.Error); ok {
				switch stripeErr.Code {
				case stripe.ErrorCodeCardDeclined:
					chargeresponse.ErrorMessage = "Card declined"
				case stripe.ErrorCodeExpiredCard:
					chargeresponse.ErrorMessage = "Card expired"
				case stripe.ErrorCodeIncorrectCVC:
					chargeresponse.ErrorMessage = "Incorrect CVC"
				case stripe.ErrorCodeInvalidChargeAmount:
					chargeresponse.ErrorMessage = "Invalid Charge amount"
				case stripe.ErrorCodeInvalidSourceUsage:
					chargeresponse.ErrorMessage = "Invalid Card type"
				case stripe.ErrorCodeResourceMissing:
					chargeresponse.ErrorMessage = "Invalid card token"
				default:
					chargeresponse.ErrorMessage = "Failed to create charge"
				}
			} else {
				chargeresponse.ErrorMessage = "Other error occurred"
			}

			jsonResponse, _ := json.Marshal(chargeresponse)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}
		chargeresponse := chargeResponse{Success: true}
		chargeresponse.Id = c.ID
		chargeresponse.SuccessMessage = fmt.Sprintf("Step 1 of card charge for amount %v successfull", c.Amount)
		jsonResponse, _ := json.Marshal(chargeresponse)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func init() {
	stripe.Key = os.Getenv("STRIPE_KEY")
}

func main() {

	r := mux.NewRouter()

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

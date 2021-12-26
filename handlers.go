package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
	"github.com/stripe/stripe-go/v72/refund"
)

func ListChargesHandler(w http.ResponseWriter, r *http.Request) {

	listchargeresponse := listchargeResponse{}
	params := &stripe.ChargeListParams{}
	i := charge.List(params)
	for i.Next() {
		c := i.Charge()
		listchargeresponse.Ids = append(listchargeresponse.Ids, c.ID)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listchargeresponse); err != nil {
		log.Panicf("JSON Encoding error: %v", err)
	}
}

func CaptureChargeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := r.FormValue("id")
	pathId := vars["chargeId"]
	if Id != pathId || Id == "" {
		chargeresponse := chargeResponse{Success: false}
		invalidIdChecker(w, Id, pathId, chargeresponse)
		return
	}
	c, err := charge.Capture(
		Id,
		nil,
	)
	if err != nil {
		chargeresponse := chargeResponse{Success: false}
		chargeresponse = errorResponseSetter(err, chargeresponse)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(chargeresponse); err != nil {
			log.Panicf("JSON Encoding error: %v", err)
		}
		return
	}
	chargeresponse := chargeResponse{Success: true}
	chargeresponse.Id = c.ID
	chargeresponse.SuccessMessage = fmt.Sprintf("Charged amount %v successfully, transaction id is: %v", c.Amount, c.ID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(chargeresponse); err != nil {
		log.Panicf("JSON Encoding error: %v", err)
	}
}

func CreateRefundHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	Id := r.FormValue("id")
	pathId := vars["chargeId"]
	if Id != pathId || Id == "" {
		chargeresponse := chargeResponse{Success: false}
		invalidIdChecker(w, Id, pathId, chargeresponse)
		return
	}
	params := &stripe.RefundParams{
		Charge: stripe.String(Id),
	}
	refund_data, err := refund.New(params)
	if err != nil {
		chargeresponse := chargeResponse{Success: false}
		if _, ok := err.(*stripe.Error); ok {
			chargeresponse.ErrorMessage = "Charge already refunded or invalid transaction id"
		} else {
			chargeresponse.ErrorMessage = "Other error occurred"
		}

		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(chargeresponse); err != nil {
			log.Panicf("JSON Encoding error: %v", err)
		}
		return
	}
	chargeresponse := chargeResponse{Success: true}
	chargeresponse.Id = refund_data.ID
	chargeresponse.SuccessMessage = fmt.Sprintf("Charged amount %v successfully refunded, refund id is: %v", refund_data.Amount, refund_data.ID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(chargeresponse); err != nil {
		log.Panicf("JSON Encoding error: %v", err)
	}
}

func CreateChargeHandler(w http.ResponseWriter, r *http.Request) {

	amount_string := r.FormValue("amount")
	sourceToken := r.FormValue("sourceToken")
	amount, err := strconv.ParseInt(amount_string, 10, 64)
	if err != nil {
		chargeresponse := chargeResponse{Success: false, ErrorMessage: "Invalid amount"}

		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(chargeresponse); err != nil {
			log.Panicf("JSON Encoding error: %v", err)
		}
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
			chargeresponse = errorResponseSetter(err, chargeresponse)
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(chargeresponse); err != nil {
				log.Panicf("JSON Encoding error: %v", err)
			}
			return
		}
		chargeresponse := chargeResponse{Success: true}
		chargeresponse.Id = c.ID
		chargeresponse.SuccessMessage = fmt.Sprintf("Step 1 of card charge for amount %v successfull", c.Amount)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(chargeresponse); err != nil {
			log.Panicf("JSON Encoding error: %v", err)
		}
	}
}

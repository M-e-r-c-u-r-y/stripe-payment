package main

import (
	"bytes"
	"log"
	"net/http"
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

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

type LogMiddleware struct {
	logger *log.Logger
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func TrackLoanApplication(response http.ResponseWriter, request *http.Request) {
	rspJson := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(rspJson)
	jsonEncoder.SetEscapeHTML(false)
	if LoanDetails.LoanID.Status == "PAID" {
		LoanDetails.LoanID.Status = "NO MORE TERMS TO PAY"
		err := jsonEncoder.Encode(LoanDetails.LoanID.Status)
		if err != nil {
			SendHttpResponse(400, response, nil)
		}
	} else {
		err := jsonEncoder.Encode(LoanDetails)
		if err != nil {
			SendHttpResponse(400, response, nil)
		}
	}
	SendHttpResponse(200, response, rspJson.Bytes())
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func ApproveLaon(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Request reached to approve the laon")
	if LoanDetails.LoanID.Status == "pending" {
		LoanDetails.LoanID.Status = "APPROVED"
	}
	fmt.Println("Loan Approved for LoanID-- %v", LoanDetails.LoanID.LoanId)
	rspJson := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(rspJson)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(LoanDetails.LoanID)
	if err != nil {
		SendHttpResponse(400, response, nil)
	}
	SendHttpResponse(200, response, rspJson.Bytes())
}

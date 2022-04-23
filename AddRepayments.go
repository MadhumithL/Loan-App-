package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PayAmount struct {
	Amount int
}
type Status struct {
	Status string `json:"status""`
}

var amountRemaining int

func RepayAmount(response http.ResponseWriter, request *http.Request) {
	var amount PayAmount
	err := json.NewDecoder(request.Body).Decode(&amount)
	if err != nil {
		SendHttpResponse(400, response, nil)
	}
	var statusPaid Status
	if LoanDetails.LoanID.Status == "PAID" && amount.Amount > 0 || LoanDetails.LoanID.Status == "NO MORE TERMS TO PAY" && amount.Amount > 0 {
		statusPaid.Status = "NO MORE PENDING TERMS, MONEY WILL CREDIT BACK"
	} else {
		for i := 0; i < createLoanRequest.Term; i++ {
			if LoanDetails.LoanID.Status == "APPROVED" && LoanDetails.Terms[i].Status == "PAID" {
				continue
			} else if LoanDetails.LoanID.Status == "APPROVED" && LoanDetails.Terms[i].Status == "pending" && amount.Amount >= LoanDetails.Terms[i].EMI {
				LoanDetails.Terms[i].Status = "PAID"
				statusPaid.Status = "PAID"
				amountRemaining = amount.Amount - LoanDetails.Terms[i].EMI + amountRemaining
				break
			} else if LoanDetails.LoanID.Status == "APPROVED" && LoanDetails.Terms[i].Status == "pending" && amount.Amount+amountRemaining >= LoanDetails.Terms[i].EMI {
				fmt.Println("INTO")
				LoanDetails.Terms[i].Status = "PAID"
				statusPaid.Status = "PAID"
				amountRemaining = amount.Amount - LoanDetails.Terms[i].EMI + amountRemaining
				break
			}
		}
	}
	rspJson := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(rspJson)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(statusPaid)
	if LoanDetails.Terms[createLoanRequest.Term-1].Status == "PAID" {
		LoanDetails.LoanID.Status = "PAID"
	}
	if err != nil {
		SendHttpResponse(400, response, nil)
	}
	SendHttpResponse(200, response, rspJson.Bytes())
}

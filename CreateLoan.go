package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type createLoanReqStruct struct {
	CustomerId string
	LoanAmount int
	Term       int
}

type Terms struct {
	Status   string `json:"status"`
	TermDate string `json:"date"`
	EMI      int    `json:"amount"`
}
type LoanDetail struct {
	LoanId int    `json:"loanId"`
	Status string `json:"loan_status"`
}

type LoanDetailResp struct {
	LoanID LoanDetail `json:"loan_details"`
	Terms  []Terms    `json:"Terms"`
}

var CustomerMap = make(map[string]LoanDetailResp)

var LoanDetails LoanDetailResp
var createLoanRequest createLoanReqStruct

func CreateLoan(response http.ResponseWriter, request *http.Request) {
	err := json.NewDecoder(request.Body).Decode(&createLoanRequest)
	if err != nil {
		SendHttpResponse(503, response, nil)
	}

	createLoanResponse := make([]Terms, createLoanRequest.Term)
	current := time.Now()
	eachTermValue := createLoanRequest.LoanAmount / createLoanRequest.Term
	for i := 0; i < createLoanRequest.Term; i++ {
		createLoanResponse[i].Status = "pending"
		createLoanResponse[i].EMI = eachTermValue
		current = current.AddDate(0, 0, 7)
		createLoanResponse[i].TermDate = current.String()
	}
	LoanDetails.LoanID.LoanId = rand.Intn(9999999999)
	LoanDetails.LoanID.Status = "pending"
	LoanDetails.Terms = createLoanResponse
	CustomerMap[createLoanRequest.CustomerId] = LoanDetails
	rspJson := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(rspJson)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(LoanDetails)
	SendHttpResponse(200, response, rspJson.Bytes())
}

func SendHttpResponse(statusCode int, response http.ResponseWriter, rspJson []byte) {
	response.WriteHeader(statusCode)
	if rspJson != nil {
		response.Header().Add("Content-Type", "application/json")
		response.Write(rspJson)
	}
}

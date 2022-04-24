package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateLoanAPI(t *testing.T) {

	var jsonStr = []byte(`{
    "loanAmount": 10000,
    "term": 3
}`)
	req, err := http.NewRequest("POST", "/createloan", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateLoan)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"loan_details":{"loanId":2505480089,"loan_status":"pending"},"Terms":[{"status":"pending","date":"2022-04-30 21:51:20.72131 +0530 IST","amount":3333},{"status":"pending","date":"2022-05-07 21:51:20.72131 +0530 IST","amount":3333},{"status":"pending","date":"2022-05-14 21:51:20.72131 +0530 IST","amount":3333}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAdminApproveLoan(t *testing.T) {
	req, err := http.NewRequest("UPDATE", "/admin/approve/2505480089", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ApproveLaon)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"loanId":2505480089,"loan_status":"APPROVED"}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestTrackLoan(t *testing.T) {
	req, err := http.NewRequest("GET", "/track/2505480089", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TrackLoanApplication)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect, after first payment
	expected := `{"loan_details":{"loanId":2505480089,"loan_status":"APPROVED"},"Terms":[{"status":"PAID","date":"2022-04-30 21:51:20.72131 +0530 IST","amount":3333},{"status":"pending","date":"2022-05-07 21:51:20.72131 +0530 IST","amount":3333},{"status":"pending","date":"2022-05-14 21:51:20.72131 +0530 IST","amount":3333}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestTrackLaonAfterAllEMIsPAID(t *testing.T) {
	req, err := http.NewRequest("GET", "/track/2505480089", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TrackLoanApplication)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{loan_details:{"loan_id":2505480089, "loan_status":"PAID"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRepaymentFuncnality(t *testing.T) {
	req, err := http.NewRequest("GET", "/pay/2505480089", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RepayAmount)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"status":"PAID"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRepaymentFuncnalityIfAllEMIsPaidAndHavingExtraAmountRemaining(t *testing.T) {
	req, err := http.NewRequest("GET", "/pay/2505480089", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RepayAmount)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"status":"NO MORE PENDING TERMS, MONEY WILL CREDIT BACK"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

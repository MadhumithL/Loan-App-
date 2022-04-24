package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
	"os"
)

func getAppRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/createloan", AddContext(CreateLoan)).Methods("POST")		// creates a loan_id, returns with terms scheduled
	router.HandleFunc("/admin/approve/{loanId}", AddContext(ApproveLaon)).Methods("PUT")	//updates the loan status to approves
	router.HandleFunc("/track/{loanId}", AddContext(TrackLoanApplication)).Methods("GET")	//tracks the loanstatus and EMIs paid
	router.HandleFunc("/pay/{loanId}", AddContext(RepayAmount)).Methods("POST")		//change the status of the emi to paid
	return router
}

var host = "localhost"
var port = "8080"

func startServer(r *mux.Router) {
	go func() {
		err := http.ListenAndServe("localhost:8080", r)
		if err != nil {
			fmt.Println("Starting server failed for %v,%v, with Error : %v", host, port, err)
			os.Exit(1)
		}
	}()
	return
}

func main() {
	ar := getAppRouter()
	startServer(ar)
	select {}
}

func AddContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println("Request failed")
		}
		fmt.Println("Request Received----", string(reqDump))
		next.ServeHTTP(res, req)
		fmt.Println("Response sent %s,%s", req.Method, req.URL.String())
	})
}

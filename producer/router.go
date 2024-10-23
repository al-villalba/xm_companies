package main

import (
	"net/http"
	"producer/common"
	"producer/company"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Register REST routes handlers
	r.Handle("/company", common.MwAuthenticate(common.MwAuthorizeWriter(http.HandlerFunc(company.PostCompanyHandler)))).Methods("POST")
	r.Handle("/company", common.MwAuthenticate(http.HandlerFunc(company.PatchCompanyHandler))).Methods("PUT")
	r.HandleFunc("/company/{uuid}", company.GetCompanyHandler).Methods("GET")
	r.Handle("/company/{uuid}", common.MwAuthenticate(http.HandlerFunc(company.DeleteCompanyHandler))).Methods("DELETE")

	return r
}

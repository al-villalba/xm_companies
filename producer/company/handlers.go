package company

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"producer/common"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// serve POST /company
func PostCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// read income data
	var company *Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		logrus.WithError(err).Error("parsing incomming post")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := company.validate(); err != nil {
		logrus.WithError(err).Error("parsing incomming post")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// save into db
	err = company.save()
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok {
			if mysqlErr.Number == 1062 {
				http.Error(w, "Duplicate entry", http.StatusConflict)
				return
			}
		}
		logrus.WithError(err).Error("saving into db")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// produce event
	message, err := company.serialiseMessage("create")
	if err != nil {
		logrus.WithError(err).Error("serialising company")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = common.ProduceEvent(message)
	if err != nil {
		logrus.WithError(err).Error("produce event failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Write output
	common.ResponseJson(w, company, http.StatusCreated)
}

// serve PUT /company
func PatchCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// read income data
	var company *Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		logrus.WithError(err).Error("parsing incomming post")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := company.validate(); err != nil {
		logrus.WithError(err).Error("parsing incomming post")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update db
	err = company.update()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		logrus.WithError(err).Error("updating db")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// produce event
	message, err := company.serialiseMessage("patch")
	if err != nil {
		logrus.WithError(err).Error("serialising company")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = common.ProduceEvent(message)
	if err != nil {
		logrus.WithError(err).Error("produce event failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Write output
	common.ResponseJson(w, company, http.StatusOK)
}

// serve GET /company/{uuid}
func GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// get uuid param
	vars := mux.Vars(r)
	id, ok := vars["uuid"]
	if !ok || !isValidUUID(id) {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	company, err := fetchCompany(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		logrus.WithError(err).Error("Fetching company from db")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// write output
	common.ResponseJson(w, company, http.StatusOK)
}

// serve DELETE /company/{uuid}
func DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// get uuid param
	vars := mux.Vars(r)
	id, ok := vars["uuid"]
	if !ok || !isValidUUID(id) {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	company, err := fetchCompany(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		logrus.WithError(err).Error("Fetching company from db")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = deleteCompany(id)
	if err != nil {
		logrus.WithError(err).Error("Deleting company from db")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// produce event
	message, err := company.serialiseMessage("delete")
	if err != nil {
		logrus.WithError(err).Error("serialising company")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = common.ProduceEvent(message)
	if err != nil {
		logrus.WithError(err).Error("produce event failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// write output
	w.WriteHeader(http.StatusNoContent)
}

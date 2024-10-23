package company

import (
	"encoding/json"
	"errors"
	"producer/common"
	"slices"

	"github.com/google/uuid"
)

type Company struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Description  string `json:"description" db:"description"`
	AmtEmployees int    `json:"amt_employees" db:"amt_employees"`
	Registered   bool   `json:"registered" db:"registered"`
	Type         string `json:"type" db:"type"`
}

// Insert into db
func (company *Company) save() error {
	query := `INSERT INTO companies VALUES (UUID_TO_BIN(:id), :name, :description, :amt_employees, :registered, :type)`
	_, err := common.GetDatabase().NamedExec(query, company)
	if err != nil {
		return err
	}

	return nil
}

// Update in the db
func (company *Company) update() error {
	// check company exists
	_, err := fetchCompany(company.ID)
	if err != nil {
		// not found falls here
		return err
	}

	query := `UPDATE companies SET id = UUID_TO_BIN(:id), name = :name, description = :description, amt_employees = :amt_employees, registered = :registered, type = :type WHERE id = UUID_TO_BIN(:id)`
	_, err = common.GetDatabase().NamedExec(query, company)
	if err != nil {
		return err
	}

	return nil
}

// prepare message for kafka
func (company *Company) serialiseMessage(action string) ([]byte, error) {
	message := struct {
		Action  string
		Company Company
	}{action, *company}

	serialised, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return serialised, nil
}

// Fetch from db
func fetchCompany(id string) (Company, error) {
	query := "SELECT BIN_TO_UUID(id) as id, name, description, amt_employees, registered, type FROM companies WHERE id = UUID_TO_BIN(?)"
	company := Company{}
	err := common.GetDatabase().Get(&company, query, id)

	return company, err
}

// Delete from db
func deleteCompany(id string) error {
	query := "DELETE FROM companies WHERE id = UUID_TO_BIN(?)"
	_, err := common.GetDatabase().Exec(query, id)

	return err
}

func (company *Company) validate() error {
	// id
	if !isValidUUID(company.ID) {
		return errors.New("invalid uuid")
	}
	// Name
	if len(company.Name) == 0 || len(company.Name) > 15 {
		return errors.New("invalid name")
	}
	// Description
	if len(company.Name) > 3000 {
		return errors.New("invalid description")
	}
	// No constraints for AmtEmployees and Registered
	// Type
	constraint := []string{"Corporations", "NonProfit", "Cooperative", "Sole Proprietorship"}
	if !slices.Contains(constraint, company.Type) {
		return errors.New("invalid type")
	}

	return nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

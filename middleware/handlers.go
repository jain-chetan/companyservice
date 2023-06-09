package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jain-chetan/companyservice/auth"
	database "github.com/jain-chetan/companyservice/database"
	models "github.com/jain-chetan/companyservice/model"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SwaggerHandler() http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"), // Replace with your server URL
	)
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to decode the request body . %v", err)
		return
	}
	token, err := auth.CreateToken(user)
	if err != nil {
		log.Printf("Unable to create token . %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(token)
}

// @Summary Create a new company
// @Description Create a new company with the specified details
// @Tags company
// @Accept json
// @Produce json
// @Param company body models.Company true "Company object that needs to be created"
// @Success 201 {object} models.Company
// @Failure 400
// @Router /companies [post]
func CreateCompany(w http.ResponseWriter, r *http.Request) {

	_, err := auth.ValidateToken(r.Header)
	if err != nil {
		res := models.Response{
			Code:    400,
			Message: "Not a valid token",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	var company models.Company

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to decode the request body . %v", err)
		return
	}

	if company.Name == "" || company.Employees == 0 || company.Type == "" {
		res := models.Response{
			Code:    400,
			Message: "Name, Employees, and Type are required fields",
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	unique, err := database.CheckNameUniqueness(company.Name)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if unique == 0 {
		companyID := database.CreateCompanyQuery(company)
		res := models.CreateResponse{
			ID:      companyID,
			Code:    201,
			Message: "Company inserted",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)

	} else {
		res := models.Response{
			Code:    400,
			Message: "Name not unique",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}

}

// @Summary Get a company by ID
// @Description Get a company by its ID
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} models.Company
// @Failure 400
// @Failure 404
// @Router /companies/{id} [get]
func GetCompany(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := uuid.Parse(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to convert string to int . %v", err)
		return
	}

	company, err := database.GetCompanyQuery(id)

	if err != nil {

		res := models.Response{
			Code:    404,
			Message: "Company not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(company)
}

// @Summary Update a company by ID
// @Description Update a company with the specified details
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param company body models.Company true "Updated company object"
// @Success 200
// @Failure 400
// @Router /companies/{id} [patch]
func PatchCompany(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ValidateToken(r.Header)
	if err != nil {
		res := models.Response{
			Code:    400,
			Message: "Not a valid token",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	params := mux.Vars(r)

	id, err := uuid.Parse(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to convert string to int . %v", err)
	}

	var company models.Company

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to decode the request body.  %v", err)
	}

	_ = database.PatchCompanyQuery(id, company)

	res := models.Response{
		Code:    200,
		Message: "Updated Successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

// @Summary Delete a company by ID
// @Description Delete a company by its ID
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200
// @Failure 400
// @Router /companies/{id} [delete]
func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ValidateToken(r.Header)
	if err != nil {
		res := models.Response{
			Code:    400,
			Message: "Not a valid token",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to convert the string into int.  %v", err)
		return
	}
	_, err = database.GetCompanyQuery(id)

	if err != nil {

		res := models.Response{
			Code:    404,
			Message: "Company not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	_ = database.DeleteCompanyQuery(id)
	res := models.Response{
		Code:    200,
		Message: "Deleted Successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

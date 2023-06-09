package router

import (
	"net/http"

	middleware "github.com/jain-chetan/companyservice/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	// Serve the Swagger UI
	router.PathPrefix("/docs").Handler(http.StripPrefix("/docs", middleware.SwaggerHandler()))
	router.HandleFunc("/createtoken", middleware.CreateToken).Methods("POST")
	router.HandleFunc("/companies", middleware.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", middleware.PatchCompany).Methods("PATCH")
	router.HandleFunc("/companies/{id}", middleware.GetCompany).Methods("GET")
	router.HandleFunc("/companies/{id}", middleware.DeleteCompany).Methods("DELETE")

	return router
}

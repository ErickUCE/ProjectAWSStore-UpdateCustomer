package routes

import (
	"ProjectAWSStore-UpdateCustomer/controllers"

	"github.com/gorilla/mux"
)

// Configurar las rutas
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Rutas principales
	router.HandleFunc("/customers/{id}", controllers.UpdateCustomer).Methods("PUT")

	return router
}

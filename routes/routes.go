package routes

import (
	"ProjectAWSStore-UpdateCustomer/controllers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes configura las rutas de la API
func SetupRoutes(db *mongo.Database) *mux.Router {
	router := mux.NewRouter()

	// ✅ Inicializar la colección en el controlador (NO devuelve nada, solo se ejecuta)
	controllers.SetCustomerCollection(db)

	// ✅ Configurar controladores
	router.HandleFunc("/customers/{id}", controllers.UpdateCustomer).Methods("PUT")
	router.HandleFunc("/sync-update", controllers.SyncUpdateCustomer).Methods("POST") // 🔥 Usa `SyncUpdateCustomer`

	return router
}

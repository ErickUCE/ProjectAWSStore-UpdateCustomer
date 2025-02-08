package main

import (
	"fmt"
	"log"
	"net/http"

	"ProjectAWSStore-UpdateCustomer/config"
	"ProjectAWSStore-UpdateCustomer/routes"
)

func main() {
	fmt.Println("🚀 Iniciando UpdateCustomerService en Golang...")

	// ✅ Conectar a MongoDB antes de levantar el servidor
	config.ConnectDB()

	// ✅ Configurar rutas
	router := routes.SetupRoutes()

	// ✅ Iniciar el servidor en el puerto 8083
	fmt.Println("✅ Servidor corriendo en el puerto 8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}

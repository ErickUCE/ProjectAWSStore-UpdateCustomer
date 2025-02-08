package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ProjectAWSStore-UpdateCustomer/config"
	"ProjectAWSStore-UpdateCustomer/routes"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🚀 Iniciando UpdateCustomerService en Golang...")

	// 📌 Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Advertencia: No se pudo cargar el archivo .env, verificando variables de entorno...")
	}

	// 📌 Conectar a MongoDB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("❌ Error al conectar con MongoDB:", err)
	}
	fmt.Println("✅ Conexión exitosa a MongoDB")

	// ✅ Configurar rutas y pasar la base de datos correctamente
	router := routes.SetupRoutes(db)

	// 📌 Obtener el puerto desde `.env`
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // 🔥 Valor por defecto
	}

	// ✅ Iniciar el servidor en el puerto definido en `.env`
	fmt.Println("✅ Servidor corriendo en el puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

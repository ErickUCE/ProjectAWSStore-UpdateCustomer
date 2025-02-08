package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var client *mongo.Client

func ConnectDB() {
	fmt.Println("📌 Ejecutando ConnectDB...")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ No se pudo cargar el archivo .env, verificando variables de entorno...")
	}

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	fmt.Println("🔗 MONGO_URI:", uri)
	fmt.Println("🛢️ Base de datos seleccionada:", dbName)

	if uri == "" || dbName == "" {
		log.Fatal("❌ ERROR: MONGO_URI o MONGO_DB_NAME están vacías.")
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ Error conectando a MongoDB:", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ Error haciendo ping a MongoDB:", err)
	}

	fmt.Println("✅ Conexión exitosa a MongoDB")
	DB = client.Database(dbName)
}

// Obtener colección
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("❌ Error: La base de datos no está inicializada. Llama a ConnectDB() primero.")
	}
	return DB.Collection(collectionName)
}

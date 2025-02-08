package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"ProjectAWSStore-UpdateCustomer/config"
	"ProjectAWSStore-UpdateCustomer/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Variable global para la colección de clientes
var customerCollection *mongo.Collection

// 📌 Función para establecer la colección en el controlador
func SetCustomerCollection(db *mongo.Database) {
	customerCollection = db.Collection("customers")
}

// 📌 **Actualizar un cliente**
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"] // 🔥 Usar ID para la consulta

	var updatedCustomer models.Customer
	err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de actualización para:", updatedCustomer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "❌ Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Actualizar cliente en MongoDB
	_, err = customerCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updatedCustomer},
	)
	if err != nil {
		http.Error(w, "❌ Error al actualizar cliente", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente actualizado correctamente:", updatedCustomer.Email)

	// 🔥 **Sincronizar con `ReadCustomer` y `CreateCustomer`**
	syncUpdateWithMicroservices(updatedCustomer)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCustomer)
}

// 📌 **Sincronizar actualización de clientes con `ReadCustomer` y `CreateCustomer`**
func syncUpdateWithMicroservices(updatedCustomer models.Customer) {
	services := []string{
		os.Getenv("READ_CUSTOMER_SERVICE"),
		os.Getenv("CREATE_CUSTOMER_SERVICE"),
	}

	customerJSON, _ := json.Marshal(updatedCustomer)

	for _, service := range services {
		url := service + "/sync-update"
		fmt.Println("🔄 Enviando sincronización a:", url)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(customerJSON))
		if err != nil {
			fmt.Println("❌ Error creando solicitud HTTP:", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("❌ Error enviando solicitud a", url, ":", err)
			continue
		}

		fmt.Println("✅ Sincronización exitosa con:", url, " Status:", resp.Status)
		resp.Body.Close()
	}
}

// 📌 **Sincronizar actualización de clientes desde `ReadCustomer`**
func SyncUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var updatedCustomer models.Customer
	err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización para:", updatedCustomer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Actualizar cliente en MongoDB
	_, err = customerCollection.UpdateOne(
		context.TODO(),
		bson.M{"email": updatedCustomer.Email},
		bson.M{"$set": updatedCustomer},
	)
	if err != nil {
		http.Error(w, "❌ Error al sincronizar actualización", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente sincronizado correctamente:", updatedCustomer.Email)
	w.WriteHeader(http.StatusOK)
}

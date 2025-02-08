package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ProjectAWSStore-UpdateCustomer/config"
	"ProjectAWSStore-UpdateCustomer/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var customerCollection = config.GetCollection("customers")

// 📌 Actualizar un cliente y sincronizar en `CreateCustomer` y `ReadCustomer`
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "❌ ID inválido", http.StatusBadRequest)
		return
	}

	var updatedCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	updatedCustomer.UpdatedAt = time.Now()

	_, err = customerCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updatedCustomer},
	)
	if err != nil {
		http.Error(w, "❌ Error actualizando cliente", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente actualizado:", updatedCustomer.Email)

	// 🔄 Sincronizar actualización en `CreateCustomer` y `ReadCustomer`
	go syncUpdateWithOtherServices(updatedCustomer)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCustomer)
}

// 📌 Función para sincronizar la actualización con los otros microservicios
func syncUpdateWithOtherServices(customer models.Customer) {
	instances := []string{
		"http://localhost:8081/sync-update", // CreateCustomer
		"http://localhost:8082/sync-update", // ReadCustomer
	}

	jsonData, err := json.Marshal(customer)
	if err != nil {
		fmt.Println("❌ Error serializando cliente:", err)
		return
	}

	for _, instance := range instances {
		resp, err := http.Post(instance, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("❌ Error notificando a", instance, ":", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("✅ Cliente sincronizado con", instance)
		} else {
			fmt.Println("⚠️ No se pudo sincronizar cliente con", instance, "Código:", resp.StatusCode)
		}
	}
}

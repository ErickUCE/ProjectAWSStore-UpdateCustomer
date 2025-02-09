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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Variable global para la colección de clientes
var customerCollection *mongo.Collection

// 📌 Función para establecer la colección en el controlador
func SetCustomerCollection(db *mongo.Database) {
	customerCollection = db.Collection("customers")
}

// 📌 **Actualizar un cliente en `UpdateCustomer`**
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// 🔄 Convertir `ID` de string a `primitive.ObjectID`
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "❌ ID inválido", http.StatusBadRequest)
		return
	}

	var updatedCustomer models.Customer
	err = json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de actualización para:", updatedCustomer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Buscar el cliente por `_id` y actualizarlo
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updatedCustomer}

	result, err := customerCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "❌ Error al actualizar cliente", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "⚠️ Cliente no encontrado en UpdateCustomer", http.StatusNotFound)
		return
	}

	fmt.Println("✅ Cliente actualizado correctamente en UpdateCustomer:", updatedCustomer.Email)

	// ✅ Convertir `ObjectID` a string antes de sincronizar
	updatedCustomer.ID = objID.Hex() // 🚀 Convertir `ObjectID` a string para sincronización

	// 🔄 **Sincronizar con ReadCustomer y CreateCustomer**
	go syncUpdateWithMicroservices(updatedCustomer)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCustomer)
}

// 📌 **Sincronizar actualización de clientes con `ReadCustomer` y `CreateCustomer`**
func syncUpdateWithMicroservices(updatedCustomer models.Customer) {
	services := []string{
		os.Getenv("READ_CUSTOMER_SERVICE") + "/sync-update",   // URL de ReadCustomer
		os.Getenv("CREATE_CUSTOMER_SERVICE") + "/sync-update", // URL de CreateCustomer
	}

	customerJSON, _ := json.Marshal(updatedCustomer)

	for _, service := range services {
		if service == "" {
			fmt.Println("⚠️ Servicio no definido, omitiendo...")
			continue
		}

		fmt.Println("🔄 Enviando sincronización a:", service)

		req, err := http.NewRequest("POST", service, bytes.NewBuffer(customerJSON))
		if err != nil {
			fmt.Println("❌ Error creando solicitud HTTP:", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("❌ Error enviando solicitud a", service, ":", err)
			continue
		}

		fmt.Println("✅ Sincronización exitosa con:", service, " Status:", resp.Status)
		resp.Body.Close()
	}
}

// 📌 **Sincronizar actualización de clientes desde `ReadCustomer`**
// 📌 **Sincronizar actualización de clientes desde `ReadCustomer` y `CreateCustomer`**
// 📌 **Sincronizar actualización de clientes desde `UpdateCustomer`**
func SyncUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var updatedCustomer models.Customer

	// 📌 Decodificar el JSON recibido
	err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		fmt.Println("❌ Error al decodificar JSON:", err)
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización para:", updatedCustomer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "❌ Database not initialized", http.StatusInternalServerError)
		return
	}

	// 📌 Verificar si `ID` está presente en el JSON recibido
	if updatedCustomer.ID == "" || updatedCustomer.ID == "000000000000000000000000" {
		fmt.Println("⚠️ Error: `ID` vacío en la sincronización de actualización.")
		http.Error(w, "⚠️ Error: `ID` vacío en la sincronización", http.StatusBadRequest)
		return
	}

	// ✅ Convertir `ID` de string a `primitive.ObjectID`
	objID, err := primitive.ObjectIDFromHex(updatedCustomer.ID)
	if err != nil {
		fmt.Println("⚠️ ID inválido en sincronización:", updatedCustomer.ID)
		http.Error(w, "⚠️ ID inválido en sincronización", http.StatusBadRequest)
		return
	}

	// ✅ Crear el filtro para buscar por `_id`
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updatedCustomer}

	// 📌 Intentar actualizar el cliente en la base de datos
	result, err := customerCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("❌ Error al actualizar cliente en MongoDB:", err)
		http.Error(w, "❌ Error al sincronizar actualización", http.StatusInternalServerError)
		return
	}

	// 📌 Verificar si realmente se encontró y actualizó el cliente
	if result.MatchedCount == 0 {
		fmt.Println("⚠️ Cliente no encontrado en la base de datos durante sincronización.")
		http.Error(w, "⚠️ Cliente no encontrado en la base de datos durante sincronización.", http.StatusNotFound)
		return
	}

	// ✅ Cliente actualizado correctamente
	fmt.Println("✅ Cliente sincronizado correctamente:", updatedCustomer.Email)
	w.WriteHeader(http.StatusOK)
}

// 📌 **Sincronizar clientes desde `CreateCustomer`**
func SyncCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización desde CreateCustomer:", customer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Verificar si el cliente ya existe en UpdateCustomer
	var existingCustomer models.Customer
	err = customerCollection.FindOne(context.TODO(), bson.M{"email": customer.Email}).Decode(&existingCustomer)
	if err == nil {
		fmt.Println("⚠️ Cliente ya existe en UpdateCustomer:", customer.Email)
		w.WriteHeader(http.StatusOK)
		return
	}

	// ✅ Insertar nuevo cliente
	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "❌ Error al sincronizar cliente", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente sincronizado correctamente en UpdateCustomer:", customer.Email)
	w.WriteHeader(http.StatusCreated)
}

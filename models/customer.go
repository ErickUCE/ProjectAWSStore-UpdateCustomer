package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer representa el esquema de la colección en MongoDB
type Customer struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone_number" bson:"phone_number"`
	Address   string             `json:"address" bson:"address"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// 🔥 Serializar ObjectID como String en JSON
func (c *Customer) MarshalJSON() ([]byte, error) {
	type Alias Customer
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    c.ID.Hex(), // 🔥 Convertir ObjectID a string
		Alias: (*Alias)(c),
	})
}

// 🔥 Convertir String a ObjectID al deserializar JSON
func (c *Customer) UnmarshalJSON(data []byte) error {
	type Alias Customer
	aux := &struct {
		ID string `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 🔄 Convertir `ID` de string a `ObjectID`
	objID, err := primitive.ObjectIDFromHex(aux.ID)
	if err != nil {
		return err
	}
	c.ID = objID

	return nil
}

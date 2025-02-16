package models

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransaksiGuru struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GuruID    primitive.ObjectID `bson:"guru_id" json:"guru_id"`
	GuruName  string             `bson:"guru_name" json:"guru_name"`
	Amount    float64            `bson:"amount" json:"amount"`
	CreatedAt string             `bson:"created_at" json:"created_at"` // ⬅️ Ubah dari Date ke string format WIB
	Notes     string             `bson:"notes" json:"notes"`
}

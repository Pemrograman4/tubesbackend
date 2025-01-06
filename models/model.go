package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type Course struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Duration    int                `bson:"duration" json:"duration"`
    Cost        float64            `bson:"cost" json:"cost"`
    Description string             `bson:"description" json:"description"`
    CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`

}

type Siswa struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FullName    string             `bson:"fullname,omitempty" json:"fullname,omitempty"`
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`
	PhoneNumber string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email       string             `bson:"email,omitempty" json:"email,omitempty"`
	Status      string             `bson:"status,omitempty" json:"status,omitempty"` // "aktif" or "nonaktif"
}
	
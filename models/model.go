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

type Guru struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FullName      string             `bson:"fullname,omitempty" json:"fullname,omitempty"`
	Address       string             `bson:"address,omitempty" json:"address,omitempty"`
	PhoneNumber   string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty"`
	SchoolSubject string             `bson:"school_subject,omitempty" json:"school_subject,omitempty"`
	Status        string             `bson:"status,omitempty" json:"status,omitempty"` // "aktif" or "nonaktif"
	JoinedAt      primitive.DateTime `bson:"joined_at,omitempty" json:"joined_at,omitempty"`
}

type Tagihan struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SiswaID      primitive.ObjectID `bson:"siswa_id" json:"siswa_id"`
	CourseID     primitive.ObjectID `bson:"course_id" json:"course_id"`
	Amount       float64            `bson:"amount" json:"amount"` // Total tagihan
	DueDate      primitive.DateTime `bson:"due_date" json:"due_date"` // Tanggal jatuh tempo
	Paid         bool               `bson:"paid" json:"paid"` // Status pembayaran
	PaidAt       primitive.DateTime `bson:"paid_at,omitempty" json:"paid_at,omitempty"` // Jika sudah dibayar
	CreatedAt    primitive.DateTime `bson:"created_at" json:"created_at"`
}
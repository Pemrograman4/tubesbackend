package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`   // Contoh: "admin", "user"
	Status    string             `bson:"status"` // Contoh: "active", "inactive"
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	Schedule  string             `json:"schedule"`
}

type Course struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Duration    int                `bson:"duration" json:"duration"`
	Cost        float64            `bson:"cost" json:"cost"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
	Schedule    string             `bson:"schedule" json:"schedule"` // Tambahkan ini
}

// Struct untuk mapping data Schedule dengan nama kursus dan multiple dates
type Schedule struct {
	Time  []string `json:"time" bson:"time"`
	Dates []string `json:"dates" bson:"dates"`
}

// CourseRegistration merepresentasikan data pendaftaran kursus
type Registration struct {
	CourseId    string   `json:"courseId" bson:"courseId"`
	StudentName string   `json:"studentName" bson:"studentName"`
	Email       string   `json:"email" bson:"email"`
	Phonenumber string   `json:"phonenumber" bson:"phonenumber"`
	Status      string   `json:"status" bson:"status"`
	Courses     []string `json:"courses" bson:"courses"`
}

type Siswa struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FullName    string             `bson:"fullname,omitempty" json:"fullname,omitempty"`
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`
	PhoneNumber string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email       string             `bson:"email,omitempty" json:"email,omitempty"`
	Status      string             `bson:"status,omitempty" json:"status,omitempty"` // "aktif" or "nonaktif"
}
type TransaksiSiswa struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SiswaID primitive.ObjectID `bson:"siswa_id" json:"siswa_id"`
	UserID  primitive.ObjectID `bson:"user_id" json:"user_id"`
	Item    string             `bson:"item" json:"item"`
	Harga   float64            `bson:"harga" json:"harga"`
	Tanggal primitive.DateTime `bson:"tanggal" json:"tanggal"`
	Status  string             `bson:"status" json:"status"`
}

type Guru struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FullName      string             `bson:"fullname,omitempty" json:"fullname,omitempty"`
	Address       string             `bson:"address,omitempty" json:"address,omitempty"`
	PhoneNumber   string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty"`
	SchoolSubject string             `bson:"school_subject,omitempty" json:"school_subject,omitempty"`
	Status        string             `bson:"status,omitempty" json:"status,omitempty"` // "aktif" atau "nonaktif"
}

type Tagihan struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	SiswaID    primitive.ObjectID  `bson:"siswa_id" json:"siswa_id"`
	SiswaNama  string              `bson:"siswa_nama" json:"siswa_nama"`
	SiswaEmail string              `bson:"siswa_email" json:"siswa_email"`
	CourseID   primitive.ObjectID  `bson:"course_id" json:"course_id"`
	CourseName string              `bson:"course_name" json:"course_name"`
	Amount     float64             `bson:"amount" json:"amount"`
	DueDate    primitive.DateTime  `bson:"due_date" json:"due_date"`
	Paid       bool                `bson:"paid" json:"paid"`
	Status     string              `bson:"status" json:"status"` // Tambahkan status di database
	PaidAt     *primitive.DateTime `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	CreatedAt  primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt  primitive.DateTime  `bson:"updated_at" json:"updated_at"`
}

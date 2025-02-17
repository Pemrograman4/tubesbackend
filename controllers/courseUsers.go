package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Struct Registration untuk mapping data JSON dan MongoDB
type Registration struct {
	CourseId    string   `json:"courseId" bson:"courseId"`
	StudentName string   `json:"studentName" bson:"studentName"`
	Email       string   `json:"email" bson:"email"`
	Phonenumber string   `json:"phonenumber" bson:"phonenumber"`
	Status      string   `json:"status" bson:"status"`
	Courses     []string `json:"courses" bson:"courses"`
}

// Struct CourseUsers menangani pendaftaran kursus
type CourseUsers struct {
	DB *mongo.Database
}

// Constructor NewCourseUsers
func NewCourseUsers(db *mongo.Database) *CourseUsers {
	return &CourseUsers{DB: db}
}

// RegisterCourse untuk mendaftarkan kursus
func (cu *CourseUsers) RegisterCourse(c *gin.Context) {
	var registration Registration

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi semua field wajib diisi
	if registration.CourseId == "" || registration.StudentName == "" || registration.Email == "" ||
		registration.Phonenumber == "" || registration.Status == "" || len(registration.Courses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Menyimpan ke MongoDB
	collection := cu.DB.Collection("course_registrations")
	_, err := collection.InsertOne(context.TODO(), registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pendaftaran kursus berhasil"})
}

// GetAllCourseRegistrations untuk mengambil semua pendaftaran kursus
func (cu *CourseUsers) GetAllCourseRegistrations(c *gin.Context) {
	collection := cu.DB.Collection("course_registrations")

	cursor, err := collection.Find(context.TODO(), bson.M{}) // Mengambil semua data tanpa filter
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var registrations []struct {
		CourseId    string   `json:"courseId" bson:"courseId"`
		StudentName string   `json:"studentName" bson:"studentName"`
		Email       string   `json:"email" bson:"email"`
		Phonenumber string   `json:"phonenumber" bson:"phonenumber"`
		Status      string   `json:"status" bson:"status"`
		Courses     []string `json:"courses" bson:"courses"`
	}

	for cursor.Next(context.TODO()) {
		var registration struct {
			CourseId    string   `json:"courseId" bson:"courseId"`
			StudentName string   `json:"studentName" bson:"studentName"`
			Email       string   `json:"email" bson:"email"`
			Phonenumber string   `json:"phonenumber" bson:"phonenumber"`
			Status      string   `json:"status" bson:"status"`
			Courses     []string `json:"courses" bson:"courses"`
		}
		if err := cursor.Decode(&registration); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		registrations = append(registrations, registration)
	}

	c.JSON(http.StatusOK, registrations)
}

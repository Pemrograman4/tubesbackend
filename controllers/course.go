package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/organisasi/tubesbackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CourseController mengelola endpoint kursus
type CourseController struct {
	DB *mongo.Database
}

// NewCourseController membuat instance CourseController
func NewCourseController(db *mongo.Database) *CourseController {
	return &CourseController{DB: db}
}

// CreateCourse membuat course baru
func (cc *CourseController) CreateCourse(c *gin.Context) {
	var course struct {
		Name        string  `json:"name"`
		Duration    int     `json:"duration"`
		Cost        float64 `json:"cost"`
		Description string  `json:"description"`
		Schedule    string  `json:"schedule"` // Field Schedule ditambahkan
	}

	// Validasi input JSON
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validasi field yang wajib diisi
	if course.Name == "" || course.Schedule == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Schedule are required"})
		return
	}

	// Create a new Course struct based on input
	newCourse := models.Course{
		ID:          primitive.NewObjectID(),
		Name:        course.Name,
		Duration:    course.Duration,
		Cost:        course.Cost,
		Description: course.Description,
		Schedule:    course.Schedule,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()), // Set creation time
	}

	// Insert the new course into the MongoDB collection
	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, newCourse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully"})
}

func (cc *CourseController) GetCourses(c *gin.Context) {
	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}
	defer cursor.Close(ctx)

	var courses []bson.M
	if err = cursor.All(ctx, &courses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// FindCourseById mendapatkan kursus berdasarkan ID
func (cc *CourseController) FindCourseById(c *gin.Context) {
	id := c.Param("id")

	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var course bson.M
	var err error

	// Coba cari berdasarkan custom field "id" (string)
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&course)
	if err != nil {
		// Jika gagal, coba cari berdasarkan "_id" (ObjectID)
		objID, objErr := primitive.ObjectIDFromHex(id)
		if objErr == nil {
			err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&course)
		}

		// Jika tetap tidak ditemukan, kembalikan 404
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
	}

	c.JSON(http.StatusOK, course)
}

func (cc *CourseController) UpdateCourseById(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5504")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Pastikan ID valid untuk ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Buat struktur untuk menerima data yang diupdate
	var updatedCourse struct {
		Name        string  `json:"name"`
		Duration    int     `json:"duration"`
		Cost        float64 `json:"cost"`
		Description string  `json:"description"`
		Schedule    string  `json:"schedule"` // Kolom schedule yang bisa diupdate
	}

	// Bind data dari JSON request body ke struktur updatedCourse
	if err := c.ShouldBindJSON(&updatedCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi bahwa schedule tidak kosong
	if updatedCourse.Schedule == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule is required"})
		return
	}

	// Ambil koleksi "courses" dari database
	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Membuat data yang akan diupdate
	update := bson.M{
		"$set": bson.M{
			"name":        updatedCourse.Name,
			"duration":    updatedCourse.Duration,
			"cost":        updatedCourse.Cost,
			"description": updatedCourse.Description,
			"schedule":    updatedCourse.Schedule, // Update schedule
		},
	}

	// Melakukan update pada dokumen yang sesuai dengan ObjectID
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	// Cek apakah kursus ditemukan
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Jika berhasil, kirimkan respon sukses
	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// DeleteCourse menghapus data kursus berdasarkan ID
func (cc *CourseController) DeleteCourse(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func (cc *CourseController) GetNextCourseId(c *gin.Context) {
	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count courses"})
		return
	}

	nextId := fmt.Sprintf("C-%04d", count+1)
	c.JSON(http.StatusOK, gin.H{"nextId": nextId})
}

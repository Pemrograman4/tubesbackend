package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseController struct {
	DB *mongo.Database
}

// CreateCourse menambahkan data kursus baru
func (cc *CourseController) CreateCourse(c *gin.Context) {
	var course struct {
		Name        string  `json:"name"`
		Duration    int     `json:"duration"`
		Cost        float64 `json:"cost"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

// GetCourses mendapatkan daftar kursus
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

// UpdateCourse memperbarui data kursus berdasarkan ID
func (cc *CourseController) UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var course struct {
		Name        string  `json:"name"`
		Duration    int     `json:"duration"`
		Cost        float64 `json:"cost"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": course}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// DeleteCourse menghapus data kursus berdasarkan ID
func (cc *CourseController) DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	// Validasi apakah ID adalah ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Jika ID bukan ObjectID, validasi sebagai string biasa
		result, deleteErr := cc.DB.Collection("courses").DeleteOne(
			c.Request.Context(),
			bson.M{"id": id}, // Gunakan kolom "id" yang sesuai di MongoDB
		)

		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
		return
	}

	// Jika ID valid sebagai ObjectID, lanjutkan proses penghapusan
	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// GetLatestCourseId mendapatkan ID kursus terbaru
func (cc *CourseController) GetLatestCourseId(c *gin.Context) {
	collection := cc.DB.Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var course struct {
		ID        primitive.ObjectID `bson:"_id"`
		CreatedAt time.Time          `bson:"createdAt"`
	}

	err := collection.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "No courses found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch latest course ID"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"latestId": course.ID.Hex()})
}

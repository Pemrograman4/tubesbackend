package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

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

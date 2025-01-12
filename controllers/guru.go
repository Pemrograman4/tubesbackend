package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/organisasi/tubesbackend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GuruController struct {
	DB *mongo.Database
}

// GetAllGuru retrieves all Guru records.
func (ctrl *GuruController) GetAllGuru(c *gin.Context) {
	var gurus []models.Guru
	cursor, err := ctrl.DB.Collection("gurus").Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var guru models.Guru
		if err := cursor.Decode(&guru); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding data"})
			return
		}
		gurus = append(gurus, guru)
	}

	c.JSON(http.StatusOK, gurus)
}

// CreateGuru creates a new Guru record.
func (ctrl *GuruController) CreateGuru(c *gin.Context) {
	var guru models.Guru
	if err := c.ShouldBindJSON(&guru); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	guru.ID = primitive.NewObjectID()
	guru.JoinedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := ctrl.DB.Collection("gurus").InsertOne(context.TODO(), guru)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Guru"})
		return
	}

	c.JSON(http.StatusCreated, guru)
}

// GetGuruByID retrieves a Guru by ID.
func (ctrl *GuruController) GetGuruByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var guru models.Guru
	err = ctrl.DB.Collection("gurus").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&guru)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guru not found"})
		return
	}

	c.JSON(http.StatusOK, guru)
}

// UpdateGuru updates an existing Guru record.
func (ctrl *GuruController) UpdateGuru(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updateData models.Guru
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = ctrl.DB.Collection("gurus").UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updateData},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Guru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guru updated successfully"})
}

// DeleteGuru deletes a Guru record by ID.
func (ctrl *GuruController) DeleteGuru(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = ctrl.DB.Collection("gurus").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Guru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guru deleted successfully"})
}

package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/organisasi/tubesbackend/models"
)

type TagihanController struct {
	DB *mongo.Database
}

// GetTagihan mendapatkan daftar tagihan
func (sc *TagihanController) GetTagihan(c *gin.Context) {
	fmt.Println("GetTagihan called") // Tambahkan log
	collection := sc.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tagihan: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var tagihanList []models.Tagihan
	if err = cursor.All(ctx, &tagihanList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse tagihan data: " + err.Error()})
		return
	}

	if len(tagihanList) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tagihan found"})
		return
	}

	c.JSON(http.StatusOK, tagihanList)
}

// CreateTagihan creates a new Tagihan record.
func (ctrl *TagihanController) CreateTagihan(c *gin.Context) {
	var tagihanInput struct {
		SiswaID  string  `json:"siswa_id"`
		CourseID string  `json:"course_id"`
		Amount   float64 `json:"amount"`
		DueDate  string  `json:"due_date"` // Format ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)
		Status   string  `json:"status"`   // Tambahkan field Status
	}

	// Validate input
	if err := c.ShouldBindJSON(&tagihanInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Convert string IDs to MongoDB types
	siswaID, err := primitive.ObjectIDFromHex(tagihanInput.SiswaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiswaID"})
		return
	}

	courseID, err := primitive.ObjectIDFromHex(tagihanInput.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CourseID"})
		return
	}

	// Parse due date
	dueDateTime, err := time.Parse(time.RFC3339, tagihanInput.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DueDate format"})
		return
	}

	// Convert to primitive.DateTime
	dueDate := primitive.NewDateTimeFromTime(dueDateTime)

	// Create a new Tagihan object
	tagihan := models.Tagihan{
		ID:        primitive.NewObjectID(),
		SiswaID:   siswaID,
		CourseID:  courseID,
		Amount:    tagihanInput.Amount,
		DueDate:   dueDate,
		Paid:      false,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Status:    tagihanInput.Status, // Tambahkan nilai Status
	}

	// Insert into database
	_, err = ctrl.DB.Collection("tagihans").InsertOne(context.TODO(), tagihan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Tagihan"})
		return
	}

	c.JSON(http.StatusCreated, tagihan)
}

// GetTagihanByID mengambil data siswa berdasarkan ID
func (sc *TagihanController) GetTagihanByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := sc.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tagihan models.Tagihan
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&tagihan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tagihan not found"})
		return
	}

	c.JSON(http.StatusOK, tagihan)
}

// UpdateTagihan memperbarui data tagihan berdasarkan ID
func (ctrl *TagihanController) UpdateTagihan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var tagihanInput struct {
		SiswaID  string  `json:"siswa_id"`
		CourseID string  `json:"course_id"`
		Amount   float64 `json:"amount"`
		DueDate  string  `json:"due_date"` // Format ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)
		Status   string  `json:"status"`
		Paid     bool    `json:"paid"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&tagihanInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Convert string IDs to MongoDB types
	siswaID, err := primitive.ObjectIDFromHex(tagihanInput.SiswaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiswaID"})
		return
	}

	courseID, err := primitive.ObjectIDFromHex(tagihanInput.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CourseID"})
		return
	}

	// Parse due date
	dueDateTime, err := time.Parse(time.RFC3339, tagihanInput.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DueDate format"})
		return
	}

	// Convert to primitive.DateTime
	dueDate := primitive.NewDateTimeFromTime(dueDateTime)

	collection := ctrl.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"siswa_id":   siswaID,
			"course_id":  courseID,
			"amount":     tagihanInput.Amount,
			"due_date":   dueDate,
			"status":     tagihanInput.Status,
			"paid":       tagihanInput.Paid,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tagihan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan updated successfully"})
}

// DeleteTagihan menghapus data tagihan berdasarkan ID
func (ctrl *TagihanController) DeleteTagihan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := ctrl.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tagihan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan deleted successfully"})
}



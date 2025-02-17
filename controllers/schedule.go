package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/organisasi/tubesbackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Struct untuk mapping data Schedule dengan nama kursus dan multiple dates
type Schedule struct {
	CourseId string   `json:"courseId" bson:"courseId"`
	Name     string   `json:"name" bson:"name"` // Nama kursus ditambahkan di sini
	Time     []string `json:"time" bson:"time"`
	Dates    []string `json:"dates" bson:"dates"`
}

// Struct ScheduleController untuk menangani jadwal kursus
type ScheduleController struct {
	DB *mongo.Database
}

// Constructor NewScheduleController
func NewScheduleController(db *mongo.Database) *ScheduleController {
	return &ScheduleController{DB: db}
}

// AddSchedule untuk menambahkan jadwal kursus dengan beberapa tanggal
func (sc *ScheduleController) AddSchedule(c *gin.Context) {
	var schedule Schedule

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi field wajib
	if schedule.CourseId == "" || len(schedule.Time) == 0 || len(schedule.Dates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CourseId, Time, and Dates are required"})
		return
	}

	// Menyimpan jadwal ke MongoDB
	collection := sc.DB.Collection("course_schedules")

	// Mengecek apakah courseId sudah ada
	var existingSchedule Schedule
	err := collection.FindOne(context.TODO(), bson.M{"courseId": schedule.CourseId}).Decode(&existingSchedule)

	if err == nil {
		// Jika sudah ada, update dengan menambahkan tanggal baru
		_, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"courseId": schedule.CourseId},
			bson.M{"$push": bson.M{"dates": bson.M{"$each": schedule.Dates}}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if err == mongo.ErrNoDocuments {
		// Jika belum ada, buat jadwal baru
		_, err := collection.InsertOne(context.TODO(), schedule)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course schedule added successfully"})
}

// GetAllSchedules untuk mengambil semua jadwal kursus
func (sc *ScheduleController) GetAllSchedules(c *gin.Context) {
	collection := sc.DB.Collection("course_schedules")

	cursor, err := collection.Find(context.TODO(), bson.M{}) // Mengambil semua data tanpa filter
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var schedules []Schedule

	for cursor.Next(context.TODO()) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, schedules)
}

// GetScheduleByCourseId untuk mengambil jadwal berdasarkan CourseId
func (sc *ScheduleController) GetScheduleByCourseId(c *gin.Context) {
	courseId := c.Param("courseId")

	// Convert courseId menjadi lowercase agar konsisten
	courseId = strings.ToLower(courseId)

	// Cari kursus berdasarkan name yang cocok
	collection := sc.DB.Collection("courses")
	var course models.Course // Menggunakan models.Course
	err := collection.FindOne(context.TODO(), bson.M{"name": courseId}).Decode(&course)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Jika kursus ditemukan, ambil jadwalnya
	scheduleCollection := sc.DB.Collection("course_schedules")
	var schedule models.Schedule // Pastikan ada model Schedule
	err = scheduleCollection.FindOne(context.TODO(), bson.M{"courseId": course.ID.Hex()}).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Gabungkan informasi kursus dengan jadwal
	c.JSON(http.StatusOK, gin.H{
		"courseId": course.ID.Hex(),
		"name":     course.Name,
		"schedule": schedule,
	})
}

// UpdateSchedule untuk memperbarui jadwal berdasarkan CourseId
func (sc *ScheduleController) UpdateSchedule(c *gin.Context) {
	courseId := c.Param("courseId")
	var updatedSchedule Schedule

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&updatedSchedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi
	if len(updatedSchedule.Time) == 0 || len(updatedSchedule.Dates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Time and Dates are required"})
		return
	}

	// Update jadwal berdasarkan courseId
	collection := sc.DB.Collection("course_schedules")
	result := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"courseId": courseId},
		bson.M{"$set": bson.M{"time": updatedSchedule.Time, "dates": updatedSchedule.Dates}},
	)

	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Err().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}

// DeleteSchedule untuk menghapus jadwal berdasarkan CourseId
func (sc *ScheduleController) DeleteSchedule(c *gin.Context) {
	courseId := c.Param("courseId")

	// Menghapus jadwal berdasarkan courseId
	collection := sc.DB.Collection("course_schedules")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"courseId": courseId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}

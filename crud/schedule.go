package crud

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "golangproject/models"
    "gorm.io/gorm"
)

// CreateSchedule handles the creation of a new schedule
func CreateSchedule(c *gin.Context, db *gorm.DB) {
    var input struct {
        Teacher_id     uint      `json:"teacher_id" binding:"required"`
        Available_time time.Time `json:"available_time" binding:"required"`
    }

    // Bind input from the request
    if err := c.ShouldBind(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Create a new schedule record
    schedule := models.Schedule{
        Teacher_id:     input.Teacher_id,
        Available_time: input.Available_time,
        Is_booked:      false, // By default, the schedule is not booked
    }

    // Save to the database
    if err := db.Create(&schedule).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
        return
    }

    // Return success message
    c.JSON(http.StatusOK, gin.H{"message": "Schedule created successfully", "schedule": schedule})
}

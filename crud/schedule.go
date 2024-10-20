package crud

import (
    "fmt"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gin-gonic/gin"
    "golangproject/models"
    "gorm.io/gorm"
)

// CreateSchedule handles the creation of a new schedule
func CreateSchedule(c *gin.Context, db *gorm.DB) {
    var input struct {
        TeacherID     string `json:"teacher_id" binding:"required"`  // รับ TeacherID เป็นสตริง
        AvailableTime string `json:"available_time" binding:"required"`
    }
    // ตรวจสอบการ bind input จาก request
    if err := c.ShouldBind(&input); err != nil {
        fmt.Println("Bind Error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
        return
    }

    // แปลง Teacher_id จากสตริงเป็น uint
    teacherID, err := strconv.ParseUint(input.TeacherID, 10, 32)
    if err != nil {
        fmt.Println("Teacher ID Parse Error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Teacher ID ไม่ถูกต้อง"})
        return
    }

    // ลองแปลง Available_time จากสตริงเป็น time.Time รองรับ ISO 8601
    parsedTime, err := time.Parse(time.RFC3339, input.AvailableTime)
    if err != nil {
        fmt.Println("Time Parse Error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบเวลาที่ส่งมาไม่ถูกต้อง"})
        return
    }

    // สร้างเรคคอร์ด schedule ใหม่
    schedule := models.Schedule{
        Teacher_id:     uint(teacherID),
        Available_time: parsedTime,
        Is_booked:      false,
    }

    // บันทึกข้อมูลลงฐานข้อมูล
    if err := db.Create(&schedule).Error; err != nil {
        fmt.Println("Database Error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างตารางได้"})
        return
    }

    // ส่งข้อความความสำเร็จกลับไป
    c.JSON(http.StatusOK, gin.H{"message": "สร้างตารางสำเร็จ", "schedule": schedule})
}


func GetSchedule(c *gin.Context, db *gorm.DB) {
    // รับ teacher_id จาก URL
    teacherID := c.Param("teacher_id")

    // ตรวจสอบว่า teacher_id เป็นตัวเลขหรือไม่
    if _, err := strconv.Atoi(teacherID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "teacher_id ไม่ถูกต้อง"})
        return
    }

    var schedules []models.Schedule

    // ดึงตารางเวลาทั้งหมดของอาจารย์ตาม teacher_id
    if err := db.Where("teacher_id = ?", teacherID).Find(&schedules).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลตารางได้"})
        return
    }

    // ตรวจสอบว่ามีตารางเวลาหรือไม่
    if len(schedules) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "ไม่มีตารางเวลาสำหรับอาจารย์นี้"})
        return
    }

    // ส่งข้อมูลตารางทั้งหมดกลับไปยัง frontend พร้อมจำนวนตารางเวลา
    c.JSON(http.StatusOK, gin.H{
        "count":     len(schedules),
        "schedules": schedules,
    })
}

// UpdateSchedule อัปเดตข้อมูลตารางเวลาของอาจารย์
func UpdateSchedule(c *gin.Context, db *gorm.DB) {
    scheduleID := c.Param("id")  // รับ ID ของตารางเวลาจาก URL

    // ตรวจสอบว่า scheduleID เป็นตัวเลขหรือไม่
    var schedule models.Schedule
    if err := db.First(&schedule, scheduleID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบตารางเวลานี้"})
        return
    }

    // รับข้อมูลใหม่จาก JSON ที่ส่งมาจาก client
    var input struct {
        AvailableTime string `json:"available_time"`  // ฟิลด์ available_time ที่จะอัปเดต
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
        return
    }

    // แปลง available_time เป็น time.Time
    parsedTime, err := time.Parse("2006-01-02T15:04:05", input.AvailableTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบเวลาที่ส่งมาไม่ถูกต้อง"})
        return
    }

    // อัปเดตข้อมูลตารางเวลา
    schedule.Available_time = parsedTime

    // บันทึกการเปลี่ยนแปลงลงในฐานข้อมูล
    if err := db.Save(&schedule).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตตารางเวลาได้"})
        return
    }

    // ส่งข้อมูลที่อัปเดตกลับไปยัง client
    c.JSON(http.StatusOK, gin.H{"message": "อัปเดตตารางเวลาสำเร็จ", "schedule": schedule})
}

func DeleteSchedule(c *gin.Context, db *gorm.DB) {
    scheduleID := c.Param("id")  // รับ ID ของตารางเวลาจาก URL

    // ตรวจสอบว่า scheduleID เป็นตัวเลขหรือไม่
    var schedule models.Schedule
    if err := db.First(&schedule, scheduleID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบตารางเวลานี้"})
        return
    }

    // ลบตารางเวลา
    if err := db.Delete(&schedule).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบตารางเวลาได้"})
        return
    }

    // ส่งข้อความยืนยันว่าลบสำเร็จ
    c.JSON(http.StatusOK, gin.H{"message": "ลบตารางเวลาสำเร็จ"})
}
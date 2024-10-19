package crud

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "golangproject/models"
)

// CreateBooking สร้างข้อมูลการจองใหม่
func CreateBooking(c *gin.Context, db *gorm.DB) {
    var input struct {
        ScheduleID uint   `json:"schedule_id" binding:"required"`  // รับ ScheduleID แทนการใช้ TeacherID โดยตรง
        Subject    string `json:"subject" binding:"required"`
    }

    // ตรวจสอบข้อมูลที่ส่งเข้ามา
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
        return
    }

    // ดึง StudentID จาก session หรือ token ของผู้ใช้ที่ล็อกอิน
    studentID, exists := c.Get("userID")  // สมมุติว่าคุณเก็บ StudentID ไว้ใน context หรือ session
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่ได้ล็อกอิน"})
        return
    }

    // ดึงข้อมูลตารางเวลาจาก Schedule เพื่อตรวจสอบ
    var schedule models.Schedule
    if err := db.First(&schedule, input.ScheduleID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบตารางเวลานี้"})
        return
    }

    // ตรวจสอบว่าถูกจองแล้วหรือยัง
    if schedule.Is_booked {
        c.JSON(http.StatusBadRequest, gin.H{"error": "เวลานี้ถูกจองแล้ว"})
        return
    }

    // สร้างการจองใหม่
    booking := models.Booking{
        Student_id:  studentID.(uint),       // ใช้ StudentID จากผู้ใช้ที่ล็อกอิน
        Teacher_id:  schedule.Teacher_id,    // ใช้ TeacherID จากตาราง Schedule
        Time:        schedule.Available_time,  // ใช้เวลาในตาราง Schedule
        Subject:     input.Subject,
        Status:      models.Pending,         // ตั้งค่าสถานะเริ่มต้นเป็น Pending
    }

    // บันทึกการจองลงฐานข้อมูล
    if err := db.Create(&booking).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างการจองได้"})
        return
    }

    // ส่งข้อมูลการจองกลับไป
    c.JSON(http.StatusOK, gin.H{"message": "สร้างการจองสำเร็จ", "booking": booking})
}

// UpdateBookingStatus อัปเดตสถานะการจอง
func UpdateBookingStatus(c *gin.Context, db *gorm.DB) {
    bookingID := c.Param("id")

    var input struct {
        Status models.BookingStatus `json:"status" binding:"required"`
    }

    // ตรวจสอบข้อมูลที่ส่งมา
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
        return
    }

    // ค้นหาข้อมูลการจองตาม bookingID
    var booking models.Booking
    if err := db.First(&booking, bookingID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลการจองนี้"})
        return
    }

    // ตรวจสอบสถานะที่ส่งมา
    var message string
    if input.Status == models.Accepted {
        // อัปเดตตารางเวลา (Schedule) ว่าเวลานี้ถูกจองแล้ว
        var schedule models.Schedule
        if err := db.Where("teacher_id = ? AND available_time = ?", booking.Teacher_id, booking.Time).First(&schedule).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่พบตารางเวลาที่ต้องการอัปเดต"})
            return
        }

        // อัปเดตสถานะการจองของ Schedule เป็น true
        schedule.Is_booked = true
        if err := db.Save(&schedule).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตสถานะของตารางเวลาได้"})
            return
        }

        // ข้อความการแจ้งเตือนสำหรับสถานะ accepted
        message = "การจองของคุณถูกยืนยันแล้ว"

    } else if input.Status == models.Rejected {
        // ข้อความการแจ้งเตือนสำหรับสถานะ rejected
        message = "การจองของคุณถูกปฏิเสธแล้ว"
    }

    // หากสถานะเป็น accepted หรือ rejected ให้สร้างการแจ้งเตือน
    if input.Status == models.Accepted || input.Status == models.Rejected {
        notification := models.Notification{
            UserID:  booking.Student_id,
            Message: message,
            Is_read:  false,
        }

        // บันทึกการแจ้งเตือนลงในฐานข้อมูล
        if err := db.Create(&notification).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างการแจ้งเตือนได้"})
            return
        }
    }

    // อัปเดตสถานะการจองในตาราง Booking
    booking.Status = input.Status
    if err := db.Save(&booking).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตสถานะการจองได้"})
        return
    }

    // ส่งข้อมูลการจองที่ถูกอัปเดตกลับไป
    c.JSON(http.StatusOK, gin.H{"message": "อัปเดตสถานะการจองสำเร็จ", "booking": booking})
}

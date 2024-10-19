package main

import (
	"log"
	"net/http"
	"os"

	"golangproject/crud"
	"golangproject/controllers"
	"golangproject/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
    dsn := "meet:1234@tcp(127.0.0.1:3306)/meet?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate the User model
    db.AutoMigrate(&models.User{})
    db.AutoMigrate(&models.Schedule{})
    db.AutoMigrate(&models.Booking{})
    db.AutoMigrate(&models.Notification{})

    router := gin.Default()

    // HTML templates
    router.LoadHTMLGlob("templates/*")

    // Routes for registration and login
    router.GET("/register", func(c *gin.Context) {
        c.HTML(http.StatusOK, "register.html", nil)
    })

    router.POST("/register", func(c *gin.Context) {
        controllers.Register(c, db)
    })

    router.GET("/login", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
    })

    router.POST("/login", func(c *gin.Context) {
        controllers.Login(c, db)
    })

	router.GET("/create", func(c * gin.Context){
		c.HTML(http.StatusOK, "create_schedule.html", nil)
	})

	router.POST("/create", func(c *gin.Context) {
		crud.CreateSchedule(c, db)
	})
    router.GET("/schedule/:teacher_id", func(c *gin.Context) {
        c.HTML(http.StatusOK, "schedule.html", nil)
    })
    router.GET("/schedule/:teacher_id/api", func(c *gin.Context) {
        crud.GetSchedule(c, db)
    })

    router.PUT("/schedule/:id", func(c *gin.Context) {
        crud.UpdateSchedule(c, db)
    })

    router.DELETE("/schedule/:id", func(c *gin.Context) {
        crud.DeleteSchedule(c, db)
    })

    router.POST("/bookings",controllers.AuthMiddleware(), func(c *gin.Context) {
        crud.CreateBooking(c, db)
    })

    router.PUT("/bookings/:id/status", func(c *gin.Context) {
        crud.UpdateBookingStatus(c, db)
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    router.Run(":" + port)
}
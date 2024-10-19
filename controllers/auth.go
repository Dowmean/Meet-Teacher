// controllers.go
package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "golangproject/models"
)

var jwtKey = []byte("secret_key")

// Claims represents the JWT claims
type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}

// Register handles user registration
func Register(c *gin.Context, db *gorm.DB) {
    var input struct {
        Name     string `form:"name" binding:"required"`
        Email    string `form:"email" binding:"required"`
        Password string `form:"password" binding:"required"`
        Is_role  string `form:"role" binding:"required"`  // รับค่า role จากแบบฟอร์ม
    }

    // Bind input from registration form
    if err := c.ShouldBind(&input); err != nil {
        // c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Invalid data"})
        return
    }

    // Hash the password before storing it
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Error hashing password"})
        return
    }

    // Create a new user record
    user := models.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: string(hashedPassword),
        Is_role:  models.Role(input.Is_role),  // บันทึกบทบาทของผู้ใช้
    }

    // Save user to the database
    if err := db.Create(&user).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "User already exists"})
        return
    }

    // On success, redirect to login page
    c.HTML(http.StatusOK, "login.html", gin.H{"message": "Registration successful! Please log in."})
}
// Login handles user authentication and JWT generation
func Login(c *gin.Context, db *gorm.DB) {
    var input struct {
        Email    string `form:"email" binding:"required"`
        Password string `form:"password" binding:"required"`
    }

    // Bind login input
    if err := c.ShouldBind(&input); err != nil {
        c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid data"})
        return
    }

    // Find user by email
    var user models.User
    if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
        return
    }

    // Verify the password with the stored hash
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
        return
    }

    // JWT token generation with userID and email
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,   // เก็บ userID ลงใน JWT token
        Email:  user.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    // Create the token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Error generating token"})
        return
    }

    // Set the JWT token in a cookie (cookie อายุ 24 ชั่วโมง)
    c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)

    // Send response after successful login
    c.HTML(http.StatusOK, "dashboard.html", gin.H{"message": "Login successful!"})
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ดึง token จากคุกกี้
        tokenString, err := c.Cookie("token")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, please log in"})
            c.Abort()
            return
        }

        // Parse token
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, invalid token"})
            c.Abort()
            return
        }

        // เก็บ userID ลงใน context เพื่อใช้ในคำขออื่นๆ
        c.Set("userID", claims.UserID)

        c.Next()
    }
}
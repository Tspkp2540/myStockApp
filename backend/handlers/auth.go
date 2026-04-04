package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"restaurant-stock/database"
	"restaurant-stock/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const sessionDuration = 24 * time.Hour
const cookieName = "session_token"

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Login handles POST /api/auth/login
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุชื่อผู้ใช้และรหัสผ่าน"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "บัญชีถูกระงับ"})
		return
	}

	if !checkPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้าง session ได้"})
		return
	}

	session := models.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	database.DB.Create(&session)

	c.SetCookie(cookieName, token, int(sessionDuration.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Logout handles POST /api/auth/logout
func Logout(c *gin.Context) {
	token, err := c.Cookie(cookieName)
	if err == nil && token != "" {
		database.DB.Where("token = ?", token).Delete(&models.Session{})
	}
	c.SetCookie(cookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "ออกจากระบบแล้ว"})
}

// GetMe handles GET /api/auth/me
func GetMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่ได้เข้าสู่ระบบ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ==================== User Management (Admin) ====================

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Order("id").Find(&users)
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var input struct {
		Username string      `json:"username" binding:"required"`
		Password string      `json:"password" binding:"required,min=4"`
		FullName string      `json:"full_name"`
		Role     models.Role `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role == "" {
		input.Role = models.RoleEmployee
	}
	if input.Role != models.RoleAdmin && input.Role != models.RoleEmployee {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role ต้องเป็น admin หรือ employee"})
		return
	}

	hashed, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเข้ารหัสรหัสผ่านได้"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: hashed,
		FullName: input.FullName,
		Role:     input.Role,
		Active:   true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ชื่อผู้ใช้นี้มีอยู่แล้ว"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้"})
		return
	}

	var input struct {
		FullName string      `json:"full_name"`
		Role     models.Role `json:"role"`
		Active   *bool       `json:"active"`
		Password string      `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Role != "" {
		if input.Role != models.RoleAdmin && input.Role != models.RoleEmployee {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role ต้องเป็น admin หรือ employee"})
			return
		}
		user.Role = input.Role
	}
	if input.Active != nil {
		user.Active = *input.Active
	}
	if input.Password != "" {
		hashed, err := hashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเข้ารหัสรหัสผ่านได้"})
			return
		}
		user.Password = hashed
	}

	database.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	if err := database.DB.Delete(&models.User{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Also remove their sessions
	database.DB.Where("user_id = ?", c.Param("id")).Delete(&models.Session{})
	c.JSON(http.StatusOK, gin.H{"message": "ลบผู้ใช้แล้ว"})
}

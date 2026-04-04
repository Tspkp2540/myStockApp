package handlers

import (
	"net/http"
	"time"

	"restaurant-stock/database"
	"restaurant-stock/models"

	"github.com/gin-gonic/gin"
)

// AuthRequired validates the session cookie
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(cookieName)
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "กรุณาเข้าสู่ระบบ"})
			return
		}

		var session models.Session
		if err := database.DB.Preload("User").Where("token = ?", token).First(&session).Error; err != nil {
			c.SetCookie(cookieName, "", -1, "/", "", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session ไม่ถูกต้อง"})
			return
		}

		if time.Now().After(session.ExpiresAt) {
			database.DB.Delete(&session)
			c.SetCookie(cookieName, "", -1, "/", "", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session หมดอายุ"})
			return
		}

		if !session.User.Active {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "บัญชีถูกระงับ"})
			return
		}

		c.Set("user", session.User)
		c.Set("role", session.User.Role)
		c.Next()
	}
}

// AdminOnly allows only admin role
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "ต้องเป็นผู้ดูแลระบบเท่านั้น"})
			return
		}
		c.Next()
	}
}

package main

import (
	"log"
	"os"

	"restaurant-stock/database"
	"restaurant-stock/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/stock.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Cookie"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		// Public — auth endpoints
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/logout", handlers.Logout)
	}

	// Protected — require login
	auth := api.Group("")
	auth.Use(handlers.AuthRequired())
	{
		auth.GET("/auth/me", handlers.GetMe)

		// Categories
		auth.GET("/categories", handlers.GetCategories)
		auth.POST("/categories", handlers.CreateCategory)
		auth.PUT("/categories/:id", handlers.UpdateCategory)
		auth.DELETE("/categories/:id", handlers.DeleteCategory)
		auth.POST("/categories/bulk-delete", handlers.BulkDeleteCategories)

		// Units
		auth.GET("/units", handlers.GetUnits)
		auth.POST("/units", handlers.CreateUnit)
		auth.PUT("/units/:id", handlers.UpdateUnit)
		auth.DELETE("/units/:id", handlers.DeleteUnit)
		auth.POST("/units/bulk-delete", handlers.BulkDeleteUnits)

		// Ingredients
		auth.GET("/ingredients", handlers.GetIngredients)
		auth.POST("/ingredients", handlers.CreateIngredient)
		auth.PUT("/ingredients/:id", handlers.UpdateIngredient)
		auth.DELETE("/ingredients/:id", handlers.DeleteIngredient)
		auth.POST("/ingredients/bulk-delete", handlers.BulkDeleteIngredients)

		// Stock transactions
		auth.GET("/transactions", handlers.GetTransactions)
		auth.POST("/transactions", handlers.CreateTransaction)

		// Dashboard
		auth.GET("/dashboard", handlers.GetDashboard)
	}

	// Admin only
	admin := auth.Group("")
	admin.Use(handlers.AdminOnly())
	{
		admin.GET("/users", handlers.GetUsers)
		admin.POST("/users", handlers.CreateUser)
		admin.PUT("/users/:id", handlers.UpdateUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)
	}

	// Serve frontend static files if the directory exists (Railway single-container mode)
	if _, err := os.Stat("./static"); err == nil {
		r.Static("/assets", "./static/assets")
		r.StaticFile("/vite.svg", "./static/vite.svg")
		r.NoRoute(func(c *gin.Context) {
			c.File("./static/index.html")
		})
		log.Println("Serving frontend static files from ./static")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

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
		// Public — health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

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
		auth.GET("/transactions/cost-summary", handlers.GetIngredientCostSummary)

		// Dashboard
		auth.GET("/dashboard", handlers.GetDashboard)

		// Menu items (read-only for all authenticated users — frontoffice needs this)
		auth.GET("/menu-items", handlers.GetMenuItems)
		auth.GET("/menu-items/:id", handlers.GetMenuItem)

		// Front-Office — Dashboard & Sales (employee)
		auth.GET("/frontoffice/dashboard", handlers.GetFrontofficeDashboard)
		auth.POST("/frontoffice/sales", handlers.CreateSale)
		auth.GET("/frontoffice/sales", handlers.GetSales)
		auth.DELETE("/frontoffice/sales/:id", handlers.DeleteSaleWithReason)
	}

	// Admin only
	admin := auth.Group("")
	admin.Use(handlers.AdminOnly())
	{
		// Backoffice — Menu management (admin only)
		admin.POST("/menu-items", handlers.CreateMenuItem)
		admin.PUT("/menu-items/:id", handlers.UpdateMenuItem)
		admin.DELETE("/menu-items/:id", handlers.DeleteMenuItem)

		// Backoffice — Sales (admin only)
		admin.POST("/sales", handlers.CreateSale)
		admin.GET("/sales", handlers.GetSales)
		admin.GET("/sales/:id", handlers.GetSale)
		admin.DELETE("/sales/:id", handlers.DeleteSale)

		// Backoffice — Dashboard (admin only)
		admin.GET("/backoffice/dashboard", handlers.GetBackofficeDashboard)

		admin.GET("/users", handlers.GetUsers)
		admin.POST("/users", handlers.CreateUser)
		admin.PUT("/users/:id", handlers.UpdateUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		// Transaction management (admin)
		admin.DELETE("/transactions/:id", handlers.DeleteTransaction)
		admin.GET("/transactions/deleted", handlers.GetDeletedTransactions)
		admin.POST("/transactions/:id/restore", handlers.RestoreTransaction)

		// Deleted sales review (admin — backoffice)
		admin.GET("/sales/deleted", handlers.GetDeletedSales)
	}

	// Serve frontend static files if the directory exists (Railway single-container mode)
	if _, err := os.Stat("./static"); err == nil {
		r.Static("/assets", "./static/assets")
		r.StaticFile("/vite.svg", "./static/vite.svg")

		// Backoffice static files
		if _, err := os.Stat("./static-backoffice"); err == nil {
			backoffice := r.Group("/backoffice")
			backoffice.Static("/assets", "./static-backoffice/assets")
			backoffice.GET("", func(c *gin.Context) {
				c.File("./static-backoffice/index.html")
			})
			backoffice.GET("/*path", func(c *gin.Context) {
				c.File("./static-backoffice/index.html")
			})
			log.Println("Serving backoffice static files from ./static-backoffice")
		}

		// Front-office static files
		if _, err := os.Stat("./static-frontoffice"); err == nil {
			fo := r.Group("/front-office")
			fo.Static("/assets", "./static-frontoffice/assets")
			fo.GET("", func(c *gin.Context) {
				c.File("./static-frontoffice/index.html")
			})
			fo.GET("/*path", func(c *gin.Context) {
				c.File("./static-frontoffice/index.html")
			})
			log.Println("Serving frontoffice static files from ./static-frontoffice")
		}

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

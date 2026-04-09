package handlers

import (
	"math"
	"net/http"
	"time"

	"restaurant-stock/database"
	"restaurant-stock/models"

	"github.com/gin-gonic/gin"
)

// ==================== Menu Items ====================

func GetMenuItems(c *gin.Context) {
	var items []models.MenuItem
	query := database.DB.Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Ingredients.Ingredient.Unit").Order("name")
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if active := c.Query("active"); active == "true" {
		query = query.Where("active = ?", true)
	}
	query.Find(&items)

	// Calculate cost for each item
	type ItemWithCost struct {
		models.MenuItem
		CostPerUnit float64 `json:"cost_per_unit"`
	}
	var result []ItemWithCost
	for _, item := range items {
		cost := calculateMenuCost(item.ID)
		result = append(result, ItemWithCost{MenuItem: item, CostPerUnit: cost})
	}

	c.JSON(http.StatusOK, result)
}

func GetMenuItem(c *gin.Context) {
	var item models.MenuItem
	if err := database.DB.Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Ingredients.Ingredient.Unit").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบเมนู"})
		return
	}
	cost := calculateMenuCost(item.ID)
	c.JSON(http.StatusOK, gin.H{
		"menu_item":     item,
		"cost_per_unit": cost,
	})
}

func CreateMenuItem(c *gin.Context) {
	var body struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		ImageURL    string  `json:"image_url"`
		Ingredients []struct {
			IngredientID uint    `json:"ingredient_id"`
			Quantity     float64 `json:"quantity"`
		} `json:"ingredients"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.MenuItem{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		ImageURL:    body.ImageURL,
		Active:      true,
	}
	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างเมนูได้"})
		return
	}

	for _, ing := range body.Ingredients {
		if ing.IngredientID == 0 || ing.Quantity <= 0 {
			continue
		}
		mi := models.MenuIngredient{
			MenuItemID:   item.ID,
			IngredientID: ing.IngredientID,
			Quantity:     ing.Quantity,
		}
		database.DB.Create(&mi)
	}

	database.DB.Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Ingredients.Ingredient.Unit").First(&item, item.ID)
	c.JSON(http.StatusCreated, item)
}

func UpdateMenuItem(c *gin.Context) {
	var item models.MenuItem
	if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบเมนู"})
		return
	}

	var body struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
		Active      *bool   `json:"active"`
		Ingredients []struct {
			IngredientID uint    `json:"ingredient_id"`
			Quantity     float64 `json:"quantity"`
		} `json:"ingredients"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Name != "" {
		item.Name = body.Name
	}
	item.Description = body.Description
	item.Price = body.Price
	item.ImageURL = body.ImageURL
	if body.Active != nil {
		item.Active = *body.Active
	}
	database.DB.Save(&item)

	// Replace ingredients
	if body.Ingredients != nil {
		database.DB.Where("menu_item_id = ?", item.ID).Delete(&models.MenuIngredient{})
		for _, ing := range body.Ingredients {
			if ing.IngredientID == 0 || ing.Quantity <= 0 {
				continue
			}
			mi := models.MenuIngredient{
				MenuItemID:   item.ID,
				IngredientID: ing.IngredientID,
				Quantity:     ing.Quantity,
			}
			database.DB.Create(&mi)
		}
	}

	database.DB.Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Ingredients.Ingredient.Unit").First(&item, item.ID)
	c.JSON(http.StatusOK, item)
}

func DeleteMenuItem(c *gin.Context) {
	database.DB.Where("menu_item_id = ?", c.Param("id")).Delete(&models.MenuIngredient{})
	if err := database.DB.Delete(&models.MenuItem{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ลบเมนูแล้ว"})
}

// ==================== Sales ====================

func CreateSale(c *gin.Context) {
	var body struct {
		Items []struct {
			MenuItemID uint `json:"menu_item_id" binding:"required"`
			Quantity   int  `json:"quantity" binding:"required"`
		} `json:"items" binding:"required"`
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(body.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ต้องมีรายการอย่างน้อย 1 รายการ"})
		return
	}

	user := c.MustGet("user").(models.User)

	sale := models.Sale{
		UserID: user.ID,
		Note:   body.Note,
	}
	database.DB.Create(&sale)

	var totalRevenue, totalCost float64

	for _, item := range body.Items {
		if item.Quantity <= 0 {
			continue
		}
		var menuItem models.MenuItem
		if err := database.DB.Preload("Ingredients").First(&menuItem, item.MenuItemID).Error; err != nil {
			continue
		}

		cost := calculateMenuCost(menuItem.ID) * float64(item.Quantity)
		lineTotal := menuItem.Price * float64(item.Quantity)

		si := models.SaleItem{
			SaleID:     sale.ID,
			MenuItemID: menuItem.ID,
			Quantity:   item.Quantity,
			Price:      menuItem.Price,
			Cost:       cost,
		}
		database.DB.Create(&si)

		// Deduct ingredients from stock
		for _, mi := range menuItem.Ingredients {
			totalQty := mi.Quantity * float64(item.Quantity)
			var ing models.Ingredient
			if err := database.DB.First(&ing, mi.IngredientID).Error; err == nil {
				ing.Quantity -= totalQty
				if ing.Quantity < 0 {
					ing.Quantity = 0
				}
				database.DB.Save(&ing)

				// Record stock-out transaction
				txn := models.Transaction{
					IngredientID: mi.IngredientID,
					Type:         models.StockOut,
					Quantity:     totalQty,
					Note:         "ขาย: " + menuItem.Name,
					UserID:       user.ID,
				}
				database.DB.Create(&txn)
			}
		}

		totalRevenue += lineTotal
		totalCost += cost
	}

	sale.Total = math.Round(totalRevenue*100) / 100
	sale.TotalCost = math.Round(totalCost*100) / 100
	sale.Profit = math.Round((totalRevenue-totalCost)*100) / 100
	database.DB.Save(&sale)

	database.DB.Preload("Items").Preload("Items.MenuItem").Preload("User").First(&sale, sale.ID)
	c.JSON(http.StatusCreated, sale)
}

func GetSales(c *gin.Context) {
	var sales []models.Sale
	query := database.DB.Preload("Items").Preload("Items.MenuItem").Preload("User").Order("created_at DESC")
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom+"T00:00:00Z")
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("created_at <= ?", dateTo+"T23:59:59Z")
	}
	query.Limit(500).Find(&sales)
	c.JSON(http.StatusOK, sales)
}

func GetSale(c *gin.Context) {
	var sale models.Sale
	if err := database.DB.Preload("Items").Preload("Items.MenuItem").Preload("User").First(&sale, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบรายการขาย"})
		return
	}
	c.JSON(http.StatusOK, sale)
}

func DeleteSale(c *gin.Context) {
	if err := database.DB.Where("sale_id = ?", c.Param("id")).Delete(&models.SaleItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Delete(&models.Sale{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ลบรายการขายแล้ว"})
}

// ==================== Backoffice Dashboard ====================

func GetBackofficeDashboard(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	// Today's sales
	var todaySales int64
	var todayRevenue, todayCost float64
	database.DB.Model(&models.Sale{}).Where("DATE(created_at) = ?", today).Count(&todaySales)
	database.DB.Model(&models.Sale{}).Where("DATE(created_at) = ?", today).Select("COALESCE(SUM(total),0)").Scan(&todayRevenue)
	database.DB.Model(&models.Sale{}).Where("DATE(created_at) = ?", today).Select("COALESCE(SUM(total_cost),0)").Scan(&todayCost)

	// Total menu items
	var totalMenuItems int64
	database.DB.Model(&models.MenuItem{}).Where("active = ?", true).Count(&totalMenuItems)

	// Top selling items (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	type TopItem struct {
		MenuItemID uint    `json:"menu_item_id"`
		Name       string  `json:"name"`
		TotalQty   int     `json:"total_qty"`
		TotalSales float64 `json:"total_sales"`
	}
	var topItems []TopItem
	database.DB.Model(&models.SaleItem{}).
		Select("sale_items.menu_item_id, menu_items.name, SUM(sale_items.quantity) as total_qty, SUM(sale_items.price * sale_items.quantity) as total_sales").
		Joins("JOIN menu_items ON menu_items.id = sale_items.menu_item_id").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Where("sales.created_at >= ?", thirtyDaysAgo+"T00:00:00Z").
		Group("sale_items.menu_item_id, menu_items.name").
		Order("total_qty DESC").
		Limit(10).
		Scan(&topItems)

	// Recent sales
	var recentSales []models.Sale
	database.DB.Preload("Items").Preload("Items.MenuItem").Preload("User").Order("created_at DESC").Limit(10).Find(&recentSales)

	c.JSON(http.StatusOK, gin.H{
		"today_sales":      todaySales,
		"today_revenue":    todayRevenue,
		"today_cost":       todayCost,
		"today_profit":     math.Round((todayRevenue-todayCost)*100) / 100,
		"total_menu_items": totalMenuItems,
		"top_items":        topItems,
		"recent_sales":     recentSales,
	})
}

// ==================== Helpers ====================

func calculateMenuCost(menuItemID uint) float64 {
	var ingredients []models.MenuIngredient
	database.DB.Where("menu_item_id = ?", menuItemID).Find(&ingredients)

	var totalCost float64
	for _, mi := range ingredients {
		// Get average cost per unit from recent stock-in transactions
		var avgCost float64
		database.DB.Model(&models.Transaction{}).
			Where("ingredient_id = ? AND type = ?", mi.IngredientID, models.StockIn).
			Select("COALESCE(AVG(price), 0)").
			Scan(&avgCost)
		totalCost += avgCost * mi.Quantity
	}
	return math.Round(totalCost*100) / 100
}

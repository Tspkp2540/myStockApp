package handlers

import (
	"fmt"
	"net/http"
	"time"

	"restaurant-stock/database"
	"restaurant-stock/models"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Order("name").Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&cat)
	c.JSON(http.StatusCreated, cat)
}

func UpdateCategory(c *gin.Context) {
	var cat models.Category
	if err := database.DB.First(&cat, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&cat)
	c.JSON(http.StatusOK, cat)
}

func DeleteCategory(c *gin.Context) {
	if err := database.DB.Delete(&models.Category{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func BulkDeleteCategories(c *gin.Context) {
	var body struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids is required"})
		return
	}
	if err := database.DB.Delete(&models.Category{}, body.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted", "count": len(body.IDs)})
}

// ==================== Units ====================

func GetUnits(c *gin.Context) {
	var units []models.Unit
	database.DB.Order("name").Find(&units)
	c.JSON(http.StatusOK, units)
}

func CreateUnit(c *gin.Context) {
	var unit models.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&unit)
	c.JSON(http.StatusCreated, unit)
}

func UpdateUnit(c *gin.Context) {
	var unit models.Unit
	if err := database.DB.First(&unit, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&unit)
	c.JSON(http.StatusOK, unit)
}

func DeleteUnit(c *gin.Context) {
	if err := database.DB.Delete(&models.Unit{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func BulkDeleteUnits(c *gin.Context) {
	var body struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids is required"})
		return
	}
	if err := database.DB.Delete(&models.Unit{}, body.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted", "count": len(body.IDs)})
}

func GetIngredients(c *gin.Context) {
	var ingredients []models.Ingredient
	query := database.DB.Preload("Category").Preload("Unit").Order("name")
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	query.Find(&ingredients)
	c.JSON(http.StatusOK, ingredients)
}

func CreateIngredient(c *gin.Context) {
	var body struct {
		Name       string  `json:"name"`
		CategoryID uint    `json:"category_id"`
		UnitID     uint    `json:"unit_id"`
		Quantity   float64 `json:"quantity"`
		MinStock   float64 `json:"min_stock"`
		Note       string  `json:"note"`
		Price      float64 `json:"price"`
		Date       string  `json:"date"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.UnitID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unit_id is required"})
		return
	}

	ing := models.Ingredient{
		Name:       body.Name,
		CategoryID: body.CategoryID,
		UnitID:     body.UnitID,
		Quantity:   body.Quantity,
		MinStock:   body.MinStock,
	}
	database.DB.Create(&ing)

	// Auto create stock-in transaction if initial quantity > 0
	if body.Quantity > 0 {
		note := body.Note
		if note == "" {
			note = "สต็อคเริ่มต้น"
		}
		user := c.MustGet("user").(models.User)
		txn := models.Transaction{
			IngredientID: ing.ID,
			Type:         models.StockIn,
			Quantity:     body.Quantity,
			Price:        body.Price,
			TotalCost:    body.Price * body.Quantity,
			Note:         note,
			UserID:       user.ID,
		}
		if body.Date != "" {
			if parsed, err := time.Parse("2006-01-02", body.Date); err == nil {
				txn.CreatedAt = parsed
			}
		}
		database.DB.Create(&txn)
	}

	database.DB.Preload("Category").Preload("Unit").First(&ing, ing.ID)
	c.JSON(http.StatusCreated, ing)
}

func UpdateIngredient(c *gin.Context) {
	var ing models.Ingredient
	if err := database.DB.First(&ing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	if err := c.ShouldBindJSON(&ing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ing.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
		return
	}
	if ing.UnitID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unit_id is required"})
		return
	}
	database.DB.Save(&ing)
	database.DB.Preload("Category").Preload("Unit").First(&ing, ing.ID)
	c.JSON(http.StatusOK, ing)
}

func DeleteIngredient(c *gin.Context) {
	if err := database.DB.Delete(&models.Ingredient{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func BulkDeleteIngredients(c *gin.Context) {
	var body struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids is required"})
		return
	}
	if err := database.DB.Delete(&models.Ingredient{}, body.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted", "count": len(body.IDs)})
}

func GetTransactions(c *gin.Context) {
	var txns []models.Transaction
	query := database.DB.Preload("Ingredient").Preload("Ingredient.Category").Preload("Ingredient.Unit").Preload("User").Order("created_at DESC")
	if ingID := c.Query("ingredient_id"); ingID != "" {
		query = query.Where("ingredient_id = ?", ingID)
	}
	if txnType := c.Query("type"); txnType != "" {
		query = query.Where("type = ?", txnType)
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom+"T00:00:00Z")
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("created_at <= ?", dateTo+"T23:59:59Z")
	}
	query.Limit(500).Find(&txns)
	c.JSON(http.StatusOK, txns)
}

func CreateTransaction(c *gin.Context) {
	var txn models.Transaction
	if err := c.ShouldBindJSON(&txn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if txn.Type != models.StockIn && txn.Type != models.StockOut {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'in' or 'out'"})
		return
	}
	if txn.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be positive"})
		return
	}
	var ing models.Ingredient
	if err := database.DB.First(&ing, txn.IngredientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	if txn.Type == models.StockOut && ing.Quantity < txn.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}
	if txn.Type == models.StockIn {
		ing.Quantity += txn.Quantity
	} else {
		ing.Quantity -= txn.Quantity
	}
	txn.TotalCost = txn.Price * txn.Quantity
	user := c.MustGet("user").(models.User)
	txn.UserID = user.ID
	database.DB.Save(&ing)
	database.DB.Create(&txn)
	database.DB.Preload("Ingredient").Preload("Ingredient.Category").Preload("Ingredient.Unit").Preload("User").First(&txn, txn.ID)
	c.JSON(http.StatusCreated, txn)
}

func GetIngredientCostSummary(c *gin.Context) {
	type CostSummary struct {
		IngredientID uint    `json:"ingredient_id"`
		TotalIn      float64 `json:"total_in"`
		TotalOut     float64 `json:"total_out"`
		CostIn       float64 `json:"cost_in"`
		CostOut      float64 `json:"cost_out"`
	}
	var summaries []CostSummary
	database.DB.Model(&models.Transaction{}).Select(
		"ingredient_id, " +
			"SUM(CASE WHEN type = 'in' THEN quantity ELSE 0 END) as total_in, " +
			"SUM(CASE WHEN type = 'out' THEN quantity ELSE 0 END) as total_out, " +
			"SUM(CASE WHEN type = 'in' THEN total_cost ELSE 0 END) as cost_in, " +
			"SUM(CASE WHEN type = 'out' THEN total_cost ELSE 0 END) as cost_out",
	).Group("ingredient_id").Scan(&summaries)
	c.JSON(http.StatusOK, summaries)
}

func DeleteTransaction(c *gin.Context) {
	var txn models.Transaction
	if err := database.DB.Preload("Ingredient").First(&txn, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}

	user := c.MustGet("user").(models.User)

	// Reverse the stock change
	var ing models.Ingredient
	if err := database.DB.First(&ing, txn.IngredientID).Error; err == nil {
		if txn.Type == models.StockIn {
			ing.Quantity -= txn.Quantity
			if ing.Quantity < 0 {
				ing.Quantity = 0
			}
		} else {
			ing.Quantity += txn.Quantity
		}
		database.DB.Save(&ing)
	}

	// Archive to deleted_transactions
	deleted := models.DeletedTransaction{
		OriginalID:        txn.ID,
		IngredientID:      txn.IngredientID,
		Type:              txn.Type,
		Quantity:          txn.Quantity,
		Price:             txn.Price,
		TotalCost:         txn.TotalCost,
		Note:              txn.Note,
		UserID:            txn.UserID,
		OriginalCreatedAt: txn.CreatedAt,
		DeleteReason:      body.Reason,
		DeletedByID:       user.ID,
	}
	database.DB.Create(&deleted)
	database.DB.Delete(&txn)

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func GetDeletedTransactions(c *gin.Context) {
	var deleted []models.DeletedTransaction
	query := database.DB.Preload("Ingredient").Preload("Ingredient.Category").Preload("Ingredient.Unit").Preload("User").Preload("DeletedBy").Order("deleted_at DESC")
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("deleted_at >= ?", dateFrom+"T00:00:00Z")
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("deleted_at <= ?", dateTo+"T23:59:59Z")
	}
	query.Limit(500).Find(&deleted)
	c.JSON(http.StatusOK, deleted)
}

func RestoreTransaction(c *gin.Context) {
	var deleted models.DeletedTransaction
	if err := database.DB.First(&deleted, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deleted transaction not found"})
		return
	}

	// Restore stock change
	var ing models.Ingredient
	if err := database.DB.First(&ing, deleted.IngredientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ingredient no longer exists"})
		return
	}

	if deleted.Type == models.StockIn {
		ing.Quantity += deleted.Quantity
	} else {
		if ing.Quantity < deleted.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock to restore this out transaction"})
			return
		}
		ing.Quantity -= deleted.Quantity
	}
	database.DB.Save(&ing)

	// Re-create original transaction
	txn := models.Transaction{
		IngredientID: deleted.IngredientID,
		Type:         deleted.Type,
		Quantity:     deleted.Quantity,
		Price:        deleted.Price,
		TotalCost:    deleted.TotalCost,
		Note:         deleted.Note,
		UserID:       deleted.UserID,
	}
	database.DB.Create(&txn)
	database.DB.Delete(&deleted)

	c.JSON(http.StatusOK, gin.H{"message": "Restored"})
}

func GetDashboard(c *gin.Context) {
	var totalIngredients int64
	var totalCategories int64
	var lowStockItems []models.Ingredient
	var recentTxns []models.Transaction

	database.DB.Model(&models.Ingredient{}).Count(&totalIngredients)
	database.DB.Model(&models.Category{}).Count(&totalCategories)
	database.DB.Preload("Category").Preload("Unit").Where("quantity <= min_stock").Find(&lowStockItems)
	database.DB.Preload("Ingredient").Preload("Ingredient.Unit").Order("created_at DESC").Limit(10).Find(&recentTxns)

	c.JSON(http.StatusOK, gin.H{
		"total_ingredients":   totalIngredients,
		"total_categories":    totalCategories,
		"low_stock_items":     lowStockItems,
		"recent_transactions": recentTxns,
	})
}

// ==================== Menu Categories ====================

func GetMenuCategories(c *gin.Context) {
	var cats []models.MenuCategory
	database.DB.Order("name").Find(&cats)
	c.JSON(http.StatusOK, cats)
}

func CreateMenuCategory(c *gin.Context) {
	var cat models.MenuCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&cat)
	c.JSON(http.StatusCreated, cat)
}

func UpdateMenuCategory(c *gin.Context) {
	var cat models.MenuCategory
	if err := database.DB.First(&cat, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบหมวดหมู่เมนู"})
		return
	}
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&cat)
	c.JSON(http.StatusOK, cat)
}

func DeleteMenuCategory(c *gin.Context) {
	if err := database.DB.Delete(&models.MenuCategory{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ลบหมวดหมู่เมนูแล้ว"})
}

// ==================== Export Sales to Excel ====================

func ExportSalesExcel(c *gin.Context) {
	var sales []models.Sale
	query := database.DB.Preload("Items").Preload("Items.MenuItem").Preload("User").Order("created_at DESC")
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom+"T00:00:00Z")
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("created_at <= ?", dateTo+"T23:59:59Z")
	}
	query.Limit(5000).Find(&sales)

	f := excelize.NewFile()
	sheet := "Sales"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"#", "วันที่/เวลา", "รายการ", "ยอดรวม", "ต้นทุน", "กำไร", "ประเภทการขาย", "การชำระเงิน", "ผู้ขาย", "หมายเหตุ"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetRowStyle(sheet, 1, 1, headerStyle)

	for row, sale := range sales {
		r := row + 2
		f.SetCellValue(sheet, cellName(1, r), sale.ID)
		f.SetCellValue(sheet, cellName(2, r), sale.CreatedAt.Format("2006-01-02 15:04"))

		var items []string
		for _, item := range sale.Items {
			name := fmt.Sprintf("#%d", item.MenuItemID)
			if item.MenuItem.Name != "" {
				name = item.MenuItem.Name
			}
			items = append(items, fmt.Sprintf("%s x%d", name, item.Quantity))
		}
		f.SetCellValue(sheet, cellName(3, r), joinStrings(items))
		f.SetCellValue(sheet, cellName(4, r), sale.Total)
		f.SetCellValue(sheet, cellName(5, r), sale.TotalCost)
		f.SetCellValue(sheet, cellName(6, r), sale.Profit)

		saleTypeText := "หน้าร้าน"
		if sale.SaleType == models.SaleTypeDelivery {
			saleTypeText = "Delivery"
		}
		f.SetCellValue(sheet, cellName(7, r), saleTypeText)

		paymentText := "เงินสด"
		if sale.PaymentMethod == models.PaymentTransfer {
			paymentText = "เงินโอน"
		}
		f.SetCellValue(sheet, cellName(8, r), paymentText)

		userName := "-"
		if sale.User.FullName != "" {
			userName = sale.User.FullName
		} else if sale.User.Username != "" {
			userName = sale.User.Username
		}
		f.SetCellValue(sheet, cellName(9, r), userName)
		f.SetCellValue(sheet, cellName(10, r), sale.Note)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=sales_export.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างไฟล์ Excel ได้"})
	}
}

func cellName(col, row int) string {
	name, _ := excelize.CoordinatesToCellName(col, row)
	return name
}

func joinStrings(s []string) string {
	result := ""
	for i, v := range s {
		if i > 0 {
			result += ", "
		}
		result += v
	}
	return result
}

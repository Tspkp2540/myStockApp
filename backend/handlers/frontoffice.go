package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"restaurant-stock/database"
	"restaurant-stock/models"

	"github.com/gin-gonic/gin"
)

// ==================== Front-Office: Sales with reason-required delete ====================

// DeleteSaleWithReason archives a sale then removes it — requires a reason
func DeleteSaleWithReason(c *gin.Context) {
	var sale models.Sale
	if err := database.DB.Preload("Items").Preload("Items.MenuItem").First(&sale, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบรายการขาย"})
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Reason) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุเหตุผลในการลบ"})
		return
	}

	user := c.MustGet("user").(models.User)

	// Build items summary text
	var parts []string
	for _, item := range sale.Items {
		name := fmt.Sprintf("%d", item.MenuItemID)
		if item.MenuItem.Name != "" {
			name = item.MenuItem.Name
		}
		parts = append(parts, fmt.Sprintf("%s ×%d", name, item.Quantity))
	}

	archived := models.DeletedSale{
		OriginalID:        sale.ID,
		Total:             sale.Total,
		TotalCost:         sale.TotalCost,
		Profit:            sale.Profit,
		SaleType:          sale.SaleType,
		PaymentMethod:     sale.PaymentMethod,
		Note:              sale.Note,
		ItemsSummary:      strings.Join(parts, ", "),
		UserID:            sale.UserID,
		OriginalCreatedAt: sale.CreatedAt,
		DeleteReason:      strings.TrimSpace(body.Reason),
		DeletedByID:       user.ID,
	}
	database.DB.Create(&archived)

	// Remove sale items then sale
	database.DB.Where("sale_id = ?", sale.ID).Delete(&models.SaleItem{})
	database.DB.Delete(&sale)

	c.JSON(http.StatusOK, gin.H{"message": "ลบรายการขายแล้ว"})
}

// GetDeletedSales returns archived/deleted sales — for backoffice review
func GetDeletedSales(c *gin.Context) {
	var deleted []models.DeletedSale
	query := database.DB.Preload("User").Preload("DeletedBy").Order("deleted_at DESC")
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("deleted_at >= ?", dateFrom+"T00:00:00Z")
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("deleted_at <= ?", dateTo+"T23:59:59Z")
	}
	query.Limit(500).Find(&deleted)
	c.JSON(http.StatusOK, deleted)
}

// ==================== Front-Office Dashboard ====================

func GetFrontofficeDashboard(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	user := c.MustGet("user").(models.User)

	// Today's sales by this user
	var mySalesToday int64
	var myRevenueToday float64
	database.DB.Model(&models.Sale{}).Where("DATE(created_at) = ? AND user_id = ?", today, user.ID).Count(&mySalesToday)
	database.DB.Model(&models.Sale{}).Where("DATE(created_at) = ? AND user_id = ?", today, user.ID).Select("COALESCE(SUM(total),0)").Scan(&myRevenueToday)

	// Today's all sales
	var allSalesToday int64
	var allRevenueToday float64
	database.DB.Model(&models.Sale{}).Where("DATE(created_at) = ?", today).Count(&allSalesToday)
	database.DB.Model(&models.Sale{}).Where("DATE(created_at) = ?", today).Select("COALESCE(SUM(total),0)").Scan(&allRevenueToday)

	// Today's total items sold
	var totalItemsToday int64
	database.DB.Model(&models.SaleItem{}).
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Where("DATE(sales.created_at) = ?", today).
		Select("COALESCE(SUM(sale_items.quantity),0)").Scan(&totalItemsToday)

	// My recent sales
	var recentSales []models.Sale
	database.DB.Preload("Items").Preload("Items.MenuItem").Preload("User").
		Where("user_id = ?", user.ID).
		Order("created_at DESC").Limit(10).Find(&recentSales)

	// Top selling items today
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
		Where("DATE(sales.created_at) = ?", today).
		Group("sale_items.menu_item_id, menu_items.name").
		Order("total_qty DESC").
		Limit(10).
		Scan(&topItems)

	c.JSON(http.StatusOK, gin.H{
		"my_sales_today":    mySalesToday,
		"my_revenue_today":  math.Round(myRevenueToday*100) / 100,
		"all_sales_today":   allSalesToday,
		"all_revenue_today": math.Round(allRevenueToday*100) / 100,
		"total_items_today": totalItemsToday,
		"recent_sales":      recentSales,
		"top_items":         topItems,
	})
}

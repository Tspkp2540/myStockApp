package database

import (
	"log"
	"os"
	"path/filepath"

	"restaurant-stock/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Category{},
		&models.Unit{},
		&models.Ingredient{},
		&models.Transaction{},
		&models.DeletedTransaction{},
		&models.MenuItem{},
		&models.MenuIngredient{},
		&models.Sale{},
		&models.SaleItem{},
		&models.DeletedSale{},
	); err != nil {
		return err
	}

	// Clean up old 'unit' text column if it exists (from previous schema)
	if DB.Migrator().HasColumn(&models.Ingredient{}, "unit") {
		if err := DB.Migrator().DropColumn(&models.Ingredient{}, "unit"); err != nil {
			log.Println("Warning: could not drop old 'unit' column:", err)
		}
	}

	seedData()
	return nil
}

func seedData() {
	// Seed default categories
	defaultCategories := []string{
		"เนื้อสัตว์", "ผัก", "ผลไม้", "เครื่องปรุง", "อาหารแห้ง",
		"นม/ผลิตภัณฑ์นม", "เครื่องดื่ม", "อาหารทะเล", "ของแช่แข็ง", "อื่นๆ",
	}
	for _, name := range defaultCategories {
		var count int64
		DB.Model(&models.Category{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			DB.Create(&models.Category{Name: name})
		}
	}

	// Seed default units
	defaultUnits := []string{
		"กก.", "กรัม", "ลิตร", "มล.", "ชิ้น", "ถุง", "กล่อง",
		"ขวด", "แพ็ค", "ลูก", "หัว", "ต้น", "มัด", "กำ", "โหล",
	}
	for _, name := range defaultUnits {
		var count int64
		DB.Model(&models.Unit{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			DB.Create(&models.Unit{Name: name})
		}
	}

	log.Println("Seed data checked/created successfully")

	// Seed default admin user
	var adminCount int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashed, err := bcrypt.GenerateFromPassword([]byte("admin1234"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Warning: could not hash admin password:", err)
			return
		}
		DB.Create(&models.User{
			Username: "admin",
			Password: string(hashed),
			FullName: "ผู้ดูแลระบบ",
			Role:     models.RoleAdmin,
			Active:   true,
		})
		log.Println("Default admin user created (admin / admin1234)")
	}
}

package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null;uniqueIndex"`
	Password  string    `json:"-" gorm:"not null"`
	FullName  string    `json:"full_name"`
	Role      Role      `json:"role" gorm:"not null;default:employee"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"-" gorm:"not null;uniqueIndex"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Unit struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Ingredient struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	CategoryID uint      `json:"category_id"`
	Category   Category  `json:"category" gorm:"foreignKey:CategoryID"`
	UnitID     uint      `json:"unit_id"`
	Unit       Unit      `json:"unit" gorm:"foreignKey:UnitID"`
	Quantity   float64   `json:"quantity" gorm:"default:0"`
	MinStock   float64   `json:"min_stock" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TransactionType string

const (
	StockIn  TransactionType = "in"
	StockOut TransactionType = "out"
)

type Transaction struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	IngredientID uint            `json:"ingredient_id" gorm:"not null"`
	Ingredient   Ingredient      `json:"ingredient" gorm:"foreignKey:IngredientID"`
	Type         TransactionType `json:"type" gorm:"not null"`
	Quantity     float64         `json:"quantity" gorm:"not null"`
	Price        float64         `json:"price" gorm:"default:0"`
	TotalCost    float64         `json:"total_cost" gorm:"default:0"`
	Note         string          `json:"note"`
	UserID       uint            `json:"user_id"`
	User         User            `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time       `json:"created_at"`
}

type DeletedTransaction struct {
	ID                uint            `json:"id" gorm:"primaryKey"`
	OriginalID        uint            `json:"original_id"`
	IngredientID      uint            `json:"ingredient_id"`
	Ingredient        Ingredient      `json:"ingredient" gorm:"foreignKey:IngredientID"`
	Type              TransactionType `json:"type"`
	Quantity          float64         `json:"quantity"`
	Price             float64         `json:"price"`
	TotalCost         float64         `json:"total_cost"`
	Note              string          `json:"note"`
	UserID            uint            `json:"user_id"`
	User              User            `json:"user" gorm:"foreignKey:UserID"`
	OriginalCreatedAt time.Time       `json:"original_created_at"`
	DeleteReason      string          `json:"delete_reason"`
	DeletedByID       uint            `json:"deleted_by_id"`
	DeletedBy         User            `json:"deleted_by" gorm:"foreignKey:DeletedByID"`
	DeletedAt         time.Time       `json:"deleted_at"`
}

// ==================== Backoffice — Menu & Sales ====================

type MenuCategory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MenuItem struct {
	ID             uint             `json:"id" gorm:"primaryKey"`
	Name           string           `json:"name" gorm:"not null"`
	Description    string           `json:"description"`
	Price          float64          `json:"price" gorm:"not null;default:0"`
	CostPrice      float64          `json:"cost_price" gorm:"not null;default:0"`
	MenuCategoryID uint             `json:"menu_category_id"`
	MenuCategory   MenuCategory     `json:"menu_category" gorm:"foreignKey:MenuCategoryID"`
	SortOrder      int              `json:"sort_order" gorm:"default:0"`
	ImageURL       string           `json:"image_url"`
	Active         bool             `json:"active" gorm:"default:true"`
	Ingredients    []MenuIngredient `json:"ingredients" gorm:"foreignKey:MenuItemID"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type MenuIngredient struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	MenuItemID   uint       `json:"menu_item_id" gorm:"not null"`
	IngredientID uint       `json:"ingredient_id" gorm:"not null"`
	Ingredient   Ingredient `json:"ingredient" gorm:"foreignKey:IngredientID"`
	Quantity     float64    `json:"quantity" gorm:"not null"`
}

type SaleType string

const (
	SaleTypeDineIn   SaleType = "dine_in"
	SaleTypeDelivery SaleType = "delivery"
)

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "cash"
	PaymentTransfer PaymentMethod = "transfer"
)

type Sale struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	Items         []SaleItem    `json:"items" gorm:"foreignKey:SaleID"`
	Total         float64       `json:"total" gorm:"not null;default:0"`
	TotalCost     float64       `json:"total_cost" gorm:"default:0"`
	Profit        float64       `json:"profit" gorm:"default:0"`
	SaleType      SaleType      `json:"sale_type" gorm:"not null;default:dine_in"`
	PaymentMethod PaymentMethod `json:"payment_method" gorm:"not null;default:cash"`
	Note          string        `json:"note"`
	UserID        uint          `json:"user_id"`
	User          User          `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time     `json:"created_at"`
}

type SaleItem struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	SaleID     uint     `json:"sale_id" gorm:"not null"`
	MenuItemID uint     `json:"menu_item_id" gorm:"not null"`
	MenuItem   MenuItem `json:"menu_item" gorm:"foreignKey:MenuItemID"`
	Quantity   int      `json:"quantity" gorm:"not null;default:1"`
	Price      float64  `json:"price" gorm:"not null"`
	Cost       float64  `json:"cost" gorm:"default:0"`
}

// ==================== Deleted Sales (Front-Office archive) ====================

type DeletedSale struct {
	ID                uint          `json:"id" gorm:"primaryKey"`
	OriginalID        uint          `json:"original_id"`
	Total             float64       `json:"total"`
	TotalCost         float64       `json:"total_cost"`
	Profit            float64       `json:"profit"`
	SaleType          SaleType      `json:"sale_type"`
	PaymentMethod     PaymentMethod `json:"payment_method"`
	Note              string        `json:"note"`
	ItemsSummary      string        `json:"items_summary"`
	UserID            uint          `json:"user_id"`
	User              User          `json:"user" gorm:"foreignKey:UserID"`
	OriginalCreatedAt time.Time     `json:"original_created_at"`
	DeleteReason      string        `json:"delete_reason"`
	DeletedByID       uint          `json:"deleted_by_id"`
	DeletedBy         User          `json:"deleted_by" gorm:"foreignKey:DeletedByID"`
	DeletedAt         time.Time     `json:"deleted_at"`
}

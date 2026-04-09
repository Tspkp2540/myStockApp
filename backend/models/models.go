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

package food_beverages

import (
	"concierge-be/internal/food_beverages_categories"
	"time"

	"gorm.io/gorm"
)

// FoodBeverage represents a food & beverage menu item (tenant-scoped)
type FoodBeverage struct {
	ID                string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TenantID          string    `gorm:"type:varchar(36);not null;index" json:"tenantId"`
	CategoryID        string    `gorm:"type:varchar(36);not null;index" json:"categoryId"`
	ItemName          string    `gorm:"type:varchar(100);not null" json:"itemName"`
	Description       string    `gorm:"type:text" json:"description"`
	Price             float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	PreparationTime   int       `gorm:"type:int;not null" json:"preparationTime"` // in minutes
	ServiceHoursStart *string   `gorm:"type:time" json:"serviceHoursStart"`
	ServiceHoursEnd   *string   `gorm:"type:time" json:"serviceHoursEnd"`
	AllDay            bool      `gorm:"default:false" json:"allDay"`
	ImagePath         string    `gorm:"type:varchar(500)" json:"imagePath"`
	IsAvailable       bool      `gorm:"default:true" json:"isAvailable"` // Available for ordering
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Category *food_beverages_categories.FoodBeverageCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

func (FoodBeverage) TableName() string {
	return "food_beverages"
}

// CreateFoodBeverageRequest represents the request body for creating a food & beverage item
type CreateFoodBeverageRequest struct {
	TenantID          string  `json:"tenantId" binding:"required"`
	CategoryID        string  `json:"categoryId" binding:"required"`
	ItemName          string  `json:"itemName" binding:"required"`
	Description       string  `json:"description"`
	Price             float64 `json:"price" binding:"required,min=0"`
	PreparationTime   int     `json:"preparationTime" binding:"required,min=1"`
	ServiceHoursStart *string `json:"serviceHoursStart"`
	ServiceHoursEnd   *string `json:"serviceHoursEnd"`
	AllDay            *bool   `json:"allDay"`
	ImagePath         string  `json:"imagePath"`
	IsAvailable       *bool   `json:"isAvailable"`
}

// UpdateFoodBeverageRequest represents the request body for updating a food & beverage item
type UpdateFoodBeverageRequest struct {
	CategoryID        string  `json:"categoryId"`
	ItemName          string  `json:"itemName"`
	Description       string  `json:"description"`
	Price             float64 `json:"price"`
	PreparationTime   int     `json:"preparationTime"`
	ServiceHoursStart *string `json:"serviceHoursStart"`
	ServiceHoursEnd   *string `json:"serviceHoursEnd"`
	AllDay            *bool   `json:"allDay"`
	ImagePath         string  `json:"imagePath"`
	IsAvailable       *bool   `json:"isAvailable"`
}

// FoodBeverageWithCategory represents a food & beverage item with its category information
type FoodBeverageWithCategory struct {
	FoodBeverage
	CategoryName string `json:"categoryName"`
}

package food_beverages_categories

import (
	"time"

	"gorm.io/gorm"
)

// FoodBeverageCategory represents a food & beverage category (global, not tenant-scoped)
type FoodBeverageCategory struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"type:varchar(255)" json:"icon"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FoodBeverageCategory) TableName() string {
	return "food_beverages_categories"
}

// CreateFoodBeverageCategoryRequest represents the request body for creating a food & beverage category
type CreateFoodBeverageCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// UpdateFoodBeverageCategoryRequest represents the request body for updating a food & beverage category
type UpdateFoodBeverageCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

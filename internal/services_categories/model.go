package services_categories

import (
	"time"

	"gorm.io/gorm"
)

// ServiceCategory represents a service category (global, not tenant-scoped)
type ServiceCategory struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"type:varchar(255)" json:"icon"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ServiceCategory) TableName() string {
	return "services_categories"
}

// CreateServiceCategoryRequest represents the request body for creating a service category
type CreateServiceCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// UpdateServiceCategoryRequest represents the request body for updating a service category
type UpdateServiceCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

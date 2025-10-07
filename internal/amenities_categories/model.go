package amenities_categories

import (
	"time"

	"gorm.io/gorm"
)

// AmenityCategory represents a category for amenities (tenant-scoped)
type AmenityCategory struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TenantID    string    `gorm:"type:varchar(36);not null;index" json:"tenantId"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AmenityCategory) TableName() string {
	return "amenities_categories"
}

// CreateAmenityCategoryRequest represents the request body for creating a category
type CreateAmenityCategoryRequest struct {
	TenantID    string `json:"tenantId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateAmenityCategoryRequest represents the request body for updating a category
type UpdateAmenityCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}


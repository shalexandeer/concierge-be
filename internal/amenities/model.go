package amenities

import (
	"concierge-be/internal/amenities_categories"
	"time"

	"gorm.io/gorm"
)

// Amenity represents an amenity item (tenant-scoped)
type Amenity struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TenantID     string    `gorm:"type:varchar(36);not null;index" json:"tenantId"`
	CategoryID   string    `gorm:"type:varchar(36);not null;index" json:"categoryId"`
	ItemName     string    `gorm:"type:varchar(100);not null" json:"itemName"`
	Description  string    `gorm:"type:text" json:"description"`
	Stock        int       `gorm:"default:0;not null" json:"stock"`
	MinimumStock int       `gorm:"default:0;not null" json:"minimumStock"`
	Available    bool      `gorm:"default:true" json:"available"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Category *amenities_categories.AmenityCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

func (Amenity) TableName() string {
	return "amenities"
}

// CreateAmenityRequest represents the request body for creating an amenity
type CreateAmenityRequest struct {
	TenantID     string `json:"tenantId" binding:"required"`
	CategoryID   string `json:"categoryId" binding:"required"`
	ItemName     string `json:"itemName" binding:"required"`
	Description  string `json:"description"`
	Stock        int    `json:"stock"`
	MinimumStock int    `json:"minimumStock"`
	Available    *bool  `json:"available"`
}

// UpdateAmenityRequest represents the request body for updating an amenity
type UpdateAmenityRequest struct {
	CategoryID   string `json:"categoryId"`
	ItemName     string `json:"itemName"`
	Description  string `json:"description"`
	Stock        *int   `json:"stock"`
	MinimumStock *int   `json:"minimumStock"`
	Available    *bool  `json:"available"`
}

// AmenityWithCategory represents an amenity with its category information
type AmenityWithCategory struct {
	Amenity
	CategoryName string `json:"categoryName"`
}


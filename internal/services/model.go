package services

import (
	"concierge-be/internal/services_categories"
	"time"

	"gorm.io/gorm"
)

// Service represents a service item (tenant-scoped)
type Service struct {
	ID                string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TenantID          string    `gorm:"type:varchar(36);not null;index" json:"tenantId"`
	CategoryID        string    `gorm:"type:varchar(36);not null;index" json:"categoryId"`
	ServiceName       string    `gorm:"type:varchar(100);not null" json:"serviceName"`
	Description       string    `gorm:"type:text" json:"description"`
	OperatingHoursFrom *string  `gorm:"type:time" json:"operatingHoursFrom"`
	OperatingHoursTo   *string  `gorm:"type:time" json:"operatingHoursTo"`
	Available24_7     bool      `gorm:"default:false" json:"available24_7"`
	ResponseTime      *int      `gorm:"type:int" json:"responseTime"`
	IsActive          bool      `gorm:"default:true" json:"isActive"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Category *services_categories.ServiceCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

func (Service) TableName() string {
	return "services"
}

// CreateServiceRequest represents the request body for creating a service
type CreateServiceRequest struct {
	TenantID          string `json:"tenantId" binding:"required"`
	CategoryID        string `json:"categoryId" binding:"required"`
	ServiceName       string `json:"serviceName" binding:"required"`
	Description       string `json:"description"`
	OperatingHoursFrom *string `json:"operatingHoursFrom"`
	OperatingHoursTo   *string `json:"operatingHoursTo"`
	Available24_7     *bool   `json:"available24_7"`
	ResponseTime      *int    `json:"responseTime"`
	IsActive          *bool   `json:"isActive"`
}

// UpdateServiceRequest represents the request body for updating a service
type UpdateServiceRequest struct {
	CategoryID        string  `json:"categoryId"`
	ServiceName       string  `json:"serviceName"`
	Description       string  `json:"description"`
	OperatingHoursFrom *string `json:"operatingHoursFrom"`
	OperatingHoursTo   *string `json:"operatingHoursTo"`
	Available24_7     *bool   `json:"available24_7"`
	ResponseTime      *int    `json:"responseTime"`
	IsActive          *bool   `json:"isActive"`
}

// ServiceWithCategory represents a service with its category information
type ServiceWithCategory struct {
	Service
	CategoryName string `json:"categoryName"`
}

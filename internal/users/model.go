package users

import (
	"time"

	"gorm.io/gorm"
)

// User represents a global user account (not tenant-scoped)
type User struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName  string    `gorm:"type:varchar(100)" json:"fullName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// UserTenant represents the many-to-many relationship between users and tenants
type UserTenant struct {
	ID       string `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID   string `gorm:"type:varchar(36);not null;index" json:"userId"`
	TenantID string `gorm:"type:varchar(36);not null;index" json:"tenantId"`
	Role     string `gorm:"type:varchar(50);default:'member'" json:"role"`
	
	// Foreign key relationships
	User   User   `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID" json:"tenant,omitempty"`
	
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (UserTenant) TableName() string {
	return "user_tenants"
}

// Tenant represents a tenant in the system
type Tenant struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Domain      string    `gorm:"type:varchar(100);uniqueIndex" json:"domain"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tenant) TableName() string {
	return "tenants"
}

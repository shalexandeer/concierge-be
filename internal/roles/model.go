package roles

import (
	"time"
	"gorm.io/gorm"
)

// Role represents a user role in the system
type Role struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

// Role constants
const (
	ROLE_SUPER_ADMIN  = "super_admin"
	ROLE_TENANT_ADMIN = "tenant_admin"
	ROLE_USER         = "user"
)

// GetDefaultRole returns the default role for new users
func GetDefaultRole() string {
	return ROLE_USER
}

// IsValidRole checks if a role name is valid
func IsValidRole(roleName string) bool {
	switch roleName {
	case ROLE_SUPER_ADMIN, ROLE_TENANT_ADMIN, ROLE_USER:
		return true
	default:
		return false
	}
}

// GetRoleHierarchy returns the hierarchy level of a role (higher number = more permissions)
func GetRoleHierarchy(roleName string) int {
	switch roleName {
	case ROLE_SUPER_ADMIN:
		return 3
	case ROLE_TENANT_ADMIN:
		return 2
	case ROLE_USER:
		return 1
	default:
		return 0
	}
}

// CanManageRole checks if a role can manage another role
func CanManageRole(managerRole, targetRole string) bool {
	managerLevel := GetRoleHierarchy(managerRole)
	targetLevel := GetRoleHierarchy(targetRole)
	
	// Can only manage roles with lower hierarchy level
	return managerLevel > targetLevel
}

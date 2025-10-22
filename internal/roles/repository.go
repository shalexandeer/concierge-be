package roles

import (
	"errors"
	"concierge-be/database"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

// CreateRole creates a new role
func (r *Repository) CreateRole(role *Role) error {
	return r.db.Create(role).Error
}

// GetRoleByID gets a role by ID
func (r *Repository) GetRoleByID(id string) (*Role, error) {
	var role Role
	err := r.db.First(&role, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// GetRoleByName gets a role by name
func (r *Repository) GetRoleByName(name string) (*Role, error) {
	var role Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// GetAllRoles gets all roles
func (r *Repository) GetAllRoles() ([]Role, error) {
	var roles []Role
	err := r.db.Find(&roles).Error
	return roles, err
}

// UpdateRole updates a role
func (r *Repository) UpdateRole(role *Role) error {
	return r.db.Save(role).Error
}

// DeleteRole deletes a role
func (r *Repository) DeleteRole(id string) error {
	return r.db.Delete(&Role{}, "id = ?", id).Error
}

// SeedRoles creates default roles if they don't exist
func (r *Repository) SeedRoles() error {
	// Check if roles already exist
	var count int64
	r.db.Model(&Role{}).Count(&count)
	if count > 0 {
		return nil // Roles already seeded
	}

	// Create default roles
	roles := []Role{
		{
			ID:          "role-" + generateRandomString(32),
			Name:        ROLE_SUPER_ADMIN,
			Description: "Super administrator with full access to all tenants and features",
		},
		{
			ID:          "role-" + generateRandomString(32),
			Name:        ROLE_TENANT_ADMIN,
			Description: "Tenant administrator with management access within their tenant",
		},
		{
			ID:          "role-" + generateRandomString(32),
			Name:        ROLE_USER,
			Description: "Regular user with view and booking access",
		},
	}

	for _, role := range roles {
		if err := r.CreateRole(&role); err != nil {
			return err
		}
	}

	return nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}

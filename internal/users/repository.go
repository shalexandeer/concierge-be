package users

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

// User repository methods
func (r *Repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	var user User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetAllUsers(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	db := r.db.Model(&User{})

	// Get total count
	db.Count(&total)

	// Paginated query
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

func (r *Repository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r *Repository) DeleteUser(id string) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}

// UserTenant repository methods
func (r *Repository) CreateUserTenant(userTenant *UserTenant) error {
	return r.db.Create(userTenant).Error
}

func (r *Repository) GetUserTenants(userID string) ([]UserTenant, error) {
	var userTenants []UserTenant
	err := r.db.Where("user_id = ?", userID).Preload("Tenant").Find(&userTenants).Error
	return userTenants, err
}

func (r *Repository) GetTenantUsers(tenantID string) ([]UserTenant, error) {
	var userTenants []UserTenant
	err := r.db.Where("tenant_id = ?", tenantID).Preload("User").Find(&userTenants).Error
	return userTenants, err
}

func (r *Repository) GetUserTenant(userID, tenantID string) (*UserTenant, error) {
	var userTenant UserTenant
	err := r.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&userTenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user-tenant relationship not found")
		}
		return nil, err
	}
	return &userTenant, nil
}

func (r *Repository) UpdateUserTenant(userTenant *UserTenant) error {
	return r.db.Save(userTenant).Error
}

func (r *Repository) DeleteUserTenant(userID, tenantID string) error {
	return r.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Delete(&UserTenant{}).Error
}

// Tenant repository methods
func (r *Repository) CreateTenant(tenant *Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *Repository) GetTenantByID(id string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.First(&tenant, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) GetTenantByDomain(domain string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.Where("domain = ?", domain).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) GetAllTenants(page, pageSize int) ([]Tenant, int64, error) {
	var tenants []Tenant
	var total int64

	db := r.db.Model(&Tenant{})

	// Get total count
	db.Count(&total)

	// Paginated query
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&tenants).Error

	return tenants, total, err
}

func (r *Repository) UpdateTenant(tenant *Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *Repository) DeleteTenant(id string) error {
	return r.db.Delete(&Tenant{}, "id = ?", id).Error
}

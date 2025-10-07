package tenants

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

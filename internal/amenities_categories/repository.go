package amenities_categories

import (
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

// Create creates a new amenity category
func (r *Repository) Create(category *AmenityCategory) error {
	return r.db.Create(category).Error
}

// GetByID retrieves an amenity category by ID
func (r *Repository) GetByID(id string) (*AmenityCategory, error) {
	var category AmenityCategory
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByTenantID retrieves all amenity categories for a specific tenant
func (r *Repository) GetByTenantID(tenantID string) ([]AmenityCategory, error) {
	var categories []AmenityCategory
	err := r.db.Where("tenant_id = ?", tenantID).Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetAll retrieves all amenity categories
func (r *Repository) GetAll() ([]AmenityCategory, error) {
	var categories []AmenityCategory
	err := r.db.Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Update updates an existing amenity category
func (r *Repository) Update(category *AmenityCategory) error {
	return r.db.Save(category).Error
}

// Delete soft deletes an amenity category
func (r *Repository) Delete(id string) error {
	return r.db.Delete(&AmenityCategory{}, "id = ?", id).Error
}

// CheckNameExists checks if a category name already exists for a tenant
func (r *Repository) CheckNameExists(tenantID, name, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&AmenityCategory{}).Where("tenant_id = ? AND name = ?", tenantID, name)
	
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}


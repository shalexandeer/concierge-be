package services_categories

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

// CreateServiceCategory creates a new service category
func (r *Repository) CreateServiceCategory(category *ServiceCategory) error {
	return r.db.Create(category).Error
}

// GetServiceCategoryByID gets a service category by ID
func (r *Repository) GetServiceCategoryByID(id string) (*ServiceCategory, error) {
	var category ServiceCategory
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllServiceCategories gets all service categories
func (r *Repository) GetAllServiceCategories() ([]ServiceCategory, error) {
	var categories []ServiceCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

// UpdateServiceCategory updates a service category
func (r *Repository) UpdateServiceCategory(id string, category *ServiceCategory) error {
	return r.db.Model(&ServiceCategory{}).Where("id = ?", id).Updates(category).Error
}

// DeleteServiceCategory soft deletes a service category
func (r *Repository) DeleteServiceCategory(id string) error {
	return r.db.Where("id = ?", id).Delete(&ServiceCategory{}).Error
}

// CheckNameExists checks if a category name already exists (excluding the given ID)
func (r *Repository) CheckNameExists(name string, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&ServiceCategory{}).Where("name = ?", name)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

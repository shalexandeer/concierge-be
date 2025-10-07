package amenities

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

// Create creates a new amenity
func (r *Repository) Create(amenity *Amenity) error {
	return r.db.Create(amenity).Error
}

// GetByID retrieves an amenity by ID with category preloaded
func (r *Repository) GetByID(id string) (*Amenity, error) {
	var amenity Amenity
	err := r.db.Preload("Category").Where("id = ?", id).First(&amenity).Error
	if err != nil {
		return nil, err
	}
	return &amenity, nil
}

// GetByTenantID retrieves all amenities for a specific tenant
func (r *Repository) GetByTenantID(tenantID string, includeCategory bool) ([]Amenity, error) {
	var amenities []Amenity
	query := r.db.Where("tenant_id = ?", tenantID)
	
	if includeCategory {
		query = query.Preload("Category")
	}
	
	err := query.Order("item_name ASC").Find(&amenities).Error
	if err != nil {
		return nil, err
	}
	return amenities, nil
}

// GetByCategoryID retrieves all amenities for a specific category
func (r *Repository) GetByCategoryID(categoryID string, includeCategory bool) ([]Amenity, error) {
	var amenities []Amenity
	query := r.db.Where("category_id = ?", categoryID)
	
	if includeCategory {
		query = query.Preload("Category")
	}
	
	err := query.Order("item_name ASC").Find(&amenities).Error
	if err != nil {
		return nil, err
	}
	return amenities, nil
}

// GetAll retrieves all amenities
func (r *Repository) GetAll(includeCategory bool) ([]Amenity, error) {
	var amenities []Amenity
	query := r.db
	
	if includeCategory {
		query = query.Preload("Category")
	}
	
	err := query.Order("item_name ASC").Find(&amenities).Error
	if err != nil {
		return nil, err
	}
	return amenities, nil
}

// GetLowStock retrieves amenities with stock below minimum for a tenant
func (r *Repository) GetLowStock(tenantID string) ([]Amenity, error) {
	var amenities []Amenity
	err := r.db.Preload("Category").
		Where("tenant_id = ? AND stock < minimum_stock", tenantID).
		Order("item_name ASC").
		Find(&amenities).Error
	if err != nil {
		return nil, err
	}
	return amenities, nil
}

// Update updates an existing amenity
func (r *Repository) Update(amenity *Amenity) error {
	return r.db.Save(amenity).Error
}

// Delete soft deletes an amenity
func (r *Repository) Delete(id string) error {
	return r.db.Delete(&Amenity{}, "id = ?", id).Error
}

// CheckItemNameExists checks if an item name already exists for a tenant
func (r *Repository) CheckItemNameExists(tenantID, itemName, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&Amenity{}).Where("tenant_id = ? AND item_name = ?", tenantID, itemName)
	
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// UpdateStock updates the stock quantity for an amenity
func (r *Repository) UpdateStock(id string, quantity int) error {
	return r.db.Model(&Amenity{}).Where("id = ?", id).Update("stock", quantity).Error
}


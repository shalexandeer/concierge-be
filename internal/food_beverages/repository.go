package food_beverages

import (
	"concierge-be/database"
	"errors"
	"time"

	"github.com/google/uuid"
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

// CreateFoodBeverage creates a new food & beverage item
func (r *Repository) CreateFoodBeverage(item *FoodBeverage) error {
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	if err := r.db.Create(item).Error; err != nil {
		return err
	}

	return nil
}

// GetFoodBeverageByID gets a food & beverage item by ID
func (r *Repository) GetFoodBeverageByID(id string) (*FoodBeverage, error) {
	var item FoodBeverage
	if err := r.db.Preload("Category").Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("food & beverage item not found")
		}
		return nil, err
	}

	return &item, nil
}

// GetAllFoodBeverages gets all food & beverage items with optional filtering
func (r *Repository) GetAllFoodBeverages(tenantID, categoryID string, page, pageSize int) ([]FoodBeverage, int64, error) {
	var items []FoodBeverage
	var total int64

	query := r.db.Preload("Category")

	// Apply tenant filter
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// Apply category filter
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Count total records
	if err := query.Model(&FoodBeverage{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetFoodBeveragesByTenant gets all food & beverage items for a specific tenant
func (r *Repository) GetFoodBeveragesByTenant(tenantID string) ([]FoodBeverage, error) {
	var items []FoodBeverage
	if err := r.db.Preload("Category").Where("tenant_id = ?", tenantID).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateFoodBeverage updates a food & beverage item
func (r *Repository) UpdateFoodBeverage(id string, req *UpdateFoodBeverageRequest) (*FoodBeverage, error) {
	var item FoodBeverage
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("food & beverage item not found")
		}
		return nil, err
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.CategoryID != "" {
		updates["category_id"] = req.CategoryID
	}
	if req.ItemName != "" {
		updates["item_name"] = req.ItemName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.PreparationTime > 0 {
		updates["preparation_time"] = req.PreparationTime
	}
	if req.ServiceHoursStart != nil {
		updates["service_hours_start"] = req.ServiceHoursStart
	}
	if req.ServiceHoursEnd != nil {
		updates["service_hours_end"] = req.ServiceHoursEnd
	}
	if req.AllDay != nil {
		updates["all_day"] = *req.AllDay
	}
	if req.ImagePath != "" {
		updates["image_path"] = req.ImagePath
	}
	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}

	updates["updated_at"] = time.Now()

	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with category information
	if err := r.db.Preload("Category").Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

// DeleteFoodBeverage soft deletes a food & beverage item
func (r *Repository) DeleteFoodBeverage(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&FoodBeverage{}).Error; err != nil {
		return err
	}

	return nil
}

// CheckFoodBeverageExists checks if a food & beverage item exists by name within a tenant
func (r *Repository) CheckFoodBeverageExists(tenantID, itemName string, excludeID ...string) (bool, error) {
	var count int64
	query := r.db.Model(&FoodBeverage{}).Where("tenant_id = ? AND item_name = ?", tenantID, itemName)
	
	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}
	
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

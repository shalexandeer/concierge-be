package food_beverages_categories

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

// CreateCategory creates a new food & beverage category
func (r *Repository) CreateCategory(category *FoodBeverageCategory) error {
	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	if err := r.db.Create(category).Error; err != nil {
		return err
	}

	return nil
}

// GetCategoryByID gets a food & beverage category by ID
func (r *Repository) GetCategoryByID(id string) (*FoodBeverageCategory, error) {
	var category FoodBeverageCategory
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

// GetAllCategories gets all food & beverage categories
func (r *Repository) GetAllCategories() ([]FoodBeverageCategory, error) {
	var categories []FoodBeverageCategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// UpdateCategory updates a food & beverage category
func (r *Repository) UpdateCategory(id string, req *UpdateFoodBeverageCategoryRequest) (*FoodBeverageCategory, error) {
	var category FoodBeverageCategory
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}

	updates["updated_at"] = time.Now()

	if err := r.db.Model(&category).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

// DeleteCategory soft deletes a food & beverage category
func (r *Repository) DeleteCategory(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&FoodBeverageCategory{}).Error; err != nil {
		return err
	}

	return nil
}

// CheckCategoryExists checks if a category exists by name
func (r *Repository) CheckCategoryExists(name string, excludeID ...string) (bool, error) {
	var count int64
	query := r.db.Model(&FoodBeverageCategory{}).Where("name = ?", name)
	
	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}
	
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

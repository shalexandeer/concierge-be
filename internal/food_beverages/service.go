package food_beverages

import (
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// CreateFoodBeverage creates a new food & beverage item
func (s *Service) CreateFoodBeverage(req *CreateFoodBeverageRequest) (*FoodBeverage, error) {
	// Check if item name already exists within the tenant
	exists, err := s.repo.CheckFoodBeverageExists(req.TenantID, req.ItemName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("food & beverage item with this name already exists in this tenant")
	}

	// Validate service hours logic
	if req.AllDay != nil && *req.AllDay {
		// If all day, clear specific hours
		req.ServiceHoursStart = nil
		req.ServiceHoursEnd = nil
	}

	item := &FoodBeverage{
		TenantID:          req.TenantID,
		CategoryID:        req.CategoryID,
		ItemName:          req.ItemName,
		Description:       req.Description,
		Price:             req.Price,
		PreparationTime:   req.PreparationTime,
		ServiceHoursStart: req.ServiceHoursStart,
		ServiceHoursEnd:   req.ServiceHoursEnd,
		AllDay:            req.AllDay != nil && *req.AllDay,
		ImagePath:         req.ImagePath,
		IsAvailable:       req.IsAvailable != nil && *req.IsAvailable,
	}

	if err := s.repo.CreateFoodBeverage(item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetFoodBeverageByID gets a food & beverage item by ID
func (s *Service) GetFoodBeverageByID(id string) (*FoodBeverage, error) {
	return s.repo.GetFoodBeverageByID(id)
}

// GetAllFoodBeverages gets all food & beverage items with optional filtering
func (s *Service) GetAllFoodBeverages(tenantID, categoryID string, page, pageSize int) ([]FoodBeverage, int64, error) {
	return s.repo.GetAllFoodBeverages(tenantID, categoryID, page, pageSize)
}

// GetFoodBeveragesByTenant gets all food & beverage items for a specific tenant
func (s *Service) GetFoodBeveragesByTenant(tenantID string) ([]FoodBeverage, error) {
	return s.repo.GetFoodBeveragesByTenant(tenantID)
}

// UpdateFoodBeverage updates a food & beverage item
func (s *Service) UpdateFoodBeverage(id string, req *UpdateFoodBeverageRequest) (*FoodBeverage, error) {
	// Check if item exists
	existingItem, err := s.repo.GetFoodBeverageByID(id)
	if err != nil {
		return nil, err
	}

	// Check if new name conflicts with existing items in the same tenant
	if req.ItemName != "" && req.ItemName != existingItem.ItemName {
		exists, err := s.repo.CheckFoodBeverageExists(existingItem.TenantID, req.ItemName, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("food & beverage item with this name already exists in this tenant")
		}
	}

	// Validate service hours logic
	if req.AllDay != nil && *req.AllDay {
		// If all day, clear specific hours
		req.ServiceHoursStart = nil
		req.ServiceHoursEnd = nil
	}

	return s.repo.UpdateFoodBeverage(id, req)
}

// DeleteFoodBeverage deletes a food & beverage item
func (s *Service) DeleteFoodBeverage(id string) error {
	// Check if item exists
	_, err := s.repo.GetFoodBeverageByID(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteFoodBeverage(id)
}

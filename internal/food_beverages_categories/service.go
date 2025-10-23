package food_beverages_categories

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

// CreateCategory creates a new food & beverage category
func (s *Service) CreateCategory(req *CreateFoodBeverageCategoryRequest) (*FoodBeverageCategory, error) {
	// Check if category name already exists
	exists, err := s.repo.CheckCategoryExists(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category with this name already exists")
	}

	category := &FoodBeverageCategory{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}

	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID gets a food & beverage category by ID
func (s *Service) GetCategoryByID(id string) (*FoodBeverageCategory, error) {
	return s.repo.GetCategoryByID(id)
}

// GetAllCategories gets all food & beverage categories
func (s *Service) GetAllCategories() ([]FoodBeverageCategory, error) {
	return s.repo.GetAllCategories()
}

// UpdateCategory updates a food & beverage category
func (s *Service) UpdateCategory(id string, req *UpdateFoodBeverageCategoryRequest) (*FoodBeverageCategory, error) {
	// Check if category exists
	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	// Check if new name conflicts with existing categories
	if req.Name != "" {
		exists, err := s.repo.CheckCategoryExists(req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("category with this name already exists")
		}
	}

	return s.repo.UpdateCategory(id, req)
}

// DeleteCategory deletes a food & beverage category
func (s *Service) DeleteCategory(id string) error {
	// Check if category exists
	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return err
	}

	// TODO: Check if category is being used by any food & beverage items
	// If yes, prevent deletion or cascade delete

	return s.repo.DeleteCategory(id)
}

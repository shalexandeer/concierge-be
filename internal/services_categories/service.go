package services_categories

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ServiceCategoryService struct {
	repo *Repository
}

func NewService() *ServiceCategoryService {
	return &ServiceCategoryService{
		repo: NewRepository(),
	}
}

// CreateServiceCategory creates a new service category
func (s *ServiceCategoryService) CreateServiceCategory(req *CreateServiceCategoryRequest) (*ServiceCategory, error) {
	// Check if name already exists
	exists, err := s.repo.CheckNameExists(req.Name, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category name already exists")
	}

	category := &ServiceCategory{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateServiceCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetServiceCategoryByID gets a service category by ID
func (s *ServiceCategoryService) GetServiceCategoryByID(id string) (*ServiceCategory, error) {
	return s.repo.GetServiceCategoryByID(id)
}

// GetAllServiceCategories gets all service categories
func (s *ServiceCategoryService) GetAllServiceCategories() ([]ServiceCategory, error) {
	return s.repo.GetAllServiceCategories()
}

// UpdateServiceCategory updates a service category
func (s *ServiceCategoryService) UpdateServiceCategory(id string, req *UpdateServiceCategoryRequest) (*ServiceCategory, error) {
	// Check if name already exists (excluding current category)
	exists, err := s.repo.CheckNameExists(req.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category name already exists")
	}

	// Get existing category
	category, err := s.repo.GetServiceCategoryByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	category.UpdatedAt = time.Now()

	if err := s.repo.UpdateServiceCategory(id, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteServiceCategory deletes a service category
func (s *ServiceCategoryService) DeleteServiceCategory(id string) error {
	// Check if category exists
	_, err := s.repo.GetServiceCategoryByID(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteServiceCategory(id)
}

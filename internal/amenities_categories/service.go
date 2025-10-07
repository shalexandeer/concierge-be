package amenities_categories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// CreateCategory creates a new amenity category
func (s *Service) CreateCategory(req *CreateAmenityCategoryRequest) (*AmenityCategory, error) {
	// Check if category name already exists for this tenant
	exists, err := s.repo.CheckNameExists(req.TenantID, req.Name, "")
	if err != nil {
		return nil, fmt.Errorf("failed to check category name: %w", err)
	}
	if exists {
		return nil, errors.New("category name already exists for this tenant")
	}

	category := &AmenityCategory{
		ID:          uuid.New().String(),
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetCategoryByID retrieves a category by ID
func (s *Service) GetCategoryByID(id string) (*AmenityCategory, error) {
	return s.repo.GetByID(id)
}

// GetCategoriesByTenantID retrieves all categories for a tenant
func (s *Service) GetCategoriesByTenantID(tenantID string) ([]AmenityCategory, error) {
	return s.repo.GetByTenantID(tenantID)
}

// GetAllCategories retrieves all categories
func (s *Service) GetAllCategories() ([]AmenityCategory, error) {
	return s.repo.GetAll()
}

// UpdateCategory updates an existing category
func (s *Service) UpdateCategory(id string, req *UpdateAmenityCategoryRequest) (*AmenityCategory, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// If name is being updated, check for duplicates
	if req.Name != "" && req.Name != category.Name {
		exists, err := s.repo.CheckNameExists(category.TenantID, req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check category name: %w", err)
		}
		if exists {
			return nil, errors.New("category name already exists for this tenant")
		}
		category.Name = req.Name
	}

	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.repo.Update(category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *Service) DeleteCategory(id string) error {
	// Check if category exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}


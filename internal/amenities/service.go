package amenities

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

// CreateAmenity creates a new amenity
func (s *Service) CreateAmenity(req *CreateAmenityRequest) (*Amenity, error) {
	// Check if item name already exists for this tenant
	exists, err := s.repo.CheckItemNameExists(req.TenantID, req.ItemName, "")
	if err != nil {
		return nil, fmt.Errorf("failed to check item name: %w", err)
	}
	if exists {
		return nil, errors.New("item name already exists for this tenant")
	}

	available := true
	if req.Available != nil {
		available = *req.Available
	}

	amenity := &Amenity{
		ID:           uuid.New().String(),
		TenantID:     req.TenantID,
		CategoryID:   req.CategoryID,
		ItemName:     req.ItemName,
		Description:  req.Description,
		Stock:        req.Stock,
		MinimumStock: req.MinimumStock,
		Available:    available,
	}

	if err := s.repo.Create(amenity); err != nil {
		return nil, fmt.Errorf("failed to create amenity: %w", err)
	}

	// Reload with category
	return s.repo.GetByID(amenity.ID)
}

// GetAmenityByID retrieves an amenity by ID
func (s *Service) GetAmenityByID(id string) (*Amenity, error) {
	return s.repo.GetByID(id)
}

// GetAmenitiesByTenantID retrieves all amenities for a tenant
func (s *Service) GetAmenitiesByTenantID(tenantID string) ([]Amenity, error) {
	return s.repo.GetByTenantID(tenantID, true)
}

// GetAmenitiesByCategoryID retrieves all amenities for a category
func (s *Service) GetAmenitiesByCategoryID(categoryID string) ([]Amenity, error) {
	return s.repo.GetByCategoryID(categoryID, true)
}

// GetAllAmenities retrieves all amenities
func (s *Service) GetAllAmenities() ([]Amenity, error) {
	return s.repo.GetAll(true)
}

// GetLowStockAmenities retrieves amenities with stock below minimum
func (s *Service) GetLowStockAmenities(tenantID string) ([]Amenity, error) {
	return s.repo.GetLowStock(tenantID)
}

// UpdateAmenity updates an existing amenity
func (s *Service) UpdateAmenity(id string, req *UpdateAmenityRequest) (*Amenity, error) {
	amenity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("amenity not found: %w", err)
	}

	// If item name is being updated, check for duplicates
	if req.ItemName != "" && req.ItemName != amenity.ItemName {
		exists, err := s.repo.CheckItemNameExists(amenity.TenantID, req.ItemName, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check item name: %w", err)
		}
		if exists {
			return nil, errors.New("item name already exists for this tenant")
		}
		amenity.ItemName = req.ItemName
	}

	if req.CategoryID != "" {
		amenity.CategoryID = req.CategoryID
	}

	if req.Description != "" {
		amenity.Description = req.Description
	}

	if req.Stock != nil {
		amenity.Stock = *req.Stock
	}

	if req.MinimumStock != nil {
		amenity.MinimumStock = *req.MinimumStock
	}

	if req.Available != nil {
		amenity.Available = *req.Available
	}

	if err := s.repo.Update(amenity); err != nil {
		return nil, fmt.Errorf("failed to update amenity: %w", err)
	}

	// Reload with category
	return s.repo.GetByID(id)
}

// UpdateStock updates the stock quantity for an amenity
func (s *Service) UpdateStock(id string, quantity int) (*Amenity, error) {
	if quantity < 0 {
		return nil, errors.New("stock quantity cannot be negative")
	}

	// Check if amenity exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("amenity not found: %w", err)
	}

	if err := s.repo.UpdateStock(id, quantity); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	// Reload with category
	return s.repo.GetByID(id)
}

// DeleteAmenity deletes an amenity
func (s *Service) DeleteAmenity(id string) error {
	// Check if amenity exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("amenity not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete amenity: %w", err)
	}

	return nil
}


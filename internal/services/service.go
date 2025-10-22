package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ServiceService struct {
	repo *Repository
}

func NewService() *ServiceService {
	return &ServiceService{
		repo: NewRepository(),
	}
}

// CreateService creates a new service
func (s *ServiceService) CreateService(req *CreateServiceRequest) (*Service, error) {
	// Check if service name already exists for this tenant
	exists, err := s.repo.CheckServiceNameExists(req.ServiceName, req.TenantID, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("service name already exists for this tenant")
	}

	// Handle 24/7 availability logic
	var operatingHoursFrom, operatingHoursTo *string
	if req.Available24_7 != nil && *req.Available24_7 {
		from := "00:00"
		to := "23:59"
		operatingHoursFrom = &from
		operatingHoursTo = &to
	} else {
		operatingHoursFrom = req.OperatingHoursFrom
		operatingHoursTo = req.OperatingHoursTo
	}

	service := &Service{
		ID:                uuid.New().String(),
		TenantID:          req.TenantID,
		CategoryID:        req.CategoryID,
		ServiceName:       req.ServiceName,
		Description:       req.Description,
		OperatingHoursFrom: operatingHoursFrom,
		OperatingHoursTo:   operatingHoursTo,
		Available24_7:     req.Available24_7 != nil && *req.Available24_7,
		ResponseTime:      req.ResponseTime,
		IsActive:          req.IsActive != nil && *req.IsActive,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.CreateService(service); err != nil {
		return nil, err
	}

	return service, nil
}

// GetServiceByID gets a service by ID
func (s *ServiceService) GetServiceByID(id string) (*Service, error) {
	return s.repo.GetServiceByID(id)
}

// GetAllServices gets all services (for super admin)
func (s *ServiceService) GetAllServices() ([]Service, error) {
	return s.repo.GetAllServices()
}

// GetServicesByTenantID gets services by tenant ID
func (s *ServiceService) GetServicesByTenantID(tenantID string) ([]Service, error) {
	return s.repo.GetServicesByTenantID(tenantID)
}

// GetServicesByTenantIDs gets services by multiple tenant IDs
func (s *ServiceService) GetServicesByTenantIDs(tenantIDs []string) ([]Service, error) {
	return s.repo.GetServicesByTenantIDs(tenantIDs)
}

// GetServicesByCategoryID gets services by category ID
func (s *ServiceService) GetServicesByCategoryID(categoryID string) ([]Service, error) {
	return s.repo.GetServicesByCategoryID(categoryID)
}

// GetServicesByCategoryIDAndTenantIDs gets services by category ID and tenant IDs
func (s *ServiceService) GetServicesByCategoryIDAndTenantIDs(categoryID string, tenantIDs []string) ([]Service, error) {
	return s.repo.GetServicesByCategoryIDAndTenantIDs(categoryID, tenantIDs)
}

// UpdateService updates a service
func (s *ServiceService) UpdateService(id string, req *UpdateServiceRequest) (*Service, error) {
	// Get existing service to check tenant
	existingService, err := s.repo.GetServiceByID(id)
	if err != nil {
		return nil, err
	}

	// Check if service name already exists for this tenant (excluding current service)
	if req.ServiceName != "" {
		exists, err := s.repo.CheckServiceNameExists(req.ServiceName, existingService.TenantID, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("service name already exists for this tenant")
		}
	}

	// Handle 24/7 availability logic
	var operatingHoursFrom, operatingHoursTo *string
	if req.Available24_7 != nil && *req.Available24_7 {
		from := "00:00"
		to := "23:59"
		operatingHoursFrom = &from
		operatingHoursTo = &to
	} else if req.OperatingHoursFrom != nil && req.OperatingHoursTo != nil {
		operatingHoursFrom = req.OperatingHoursFrom
		operatingHoursTo = req.OperatingHoursTo
	}

	// Update fields
	if req.CategoryID != "" {
		existingService.CategoryID = req.CategoryID
	}
	if req.ServiceName != "" {
		existingService.ServiceName = req.ServiceName
	}
	if req.Description != "" {
		existingService.Description = req.Description
	}
	if operatingHoursFrom != nil {
		existingService.OperatingHoursFrom = operatingHoursFrom
	}
	if operatingHoursTo != nil {
		existingService.OperatingHoursTo = operatingHoursTo
	}
	if req.Available24_7 != nil {
		existingService.Available24_7 = *req.Available24_7
	}
	if req.ResponseTime != nil {
		existingService.ResponseTime = req.ResponseTime
	}
	if req.IsActive != nil {
		existingService.IsActive = *req.IsActive
	}
	existingService.UpdatedAt = time.Now()

	if err := s.repo.UpdateService(id, existingService); err != nil {
		return nil, err
	}

	return existingService, nil
}

// DeleteService deletes a service
func (s *ServiceService) DeleteService(id string) error {
	// Check if service exists
	_, err := s.repo.GetServiceByID(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteService(id)
}

package services

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

// CreateService creates a new service
func (r *Repository) CreateService(service *Service) error {
	return r.db.Create(service).Error
}

// GetServiceByID gets a service by ID
func (r *Repository) GetServiceByID(id string) (*Service, error) {
	var service Service
	err := r.db.Preload("Category").Where("id = ?", id).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// GetAllServices gets all services with category information
func (r *Repository) GetAllServices() ([]Service, error) {
	var services []Service
	err := r.db.Preload("Category").Find(&services).Error
	return services, err
}

// GetServicesByTenantID gets services by tenant ID
func (r *Repository) GetServicesByTenantID(tenantID string) ([]Service, error) {
	var services []Service
	err := r.db.Preload("Category").Where("tenant_id = ?", tenantID).Find(&services).Error
	return services, err
}

// GetServicesByTenantIDs gets services by multiple tenant IDs
func (r *Repository) GetServicesByTenantIDs(tenantIDs []string) ([]Service, error) {
	var services []Service
	err := r.db.Preload("Category").Where("tenant_id IN ?", tenantIDs).Find(&services).Error
	return services, err
}

// GetServicesByCategoryID gets services by category ID
func (r *Repository) GetServicesByCategoryID(categoryID string) ([]Service, error) {
	var services []Service
	err := r.db.Preload("Category").Where("category_id = ?", categoryID).Find(&services).Error
	return services, err
}

// GetServicesByCategoryIDAndTenantIDs gets services by category ID and tenant IDs
func (r *Repository) GetServicesByCategoryIDAndTenantIDs(categoryID string, tenantIDs []string) ([]Service, error) {
	var services []Service
	err := r.db.Preload("Category").Where("category_id = ? AND tenant_id IN ?", categoryID, tenantIDs).Find(&services).Error
	return services, err
}

// UpdateService updates a service
func (r *Repository) UpdateService(id string, service *Service) error {
	return r.db.Model(&Service{}).Where("id = ?", id).Updates(service).Error
}

// DeleteService soft deletes a service
func (r *Repository) DeleteService(id string) error {
	return r.db.Where("id = ?", id).Delete(&Service{}).Error
}

// CheckServiceNameExists checks if a service name already exists for a tenant
func (r *Repository) CheckServiceNameExists(serviceName string, tenantID string, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&Service{}).Where("service_name = ? AND tenant_id = ?", serviceName, tenantID)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

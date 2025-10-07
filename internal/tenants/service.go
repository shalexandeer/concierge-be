package tenants

import (
	"crypto/rand"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// generateUUID generates a new UUID
func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Tenant service methods
func (s *Service) CreateTenant(tenant *Tenant) error {
	// Generate UUID if not set
	if tenant.ID == "" {
		tenant.ID = generateUUID()
	}
	return s.repo.CreateTenant(tenant)
}

func (s *Service) GetTenantByID(id string) (*Tenant, error) {
	return s.repo.GetTenantByID(id)
}

func (s *Service) GetTenantByDomain(domain string) (*Tenant, error) {
	return s.repo.GetTenantByDomain(domain)
}

func (s *Service) GetAllTenants(page, pageSize int) ([]Tenant, int64, error) {
	return s.repo.GetAllTenants(page, pageSize)
}

func (s *Service) UpdateTenant(tenant *Tenant) error {
	return s.repo.UpdateTenant(tenant)
}

func (s *Service) DeleteTenant(id string) error {
	return s.repo.DeleteTenant(id)
}

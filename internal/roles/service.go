package roles

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

// CreateRole creates a new role
func (s *Service) CreateRole(role *Role) error {
	// Validate role name
	if !IsValidRole(role.Name) {
		return fmt.Errorf("invalid role name: %s", role.Name)
	}

	// Check if role already exists
	existingRole, err := s.repo.GetRoleByName(role.Name)
	if err == nil && existingRole != nil {
		return fmt.Errorf("role with name %s already exists", role.Name)
	}

	// Generate UUID if not set
	if role.ID == "" {
		role.ID = generateUUID()
	}

	return s.repo.CreateRole(role)
}

// GetRoleByID gets a role by ID
func (s *Service) GetRoleByID(id string) (*Role, error) {
	return s.repo.GetRoleByID(id)
}

// GetRoleByName gets a role by name
func (s *Service) GetRoleByName(name string) (*Role, error) {
	return s.repo.GetRoleByName(name)
}

// GetAllRoles gets all roles
func (s *Service) GetAllRoles() ([]Role, error) {
	return s.repo.GetAllRoles()
}

// UpdateRole updates a role
func (s *Service) UpdateRole(role *Role) error {
	// Validate role name
	if !IsValidRole(role.Name) {
		return fmt.Errorf("invalid role name: %s", role.Name)
	}

	return s.repo.UpdateRole(role)
}

// DeleteRole deletes a role
func (s *Service) DeleteRole(id string) error {
	// Check if role exists
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return err
	}

	// Prevent deletion of system roles
	if role.Name == ROLE_SUPER_ADMIN || role.Name == ROLE_TENANT_ADMIN || role.Name == ROLE_USER {
		return fmt.Errorf("cannot delete system role: %s", role.Name)
	}

	return s.repo.DeleteRole(id)
}

// SeedRoles creates default roles if they don't exist
func (s *Service) SeedRoles() error {
	return s.repo.SeedRoles()
}

// ValidateRoleAssignment validates if a user can assign a specific role
func (s *Service) ValidateRoleAssignment(assignerRole, targetRole string) error {
	if !CanManageRole(assignerRole, targetRole) {
		return fmt.Errorf("role %s cannot assign role %s", assignerRole, targetRole)
	}
	return nil
}

// GetDefaultRole returns the default role for new users
func (s *Service) GetDefaultRole() (*Role, error) {
	return s.repo.GetRoleByName(GetDefaultRole())
}

// generateUUID generates a new UUID
func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

package users

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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

// User service methods
func (s *Service) CreateUser(user *User) error {
	// Generate UUID if not set
	if user.ID == "" {
		user.ID = generateUUID()
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.repo.CreateUser(user)
}

func (s *Service) GetUserByID(id string) (*User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) GetUserByUsername(username string) (*User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *Service) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *Service) GetAllUsers(page, pageSize int) ([]User, int64, error) {
	return s.repo.GetAllUsers(page, pageSize)
}

func (s *Service) UpdateUser(user *User) error {
	// If password is being updated, hash it
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return s.repo.UpdateUser(user)
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *Service) VerifyPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// UserTenant service methods
func (s *Service) AddUserToTenant(userID, tenantID, role string) error {
	userTenant := &UserTenant{
		ID:       generateUUID(),
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
	}
	return s.repo.CreateUserTenant(userTenant)
}

func (s *Service) GetUserTenants(userID string) ([]UserTenant, error) {
	return s.repo.GetUserTenants(userID)
}

func (s *Service) GetTenantUsers(tenantID string) ([]UserTenant, error) {
	return s.repo.GetTenantUsers(tenantID)
}

func (s *Service) GetUserTenant(userID, tenantID string) (*UserTenant, error) {
	return s.repo.GetUserTenant(userID, tenantID)
}

func (s *Service) UpdateUserTenantRole(userID, tenantID, role string) error {
	userTenant, err := s.repo.GetUserTenant(userID, tenantID)
	if err != nil {
		return err
	}
	userTenant.Role = role
	return s.repo.UpdateUserTenant(userTenant)
}

func (s *Service) RemoveUserFromTenant(userID, tenantID string) error {
	return s.repo.DeleteUserTenant(userID, tenantID)
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

// Helper methods
func (s *Service) IsUserInTenant(userID, tenantID string) (bool, error) {
	_, err := s.repo.GetUserTenant(userID, tenantID)
	if err != nil {
		if err.Error() == "user-tenant relationship not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Service) GetUserRoleInTenant(userID, tenantID string) (string, error) {
	userTenant, err := s.repo.GetUserTenant(userID, tenantID)
	if err != nil {
		return "", err
	}
	return userTenant.Role, nil
}

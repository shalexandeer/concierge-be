package users

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"concierge-be/internal/roles"
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

// CreateUserWithRole creates a user with role validation
func (s *Service) CreateUserWithRole(user *User, assignerRole string) error {
	// If no role specified, assign default user role
	if user.RoleID == "" {
		userRole, err := s.repo.GetRoleByName("user")
		if err != nil {
			return fmt.Errorf("failed to get default user role: %v", err)
		}
		user.RoleID = userRole.ID
	}

	// Get the role being assigned
	role, err := s.repo.GetRoleByID(user.RoleID)
	if err != nil {
		return fmt.Errorf("invalid role: %v", err)
	}

	// Validate role assignment
	if assignerRole != "" {
		if !CanManageRole(assignerRole, role.Name) {
			return fmt.Errorf("role %s cannot assign role %s", assignerRole, role.Name)
		}
	}

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

// AssignDefaultRoleToUsersWithoutRole assigns the default 'user' role to users who don't have a role
func (s *Service) AssignDefaultRoleToUsersWithoutRole() error {
	// Get the default user role
	userRole, err := s.repo.GetRoleByName("user")
	if err != nil {
		return fmt.Errorf("failed to get default user role: %v", err)
	}

	// Update all users without a role
	return s.repo.AssignDefaultRoleToUsersWithoutRole(userRole.ID)
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

// GetUserWithRole gets a user with role information
func (s *Service) GetUserWithRole(id string) (*User, error) {
	return s.repo.GetUserWithRole(id)
}

// GetDefaultRole gets the default role for new users
func (s *Service) GetDefaultRole() (*roles.Role, error) {
	return s.repo.GetDefaultRole()
}

// GetUserTenants gets all tenants for a user
func (s *Service) GetUserTenants(userID string) ([]UserTenant, error) {
	return s.repo.GetUserTenants(userID)
}

// CanManageRole checks if a role can manage another role
func CanManageRole(managerRole, targetRole string) bool {
	// Define role hierarchy
	roleHierarchy := map[string]int{
		"super_admin":  3,
		"tenant_admin": 2,
		"user":         1,
	}

	managerLevel, managerExists := roleHierarchy[managerRole]
	targetLevel, targetExists := roleHierarchy[targetRole]

	if !managerExists || !targetExists {
		return false
	}

	// Can only manage roles with lower hierarchy level
	return managerLevel > targetLevel
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

// GetRoleByID gets a role by ID
func (s *Service) GetRoleByID(roleID string) (*roles.Role, error) {
	return s.repo.GetRoleByID(roleID)
}

// GetUserTenantIDs gets all tenant IDs for a user
func (s *Service) GetUserTenantIDs(userID string) ([]string, error) {
	userTenants, err := s.repo.GetUserTenants(userID)
	if err != nil {
		return nil, err
	}
	
	var tenantIDs []string
	for _, userTenant := range userTenants {
		tenantIDs = append(tenantIDs, userTenant.TenantID)
	}
	
	return tenantIDs, nil
}

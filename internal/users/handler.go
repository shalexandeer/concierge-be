package users

import (
	"net/http"
	"strconv"

	"concierge-be/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// CreateUser creates a new user
func (h *Handler) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Clear password from response
	user.Password = ""
	utils.SuccessResponse(c, user)
}

// GetUser gets a user by ID
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// Clear password from response
	user.Password = ""
	utils.SuccessResponse(c, user)
}

// GetAllUsers gets all users with pagination
func (h *Handler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.service.GetAllUsers(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Clear passwords from response
	for i := range users {
		users[i].Password = ""
	}

	utils.SuccessResponseWithPagination(c, users, page, pageSize, int(total))
}

// UpdateUser updates a user
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	user.ID = id
	if err := h.service.UpdateUser(&user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Clear password from response
	user.Password = ""
	utils.SuccessResponse(c, user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "User deleted successfully"})
}

// CreateTenant creates a new tenant
func (h *Handler) CreateTenant(c *gin.Context) {
	var tenant Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if err := h.service.CreateTenant(&tenant); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, tenant)
}

// GetTenant gets a tenant by ID
func (h *Handler) GetTenant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tenant ID is required")
		return
	}

	tenant, err := h.service.GetTenantByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, tenant)
}

// GetAllTenants gets all tenants with pagination
func (h *Handler) GetAllTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	tenants, total, err := h.service.GetAllTenants(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithPagination(c, tenants, page, pageSize, int(total))
}

// UpdateTenant updates a tenant
func (h *Handler) UpdateTenant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tenant ID is required")
		return
	}

	var tenant Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	tenant.ID = id
	if err := h.service.UpdateTenant(&tenant); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, tenant)
}

// DeleteTenant deletes a tenant
func (h *Handler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tenant ID is required")
		return
	}

	if err := h.service.DeleteTenant(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Tenant deleted successfully"})
}

// AddUserToTenant adds a user to a tenant
func (h *Handler) AddUserToTenant(c *gin.Context) {
	var req struct {
		UserID   string `json:"userId" binding:"required"`
		TenantID string `json:"tenantId" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if req.Role == "" {
		req.Role = "member"
	}

	if err := h.service.AddUserToTenant(req.UserID, req.TenantID, req.Role); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "User added to tenant successfully"})
}

// GetUserTenants gets all tenants for a user
func (h *Handler) GetUserTenants(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userTenants, err := h.service.GetUserTenants(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, userTenants)
}

// GetTenantUsers gets all users for a tenant
func (h *Handler) GetTenantUsers(c *gin.Context) {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tenant ID is required")
		return
	}

	userTenants, err := h.service.GetTenantUsers(tenantID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, userTenants)
}

// RemoveUserFromTenant removes a user from a tenant
func (h *Handler) RemoveUserFromTenant(c *gin.Context) {
	userID := c.Param("userId")
	tenantID := c.Param("tenantId")

	if userID == "" || tenantID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID and Tenant ID are required")
		return
	}

	if err := h.service.RemoveUserFromTenant(userID, tenantID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "User removed from tenant successfully"})
}

package tenants

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

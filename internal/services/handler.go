package services

import (
	"concierge-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *ServiceService
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// CreateService handles POST /api/v1/services
func (h *Handler) CreateService(c *gin.Context) {
	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	service, err := h.service.CreateService(&req)
	if err != nil {
		if err.Error() == "service name already exists for this tenant" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, utils.Response{
		Code:    http.StatusCreated,
		Message: "Success",
		Data:    service,
	})
}

// GetService handles GET /api/v1/services/:id
func (h *Handler) GetService(c *gin.Context) {
	id := c.Param("id")

	service, err := h.service.GetServiceByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Service not found")
		return
	}

	utils.SuccessResponse(c, service)
}

// GetAllServices handles GET /api/v1/services
// Public endpoint - no authentication required
func (h *Handler) GetAllServices(c *gin.Context) {
	categoryID := c.Query("categoryId")
	tenantID := c.Query("tenantId")

	var services []Service
	var err error

	if categoryID != "" && tenantID != "" {
		// Filter by both category and tenant
		services, err = h.service.GetServicesByCategoryIDAndTenantIDs(categoryID, []string{tenantID})
	} else if categoryID != "" {
		// Filter by category only
		services, err = h.service.GetServicesByCategoryID(categoryID)
	} else if tenantID != "" {
		// Filter by tenant only
		services, err = h.service.GetServicesByTenantID(tenantID)
	} else {
		// Get all services
		services, err = h.service.GetAllServices()
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, services)
}

// UpdateService handles PUT /api/v1/services/:id
func (h *Handler) UpdateService(c *gin.Context) {
	id := c.Param("id")

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	service, err := h.service.UpdateService(id, &req)
	if err != nil {
		if err.Error() == "service name already exists for this tenant" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, service)
}

// DeleteService handles DELETE /api/v1/services/:id
func (h *Handler) DeleteService(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteService(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Service deleted successfully"})
}

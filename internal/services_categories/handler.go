package services_categories

import (
	"concierge-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *ServiceCategoryService
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// CreateServiceCategory handles POST /api/v1/services-categories
func (h *Handler) CreateServiceCategory(c *gin.Context) {
	var req CreateServiceCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.CreateServiceCategory(&req)
	if err != nil {
		if err.Error() == "category name already exists" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, utils.Response{
		Code:    http.StatusCreated,
		Message: "Success",
		Data:    category,
	})
}

// GetServiceCategory handles GET /api/v1/services-categories/:id
func (h *Handler) GetServiceCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := h.service.GetServiceCategoryByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Service category not found")
		return
	}

	utils.SuccessResponse(c, category)
}

// GetAllServiceCategories handles GET /api/v1/services-categories
func (h *Handler) GetAllServiceCategories(c *gin.Context) {
	categories, err := h.service.GetAllServiceCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, categories)
}

// UpdateServiceCategory handles PUT /api/v1/services-categories/:id
func (h *Handler) UpdateServiceCategory(c *gin.Context) {
	id := c.Param("id")

	var req UpdateServiceCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.UpdateServiceCategory(id, &req)
	if err != nil {
		if err.Error() == "category name already exists" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, category)
}

// DeleteServiceCategory handles DELETE /api/v1/services-categories/:id
func (h *Handler) DeleteServiceCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteServiceCategory(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Service category deleted successfully"})
}

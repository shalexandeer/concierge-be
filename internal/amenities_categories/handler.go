package amenities_categories

import (
	"concierge-be/utils"
	"net/http"

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

// CreateCategory handles POST /api/v1/amenities-categories
func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateAmenityCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.CreateCategory(&req)
	if err != nil {
		if err.Error() == "category name already exists for this tenant" {
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

// GetCategory handles GET /api/v1/amenities-categories/:id
func (h *Handler) GetCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	utils.SuccessResponse(c, category)
}

// GetAllCategories handles GET /api/v1/amenities-categories
// Optionally filtered by tenantId query parameter
func (h *Handler) GetAllCategories(c *gin.Context) {
	tenantID := c.Query("tenantId")

	var categories []AmenityCategory
	var err error

	if tenantID != "" {
		categories, err = h.service.GetCategoriesByTenantID(tenantID)
	} else {
		categories, err = h.service.GetAllCategories()
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, categories)
}

// UpdateCategory handles PUT /api/v1/amenities-categories/:id
func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAmenityCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.UpdateCategory(id, &req)
	if err != nil {
		if err.Error() == "category name already exists for this tenant" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, category)
}

// DeleteCategory handles DELETE /api/v1/amenities-categories/:id
func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteCategory(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Category deleted successfully"})
}


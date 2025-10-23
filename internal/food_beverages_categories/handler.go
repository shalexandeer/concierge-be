package food_beverages_categories

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

// CreateCategory creates a new food & beverage category
func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateFoodBeverageCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	category, err := h.service.CreateCategory(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, category)
}

// GetCategory gets a food & beverage category by ID
func (h *Handler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	utils.SuccessResponse(c, category)
}

// GetAllCategories gets all food & beverage categories
func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, categories)
}

// UpdateCategory updates a food & beverage category
func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	var req UpdateFoodBeverageCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	category, err := h.service.UpdateCategory(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, category)
}

// DeleteCategory deletes a food & beverage category
func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Category deleted successfully"})
}

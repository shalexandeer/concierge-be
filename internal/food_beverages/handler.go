package food_beverages

import (
	"concierge-be/utils"
	"net/http"
	"strconv"

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

// CreateFoodBeverage creates a new food & beverage item
func (h *Handler) CreateFoodBeverage(c *gin.Context) {
	var req CreateFoodBeverageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	item, err := h.service.CreateFoodBeverage(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, item)
}

// GetFoodBeverage gets a food & beverage item by ID
func (h *Handler) GetFoodBeverage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Food & beverage item ID is required")
		return
	}

	item, err := h.service.GetFoodBeverageByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Food & beverage item not found")
		return
	}

	utils.SuccessResponse(c, item)
}

// GetAllFoodBeverages gets all food & beverage items with optional filtering
func (h *Handler) GetAllFoodBeverages(c *gin.Context) {
	// Get query parameters
	tenantID := c.Query("tenantId")
	categoryID := c.Query("categoryId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := h.service.GetAllFoodBeverages(tenantID, categoryID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := gin.H{
		"data": items,
		"pagination": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": totalPages,
		},
	}

	utils.SuccessResponse(c, response)
}

// UpdateFoodBeverage updates a food & beverage item
func (h *Handler) UpdateFoodBeverage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Food & beverage item ID is required")
		return
	}

	var req UpdateFoodBeverageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	item, err := h.service.UpdateFoodBeverage(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, item)
}

// DeleteFoodBeverage deletes a food & beverage item
func (h *Handler) DeleteFoodBeverage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Food & beverage item ID is required")
		return
	}

	if err := h.service.DeleteFoodBeverage(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Food & beverage item deleted successfully"})
}

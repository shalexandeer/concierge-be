package amenities

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

// CreateAmenity handles POST /api/v1/amenities
func (h *Handler) CreateAmenity(c *gin.Context) {
	var req CreateAmenityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amenity, err := h.service.CreateAmenity(&req)
	if err != nil {
		if err.Error() == "item name already exists for this tenant" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, utils.Response{
		Code:    http.StatusCreated,
		Message: "Success",
		Data:    amenity,
	})
}

// GetAmenity handles GET /api/v1/amenities/:id
func (h *Handler) GetAmenity(c *gin.Context) {
	id := c.Param("id")

	amenity, err := h.service.GetAmenityByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Amenity not found")
		return
	}

	utils.SuccessResponse(c, amenity)
}

// GetAllAmenities handles GET /api/v1/amenities
// Supports filtering by tenantId, categoryId, and lowStock query parameters
func (h *Handler) GetAllAmenities(c *gin.Context) {
	tenantID := c.Query("tenantId")
	categoryID := c.Query("categoryId")
	lowStock := c.Query("lowStock")

	var amenities []Amenity
	var err error

	// Priority: lowStock > categoryId > tenantId > all
	if lowStock == "true" && tenantID != "" {
		amenities, err = h.service.GetLowStockAmenities(tenantID)
	} else if categoryID != "" {
		amenities, err = h.service.GetAmenitiesByCategoryID(categoryID)
	} else if tenantID != "" {
		amenities, err = h.service.GetAmenitiesByTenantID(tenantID)
	} else {
		amenities, err = h.service.GetAllAmenities()
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, amenities)
}

// UpdateAmenity handles PUT /api/v1/amenities/:id
func (h *Handler) UpdateAmenity(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAmenityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amenity, err := h.service.UpdateAmenity(id, &req)
	if err != nil {
		if err.Error() == "item name already exists for this tenant" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, amenity)
}

// UpdateStock handles PATCH /api/v1/amenities/:id/stock
func (h *Handler) UpdateStock(c *gin.Context) {
	id := c.Param("id")
	
	stockStr := c.Query("quantity")
	if stockStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "quantity query parameter is required")
		return
	}

	quantity, err := strconv.Atoi(stockStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "quantity must be a valid integer")
		return
	}

	amenity, err := h.service.UpdateStock(id, quantity)
	if err != nil {
		if err.Error() == "stock quantity cannot be negative" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, amenity)
}

// DeleteAmenity handles DELETE /api/v1/amenities/:id
func (h *Handler) DeleteAmenity(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteAmenity(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Amenity deleted successfully"})
}


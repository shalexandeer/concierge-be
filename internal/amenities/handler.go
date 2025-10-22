package amenities

import (
	"concierge-be/internal/users"
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
// Automatically filters by user's tenants (except for super admin)
func (h *Handler) GetAllAmenities(c *gin.Context) {
	// Get current user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get current user with role information
	userService := users.NewService()
	currentUser, err := userService.GetUserWithRole(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user information")
		return
	}

	tenantID := c.Query("tenantId")
	categoryID := c.Query("categoryId")
	lowStock := c.Query("lowStock")

	var amenities []Amenity
	var err2 error

	// Check if user is super admin - if so, return all data
	if currentUser.Role.Name == "super_admin" {
		// Super admin can see all amenities
		if lowStock == "true" && tenantID != "" {
			amenities, err2 = h.service.GetLowStockAmenities(tenantID)
		} else if categoryID != "" {
			amenities, err2 = h.service.GetAmenitiesByCategoryID(categoryID)
		} else if tenantID != "" {
			amenities, err2 = h.service.GetAmenitiesByTenantID(tenantID)
		} else {
			amenities, err2 = h.service.GetAllAmenities()
		}
	} else {
		// Non-super admin users: filter by their tenants
		userTenantIDs, err := userService.GetUserTenantIDs(userID.(string))
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user tenants")
			return
		}

		if len(userTenantIDs) == 0 {
			// User has no tenants, return empty array
			amenities = []Amenity{}
		} else {
			// Filter by user's tenants
			if lowStock == "true" {
				// For low stock, check all user's tenants
				amenities, err2 = h.service.GetLowStockAmenitiesByTenantIDs(userTenantIDs)
			} else if categoryID != "" {
				amenities, err2 = h.service.GetAmenitiesByCategoryIDAndTenantIDs(categoryID, userTenantIDs)
			} else {
				amenities, err2 = h.service.GetAmenitiesByTenantIDs(userTenantIDs)
			}
		}
	}

	if err2 != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err2.Error())
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


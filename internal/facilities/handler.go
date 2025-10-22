package facilities

import (
	"net/http"
	"strconv"
	"concierge-be/internal/users"
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

// Facility handlers

// CreateFacility creates a new facility
func (h *Handler) CreateFacility(c *gin.Context) {
	var req CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	facility, err := h.service.CreateFacility(&req)
	if err != nil {
		if err.Error() == "facility name already exists for this tenant" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, facility)
}

// GetFacility gets a facility by ID
func (h *Handler) GetFacility(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Facility ID is required")
		return
	}

	facility, err := h.service.GetFacilityByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Facility not found")
		return
	}

	utils.SuccessResponse(c, facility)
}

// GetAllFacilities gets all facilities with pagination
// Automatically filters by user's tenants (except for super admin)
func (h *Handler) GetAllFacilities(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var facilities []FacilityResponse
	var total int64
	var err2 error

	// Check if user is super admin - if so, return all data
	if currentUser.Role.Name == "super_admin" {
		// Super admin can see all facilities
		facilities, total, err2 = h.service.GetAllFacilities(page, pageSize)
	} else {
		// Non-super admin users: filter by their tenants
		userTenantIDs, err := userService.GetUserTenantIDs(userID.(string))
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user tenants")
			return
		}

		if len(userTenantIDs) == 0 {
			// User has no tenants, return empty array
			facilities = []FacilityResponse{}
			total = 0
		} else {
			// Filter by user's tenants
			facilities, total, err2 = h.service.GetFacilitiesByTenantIDs(userTenantIDs, page, pageSize)
		}
	}

	if err2 != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err2.Error())
		return
	}

	utils.SuccessResponseWithNestedPagination(c, facilities, page, pageSize, int(total))
}

// GetFacilitiesByTenant gets facilities for a specific tenant
func (h *Handler) GetFacilitiesByTenant(c *gin.Context) {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tenant ID is required")
		return
	}

	facilities, err := h.service.GetFacilitiesByTenantID(tenantID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, facilities)
}

// UpdateFacility updates a facility
func (h *Handler) UpdateFacility(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Facility ID is required")
		return
	}

	var req UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	facility, err := h.service.UpdateFacility(id, &req)
	if err != nil {
		if err.Error() == "facility name already exists for this tenant" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, facility)
}

// DeleteFacility deletes a facility
func (h *Handler) DeleteFacility(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Facility ID is required")
		return
	}

	if err := h.service.DeleteFacility(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Facility deleted successfully"})
}

// Booking handlers

// CreateBooking creates a new booking
func (h *Handler) CreateBooking(c *gin.Context) {
	facilityID := c.Param("id")
	if facilityID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Facility ID is required")
		return
	}

	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Set facility ID from URL parameter
	req.FacilityID = facilityID

	booking, err := h.service.CreateBooking(&req)
	if err != nil {
		if err.Error() == "facility is not available for the requested time" ||
			err.Error() == "cannot book in the past" ||
			err.Error() == "start time must be before end time" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, booking)
}

// GetBooking gets a booking by ID
func (h *Handler) GetBooking(c *gin.Context) {
	id := c.Param("bookingId")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Booking ID is required")
		return
	}

	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Booking not found")
		return
	}

	utils.SuccessResponse(c, booking)
}

// GetFacilityBookings gets all bookings for a facility
func (h *Handler) GetFacilityBookings(c *gin.Context) {
	facilityID := c.Param("id")
	if facilityID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Facility ID is required")
		return
	}

	bookings, err := h.service.GetFacilityBookingHistory(facilityID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, bookings)
}

// GetUpcomingBookings gets upcoming bookings for a facility
func (h *Handler) GetUpcomingBookings(c *gin.Context) {
	facilityID := c.Param("id")
	if facilityID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Facility ID is required")
		return
	}

	bookings, err := h.service.GetUpcomingBookings(facilityID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, bookings)
}

// DeleteBooking deletes a booking
func (h *Handler) DeleteBooking(c *gin.Context) {
	id := c.Param("bookingId")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Booking ID is required")
		return
	}

	if err := h.service.DeleteBooking(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Booking deleted successfully"})
}

// GetBookingsByTenant gets all bookings for a tenant with pagination
func (h *Handler) GetBookingsByTenant(c *gin.Context) {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tenant ID is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	bookings, total, err := h.service.GetBookingsByTenantID(tenantID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithPagination(c, bookings, page, pageSize, int(total))
}

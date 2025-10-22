package facilities

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// generateUUID generates a new UUID
func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Facility service methods

func (s *Service) CreateFacility(req *CreateFacilityRequest) (*FacilityResponse, error) {
	// Check if facility name already exists for this tenant
	exists, err := s.repo.CheckFacilityNameExists(req.TenantID, req.FacilityName, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("facility name already exists for this tenant")
	}

	// Create facility
	// Convert EquipmentList to JSON string
	equipmentJSON, err := json.Marshal(req.Equipment)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal equipment: %v", err)
	}

	facility := &Facility{
		ID:           generateUUID(),
		TenantID:     req.TenantID,
		FacilityName: req.FacilityName,
		Capacity:     req.Capacity,
		RatePerHour:  req.RatePerHour,
		Status:       "available",
		Equipment:    string(equipmentJSON),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(facility); err != nil {
		return nil, err
	}

	// Convert back to response format
	return s.convertFacilityToResponse(facility)
}

func (s *Service) GetFacilityByID(id string) (*FacilityResponse, error) {
	facility, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.convertFacilityToResponse(facility)
}

func (s *Service) GetFacilitiesByTenantID(tenantID string) ([]FacilityResponse, error) {
	facilities, err := s.repo.GetByTenantID(tenantID)
	if err != nil {
		return nil, err
	}

	// Convert each facility to response format
	var responses []FacilityResponse
	for _, facility := range facilities {
		response, err := s.convertFacilityToResponse(&facility)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

func (s *Service) GetAllFacilities(page, pageSize int) ([]FacilityResponse, int64, error) {
	facilities, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Convert each facility to response format
	var responses []FacilityResponse
	for _, facility := range facilities {
		response, err := s.convertFacilityToResponse(&facility)
		if err != nil {
			return nil, 0, err
		}
		responses = append(responses, *response)
	}

	return responses, total, nil
}

func (s *Service) UpdateFacility(id string, req *UpdateFacilityRequest) (*Facility, error) {
	facility, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if facility name already exists for this tenant (if name is being updated)
	if req.FacilityName != "" && req.FacilityName != facility.FacilityName {
		exists, err := s.repo.CheckFacilityNameExists(facility.TenantID, req.FacilityName, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("facility name already exists for this tenant")
		}
		facility.FacilityName = req.FacilityName
	}

	if req.Capacity != nil {
		facility.Capacity = *req.Capacity
	}
	if req.RatePerHour != nil {
		facility.RatePerHour = *req.RatePerHour
	}
	if req.Equipment != nil {
		equipmentJSON, err := json.Marshal(*req.Equipment)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal equipment: %v", err)
		}
		facility.Equipment = string(equipmentJSON)
	}

	facility.UpdatedAt = time.Now()

	if err := s.repo.Update(facility); err != nil {
		return nil, err
	}

	return facility, nil
}

func (s *Service) DeleteFacility(id string) error {
	return s.repo.Delete(id)
}

// Booking service methods

func (s *Service) CreateBooking(req *CreateBookingRequest) (*Booking, error) {
	// Validate time range
	if req.StartDateTime.After(req.EndDateTime) {
		return nil, fmt.Errorf("start time must be before end time")
	}

	if req.StartDateTime.Before(time.Now()) {
		return nil, fmt.Errorf("cannot book in the past")
	}

	// Check if facility exists
	_, err := s.repo.GetByID(req.FacilityID)
	if err != nil {
		return nil, fmt.Errorf("facility not found")
	}

	// Check if facility is available for the requested time
	available, err := s.repo.CheckFacilityAvailability(req.FacilityID, req.StartDateTime, req.EndDateTime)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, fmt.Errorf("facility is not available for the requested time")
	}

	// Create booking
	booking := &Booking{
		ID:            generateUUID(),
		TenantID:      req.TenantID,
		FacilityID:    req.FacilityID,
		GuestName:     req.GuestName,
		StartDateTime: req.StartDateTime,
		EndDateTime:   req.EndDateTime,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateBooking(booking); err != nil {
		return nil, err
	}

	// Update facility status to booked if this is the next upcoming booking
	s.updateFacilityStatus(req.FacilityID)

	return booking, nil
}

func (s *Service) GetBookingByID(id string) (*Booking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *Service) GetFacilityBookings(facilityID string) ([]Booking, error) {
	return s.repo.GetFacilityBookings(facilityID)
}

func (s *Service) GetFacilityBookingHistory(facilityID string) ([]Booking, error) {
	return s.repo.GetFacilityBookingHistory(facilityID)
}

func (s *Service) GetUpcomingBookings(facilityID string) ([]Booking, error) {
	return s.repo.GetUpcomingBookings(facilityID)
}

func (s *Service) DeleteBooking(id string) error {
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteBooking(id); err != nil {
		return err
	}

	// Update facility status after booking deletion
	s.updateFacilityStatus(booking.FacilityID)

	return nil
}

func (s *Service) GetBookingsByTenantID(tenantID string, page, pageSize int) ([]Booking, int64, error) {
	return s.repo.GetBookingsByTenantID(tenantID, page, pageSize)
}

// updateFacilityStatus updates the facility status based on current bookings
func (s *Service) updateFacilityStatus(facilityID string) {
	now := time.Now()
	
	// Check if there are any active bookings (current time is between start and end)
	upcomingBookings, err := s.repo.GetUpcomingBookings(facilityID)
	if err != nil {
		return
	}

	// Check if there's an active booking
	hasActiveBooking := false
	for _, booking := range upcomingBookings {
		if now.After(booking.StartDateTime) && now.Before(booking.EndDateTime) {
			hasActiveBooking = true
			break
		}
	}

	// Update facility status
	status := "available"
	if hasActiveBooking {
		status = "booked"
	}

	s.repo.UpdateStatus(facilityID, status)
}

// convertFacilityToResponse converts a Facility with JSON string to a response format
func (s *Service) convertFacilityToResponse(facility *Facility) (*FacilityResponse, error) {
	// Parse equipment from JSON string to EquipmentList
	var equipment EquipmentList
	if len(facility.Equipment) > 0 {
		if err := json.Unmarshal([]byte(facility.Equipment), &equipment); err != nil {
			return nil, fmt.Errorf("failed to unmarshal equipment: %v", err)
		}
	}

	// Create a response with the parsed equipment
	response := &FacilityResponse{
		ID:           facility.ID,
		TenantID:     facility.TenantID,
		FacilityName: facility.FacilityName,
		Capacity:     facility.Capacity,
		RatePerHour:  facility.RatePerHour,
		Status:       facility.Status,
		Equipment:    equipment,
		CreatedAt:    facility.CreatedAt,
		UpdatedAt:    facility.UpdatedAt,
		Bookings:     facility.Bookings,
	}

	return response, nil
}

package facilities

import (
	"concierge-be/database"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

// Facility CRUD operations

// Create creates a new facility
func (r *Repository) Create(facility *Facility) error {
	return r.db.Create(facility).Error
}

// GetByID retrieves a facility by ID
func (r *Repository) GetByID(id string) (*Facility, error) {
	var facility Facility
	err := r.db.Where("id = ?", id).First(&facility).Error
	if err != nil {
		return nil, err
	}
	return &facility, nil
}

// GetByTenantID retrieves all facilities for a specific tenant
func (r *Repository) GetByTenantID(tenantID string) ([]Facility, error) {
	var facilities []Facility
	err := r.db.Where("tenant_id = ?", tenantID).Order("facility_name ASC").Find(&facilities).Error
	if err != nil {
		return nil, err
	}
	return facilities, nil
}

// GetAll retrieves all facilities with pagination
func (r *Repository) GetAll(page, pageSize int) ([]Facility, int64, error) {
	var facilities []Facility
	var total int64

	// Count total records
	if err := r.db.Model(&Facility{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("facility_name ASC").Find(&facilities).Error
	if err != nil {
		return nil, 0, err
	}

	return facilities, total, nil
}

// Update updates an existing facility
func (r *Repository) Update(facility *Facility) error {
	return r.db.Save(facility).Error
}

// Delete soft deletes a facility
func (r *Repository) Delete(id string) error {
	return r.db.Delete(&Facility{}, "id = ?", id).Error
}

// CheckFacilityNameExists checks if a facility name already exists for a tenant
func (r *Repository) CheckFacilityNameExists(tenantID, facilityName, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&Facility{}).Where("tenant_id = ? AND facility_name = ?", tenantID, facilityName)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus updates the status of a facility
func (r *Repository) UpdateStatus(id, status string) error {
	return r.db.Model(&Facility{}).Where("id = ?", id).Update("status", status).Error
}

// Booking CRUD operations

// CreateBooking creates a new booking
func (r *Repository) CreateBooking(booking *Booking) error {
	return r.db.Create(booking).Error
}

// GetBookingByID retrieves a booking by ID
func (r *Repository) GetBookingByID(id string) (*Booking, error) {
	var booking Booking
	err := r.db.Preload("Facility").Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// GetFacilityBookings retrieves all bookings for a specific facility
func (r *Repository) GetFacilityBookings(facilityID string) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("facility_id = ?", facilityID).
		Order("start_date_time DESC").
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

// GetFacilityBookingHistory retrieves all bookings for a facility (past and future)
func (r *Repository) GetFacilityBookingHistory(facilityID string) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Preload("Facility").
		Where("facility_id = ?", facilityID).
		Order("start_date_time DESC").
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

// GetUpcomingBookings retrieves future bookings for a facility
func (r *Repository) GetUpcomingBookings(facilityID string) ([]Booking, error) {
	var bookings []Booking
	now := time.Now()
	err := r.db.Where("facility_id = ? AND start_date_time > ?", facilityID, now).
		Order("start_date_time ASC").
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

// CheckFacilityAvailability checks if a facility is available for a given time range
func (r *Repository) CheckFacilityAvailability(facilityID string, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&Booking{}).
		Where("facility_id = ? AND deleted_at IS NULL AND ((start_date_time < ? AND end_date_time > ?) OR (start_date_time < ? AND end_date_time > ?) OR (start_date_time >= ? AND end_date_time <= ?))",
			facilityID, endTime, startTime, startTime, endTime, startTime, endTime).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// DeleteBooking soft deletes a booking
func (r *Repository) DeleteBooking(id string) error {
	return r.db.Delete(&Booking{}, "id = ?", id).Error
}

// GetBookingsByTenantID retrieves all bookings for a specific tenant
func (r *Repository) GetBookingsByTenantID(tenantID string, page, pageSize int) ([]Booking, int64, error) {
	var bookings []Booking
	var total int64

	// Count total records
	if err := r.db.Model(&Booking{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.Preload("Facility").
		Where("tenant_id = ?", tenantID).
		Offset(offset).
		Limit(pageSize).
		Order("start_date_time DESC").
		Find(&bookings).Error
	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

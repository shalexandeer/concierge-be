package facilities

import (
	"time"

	"gorm.io/gorm"
)

// Equipment represents equipment available in a facility
type Equipment struct {
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

// EquipmentList is a custom type for JSON array of Equipment
type EquipmentList []Equipment

// Facility represents a hotel facility (tenant-scoped)
type Facility struct {
	ID           string        `gorm:"type:varchar(36);primaryKey" json:"id"`
	TenantID     string        `gorm:"type:varchar(36);not null;index" json:"tenantId"`
	FacilityName string        `gorm:"type:varchar(100);not null" json:"facilityName"`
	Capacity     int           `gorm:"not null" json:"capacity"`
	RatePerHour  float64       `gorm:"type:decimal(10,2);not null" json:"ratePerHour"`
	Status       string        `gorm:"type:varchar(20);default:'available'" json:"status"`
	Equipment    string `gorm:"type:json" json:"equipment"`
	ImagePath    string        `gorm:"type:varchar(500)" json:"imagePath"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Bookings []Booking `gorm:"foreignKey:FacilityID;references:ID" json:"bookings,omitempty"`
}

func (Facility) TableName() string {
	return "facilities"
}

// Booking represents a facility booking
type Booking struct {
	ID            string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TenantID      string    `gorm:"type:varchar(36);not null;index" json:"tenantId"`
	FacilityID    string    `gorm:"type:varchar(36);not null;index" json:"facilityId"`
	GuestName     string    `gorm:"type:varchar(100);not null" json:"guestName"`
	StartDateTime time.Time `gorm:"not null" json:"startDateTime"`
	EndDateTime   time.Time `gorm:"not null" json:"endDateTime"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Facility *Facility `gorm:"foreignKey:FacilityID;references:ID" json:"facility,omitempty"`
}

func (Booking) TableName() string {
	return "facility_bookings"
}

// CreateFacilityRequest represents the request body for creating a facility
type CreateFacilityRequest struct {
	TenantID     string        `json:"tenantId" binding:"required"`
	FacilityName string        `json:"facilityName" binding:"required"`
	Capacity     int           `json:"capacity" binding:"required,min=1"`
	RatePerHour  float64       `json:"ratePerHour" binding:"required,min=0"`
	Equipment    EquipmentList `json:"equipment"`
	ImagePath    string        `json:"imagePath"`
}

// UpdateFacilityRequest represents the request body for updating a facility
type UpdateFacilityRequest struct {
	FacilityName string         `json:"facilityName"`
	Capacity     *int           `json:"capacity"`
	RatePerHour  *float64       `json:"ratePerHour"`
	Equipment    *EquipmentList `json:"equipment"`
	ImagePath    string         `json:"imagePath"`
}

// CreateBookingRequest represents the request body for creating a booking
type CreateBookingRequest struct {
	TenantID      string    `json:"tenantId" binding:"required"`
	FacilityID    string    `json:"facilityId" binding:"required"`
	GuestName     string    `json:"guestName" binding:"required"`
	StartDateTime time.Time `json:"startDateTime" binding:"required"`
	EndDateTime   time.Time `json:"endDateTime" binding:"required"`
}

// FacilityWithBookings represents a facility with its booking information
type FacilityWithBookings struct {
	Facility
	UpcomingBookings []Booking `json:"upcomingBookings"`
	PastBookings     []Booking `json:"pastBookings"`
}

// FacilityResponse represents a facility response with parsed equipment
type FacilityResponse struct {
	ID           string        `json:"id"`
	TenantID     string        `json:"tenantId"`
	FacilityName string        `json:"facilityName"`
	Capacity     int           `json:"capacity"`
	RatePerHour  float64       `json:"ratePerHour"`
	Status       string        `json:"status"`
	Equipment    EquipmentList `json:"equipment"`
	ImagePath    string        `json:"imagePath"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	Bookings     []Booking     `json:"bookings,omitempty"`
}

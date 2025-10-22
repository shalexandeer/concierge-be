-- Create facilities table
CREATE TABLE facilities (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  facility_name VARCHAR(100) NOT NULL,
  capacity INT NOT NULL,
  rate_per_hour DECIMAL(10,2) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'available',
  equipment JSONB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  INDEX(tenant_id)
);

-- Create facility_bookings table
CREATE TABLE facility_bookings (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  facility_id VARCHAR(36) NOT NULL,
  guest_name VARCHAR(100) NOT NULL,
  start_date_time TIMESTAMP NOT NULL,
  end_date_time TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  FOREIGN KEY (facility_id) REFERENCES facilities(id),
  INDEX(tenant_id),
  INDEX(facility_id),
  INDEX(start_date_time),
  INDEX(end_date_time)
);

-- Insert sample data
INSERT INTO facilities (id, tenant_id, facility_name, capacity, rate_per_hour, status, equipment) VALUES
('facility-1', 'tenant-1', 'Conference Room A', 20, 500000.00, 'available', '[{"name": "Projector", "quantity": 1}, {"name": "Whiteboard", "quantity": 1}, {"name": "Microphone", "quantity": 2}]'),
('facility-2', 'tenant-1', 'Meeting Room B', 8, 250000.00, 'available', '[{"name": "TV Screen", "quantity": 1}, {"name": "Whiteboard", "quantity": 1}]'),
('facility-3', 'tenant-1', 'Ballroom', 100, 2000000.00, 'available', '[{"name": "Stage", "quantity": 1}, {"name": "Sound System", "quantity": 1}, {"name": "Lighting", "quantity": 1}]'),
('facility-4', 'tenant-2', 'Gymnasium', 50, 300000.00, 'available', '[{"name": "Treadmill", "quantity": 5}, {"name": "Dumbbells", "quantity": 10}]'),
('facility-5', 'tenant-2', 'Swimming Pool', 30, 150000.00, 'available', '[{"name": "Pool Chairs", "quantity": 20}, {"name": "Umbrella", "quantity": 10}]');

-- Insert sample bookings
INSERT INTO facility_bookings (id, tenant_id, facility_id, guest_name, start_date_time, end_date_time) VALUES
('booking-1', 'tenant-1', 'facility-1', 'John Doe', '2024-01-15 09:00:00', '2024-01-15 11:00:00'),
('booking-2', 'tenant-1', 'facility-1', 'Jane Smith', '2024-01-15 14:00:00', '2024-01-15 16:00:00'),
('booking-3', 'tenant-1', 'facility-2', 'Bob Johnson', '2024-01-16 10:00:00', '2024-01-16 12:00:00'),
('booking-4', 'tenant-1', 'facility-3', 'Alice Brown', '2024-01-20 18:00:00', '2024-01-20 22:00:00'),
('booking-5', 'tenant-2', 'facility-4', 'Charlie Wilson', '2024-01-17 07:00:00', '2024-01-17 09:00:00');

# API Usage Examples

Quick reference for testing the Concierge BE API.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication Endpoints

### Register New User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepass123",
    "fullName": "John Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepass123"
  }'
```

### Get Current User (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Current User (Protected)
```bash
curl -X PUT http://localhost:8080/api/v1/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "email": "newemail@example.com",
    "fullName": "John Updated Doe"
  }'
```

## User Management Endpoints

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jane",
    "email": "jane@example.com",
    "password": "password123",
    "fullName": "Jane Smith"
  }'
```

### Get User by ID
```bash
curl -X GET http://localhost:8080/api/v1/users/USER_ID \
  -H "Content-Type: application/json"
```

### Get All Users (with pagination)
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&pageSize=10" \
  -H "Content-Type: application/json"
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "email": "updated@example.com",
    "fullName": "Updated Name"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/USER_ID \
  -H "Content-Type: application/json"
```

## Tenant Management Endpoints

### Create Tenant
```bash
curl -X POST http://localhost:8080/api/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Startup Inc",
    "description": "A cutting-edge technology company",
    "domain": "techstartup.example.com",
    "isActive": true
  }'
```

### Get Tenant by ID
```bash
curl -X GET http://localhost:8080/api/v1/tenants/TENANT_ID \
  -H "Content-Type: application/json"
```

### Get All Tenants (with pagination)
```bash
curl -X GET "http://localhost:8080/api/v1/tenants?page=1&pageSize=10" \
  -H "Content-Type: application/json"
```

### Update Tenant
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/TENANT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Company Name",
    "description": "Updated description",
    "isActive": true
  }'
```

### Delete Tenant
```bash
curl -X DELETE http://localhost:8080/api/v1/tenants/TENANT_ID \
  -H "Content-Type: application/json"
```

## User-Tenant Relationship Endpoints

### Add User to Tenant
```bash
curl -X POST http://localhost:8080/api/v1/user-tenants \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "USER_UUID",
    "tenantId": "TENANT_UUID",
    "role": "admin"
  }'
```

### Get All Tenants for a User
```bash
curl -X GET http://localhost:8080/api/v1/user-tenants/users/USER_ID \
  -H "Content-Type: application/json"
```

### Get All Users in a Tenant
```bash
curl -X GET http://localhost:8080/api/v1/user-tenants/tenants/TENANT_ID \
  -H "Content-Type: application/json"
```

### Remove User from Tenant
```bash
curl -X DELETE http://localhost:8080/api/v1/user-tenants/users/USER_ID/tenants/TENANT_ID \
  -H "Content-Type: application/json"
```

## Health Check

### Check Service Status
```bash
curl -X GET http://localhost:8080/api/v1/health
```

## Running the Server

### Development Mode
```bash
./concierge-be -e development
```

### Production Mode
```bash
./concierge-be -e production
```

## Common Response Formats

### Success Response
```json
{
  "code": 200,
  "message": "Success",
  "data": { ... }
}
```

### Success Response with Pagination
```json
{
  "code": 200,
  "message": "Success",
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 100
  }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Error message description"
}
```

## Amenity Categories Endpoints

### Create Amenity Category
```bash
curl -X POST http://localhost:8080/api/v1/amenities-categories \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "tenant-uuid-here",
    "name": "Bedding",
    "description": "Bedding and linen items"
  }'
```

### Get All Amenity Categories
```bash
# Get all categories
curl -X GET http://localhost:8080/api/v1/amenities-categories

# Get categories for a specific tenant
curl -X GET "http://localhost:8080/api/v1/amenities-categories?tenantId=tenant-uuid-here"
```

### Get Single Amenity Category
```bash
curl -X GET http://localhost:8080/api/v1/amenities-categories/category-uuid-here
```

### Update Amenity Category
```bash
curl -X PUT http://localhost:8080/api/v1/amenities-categories/category-uuid-here \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bedding & Linen",
    "description": "All bedding and linen items"
  }'
```

### Delete Amenity Category
```bash
curl -X DELETE http://localhost:8080/api/v1/amenities-categories/category-uuid-here
```

## Amenities Endpoints

### Create Amenity
```bash
curl -X POST http://localhost:8080/api/v1/amenities \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "tenant-uuid-here",
    "categoryId": "category-uuid-here",
    "itemName": "King Size Bed Sheet",
    "description": "White cotton bed sheet for king size bed",
    "stock": 50,
    "minimumStock": 10,
    "available": true
  }'
```

### Get All Amenities
```bash
# Get all amenities
curl -X GET http://localhost:8080/api/v1/amenities

# Get amenities for a specific tenant
curl -X GET "http://localhost:8080/api/v1/amenities?tenantId=tenant-uuid-here"

# Get amenities for a specific category
curl -X GET "http://localhost:8080/api/v1/amenities?categoryId=category-uuid-here"

# Get low stock amenities for a tenant
curl -X GET "http://localhost:8080/api/v1/amenities?tenantId=tenant-uuid-here&lowStock=true"
```

### Get Single Amenity
```bash
curl -X GET http://localhost:8080/api/v1/amenities/amenity-uuid-here
```

### Update Amenity
```bash
curl -X PUT http://localhost:8080/api/v1/amenities/amenity-uuid-here \
  -H "Content-Type: application/json" \
  -d '{
    "itemName": "King Size Bed Sheet - Premium",
    "description": "Premium white cotton bed sheet for king size bed",
    "stock": 45,
    "minimumStock": 15,
    "available": true
  }'
```

### Update Amenity Stock
```bash
# Update stock to a specific quantity
curl -X PATCH "http://localhost:8080/api/v1/amenities/amenity-uuid-here/stock?quantity=30"
```

### Delete Amenity
```bash
curl -X DELETE http://localhost:8080/api/v1/amenities/amenity-uuid-here
```

## Notes

- All timestamps are in ISO 8601 format
- All IDs are UUIDs (36 characters)
- JWT tokens expire after 72 hours (development) or 24 hours (production)
- Passwords are automatically hashed with bcrypt
- JSON field names use camelCase convention
- Amenities and amenity categories are tenant-scoped (require `tenantId`)
- Stock quantities must be non-negative integers
- Low stock filter only works when combined with `tenantId` parameter


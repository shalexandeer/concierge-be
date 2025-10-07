# Amenities API - Quick Start Guide

## 📊 Overview

Two new resources have been added to the Concierge BE API:
- **Amenity Categories**: Organize amenities into categories (e.g., Bedding, Toiletries)
- **Amenities**: Track inventory items with stock levels and availability

Both are **tenant-scoped** for multi-tenancy support.

## 🚀 Quick Start

### 1. Create a Category

```bash
curl -X POST http://localhost:8080/api/v1/amenities-categories \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "your-tenant-id",
    "name": "Bedding",
    "description": "Bedding and linen items"
  }'
```

### 2. Create an Amenity

```bash
curl -X POST http://localhost:8080/api/v1/amenities \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "your-tenant-id",
    "categoryId": "category-id-from-step-1",
    "itemName": "King Size Bed Sheet",
    "description": "White cotton bed sheet",
    "stock": 50,
    "minimumStock": 10,
    "available": true
  }'
```

### 3. Check Low Stock Items

```bash
curl -X GET "http://localhost:8080/api/v1/amenities?tenantId=your-tenant-id&lowStock=true"
```

### 4. Update Stock

```bash
curl -X PATCH "http://localhost:8080/api/v1/amenities/amenity-id/stock?quantity=75"
```

## 📋 All Endpoints

### Amenity Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/amenities-categories` | Create new category |
| GET | `/amenities-categories` | List all categories |
| GET | `/amenities-categories?tenantId=X` | List categories for tenant |
| GET | `/amenities-categories/:id` | Get single category |
| PUT | `/amenities-categories/:id` | Update category |
| DELETE | `/amenities-categories/:id` | Delete category (soft) |

### Amenities

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/amenities` | Create new amenity |
| GET | `/amenities` | List all amenities |
| GET | `/amenities?tenantId=X` | List amenities for tenant |
| GET | `/amenities?categoryId=X` | List amenities by category |
| GET | `/amenities?tenantId=X&lowStock=true` | List low stock items |
| GET | `/amenities/:id` | Get single amenity |
| PUT | `/amenities/:id` | Update amenity |
| PATCH | `/amenities/:id/stock?quantity=X` | Update stock only |
| DELETE | `/amenities/:id` | Delete amenity (soft) |

## 🎯 Common Use Cases

### Monitor Low Stock Items

```bash
# Get all items below minimum stock for a tenant
curl -X GET "http://localhost:8080/api/v1/amenities?tenantId=<ID>&lowStock=true"
```

### Bulk Stock Check

```bash
# Get all amenities with their current stock levels
curl -X GET "http://localhost:8080/api/v1/amenities?tenantId=<ID>" | jq '.data[] | {itemName, stock, minimumStock}'
```

### Category-based Inventory

```bash
# Get all toiletries
curl -X GET "http://localhost:8080/api/v1/amenities?categoryId=<toiletries-category-id>"
```

## 💡 Tips

1. **Always include tenantId** when creating resources
2. **Stock updates** use PATCH for efficiency
3. **Low stock filter** only works with tenantId parameter
4. **Soft deletes** mean items can be recovered from database
5. **Category preloading** automatically includes category info in amenity responses

## 🔑 Key Features

- ✅ Multi-tenant support
- ✅ Automatic low stock detection
- ✅ Category relationships
- ✅ Soft delete support
- ✅ Stock validation (no negative values)
- ✅ Duplicate name detection per tenant
- ✅ Full CRUD operations

## 📝 Response Format

### Success (200/201)
```json
{
  "code": 200,
  "message": "Success",
  "data": { ... }
}
```

### Error (400/404/409/500)
```json
{
  "code": 409,
  "message": "category name already exists for this tenant"
}
```

## 🗄️ Database Tables

- `amenities_categories`: Category definitions
- `amenities`: Amenity items with stock tracking

Both use VARCHAR(36) UUIDs and support soft deletes via `deleted_at`.

## 🧪 Testing

Run the comprehensive test suite:
```bash
./scripts/test-amenities.sh
```

View detailed test results:
```bash
cat AMENITIES_TEST_RESULTS.md
```

## 📚 Full Documentation

See `API_EXAMPLES.md` for complete curl examples and detailed endpoint documentation.


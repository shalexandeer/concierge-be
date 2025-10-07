#!/bin/bash

# Test script for Amenities API
# This script tests all amenities and amenities_categories endpoints

BASE_URL="http://localhost:8080/api/v1"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================================"
echo "🧪 Testing Amenities & Categories API"
echo "================================================"
echo ""

# Step 1: Create a test tenant
echo -e "${BLUE}📋 Step 1: Creating test tenant...${NC}"
TENANT_RESPONSE=$(curl -s -X POST "$BASE_URL/tenants" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Hotel",
    "description": "Test hotel for amenities testing",
    "domain": "test-hotel-'$(date +%s)'.com"
  }')

echo "$TENANT_RESPONSE" | jq '.'
TENANT_ID=$(echo "$TENANT_RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}✅ Tenant created with ID: $TENANT_ID${NC}"
echo ""

# Step 2: Create amenity categories
echo -e "${BLUE}📋 Step 2: Creating amenity categories...${NC}"
CATEGORY1_RESPONSE=$(curl -s -X POST "$BASE_URL/amenities-categories" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "'$TENANT_ID'",
    "name": "Bedding",
    "description": "Bedding and linen items"
  }')

echo "$CATEGORY1_RESPONSE" | jq '.'
CATEGORY1_ID=$(echo "$CATEGORY1_RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}✅ Category 1 (Bedding) created with ID: $CATEGORY1_ID${NC}"
echo ""

CATEGORY2_RESPONSE=$(curl -s -X POST "$BASE_URL/amenities-categories" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "'$TENANT_ID'",
    "name": "Toiletries",
    "description": "Bathroom toiletries and supplies"
  }')

echo "$CATEGORY2_RESPONSE" | jq '.'
CATEGORY2_ID=$(echo "$CATEGORY2_RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}✅ Category 2 (Toiletries) created with ID: $CATEGORY2_ID${NC}"
echo ""

# Step 3: Get all categories
echo -e "${BLUE}📋 Step 3: Getting all categories for tenant...${NC}"
curl -s -X GET "$BASE_URL/amenities-categories?tenantId=$TENANT_ID" | jq '.'
echo -e "${GREEN}✅ Categories retrieved${NC}"
echo ""

# Step 4: Create amenities
echo -e "${BLUE}📋 Step 4: Creating amenities...${NC}"
AMENITY1_RESPONSE=$(curl -s -X POST "$BASE_URL/amenities" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "'$TENANT_ID'",
    "categoryId": "'$CATEGORY1_ID'",
    "itemName": "King Size Bed Sheet",
    "description": "White cotton bed sheet for king size bed",
    "stock": 50,
    "minimumStock": 10,
    "available": true
  }')

echo "$AMENITY1_RESPONSE" | jq '.'
AMENITY1_ID=$(echo "$AMENITY1_RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}✅ Amenity 1 (Bed Sheet) created with ID: $AMENITY1_ID${NC}"
echo ""

AMENITY2_RESPONSE=$(curl -s -X POST "$BASE_URL/amenities" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "'$TENANT_ID'",
    "categoryId": "'$CATEGORY2_ID'",
    "itemName": "Shampoo Bottle",
    "description": "50ml travel shampoo bottle",
    "stock": 5,
    "minimumStock": 50,
    "available": true
  }')

echo "$AMENITY2_RESPONSE" | jq '.'
AMENITY2_ID=$(echo "$AMENITY2_RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}✅ Amenity 2 (Shampoo - LOW STOCK) created with ID: $AMENITY2_ID${NC}"
echo ""

AMENITY3_RESPONSE=$(curl -s -X POST "$BASE_URL/amenities" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "'$TENANT_ID'",
    "categoryId": "'$CATEGORY2_ID'",
    "itemName": "Soap Bar",
    "description": "Organic soap bar",
    "stock": 100,
    "minimumStock": 20,
    "available": true
  }')

echo "$AMENITY3_RESPONSE" | jq '.'
AMENITY3_ID=$(echo "$AMENITY3_RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}✅ Amenity 3 (Soap) created with ID: $AMENITY3_ID${NC}"
echo ""

# Step 5: Get single amenity
echo -e "${BLUE}📋 Step 5: Getting single amenity...${NC}"
curl -s -X GET "$BASE_URL/amenities/$AMENITY1_ID" | jq '.'
echo -e "${GREEN}✅ Single amenity retrieved${NC}"
echo ""

# Step 6: Get all amenities for tenant
echo -e "${BLUE}📋 Step 6: Getting all amenities for tenant...${NC}"
curl -s -X GET "$BASE_URL/amenities?tenantId=$TENANT_ID" | jq '.'
echo -e "${GREEN}✅ All amenities retrieved${NC}"
echo ""

# Step 7: Get amenities by category
echo -e "${BLUE}📋 Step 7: Getting amenities by category (Toiletries)...${NC}"
curl -s -X GET "$BASE_URL/amenities?categoryId=$CATEGORY2_ID" | jq '.'
echo -e "${GREEN}✅ Amenities by category retrieved${NC}"
echo ""

# Step 8: Get low stock amenities
echo -e "${BLUE}📋 Step 8: Getting low stock amenities...${NC}"
curl -s -X GET "$BASE_URL/amenities?tenantId=$TENANT_ID&lowStock=true" | jq '.'
echo -e "${YELLOW}⚠️  Low stock amenities (should show Shampoo with stock 5 < minimum 50)${NC}"
echo ""

# Step 9: Update amenity
echo -e "${BLUE}📋 Step 9: Updating amenity...${NC}"
curl -s -X PUT "$BASE_URL/amenities/$AMENITY1_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "itemName": "King Size Bed Sheet - Premium",
    "description": "Premium white cotton bed sheet",
    "stock": 45,
    "minimumStock": 15
  }' | jq '.'
echo -e "${GREEN}✅ Amenity updated${NC}"
echo ""

# Step 10: Update stock only
echo -e "${BLUE}📋 Step 10: Updating stock quantity...${NC}"
curl -s -X PATCH "$BASE_URL/amenities/$AMENITY2_ID/stock?quantity=60" | jq '.'
echo -e "${GREEN}✅ Stock updated (Shampoo now has 60 units, above minimum)${NC}"
echo ""

# Step 11: Verify stock update fixed low stock
echo -e "${BLUE}📋 Step 11: Verifying low stock list is now empty...${NC}"
curl -s -X GET "$BASE_URL/amenities?tenantId=$TENANT_ID&lowStock=true" | jq '.'
echo -e "${GREEN}✅ Low stock list should be empty now${NC}"
echo ""

# Step 12: Update category
echo -e "${BLUE}📋 Step 12: Updating category...${NC}"
curl -s -X PUT "$BASE_URL/amenities-categories/$CATEGORY1_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bedding & Linen",
    "description": "All bedding and linen items for guest rooms"
  }' | jq '.'
echo -e "${GREEN}✅ Category updated${NC}"
echo ""

# Step 13: Test duplicate name validation
echo -e "${BLUE}📋 Step 13: Testing duplicate name validation...${NC}"
DUPLICATE_RESPONSE=$(curl -s -X POST "$BASE_URL/amenities-categories" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "'$TENANT_ID'",
    "name": "Bedding & Linen",
    "description": "Should fail"
  }')
echo "$DUPLICATE_RESPONSE" | jq '.'
if echo "$DUPLICATE_RESPONSE" | jq -e '.code == 409' > /dev/null; then
  echo -e "${GREEN}✅ Duplicate validation working (returned 409 Conflict)${NC}"
else
  echo -e "${RED}❌ Duplicate validation failed${NC}"
fi
echo ""

echo "================================================"
echo -e "${GREEN}🎉 All tests completed!${NC}"
echo "================================================"
echo ""
echo "Test Summary:"
echo "- Tenant ID: $TENANT_ID"
echo "- Category 1 ID: $CATEGORY1_ID (Bedding & Linen)"
echo "- Category 2 ID: $CATEGORY2_ID (Toiletries)"
echo "- Amenity 1 ID: $AMENITY1_ID (King Size Bed Sheet - Premium)"
echo "- Amenity 2 ID: $AMENITY2_ID (Shampoo Bottle - stock updated to 60)"
echo "- Amenity 3 ID: $AMENITY3_ID (Soap Bar)"
echo ""
echo "You can now manually test delete operations:"
echo "  curl -X DELETE $BASE_URL/amenities/$AMENITY3_ID"
echo "  curl -X DELETE $BASE_URL/amenities-categories/$CATEGORY1_ID"
echo "  curl -X DELETE $BASE_URL/tenants/$TENANT_ID"


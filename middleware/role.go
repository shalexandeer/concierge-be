package middleware

import (
	"net/http"
	"strings"

	"concierge-be/utils"
	"github.com/gin-gonic/gin"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	RoleID   string `json:"roleId"`
	RoleName string `json:"roleName"`
	TenantID string `json:"tenantId"`
}

// RequireRole middleware that checks if the user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user info from JWT token (assuming it's set by auth middleware)
		_, exists := c.Get("userID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		roleName, exists := c.Get("roleName")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User role not found")
			c.Abort()
			return
		}

		userRole := roleName.(string)
		
		// Check if user has one of the required roles
		hasRequiredRole := false
		for _, role := range roles {
			if userRole == role {
				hasRequiredRole = true
				break
			}
		}

		if !hasRequiredRole {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin middleware that requires super admin role
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole("super_admin")
}

// RequireTenantAdmin middleware that requires tenant admin or super admin role
func RequireTenantAdmin() gin.HandlerFunc {
	return RequireRole("super_admin", "tenant_admin")
}

// RequireUser middleware that requires any authenticated user (user, tenant_admin, or super_admin)
func RequireUser() gin.HandlerFunc {
	return RequireRole("super_admin", "tenant_admin", "user")
}

// ExtractUserInfo extracts user information from JWT token and sets it in context
func ExtractUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate JWT token
		claims, err := utils.ParseJWT(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Set("roleName", claims.RoleName)
		c.Set("tenantID", claims.TenantID)

		c.Next()
	}
}

// RequireTenantAccess middleware that ensures user can only access their tenant's data
func RequireTenantAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName, exists := c.Get("roleName")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User role not found")
			c.Abort()
			return
		}

		// Super admin can access all tenants
		if roleName.(string) == "super_admin" {
			c.Next()
			return
		}

		// For tenant admin and user, check tenant access
		userTenantID, exists := c.Get("tenantID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User tenant not found")
			c.Abort()
			return
		}

		// Get tenant ID from URL parameter
		tenantID := c.Param("tenantId")
		if tenantID == "" {
			// Try alternative parameter names
			tenantID = c.Param("tenant_id")
		}

		// If no tenant ID in URL, allow access (for global operations)
		if tenantID == "" {
			c.Next()
			return
		}

		// Check if user's tenant matches the requested tenant
		if userTenantID.(string) != tenantID {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied to this tenant")
			c.Abort()
			return
		}

		c.Next()
	}
}

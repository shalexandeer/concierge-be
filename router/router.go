package router

import (
	"concierge-be/internal/amenities"
	"concierge-be/internal/amenities_categories"
	"concierge-be/internal/facilities"
	"concierge-be/internal/roles"
	"concierge-be/internal/tenants"
	"concierge-be/internal/users"
	"concierge-be/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// API 版本分组
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Service is running",
			})
		})

		// Auth routes (no authentication required)
		userHandler := users.NewHandler()
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
		}

		// User routes (protected by super admin role)
		userRoutes := v1.Group("/users")
		userRoutes.Use(middleware.JWTAuth())
		{
			userRoutes.POST("", userHandler.CreateUser)
			userRoutes.GET("/:id", userHandler.GetUser)
			userRoutes.GET("", userHandler.GetAllUsers)
			userRoutes.PUT("/:id", userHandler.UpdateUser)
			userRoutes.DELETE("/:id", userHandler.DeleteUser)
		}

		// Roles routes (protected by super admin role)
		roleHandler := roles.NewHandler()
		roleRoutes := v1.Group("/roles")
		roleRoutes.Use(middleware.ExtractUserInfo())
		roleRoutes.Use(middleware.RequireSuperAdmin())
		{
			roleRoutes.GET("", roleHandler.GetAllRoles)
			roleRoutes.GET("/:id", roleHandler.GetRole)
		}

		// Authenticated routes
		authenticated := v1.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// Current user routes
			authenticated.GET("/me", userHandler.GetCurrentUser)
			authenticated.PUT("/me", userHandler.UpdateCurrentUser)
		}

		// Tenant routes (protected by super admin role)
		tenantHandler := tenants.NewHandler()
		tenantRoutes := v1.Group("/tenants")
		tenantRoutes.Use(middleware.ExtractUserInfo())
		tenantRoutes.Use(middleware.RequireSuperAdmin())
		{
			tenantRoutes.POST("", tenantHandler.CreateTenant)
			tenantRoutes.GET("/:id", tenantHandler.GetTenant)
			tenantRoutes.GET("", tenantHandler.GetAllTenants)
			tenantRoutes.PUT("/:id", tenantHandler.UpdateTenant)
			tenantRoutes.DELETE("/:id", tenantHandler.DeleteTenant)
		}

		// User-Tenant relationship routes
		userTenantRoutes := v1.Group("/user-tenants")
		{
			userTenantRoutes.POST("", userHandler.AddUserToTenant)
			userTenantRoutes.GET("/users/:userId", userHandler.GetUserTenants)
			userTenantRoutes.GET("/tenants/:tenantId", userHandler.GetTenantUsers)
			userTenantRoutes.DELETE("/users/:userId/tenants/:tenantId", userHandler.RemoveUserFromTenant)
		}

		// Amenity Categories routes
		categoriesHandler := amenities_categories.NewHandler()
		categoriesRoutes := v1.Group("/amenities-categories")
		{
			categoriesRoutes.POST("", categoriesHandler.CreateCategory)
			categoriesRoutes.GET("/:id", categoriesHandler.GetCategory)
			categoriesRoutes.GET("", categoriesHandler.GetAllCategories)
			categoriesRoutes.PUT("/:id", categoriesHandler.UpdateCategory)
			categoriesRoutes.DELETE("/:id", categoriesHandler.DeleteCategory)
		}

		// Amenities routes (protected by authentication)
		amenitiesHandler := amenities.NewHandler()
		amenitiesRoutes := v1.Group("/amenities")
		amenitiesRoutes.Use(middleware.JWTAuth())
		{
			amenitiesRoutes.POST("", amenitiesHandler.CreateAmenity)
			amenitiesRoutes.GET("/:id", amenitiesHandler.GetAmenity)
			amenitiesRoutes.GET("", amenitiesHandler.GetAllAmenities)
			amenitiesRoutes.PUT("/:id", amenitiesHandler.UpdateAmenity)
			amenitiesRoutes.PATCH("/:id/stock", amenitiesHandler.UpdateStock)
			amenitiesRoutes.DELETE("/:id", amenitiesHandler.DeleteAmenity)
		}

		// Public facilities routes (protected by authentication)
		facilitiesHandler := facilities.NewHandler()
		publicFacilitiesRoutes := v1.Group("/facilities")
		publicFacilitiesRoutes.Use(middleware.JWTAuth())
		{
			publicFacilitiesRoutes.GET("/:id", facilitiesHandler.GetFacility)
			publicFacilitiesRoutes.GET("", facilitiesHandler.GetAllFacilities)
			publicFacilitiesRoutes.GET("/tenant/:tenantId", facilitiesHandler.GetFacilitiesByTenant)
		}

		// Facilities management routes (protected by tenant admin role)
		facilitiesManagementRoutes := v1.Group("/facilities")
		facilitiesManagementRoutes.Use(middleware.ExtractUserInfo())
		facilitiesManagementRoutes.Use(middleware.RequireTenantAdmin())
		{
			facilitiesManagementRoutes.POST("", facilitiesHandler.CreateFacility)
			facilitiesManagementRoutes.PUT("/:id", facilitiesHandler.UpdateFacility)
			facilitiesManagementRoutes.DELETE("/:id", facilitiesHandler.DeleteFacility)
		}

		// Facility Bookings routes (accessible to all authenticated users)
		bookingRoutes := v1.Group("/facilities/:id/bookings")
		bookingRoutes.Use(middleware.ExtractUserInfo())
		bookingRoutes.Use(middleware.RequireUser())
		{
			bookingRoutes.POST("", facilitiesHandler.CreateBooking)
			bookingRoutes.GET("", facilitiesHandler.GetFacilityBookings)
			bookingRoutes.GET("/upcoming", facilitiesHandler.GetUpcomingBookings)
			bookingRoutes.GET("/:bookingId", facilitiesHandler.GetBooking)
			bookingRoutes.DELETE("/:bookingId", facilitiesHandler.DeleteBooking)
		}

		// Tenant Bookings routes (protected by tenant admin role)
		tenantBookingRoutes := v1.Group("/tenant-bookings")
		tenantBookingRoutes.Use(middleware.ExtractUserInfo())
		tenantBookingRoutes.Use(middleware.RequireTenantAdmin())
		{
			tenantBookingRoutes.GET("/:tenantId", facilitiesHandler.GetBookingsByTenant)
		}
	}

	return r
}

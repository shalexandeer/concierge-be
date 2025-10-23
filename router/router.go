package router

import (
	"concierge-be/internal/amenities"
	"concierge-be/internal/amenities_categories"
	"concierge-be/internal/facilities"
	"concierge-be/internal/food_beverages"
	"concierge-be/internal/food_beverages_categories"
	"concierge-be/internal/roles"
	"concierge-be/internal/services"
	"concierge-be/internal/services_categories"
	"concierge-be/internal/tenants"
	"concierge-be/internal/upload"
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

		// Service Categories routes (public GET, protected POST/PUT/DELETE)
		serviceCategoriesHandler := services_categories.NewHandler()
		serviceCategoriesRoutes := v1.Group("/services-categories")
		{
			// Public routes (no authentication required)
			serviceCategoriesRoutes.GET("", serviceCategoriesHandler.GetAllServiceCategories)
			serviceCategoriesRoutes.GET("/:id", serviceCategoriesHandler.GetServiceCategory)
		}

		// Service Categories management routes (protected by super_admin + tenant_admin)
		serviceCategoriesManagementRoutes := v1.Group("/services-categories")
		serviceCategoriesManagementRoutes.Use(middleware.ExtractUserInfo())
		serviceCategoriesManagementRoutes.Use(middleware.RequireTenantAdmin())
		{
			serviceCategoriesManagementRoutes.POST("", serviceCategoriesHandler.CreateServiceCategory)
			serviceCategoriesManagementRoutes.PUT("/:id", serviceCategoriesHandler.UpdateServiceCategory)
			serviceCategoriesManagementRoutes.DELETE("/:id", serviceCategoriesHandler.DeleteServiceCategory)
		}

		// Services routes (public GET, protected POST/PUT/DELETE)
		servicesHandler := services.NewHandler()
		servicesRoutes := v1.Group("/services")
		{
			// Public routes (no authentication required)
			servicesRoutes.GET("", servicesHandler.GetAllServices)
			servicesRoutes.GET("/:id", servicesHandler.GetService)
		}

		// Services management routes (protected by super_admin + tenant_admin)
		servicesManagementRoutes := v1.Group("/services")
		servicesManagementRoutes.Use(middleware.ExtractUserInfo())
		servicesManagementRoutes.Use(middleware.RequireTenantAdmin())
		{
			servicesManagementRoutes.POST("", servicesHandler.CreateService)
			servicesManagementRoutes.PUT("/:id", servicesHandler.UpdateService)
			servicesManagementRoutes.DELETE("/:id", servicesHandler.DeleteService)
		}

		// Food & Beverage Categories routes (public GET, protected POST/PUT/DELETE)
		fbCategoriesHandler := food_beverages_categories.NewHandler()
		fbCategoriesRoutes := v1.Group("/food-beverages-categories")
		{
			// Public routes (no authentication required)
			fbCategoriesRoutes.GET("", fbCategoriesHandler.GetAllCategories)
			fbCategoriesRoutes.GET("/:id", fbCategoriesHandler.GetCategory)
		}

		// Food & Beverage Categories management routes (protected by super_admin + tenant_admin)
		fbCategoriesManagementRoutes := v1.Group("/food-beverages-categories")
		fbCategoriesManagementRoutes.Use(middleware.ExtractUserInfo())
		fbCategoriesManagementRoutes.Use(middleware.RequireTenantAdmin())
		{
			fbCategoriesManagementRoutes.POST("", fbCategoriesHandler.CreateCategory)
			fbCategoriesManagementRoutes.PUT("/:id", fbCategoriesHandler.UpdateCategory)
			fbCategoriesManagementRoutes.DELETE("/:id", fbCategoriesHandler.DeleteCategory)
		}

		// Food & Beverage Items routes (public GET, protected POST/PUT/DELETE)
		fbItemsHandler := food_beverages.NewHandler()
		fbItemsRoutes := v1.Group("/food-beverages")
		{
			// Public routes (no authentication required)
			fbItemsRoutes.GET("", fbItemsHandler.GetAllFoodBeverages)
			fbItemsRoutes.GET("/:id", fbItemsHandler.GetFoodBeverage)
		}

		// Food & Beverage Items management routes (protected by super_admin + tenant_admin)
		fbItemsManagementRoutes := v1.Group("/food-beverages")
		fbItemsManagementRoutes.Use(middleware.ExtractUserInfo())
		fbItemsManagementRoutes.Use(middleware.RequireTenantAdmin())
		{
			fbItemsManagementRoutes.POST("", fbItemsHandler.CreateFoodBeverage)
			fbItemsManagementRoutes.PUT("/:id", fbItemsHandler.UpdateFoodBeverage)
			fbItemsManagementRoutes.DELETE("/:id", fbItemsHandler.DeleteFoodBeverage)
		}

		// Upload routes (protected by authentication)
		uploadHandler := upload.NewHandler()
		uploadRoutes := v1.Group("/uploads")
		uploadRoutes.Use(middleware.JWTAuth())
		{
			uploadRoutes.POST("/images", uploadHandler.UploadImage)
			uploadRoutes.DELETE("/images/:filename", uploadHandler.DeleteImage)
		}

		// Public image serving route (no authentication required)
		v1.GET("/uploads/images/:filename", uploadHandler.ServeImage)
	}

	return r
}

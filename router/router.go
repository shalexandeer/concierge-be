package router

import (
	"concierge-be/internal/amenities"
	"concierge-be/internal/amenities_categories"
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

		// User routes
		userRoutes := v1.Group("/users")
		{
			userRoutes.POST("", userHandler.CreateUser)
			userRoutes.GET("/:id", userHandler.GetUser)
			userRoutes.GET("", userHandler.GetAllUsers)
			userRoutes.PUT("/:id", userHandler.UpdateUser)
			userRoutes.DELETE("/:id", userHandler.DeleteUser)
		}

		// Authenticated routes
		authenticated := v1.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// Current user routes
			authenticated.GET("/me", userHandler.GetCurrentUser)
			authenticated.PUT("/me", userHandler.UpdateCurrentUser)
		}

		// Tenant routes
		tenantHandler := tenants.NewHandler()
		tenantRoutes := v1.Group("/tenants")
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

		// Amenities routes
		amenitiesHandler := amenities.NewHandler()
		amenitiesRoutes := v1.Group("/amenities")
		{
			amenitiesRoutes.POST("", amenitiesHandler.CreateAmenity)
			amenitiesRoutes.GET("/:id", amenitiesHandler.GetAmenity)
			amenitiesRoutes.GET("", amenitiesHandler.GetAllAmenities)
			amenitiesRoutes.PUT("/:id", amenitiesHandler.UpdateAmenity)
			amenitiesRoutes.PATCH("/:id/stock", amenitiesHandler.UpdateStock)
			amenitiesRoutes.DELETE("/:id", amenitiesHandler.DeleteAmenity)
		}
	}

	return r
}

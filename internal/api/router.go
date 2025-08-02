package api

import (
	"github.com/mnadev/stok/internal/api/handlers"
	"github.com/mnadev/stok/internal/api/middleware"
	"github.com/mnadev/stok/internal/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, db *database.DB) {
	// Create handlers
	authHandler := handlers.NewAuthHandler(db)
	inventoryHandler := handlers.NewInventoryHandler(db)

	// Public routes for authentication
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// API v1 routes, protected by JWT middleware
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		// User routes
		v1.GET("/me", authHandler.GetCurrentUser)

		// Invitation routes (for account admins)
		v1.GET("/accounts/:account_id/invitations", authHandler.GetInvitationsByAccount)
		v1.POST("/accounts/:account_id/invitations", authHandler.CreateInvitation)
		v1.DELETE("/accounts/:account_id/invitations/:invitation_id", authHandler.DeleteInvitation)

		// Inventory item routes
		v1.GET("/inventory/items", inventoryHandler.GetInventoryItems)
		v1.GET("/inventory/items/low-stock", inventoryHandler.GetLowStockItems)
		v1.POST("/inventory/items", inventoryHandler.CreateInventoryItem)
		v1.GET("/inventory/items/:id", inventoryHandler.GetInventoryItem)
		v1.PUT("/inventory/items/:id", inventoryHandler.UpdateInventoryItem)
		v1.DELETE("/inventory/items/:id", inventoryHandler.DeleteInventoryItem)

		// Menu item routes
		v1.GET("/menu/items", inventoryHandler.GetMenuItems)
		v1.POST("/menu/items", inventoryHandler.CreateMenuItem)

		// Delivery routes
		v1.GET("/deliveries", inventoryHandler.GetDeliveries)
		v1.POST("/deliveries", inventoryHandler.LogDelivery)
		v1.GET("/deliveries/vendor/:vendor", inventoryHandler.GetDeliveriesByVendor)

		// Vendor routes
		v1.GET("/inventory/vendor/:vendor", inventoryHandler.GetInventoryItemsByVendor)

		// Snapshot routes for inventory counts
		v1.GET("/snapshots", inventoryHandler.GetInventorySnapshots)
		v1.POST("/snapshots", inventoryHandler.CreateInventorySnapshot)
	}
}

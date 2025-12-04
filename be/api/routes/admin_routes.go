package routes

import (
	"milestone3/be/api/middleware"
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterAdminRoutes(adminCtrl *controller.AdminController) {
	adminRoutes := r.echo.Group("admin")
	adminRoutes.Use(middleware.JWTMiddleware)
	adminRoutes.Use(middleware.LoggingMiddleware)

	//admin endpoint
	adminRoutes.GET("/dashboard", adminCtrl.AdminDashboard)
	// adminRoutes.GET("/reports", adminCtrl.AdminReport)
}
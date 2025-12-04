package routes

import (
	"milestone3/be/api/middleware"
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterUserRoutes(userCtrl *controller.UserController) {
	userRoutes := r.echo.Group("auth")
	userRoutes.Use(middleware.LoggingMiddleware)

	// auth endpoint
	userRoutes.POST("/register", userCtrl.CreateUser)
	userRoutes.POST("/login", userCtrl.LoginUser)
}

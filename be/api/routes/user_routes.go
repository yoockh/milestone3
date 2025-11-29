package routes

import (
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterUserRoutes(userCtrl *controller.UserController) {
	userRoutes := r.echo.Group("auth")

	// auth endpoint
	userRoutes.POST("/register", userCtrl.CreateUser)
	userRoutes.GET("/login", userCtrl.LoginUser)
}
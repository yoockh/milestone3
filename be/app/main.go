package main

import (
	"context"
	"log"
	"milestone3/be/api/routes"
	"milestone3/be/config"
	"milestone3/be/internal/controller"
	"milestone3/be/internal/repository"
	"milestone3/be/internal/service"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	db := config.ConnectionDb()
	validator := validator.New()
	ctx := context.Background()
	
	//dependency injection
	//repository
	userRepo := repository.NewUserRepo(db, ctx)
	
	//service
	userServ := service.NewUserService(userRepo)
	
	//controller
	userControl := controller.NewUserController(validator, userServ)
	
	//echo
	e := echo.New()
	//router
	router := routes.NewRouter(e)
	router.RegisterUserRoutes(userControl)

	address := os.Getenv("PORT")
	if err := e.Start(":" + address); err != nil {
		log.Printf("faile to start server %s", err)
	}
}
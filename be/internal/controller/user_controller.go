package controller

import (
	"milestone3/be/internal/dto"
	"milestone3/be/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	CreateUser(req dto.UserRequest) (res dto.UserResponse, err error)
	GetUserByEmail(email, password string) (accessToken string, err error)
}

type UserController struct {
	userService UserService
	validate *validator.Validate
}

func NewUserController(validate *validator.Validate, us UserService) *UserController {
	return &UserController{validate: validate, userService: us}
}

func (uc *UserController) CreateUser(c echo.Context) error {
	req := new(dto.UserRequest)

	if err := c.Bind(req); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := uc.validate.Struct(req); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	resp, err := uc.userService.CreateUser(*req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.CreatedResponse(c, "user created", resp)
}

func (uc *UserController) LoginUser(c echo.Context) error {
	req := new(dto.UserLoginRequest)

	if err := c.Bind(req); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := uc.validate.Struct(req); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	resp, err := uc.userService.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "success login", resp)
}
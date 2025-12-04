package controller

import (
	"errors"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/service"
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
	validate    *validator.Validate
}

func NewUserController(validate *validator.Validate, us UserService) *UserController {
	return &UserController{validate: validate, userService: us}
}

// CreateUser godoc
// @Summary Register new user
// @Description Register a new user account in the system
// @Tags Your Donate Rise API - Authentication
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "User registration data"
// @Success 201 {object} utils.SuccessResponseData{data=dto.UserResponse} "user created"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid payload or validation error"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/register [post]
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

// LoginUser godoc
// @Summary User login
// @Description Authenticate user and return access token
// @Tags Your Donate Rise API - Authentication
// @Accept json
// @Produce json
// @Param credentials body dto.UserLoginRequest true "User login credentials"
// @Success 200 {object} utils.SuccessResponseData{data=string} "success login"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid credentials format"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid email or password"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/login [post]
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
		if errors.Is(err, service.ErrInvalidCredential) {
			return utils.UnauthorizedResponse(c, "invalid email or password")
		}
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "success login", resp)
}

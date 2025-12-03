package controller

import (
	"milestone3/be/internal/dto"
	"milestone3/be/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AdminService interface {
	AdminDashboard() (resp dto.AdminDashboardResponse, err error)
	// AdminReport() (err error)
}

type AdminController struct {
	adminService AdminService
}

func NewAdminController(as AdminService) *AdminController {
	return &AdminController{adminService: as}
}

func (ac *AdminController) AdminDashboard(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	role := claim["role"].(string)

	if role != "admin" {
		return utils.ForbiddenResponse(c, "forbidden request")
	}
	
	resp, err := ac.adminService.AdminDashboard()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "ok", resp)
}
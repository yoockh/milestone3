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

// AdminDashboard godoc
// @Summary Get admin dashboard analytics
// @Description Get comprehensive dashboard analytics including donation stats, auction metrics, and system overview
// @Tags Your Donate Rise API - Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseData "ok"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /admin/dashboard [get]
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

// WIP
// func (ac *AdminController) AdminReport(c echo.Context) error {

// }
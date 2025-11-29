package controller

import (
	"milestone3/be/internal/service"
	"milestone3/be/internal/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FinalDonationController struct {
	svc service.FinalDonationService
}

func NewFinalDonationController(finalDonationService service.FinalDonationService) *FinalDonationController {
	return &FinalDonationController{svc: finalDonationService}
}

func (h *FinalDonationController) GetAllFinalDonations(c echo.Context) error {
	finalDonations, err := h.svc.GetAllFinalDonations()
	if err != nil {
		return utils.BadRequestResponse(c, "Failed to fetch final donations")
	}
	return utils.SuccessResponse(c, "Final donations fetched successfully", finalDonations)
}
func (h *FinalDonationController) GetAllFinalDonationsByUserID(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user id")
	}
	finalDonations, err := h.svc.GetAllFinalDonationsByUserID(userID)
	if err != nil {
		return utils.BadRequestResponse(c, "Failed to fetch final donations for the user")
	}
	return utils.SuccessResponse(c, "Final donations fetched successfully", finalDonations)
}

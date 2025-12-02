package controller

import (
	"milestone3/be/internal/dto"
	"milestone3/be/internal/service"
	"milestone3/be/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuctionSessionController struct {
	svc      service.AuctionSessionService
	validate *validator.Validate
}

func NewAuctionSessionController(s service.AuctionSessionService, validate *validator.Validate) *AuctionSessionController {
	return &AuctionSessionController{svc: s, validate: validate}
}

func (h *AuctionSessionController) CreateAuctionSession(c echo.Context) error {
	if !utils.IsAdmin(c) {
		return utils.ForbiddenResponse(c, "only admin can create auction sessions")
	}

	var payload dto.AuctionSessionDTO
	if err := c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}

	if err := h.validate.Struct(payload); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	createdSession, err := h.svc.Create(&payload)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to create auction session")
	}

	return utils.CreatedResponse(c, "auction session created successfully", createdSession)
}

func (h *AuctionSessionController) GetAuctionSessionByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction session ID")
	}

	session, err := h.svc.GetByID(id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to retrieve auction session")
	}

	return utils.SuccessResponse(c, "auction session retrieved successfully", session)
}

func (h *AuctionSessionController) GetAllAuctionSessions(c echo.Context) error {
	sessions, err := h.svc.GetAll()
	if err != nil {
		switch err {
		case service.ErrAuctionNotFound:
			return utils.NotFoundResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "failed to retrieve auction sessions")
		}
	}

	return utils.SuccessResponse(c, "auction sessions retrieved successfully", sessions)
}

func (h *AuctionSessionController) UpdateAuctionSession(c echo.Context) error {
	if !utils.IsAdmin(c) {
		return utils.ForbiddenResponse(c, "only admin can update auction sessions")
	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction session ID")
	}

	var payload dto.AuctionSessionDTO
	if err = c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}

	if err = h.validate.Struct(payload); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	updatedSession, err := h.svc.Update(id, &payload)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to update auction session")
	}

	return utils.SuccessResponse(c, "auction session updated successfully", updatedSession)
}

func (h *AuctionSessionController) DeleteAuctionSession(c echo.Context) error {
	if !utils.IsAdmin(c) {
		return utils.ForbiddenResponse(c, "only admin can delete auction sessions")
	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction session ID")
	}

	err = h.svc.Delete(id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to delete auction session")
	}

	return utils.SuccessResponse(c, "auction session deleted successfully", nil)
}

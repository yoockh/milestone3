package controller

import (
	"milestone3/be/internal/dto"
	"milestone3/be/internal/service"
	"milestone3/be/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuctionSessionController struct {
	svc      service.AuctionSessionService
	validate *validator.Validate
}

func NewAuctionSessionController(s service.AuctionSessionService, validate *validator.Validate) *AuctionSessionController {
	return &AuctionSessionController{svc: s, validate: validate}
}

func isAdminFromTokenSession(c echo.Context) bool {
	token := c.Get("user")
	if token == nil {
		return false
	}

	claims, ok := token.(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	role, ok := claims["role"].(string)
	return ok && role == "admin"
}

// CreateAuctionSession godoc
// @Summary Create new auction session
// @Description Create a new auction session with start and end times
// @Tags Your Donate Rise API - Auction Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param auctionSession body dto.AuctionSessionDTO true "Auction session data"
// @Success 201 {object} utils.SuccessResponseData "auction session created successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid payload"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions [post]
func (h *AuctionSessionController) CreateAuctionSession(c echo.Context) error {
	if !isAdminFromTokenSession(c) {
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

// GetAuctionSessionByID godoc
// @Summary Get auction session by ID
// @Description Retrieve a specific auction session by its ID
// @Tags Your Donate Rise API - Auction Sessions
// @Accept json
// @Produce json
// @Param id path int true "Auction Session ID"
// @Success 200 {object} utils.SuccessResponseData "auction session retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid auction session ID"
// @Failure 404 {object} utils.ErrorResponse "Auction session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions/{id} [get]
func (h *AuctionSessionController) GetAuctionSessionByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction session ID")
	}

	session, err := h.svc.GetByID(id)
	if err != nil {
		switch err {
		case service.ErrAuctionNotFound:
			return utils.NotFoundResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "failed to retrieve auction session")
		}
	}

	return utils.SuccessResponse(c, "auction session retrieved successfully", session)
}

// GetAllAuctionSessions godoc
// @Summary Get all auction sessions
// @Description Retrieve all auction sessions
// @Tags Your Donate Rise API - Auction Sessions
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseData "auction sessions retrieved successfully"
// @Failure 404 {object} utils.ErrorResponse "No auction sessions found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions [get]
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

// UpdateAuctionSession godoc
// @Summary Update auction session
// @Description Update an existing auction session by ID
// @Tags Your Donate Rise API - Auction Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Auction Session ID"
// @Param auctionSession body dto.AuctionSessionDTO true "Updated auction session data"
// @Success 200 {object} utils.SuccessResponseData "auction session updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid ID, payload, or session is active"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 404 {object} utils.ErrorResponse "Auction session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions/{id} [put]
func (h *AuctionSessionController) UpdateAuctionSession(c echo.Context) error {
	if !isAdminFromTokenSession(c) {
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
		switch err {
		case service.ErrAuctionNotFoundID:
			return utils.NotFoundResponse(c, err.Error())
		case service.ErrActiveSession:
			return utils.BadRequestResponse(c, err.Error())
		case service.ErrInvalidDate:
			return utils.BadRequestResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "failed to update auction session")
		}
	}

	return utils.SuccessResponse(c, "auction session updated successfully", updatedSession)
}

// DeleteAuctionSession godoc
// @Summary Delete auction session
// @Description Delete an auction session by ID
// @Tags Your Donate Rise API - Auction Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Auction Session ID"
// @Success 200 {object} utils.SuccessResponseData "auction session deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid ID or session is active"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 404 {object} utils.ErrorResponse "Auction session not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions/{id} [delete]
func (h *AuctionSessionController) DeleteAuctionSession(c echo.Context) error {
	if !isAdminFromTokenSession(c) {
		return utils.ForbiddenResponse(c, "only admin can delete auction sessions")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction session ID")
	}

	err = h.svc.Delete(id)
	if err != nil {
		switch err {
		case service.ErrAuctionNotFoundID:
			return utils.NotFoundResponse(c, err.Error())
		case service.ErrActiveSession, service.ErrInvalidDate:
			return utils.BadRequestResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "failed to delete auction session")
		}
	}

	return utils.SuccessResponse(c, "auction session deleted successfully", nil)
}

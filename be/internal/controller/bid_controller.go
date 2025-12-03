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

type BidController struct {
	svc        service.BidService
	sessionSvc service.AuctionSessionService
	validate   *validator.Validate
}

func NewBidController(s service.BidService, sessionSvc service.AuctionSessionService, validate *validator.Validate) *BidController {
	return &BidController{svc: s, sessionSvc: sessionSvc, validate: validate}
}

func getUserIDFromToken(c echo.Context) (int64, error) {
	token := c.Get("user")
	if token == nil {
		return 0, echo.NewHTTPError(401, "unauthenticated")
	}

	claims, ok := token.(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return 0, echo.NewHTTPError(401, "invalid token")
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, echo.NewHTTPError(401, "invalid token")
	}

	return int64(userIDFloat), nil
}

func isAdminFromToken(c echo.Context) bool {
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

func (h *BidController) PlaceBid(c echo.Context) error {
	sessionIDStr := c.Param("sessionID")
	itemIDStr := c.Param("itemID")

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid sessionID")
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid itemID")
	}

	var payload dto.BidDTO
	if err = c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}

	if err = h.validate.Struct(payload); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}

	session, err := h.sessionSvc.GetByID(sessionID)
	if err != nil {
		return utils.NotFoundResponse(c, "auction session not found")
	}
	err = h.svc.PlaceBid(
		sessionID,
		itemID,
		userID,
		payload.Amount,
		session.EndTime,
	)

	if err != nil {
		switch err {
		case service.ErrBidTooLow, service.ErrInvalidBidding:
			return utils.BadRequestResponse(c, err.Error())
		case service.ErrAuctionNotFound:
			return utils.NotFoundResponse(c, err.Error())
		case service.ErrInvalidAuction:
			return utils.ConflictResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "failed placing bid")
		}
	}

	return utils.SuccessResponse(c, "bid placed successfully", nil)
}

func (h *BidController) GetHighestBid(c echo.Context) error {
	sessionIDStr := c.Param("sessionID")
	itemIDStr := c.Param("itemID")

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid sessionID")
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid itemID")
	}

	highest, bidder, err := h.svc.GetHighestBid(sessionID, itemID)
	if err != nil {
		switch err {
		case service.ErrAuctionNotFound:
			return utils.NotFoundResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "failed retrieving highest bid")
		}
	}

	resp := map[string]interface{}{
		"session_id":  sessionID,
		"item_id":     itemID,
		"highest_bid": highest,
		"bidder_id":   bidder,
	}

	return utils.SuccessResponse(c, "highest bid retrieved successfully", resp)
}

func (h *BidController) SyncHighestBid(c echo.Context) error {
	if !isAdminFromToken(c) {
		return utils.ForbiddenResponse(c, "only admin can sync bids")
	}

	sessionIDStr := c.Param("sessionID")
	itemIDStr := c.Param("itemID")

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid sessionID")
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid itemID")
	}

	err = h.svc.SyncHighestBid(sessionID, itemID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "highest bid synced to database", nil)
}

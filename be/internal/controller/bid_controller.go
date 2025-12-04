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

// PlaceBid godoc
// @Summary Place bid on auction item
// @Description Place a bid on a specific auction item within an active session
// @Tags Your Donate Rise API - Bidding
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionID path int true "Auction Session ID"
// @Param itemID path int true "Auction Item ID"
// @Param bid body dto.BidDTO true "Bid amount"
// @Success 200 {object} utils.SuccessResponseData "bid placed successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid parameters or bid too low"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 404 {object} utils.ErrorResponse "Auction session or item not found"
// @Failure 409 {object} utils.ErrorResponse "Conflict - Invalid auction state"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions/{sessionID}/items/{itemID}/bid [post]
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

// GetHighestBid godoc
// @Summary Get highest bid for auction item
// @Description Retrieve the current highest bid for a specific auction item
// @Tags Your Donate Rise API - Bidding
// @Accept json
// @Produce json
// @Param sessionID path int true "Auction Session ID"
// @Param itemID path int true "Auction Item ID"
// @Success 200 {object} utils.SuccessResponseData "highest bid retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid session or item ID"
// @Failure 404 {object} utils.ErrorResponse "Auction not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions/{sessionID}/items/{itemID}/highest-bid [get]
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

// SyncHighestBid godoc
// @Summary Sync highest bid to database
// @Description Synchronize the highest bid from Redis cache to the database
// @Tags Your Donate Rise API - Bidding
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionID path int true "Auction Session ID"
// @Param itemID path int true "Auction Item ID"
// @Success 200 {object} utils.SuccessResponseData "highest bid synced to database"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid session or item ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auction/sessions/{sessionID}/items/{itemID}/sync [post]
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

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

type AuctionController struct {
	svc      service.AuctionItemService
	validate *validator.Validate
}

func NewAuctionController(s service.AuctionItemService, validate *validator.Validate) *AuctionController {
	return &AuctionController{svc: s, validate: validate}
}

func getUserIDFromTokenItem(c echo.Context) (int64, error) {
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

// Helper function to check if user is admin from JWT token
func isAdminFromTokenItem(c echo.Context) bool {
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

func (h *AuctionController) CreateAuctionItem(c echo.Context) error {
	// Check if user is admin
	if !isAdminFromTokenItem(c) {
		return utils.ForbiddenResponse(c, "only admin can create auction items")
	}

	// Get user ID from JWT token
	userID, err := getUserIDFromTokenItem(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}

	var payload dto.AuctionItemDTO
	if err := c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}

	if err := h.validate.Struct(payload); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	payload.UserID = userID
	createdItem, err := h.svc.Create(&payload)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed creating auction item")
	}

	return utils.CreatedResponse(c, "auction item created successfully", createdItem)
}

func (h *AuctionController) GetAllAuctionItems(c echo.Context) error {
	items, err := h.svc.GetAll()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed retrieving auction items")
	}
	return utils.SuccessResponse(c, "auction items retrieved successfully", items)
}

func (h *AuctionController) GetAuctionItemByID(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction item ID")
	}

	item, err := h.svc.GetByID(id)
	if err != nil {
		switch err {
		case service.ErrAuctionNotFoundID:
			return utils.NotFoundResponse(c, "auction item not found")
		default:
			return utils.InternalServerErrorResponse(c, "failed retrieving auction item")
		}
	}

	return utils.SuccessResponse(c, "auction item retrieved successfully", item)
}

func (h *AuctionController) UpdateAuctionItem(c echo.Context) error {
	if !isAdminFromTokenItem(c) {
		return utils.ForbiddenResponse(c, "only admin can update auction items")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction item ID")
	}

	var payload dto.AuctionItemDTO
	if err = c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}

	if err = h.validate.Struct(payload); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	updatedItem, err := h.svc.Update(id, &payload)
	if err != nil {
		switch err {
		case service.ErrAuctionNotFoundID:
			return utils.NotFoundResponse(c, "auction item not found")
		case service.ErrInvalidAuction:
			return utils.BadRequestResponse(c, "invalid auction item data")
		default:
			return utils.InternalServerErrorResponse(c, "failed updating auction item")
		}
	}

	return utils.SuccessResponse(c, "auction item updated successfully", updatedItem)
}

func (h *AuctionController) DeleteAuctionItem(c echo.Context) error {
	// Check if user is admin
	if !isAdminFromTokenItem(c) {
		return utils.ForbiddenResponse(c, "only admin can delete auction items")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction item ID")
	}

	err = h.svc.Delete(id)
	if err != nil {
		switch err {
		case service.ErrAuctionNotFoundID:
			return utils.NotFoundResponse(c, "auction item not found")
		case service.ErrInvalidAuction:
			return utils.BadRequestResponse(c, "cannot delete auction item")
		default:
			return utils.InternalServerErrorResponse(c, "failed deleting auction item")
		}
	}

	return utils.SuccessResponse(c, "auction item deleted successfully", nil)
}

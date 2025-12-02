package controller

import (
	"milestone3/be/config"
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

func (h *AuctionController) CreateAuctionItem(c echo.Context) error {
	t := c.Get("user")
	if t == nil {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	user := t.(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	userID := int64(claim["id"].(float64))
	role := ""
	if r, ok := claim["role"].(string); ok {
		role = r
	}

	if role == "" {
		db := config.ConnectionDb()
		var roleName string
		if err := db.Raw("SELECT role FROM users WHERE id = ?", userID).Scan(&roleName).Error; err == nil {
			role = roleName
		}
	}

	if role != "admin" {
		return utils.ForbiddenResponse(c, "only admin can create auction items")
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
	t := c.Get("user")
	if t == nil {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	user := t.(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	userID := int64(claim["id"].(float64))
	role := ""
	if r, ok := claim["role"].(string); ok {
		role = r
	}

	if role == "" {
		db := config.ConnectionDb()
		var roleName string
		if err := db.Raw("SELECT role FROM users WHERE id = ?", userID).Scan(&roleName).Error; err == nil {
			role = roleName
		}
	}

	if role != "admin" {
		return utils.ForbiddenResponse(c, "only admin can create auction items")
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
	t := c.Get("user")
	if t == nil {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	user := t.(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	role := ""
	if r, ok := claim["role"].(string); ok {
		role = r
	}

	if role != "admin" {
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

package controller

import (
	"log"
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

func getUserIDTest(c echo.Context) (int64, bool) {
	user := c.Get("user")
	if user == nil {
		return 0, false
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return 0, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	// cek beberapa nama claim umum: "id" atau "user_id"
	if raw, ok := claims["id"]; ok && raw != nil {
		switch v := raw.(type) {
		case float64:
			return int64(v), true
		case string:
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
				return id, true
			}
		}
	}
	if raw, ok := claims["user_id"]; ok && raw != nil {
		switch v := raw.(type) {
		case float64:
			return int64(v), true
		case string:
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
				return id, true
			}
		}
	}
	return 0, false
}

// helper lokal: cek apakah role claim = "admin"
func isAdminTest(c echo.Context) bool {
	user := c.Get("user")
	if user == nil {
		return false
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	if role, exists := claims["role"]; exists {
		if rs, ok := role.(string); ok && rs == "admin" {
			return true
		}
	}
	// fallback: cek claim boolean is_admin / admin
	if raw, exists := claims["is_admin"]; exists {
		switch v := raw.(type) {
		case bool:
			return v
		case float64:
			return v != 0
		}
	}
	return false
}

// ...existing code...

func (h *AuctionController) CreateAuctionItem(c echo.Context) error {
	// echo-jwt stores token under context key "user"
	user := c.Get("user")
	log.Printf("DEBUG c.Get(\"user\"): %#v", user)

	if token, ok := user.(*jwt.Token); ok {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			log.Printf("DEBUG token claims: %+v", claims)
		} else {
			log.Println("DEBUG: token has no MapClaims")
		}
	} else {
		log.Println("DEBUG: no *jwt.Token in context")
	}

	if !isAdminTest(c) {
		return utils.ForbiddenResponse(c, "only admin can create auction items")
	}

	var payload dto.AuctionItemDTO
	if err := c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}

	if err := h.validate.Struct(payload); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	userID, ok := getUserIDTest(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	payload.UserID = int64(userID)

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
		return utils.InternalServerErrorResponse(c, "failed retrieving auction item")
	}

	return utils.SuccessResponse(c, "auction item retrieved successfully", item)
}

func (h *AuctionController) UpdateAuctionItem(c echo.Context) error {
	if !isAdminTest(c) {
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
		return utils.InternalServerErrorResponse(c, "failed updating auction item")
	}
	return utils.SuccessResponse(c, "auction item updated successfully", updatedItem)
}

func (h *AuctionController) DeleteAuctionItem(c echo.Context) error {
	if !isAdminTest(c) {
		return utils.ForbiddenResponse(c, "only admin can delete auction items")
	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid auction item ID")
	}

	err = h.svc.Delete(id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed deleting auction item")
	}
	return utils.SuccessResponse(c, "auction item deleted successfully", nil)
}

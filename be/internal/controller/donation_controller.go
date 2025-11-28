package controller

import (
	"errors"
	"strconv"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/service"
	"milestone3/be/internal/utils"

	"github.com/labstack/echo/v4"
)

type DonationController struct {
	svc service.DonationService
}

func NewDonationController(s service.DonationService) *DonationController {
	return &DonationController{svc: s}
}

// helper: read user id from context (set by auth middleware)
func getUserID(c echo.Context) (uint, bool) {
	v := c.Get("user_id")
	switch t := v.(type) {
	case uint:
		return t, true
	case int:
		return uint(t), true
	case int64:
		return uint(t), true
	case float64:
		return uint(t), true
	default:
		return 0, false
	}
}

// donor/user: POST /donations
func (h *DonationController) CreateDonation(c echo.Context) error {
	var payload dto.DonationDTO
	if err := c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}

	userID, ok := getUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	payload.UserID = userID

	if err := h.svc.CreateDonation(payload); err != nil {
		// if service has validation error mapping, map accordingly; fallback internal error
		return utils.InternalServerErrorResponse(c, "failed creating donation")
	}
	return utils.CreatedResponse(c, "donation created successfully", nil)
}

// user/admin: GET /donations
// - admin: returns all donations
// - user: returns only own donations
func (h *DonationController) GetAllDonations(c echo.Context) error {
	donations, err := h.svc.GetAllDonations()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed fetching donations")
	}

	// if admin return all
	if isAdmin(c) {
		return utils.SuccessResponse(c, "donations fetched", donations)
	}

	// else filter by user
	userID, ok := getUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	filtered := make([]dto.DonationDTO, 0, len(donations))
	for _, d := range donations {
		if d.UserID == userID {
			filtered = append(filtered, d)
		}
	}
	return utils.SuccessResponse(c, "donations fetched", filtered)
}

// user/admin: GET /donations/:id (owner or admin)
func (h *DonationController) GetDonationByID(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	d, err := h.svc.GetDonationByID(uint(id64))
	if err != nil {
		if errors.Is(err, service.ErrDonationNotFound) {
			return utils.NotFoundResponse(c, "donation not found")
		}
		return utils.InternalServerErrorResponse(c, "failed fetching donation")
	}

	// permission check: owner or admin
	if !isAdmin(c) {
		userID, ok := getUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "unauthenticated")
		}
		if d.UserID != userID {
			return utils.ForbiddenResponse(c, "forbidden")
		}
	}

	return utils.SuccessResponse(c, "donation fetched", d)
}

// user/admin: PUT /donations/:id (only owner or admin can update)
func (h *DonationController) UpdateDonation(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	var payload dto.DonationDTO
	if err := c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}
	payload.ID = uint(id64)

	userID, ok := getUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	isAdm := isAdmin(c)

	if err := h.svc.UpdateDonation(payload, userID, isAdm); err != nil {
		if errors.Is(err, service.ErrDonationNotFound) {
			return utils.NotFoundResponse(c, "donation not found")
		}
		if errors.Is(err, service.ErrForbidden) {
			return utils.ForbiddenResponse(c, "forbidden")
		}
		return utils.InternalServerErrorResponse(c, "failed updating donation")
	}
	return utils.SuccessResponse(c, "donation updated", nil)
}

// user/admin: DELETE /donations/:id (only owner or admin can delete)
func (h *DonationController) DeleteDonation(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	userID, ok := getUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	isAdm := isAdmin(c)

	if err := h.svc.DeleteDonation(uint(id64), userID, isAdm); err != nil {
		if errors.Is(err, service.ErrDonationNotFound) {
			return utils.NotFoundResponse(c, "donation not found")
		}
		if errors.Is(err, service.ErrForbidden) {
			return utils.ForbiddenResponse(c, "forbidden")
		}
		return utils.InternalServerErrorResponse(c, "failed deleting donation")
	}
	return utils.NoContentResponse(c)
}

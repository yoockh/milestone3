package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
	"milestone3/be/internal/service"
	"milestone3/be/internal/utils"

	"github.com/labstack/echo/v4"
)

type DonationController struct {
	svc          service.DonationService
	privateStore repository.GCPStorageRepo
}

func NewDonationController(s service.DonationService, privateStore repository.GCPStorageRepo) *DonationController {
	return &DonationController{svc: s, privateStore: privateStore}
}

// POST /donations (fix)
func (h *DonationController) CreateDonation(c echo.Context) error {
	var payload dto.DonationDTO

	contentType := c.Request().Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {

		if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
			return utils.BadRequestResponse(c, "invalid multipart form")
		}
		form := c.Request().MultipartForm

		payload.Title = form.Value["title"][0]
		payload.Description = form.Value["description"][0]
		payload.Category = form.Value["category"][0]
		payload.Condition = form.Value["condition"][0]

		if v, ok := form.Value["status"]; ok {
			payload.Status = entity.StatusDonation(v[0])
		}

		// FILE HANDLING PRIVATE ONLY
		if fhs, ok := form.File["photos"]; ok {
			for _, fh := range fhs {
				f, err := fh.Open()
				if err != nil {
					return utils.BadRequestResponse(c, "cannot open file")
				}

				objName := fmt.Sprintf("donations/private/%d_%s", time.Now().UnixNano(), fh.Filename)
				objectName, err := h.privateStore.UploadFile(c.Request().Context(), f, objName)
				_ = f.Close()

				if err != nil {
					return utils.InternalServerErrorResponse(c, "failed upload")
				}

				// SAVE PRIVATE STORAGE (objectName)
				payload.Photos = append(payload.Photos, objectName)
			}
		}

	} else {
		if err := c.Bind(&payload); err != nil {
			return utils.BadRequestResponse(c, "invalid payload")
		}
	}

	userID, ok := utils.GetUserID(c)
	if !ok || userID == 0 {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	payload.UserID = userID

	if err := h.svc.CreateDonation(payload); err != nil {
		return utils.InternalServerErrorResponse(c, "failed creating donation")
	}

	return utils.CreatedResponse(c, "donation created successfully", nil)
}

// user/admin: GET /donations
// - admin: returns all donations
// - user: returns only own donations
func (h *DonationController) GetAllDonations(c echo.Context) error {
	userID, _ := utils.GetUserID(c) // unauthenticated => 0,false
	isAdm := utils.IsAdmin(c)

	// require auth for user-level listing; admin may call even without user_id set by middleware
	if !isAdm {
		if userID == 0 {
			return utils.UnauthorizedResponse(c, "unauthenticated")
		}
	}

	donations, err := h.svc.GetAllDonations(userID, isAdm)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed fetching donations")
	}
	return utils.SuccessResponse(c, "donations fetched", donations)
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
	if !utils.IsAdmin(c) {
		userID, ok := utils.GetUserID(c)
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

	userID, ok := utils.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	isAdm := utils.IsAdmin(c)

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

	userID, ok := utils.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	isAdm := utils.IsAdmin(c)

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

func (h *DonationController) PatchDonation(c echo.Context) error {
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

	userID, ok := utils.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthenticated")
	}
	isAdm := utils.IsAdmin(c)

	if err := h.svc.PatchDonation(payload, userID, isAdm); err != nil {
		if errors.Is(err, service.ErrDonationNotFound) {
			return utils.NotFoundResponse(c, "donation not found")
		}
		if errors.Is(err, service.ErrForbidden) {
			return utils.ForbiddenResponse(c, "forbidden")
		}
		return utils.InternalServerErrorResponse(c, "failed patching donation")
	}
	return utils.SuccessResponse(c, "donation patched", nil)
}

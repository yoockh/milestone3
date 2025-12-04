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

// CreateDonation godoc
// @Summary Create new donation
// @Description Submit a new donation with photos and details
// @Tags Your Donate Rise API - Donations
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Donation title"
// @Param description formData string true "Donation description"
// @Param category formData string true "Donation category"
// @Param condition formData string true "Item condition"
// @Param photos formData file false "Donation photos (multiple files allowed)"
// @Success 201 {object} utils.SuccessResponseData "donation created successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid payload or file upload"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /donations [post]
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

// GetAllDonations godoc
// @Summary Get all donations
// @Description Get all donations (admin sees all, users see only their own)
// @Tags Your Donate Rise API - Donations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseData "donations fetched"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /donations [get]
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

// GetDonationByID godoc
// @Summary Get donation by ID
// @Description Retrieve a specific donation by ID (owner or admin only)
// @Tags Your Donate Rise API - Donations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Donation ID"
// @Success 200 {object} utils.SuccessResponseData "donation fetched"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid donation ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Access denied"
// @Failure 404 {object} utils.ErrorResponse "Donation not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /donations/{id} [get]
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

// UpdateDonation godoc
// @Summary Update donation
// @Description Update an existing donation (owner or admin only)
// @Tags Your Donate Rise API - Donations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Donation ID"
// @Param donation body dto.DonationDTO true "Updated donation data"
// @Success 200 {object} utils.SuccessResponseData "donation updated"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid ID or payload"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Access denied"
// @Failure 404 {object} utils.ErrorResponse "Donation not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /donations/{id} [put]
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

// DeleteDonation godoc
// @Summary Delete donation
// @Description Delete a donation by ID (owner or admin only)
// @Tags Your Donate Rise API - Donations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Donation ID"
// @Success 204 "Donation deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid donation ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Access denied"
// @Failure 404 {object} utils.ErrorResponse "Donation not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /donations/{id} [delete]
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

// PatchDonation godoc
// @Summary Partially update donation
// @Description Partially update a donation by ID (owner or admin only)
// @Tags Your Donate Rise API - Donations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Donation ID"
// @Param donation body dto.DonationDTO true "Partial donation data"
// @Success 200 {object} utils.SuccessResponseData "donation patched"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid ID or payload"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Access denied"
// @Failure 404 {object} utils.ErrorResponse "Donation not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /donations/{id} [patch]
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

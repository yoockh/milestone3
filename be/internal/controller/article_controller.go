package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/repository"
	"milestone3/be/internal/service"
	"milestone3/be/internal/utils"

	"github.com/labstack/echo/v4"
)

type ArticleController struct {
	svc           service.ArticleService
	storagePublic repository.GCPStorageRepo
}

func NewArticleController(s service.ArticleService, storage repository.GCPStorageRepo) *ArticleController {
	return &ArticleController{svc: s, storagePublic: storage}
}

// GetAllArticles godoc
// @Summary Get all transparency articles
// @Description Retrieve all published weekly transparency articles with pagination
// @Tags Your Donate Rise API - Articles
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} utils.SuccessResponseData "articles fetched"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /articles [get]
func (h *ArticleController) GetAllArticles(c echo.Context) error {
	// Parse pagination params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	articles, total, err := h.svc.GetAllArticles(page, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed fetching articles")
	}

	response := map[string]interface{}{
		"articles": articles,
		"page":     page,
		"limit":    limit,
		"total":    total,
	}
	return utils.SuccessResponse(c, "articles fetched", response)
}

// GetArticleByID godoc
// @Summary Get article by ID
// @Description Retrieve a specific transparency article by its ID
// @Tags Your Donate Rise API - Articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} utils.SuccessResponseData "article fetched"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid article ID"
// @Failure 404 {object} utils.ErrorResponse "Article not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /articles/{id} [get]
func (h *ArticleController) GetArticleByID(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}
	article, err := h.svc.GetArticleByID(uint(id64))
	if err != nil {
		if err == service.ErrArticleNotFound {
			return utils.NotFoundResponse(c, "article not found")
		}
		return utils.InternalServerErrorResponse(c, "failed fetching article")
	}
	return utils.SuccessResponse(c, "article fetched", article)
}

// CreateArticle godoc
// @Summary Create new transparency article
// @Description Create a new weekly transparency article with optional image upload
// @Tags Your Donate Rise API - Articles
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Article title"
// @Param content formData string true "Article content"
// @Param week formData int true "Week number"
// @Param image formData file false "Article image (optional)"
// @Success 201 {object} utils.SuccessResponseData "article created"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid payload or image"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /articles [post]
func (h *ArticleController) CreateArticle(c echo.Context) error {
	if !utils.IsAdmin(c) {
		return utils.ForbiddenResponse(c, "admin only")
	}

	contentType := c.Request().Header.Get("Content-Type")
	var payload dto.ArticleDTO

	// if multipart/form-data (with image)
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
			return utils.BadRequestResponse(c, "invalid multipart form")
		}

		form := c.Request().MultipartForm

		// Validate required fields
		if len(form.Value["title"]) == 0 {
			return utils.BadRequestResponse(c, "title is required")
		}
		if len(form.Value["content"]) == 0 {
			return utils.BadRequestResponse(c, "content is required")
		}
		if len(form.Value["week"]) == 0 {
			return utils.BadRequestResponse(c, "week is required")
		}

		payload.Title = form.Value["title"][0]
		payload.Content = form.Value["content"][0]
		week, _ := strconv.Atoi(form.Value["week"][0])
		payload.Week = week

		// handle image (opsional)
		if fhs, ok := form.File["image"]; ok && len(fhs) > 0 {
			fh := fhs[0]

			// Validate file size (max 5MB)
			if fh.Size > 5*1024*1024 {
				return utils.BadRequestResponse(c, "image size exceeds 5MB limit")
			}

			// Validate file type (only images)
			contentType := fh.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "image/") {
				return utils.BadRequestResponse(c, "only image files are allowed")
			}

			file, err := fh.Open()
			if err != nil {
				return utils.BadRequestResponse(c, "failed open image")
			}
			defer file.Close()

			// Sanitize filename
			safeFilename := strings.ReplaceAll(fh.Filename, "..", "")
			safeFilename = strings.ReplaceAll(safeFilename, "/", "")
			safeFilename = strings.ReplaceAll(safeFilename, "\\", "")

			objName := fmt.Sprintf("articles/%d_%s", time.Now().UnixNano(), safeFilename)

			//  upload to public storage
			url, err := h.storagePublic.UploadFile(c.Request().Context(), file, objName)
			if err != nil {
				return utils.InternalServerErrorResponse(c, "failed uploading image")
			}

			payload.Image = url
		}

	} else {
		// support JSON without image
		if err := c.Bind(&payload); err != nil {
			return utils.BadRequestResponse(c, "invalid payload")
		}
	}

	// send to service
	if err := h.svc.CreateArticle(payload); err != nil {
		return utils.InternalServerErrorResponse(c, "failed creating article")
	}

	return utils.CreatedResponse(c, "article created", nil)
}

// UpdateArticle godoc
// @Summary Update existing article
// @Description Update an existing transparency article by ID
// @Tags Your Donate Rise API - Articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Article ID"
// @Param article body dto.ArticleDTO true "Updated article data"
// @Success 200 {object} utils.SuccessResponseData "article updated"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid ID or payload"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 404 {object} utils.ErrorResponse "Article not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /articles/{id} [put]
func (h *ArticleController) UpdateArticle(c echo.Context) error {
	if !utils.IsAdmin(c) {
		return utils.ForbiddenResponse(c, "admin only")
	}
	var payload dto.ArticleDTO
	if err := c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}
	payload.ID = uint(id64)
	if err := h.svc.UpdateArticle(payload); err != nil {
		if err == service.ErrArticleNotFound {
			return utils.NotFoundResponse(c, "article not found")
		}
		return utils.InternalServerErrorResponse(c, "failed updating article")
	}
	return utils.SuccessResponse(c, "article updated", nil)
}

// DeleteArticle godoc
// @Summary Delete article
// @Description Delete a transparency article by ID
// @Tags Your Donate Rise API - Articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Article ID"
// @Success 204 "Article deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid article ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin access required"
// @Failure 404 {object} utils.ErrorResponse "Article not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /articles/{id} [delete]
func (h *ArticleController) DeleteArticle(c echo.Context) error {
	if !utils.IsAdmin(c) {
		return utils.ForbiddenResponse(c, "admin only")
	}
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}
	if err := h.svc.DeleteArticle(uint(id64)); err != nil {
		if err == service.ErrArticleNotFound {
			return utils.NotFoundResponse(c, "article not found")
		}
		return utils.InternalServerErrorResponse(c, "failed deleting article")
	}
	return utils.NoContentResponse(c)
}

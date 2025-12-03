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

func (h *ArticleController) GetAllArticles(c echo.Context) error {
	articles, err := h.svc.GetAllArticles()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed fetching articles")
	}
	return utils.SuccessResponse(c, "articles fetched", articles)
}

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

		payload.Title = form.Value["title"][0]
		payload.Content = form.Value["content"][0]
		week, _ := strconv.Atoi(form.Value["week"][0])
		payload.Week = week

		// handle image (opsional)
		if fhs, ok := form.File["image"]; ok && len(fhs) > 0 {
			fh := fhs[0]

			file, err := fh.Open()
			if err != nil {
				return utils.BadRequestResponse(c, "failed open image")
			}
			defer file.Close()

			objName := fmt.Sprintf("articles/%d_%s", time.Now().UnixNano(), fh.Filename)

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

package controller

import (
	"strconv"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/service"
	"milestone3/be/internal/utils"

	"github.com/labstack/echo/v4"
)

type ArticleController struct {
	svc service.ArticleService
}

func NewArticleController(s service.ArticleService) *ArticleController {
	return &ArticleController{svc: s}
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
	var payload dto.ArticleDTO
	if err := c.Bind(&payload); err != nil {
		return utils.BadRequestResponse(c, "invalid payload")
	}
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

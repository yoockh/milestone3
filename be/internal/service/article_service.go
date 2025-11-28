package service

import (
	"errors"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/repository"

	"gorm.io/gorm"
)

type ArticleService interface {
	CreateArticle(articleDTO dto.ArticleDTO) error
	GetAllArticles() ([]dto.ArticleDTO, error)
	GetArticleByID(id uint) (dto.ArticleDTO, error)
	UpdateArticle(articleDTO dto.ArticleDTO) error
	DeleteArticle(id uint) error
}

type articleService struct {
	repo repository.ArticleRepo
}

func NewArticleService(repo repository.ArticleRepo) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) CreateArticle(articleDTO dto.ArticleDTO) error {
	article, err := dto.ArticleRequest(articleDTO)
	if err != nil {
		return err
	}
	return s.repo.CreateArticle(article)
}

func (s *articleService) GetAllArticles() ([]dto.ArticleDTO, error) {
	articles, err := s.repo.GetAllArticles()
	if err != nil {
		return nil, err
	}
	return dto.ArticleResponses(articles), nil
}

func (s *articleService) GetArticleByID(id uint) (dto.ArticleDTO, error) {
	article, err := s.repo.GetArticleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ArticleDTO{}, ErrArticleNotFound
		}
		return dto.ArticleDTO{}, err
	}
	return dto.ArticleResponse(article), nil
}

func (s *articleService) UpdateArticle(articleDTO dto.ArticleDTO) error {
	article, err := dto.ArticleRequest(articleDTO)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateArticle(article); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrArticleNotFound
		}
		return err
	}
	return nil
}

func (s *articleService) DeleteArticle(id uint) error {
	if err := s.repo.DeleteArticle(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrArticleNotFound
		}
		return err
	}
	return nil
}

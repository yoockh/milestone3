package dto

import (
	"milestone3/be/internal/entity"
	"time"
)

type ArticleDTO struct {
	ID        uint      `json:"id,omitempty" validate:"omitempty"`
	Title     string    `json:"title,omitempty" validate:"required"`
	Content   string    `json:"content,omitempty" validate:"required"`
	Week      int       `json:"week,omitempty" validate:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// ArticleRequest converts DTO to entity.Article
func ArticleRequest(a ArticleDTO) (entity.Article, error) {
	return entity.Article{
		ID:        a.ID,
		Title:     a.Title,
		Content:   a.Content,
		Week:      a.Week,
		CreatedAt: a.CreatedAt,
	}, nil
}

// ArticleResponse converts entity.Article to DTO
func ArticleResponse(m entity.Article) ArticleDTO {
	return ArticleDTO{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Content,
		Week:      m.Week,
		CreatedAt: m.CreatedAt,
	}
}

// ArticleResponses converts slice of entity.Article to slice of DTOs
func ArticleResponses(articles []entity.Article) []ArticleDTO {
	dtos := make([]ArticleDTO, len(articles))
	for i, article := range articles {
		dtos[i] = ArticleResponse(article)
	}
	return dtos
}

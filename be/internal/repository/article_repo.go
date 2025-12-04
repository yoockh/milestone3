package repository

import (
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type ArticleRepo interface {
	GetAllArticles(page, limit int) ([]entity.Article, int64, error)
	GetArticleByID(id uint) (entity.Article, error)
	// Admin functionalities
	CreateArticle(article entity.Article) error
	UpdateArticle(article entity.Article) error
	DeleteArticle(id uint) error
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) ArticleRepo {
	return &articleRepo{db: db}
}

func (r *articleRepo) CreateArticle(article entity.Article) error {
	return r.db.Create(&article).Error
}

func (r *articleRepo) GetAllArticles(page, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&articles).Error
	return articles, total, err
}

func (r *articleRepo) GetArticleByID(id uint) (entity.Article, error) {
	var article entity.Article
	err := r.db.First(&article, id).Error
	return article, err
}

func (r *articleRepo) UpdateArticle(article entity.Article) error {
	return r.db.Save(&article).Error
}

func (r *articleRepo) DeleteArticle(id uint) error {
	return r.db.Delete(&entity.Article{}, id).Error
}

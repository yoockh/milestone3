package repository

import (
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type ArticleRepo interface {
	CreateArticle(article entity.Article) error
	GetAllArticles() ([]entity.Article, error)
	GetArticleByID(id uint) (entity.Article, error)
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

func (r *articleRepo) GetAllArticles() ([]entity.Article, error) {
	var articles []entity.Article
	err := r.db.Find(&articles).Error
	return articles, err
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

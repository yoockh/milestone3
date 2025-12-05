package service

import (
	"errors"
	"testing"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestArticleService_CreateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepo(ctrl)
	articleService := NewArticleService(mockRepo)

	tests := []struct {
		name    string
		req     dto.ArticleDTO
		setup   func()
		wantErr bool
	}{
		{
			name: "successful article creation",
			req: dto.ArticleDTO{
				Title:   "Test Article",
				Content: "Test content",
			},
			setup: func() {
				mockRepo.EXPECT().CreateArticle(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository create error",
			req: dto.ArticleDTO{
				Title:   "Test Article",
				Content: "Test content",
			},
			setup: func() {
				mockRepo.EXPECT().CreateArticle(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			
			err := articleService.CreateArticle(tt.req)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArticleService_GetAllArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepo(ctrl)
	articleService := NewArticleService(mockRepo)

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful get all articles",
			setup: func() {
				articles := []entity.Article{
					{ID: 1, Title: "Article 1", Content: "Content 1"},
					{ID: 2, Title: "Article 2", Content: "Content 2"},
				}
				mockRepo.EXPECT().GetAllArticles(1, 10).Return(articles, int64(2), nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func() {
				mockRepo.EXPECT().GetAllArticles(1, 10).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			
			result, total, err := articleService.GetAllArticles(1, 10)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, int64(2), total)
			}
		})
	}
}

func TestArticleService_GetArticleByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepo(ctrl)
	articleService := NewArticleService(mockRepo)

	tests := []struct {
		name    string
		id      uint
		setup   func()
		wantErr bool
	}{
		{
			name: "successful get article by id",
			id:   1,
			setup: func() {
				article := entity.Article{ID: 1, Title: "Test Article", Content: "Test content"}
				mockRepo.EXPECT().GetArticleByID(uint(1)).Return(article, nil)
			},
			wantErr: false,
		},
		{
			name: "article not found",
			id:   999,
			setup: func() {
				mockRepo.EXPECT().GetArticleByID(uint(999)).Return(entity.Article{}, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			
			result, err := articleService.GetArticleByID(tt.id)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint(1), result.ID)
			}
		})
	}
}


func TestArticleService_UpdateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepo(ctrl)
	articleService := NewArticleService(mockRepo)

	tests := []struct {
		name    string
		req     dto.ArticleDTO
		setup   func()
		wantErr bool
	}{
		{
			name: "successful update",
			req: dto.ArticleDTO{
				ID:    1,
				Title: "Updated",
			},
			setup: func() {
				mockRepo.EXPECT().UpdateArticle(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "article not found",
			req: dto.ArticleDTO{
				ID: 999,
			},
			setup: func() {
				mockRepo.EXPECT().UpdateArticle(gomock.Any()).Return(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := articleService.UpdateArticle(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArticleService_DeleteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepo(ctrl)
	articleService := NewArticleService(mockRepo)

	tests := []struct {
		name    string
		id      uint
		setup   func()
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   1,
			setup: func() {
				mockRepo.EXPECT().DeleteArticle(uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "article not found",
			id:   999,
			setup: func() {
				mockRepo.EXPECT().DeleteArticle(uint(999)).Return(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := articleService.DeleteArticle(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

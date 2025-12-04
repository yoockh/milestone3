package mocks

import (
	"milestone3/be/internal/entity"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockArticleRepo is a mock of ArticleRepo interface.
type MockArticleRepo struct {
	ctrl     *gomock.Controller
	recorder *MockArticleRepoMockRecorder
}

// MockArticleRepoMockRecorder is the mock recorder for MockArticleRepo.
type MockArticleRepoMockRecorder struct {
	mock *MockArticleRepo
}

// NewMockArticleRepo creates a new mock instance.
func NewMockArticleRepo(ctrl *gomock.Controller) *MockArticleRepo {
	mock := &MockArticleRepo{ctrl: ctrl}
	mock.recorder = &MockArticleRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleRepo) EXPECT() *MockArticleRepoMockRecorder {
	return m.recorder
}

// CreateArticle mocks base method.
func (m *MockArticleRepo) CreateArticle(article entity.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateArticle", article)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateArticle indicates an expected call of CreateArticle.
func (mr *MockArticleRepoMockRecorder) CreateArticle(article interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateArticle", reflect.TypeOf((*MockArticleRepo)(nil).CreateArticle), article)
}

// DeleteArticle mocks base method.
func (m *MockArticleRepo) DeleteArticle(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteArticle", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteArticle indicates an expected call of DeleteArticle.
func (mr *MockArticleRepoMockRecorder) DeleteArticle(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteArticle", reflect.TypeOf((*MockArticleRepo)(nil).DeleteArticle), id)
}

// GetAllArticles mocks base method.
func (m *MockArticleRepo) GetAllArticles(page, limit int) ([]entity.Article, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllArticles", page, limit)
	ret0, _ := ret[0].([]entity.Article)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllArticles indicates an expected call of GetAllArticles.
func (mr *MockArticleRepoMockRecorder) GetAllArticles(page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllArticles", reflect.TypeOf((*MockArticleRepo)(nil).GetAllArticles), page, limit)
}

// GetArticleByID mocks base method.
func (m *MockArticleRepo) GetArticleByID(id uint) (entity.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArticleByID", id)
	ret0, _ := ret[0].(entity.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArticleByID indicates an expected call of GetArticleByID.
func (mr *MockArticleRepoMockRecorder) GetArticleByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArticleByID", reflect.TypeOf((*MockArticleRepo)(nil).GetArticleByID), id)
}

// UpdateArticle mocks base method.
func (m *MockArticleRepo) UpdateArticle(article entity.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateArticle", article)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateArticle indicates an expected call of UpdateArticle.
func (mr *MockArticleRepoMockRecorder) UpdateArticle(article interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateArticle", reflect.TypeOf((*MockArticleRepo)(nil).UpdateArticle), article)
}
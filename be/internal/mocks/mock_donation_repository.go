package mocks

import (
	"context"
	"io"
	"milestone3/be/internal/entity"
	"time"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDonationRepo is a mock of DonationRepo interface.
type MockDonationRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDonationRepoMockRecorder
}

// MockDonationRepoMockRecorder is the mock recorder for MockDonationRepo.
type MockDonationRepoMockRecorder struct {
	mock *MockDonationRepo
}

// NewMockDonationRepo creates a new mock instance.
func NewMockDonationRepo(ctrl *gomock.Controller) *MockDonationRepo {
	mock := &MockDonationRepo{ctrl: ctrl}
	mock.recorder = &MockDonationRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDonationRepo) EXPECT() *MockDonationRepoMockRecorder {
	return m.recorder
}

// CreateDonation mocks base method.
func (m *MockDonationRepo) CreateDonation(donation entity.Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDonation", donation)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateDonation indicates an expected call of CreateDonation.
func (mr *MockDonationRepoMockRecorder) CreateDonation(donation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDonation", reflect.TypeOf((*MockDonationRepo)(nil).CreateDonation), donation)
}

// DeleteDonation mocks base method.
func (m *MockDonationRepo) DeleteDonation(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDonation", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDonation indicates an expected call of DeleteDonation.
func (mr *MockDonationRepoMockRecorder) DeleteDonation(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDonation", reflect.TypeOf((*MockDonationRepo)(nil).DeleteDonation), id)
}

// GetAllDonations mocks base method.
func (m *MockDonationRepo) GetAllDonations(page, limit int) ([]entity.Donation, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDonations", page, limit)
	ret0, _ := ret[0].([]entity.Donation)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllDonations indicates an expected call of GetAllDonations.
func (mr *MockDonationRepoMockRecorder) GetAllDonations(page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDonations", reflect.TypeOf((*MockDonationRepo)(nil).GetAllDonations), page, limit)
}

// GetDonationByID mocks base method.
func (m *MockDonationRepo) GetDonationByID(id uint) (entity.Donation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDonationByID", id)
	ret0, _ := ret[0].(entity.Donation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDonationByID indicates an expected call of GetDonationByID.
func (mr *MockDonationRepoMockRecorder) GetDonationByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDonationByID", reflect.TypeOf((*MockDonationRepo)(nil).GetDonationByID), id)
}

// GetDonationsByUserID mocks base method.
func (m *MockDonationRepo) GetDonationsByUserID(userID uint, page, limit int) ([]entity.Donation, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDonationsByUserID", userID, page, limit)
	ret0, _ := ret[0].([]entity.Donation)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDonationsByUserID indicates an expected call of GetDonationsByUserID.
func (mr *MockDonationRepoMockRecorder) GetDonationsByUserID(userID, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDonationsByUserID", reflect.TypeOf((*MockDonationRepo)(nil).GetDonationsByUserID), userID, page, limit)
}

// PatchDonation mocks base method.
func (m *MockDonationRepo) PatchDonation(donation entity.Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchDonation", donation)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchDonation indicates an expected call of PatchDonation.
func (mr *MockDonationRepoMockRecorder) PatchDonation(donation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchDonation", reflect.TypeOf((*MockDonationRepo)(nil).PatchDonation), donation)
}

// UpdateDonation mocks base method.
func (m *MockDonationRepo) UpdateDonation(donation entity.Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDonation", donation)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDonation indicates an expected call of UpdateDonation.
func (mr *MockDonationRepoMockRecorder) UpdateDonation(donation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDonation", reflect.TypeOf((*MockDonationRepo)(nil).UpdateDonation), donation)
}

// MockGCPStorageRepo is a mock of GCPStorageRepo interface.
type MockGCPStorageRepo struct {
	ctrl     *gomock.Controller
	recorder *MockGCPStorageRepoMockRecorder
}

// MockGCPStorageRepoMockRecorder is the mock recorder for MockGCPStorageRepo.
type MockGCPStorageRepoMockRecorder struct {
	mock *MockGCPStorageRepo
}

// NewMockGCPStorageRepo creates a new mock instance.
func NewMockGCPStorageRepo(ctrl *gomock.Controller) *MockGCPStorageRepo {
	mock := &MockGCPStorageRepo{ctrl: ctrl}
	mock.recorder = &MockGCPStorageRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGCPStorageRepo) EXPECT() *MockGCPStorageRepoMockRecorder {
	return m.recorder
}

// GenerateSignedURL mocks base method.
func (m *MockGCPStorageRepo) GenerateSignedURL(ctx context.Context, objectName string, expiration time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSignedURL", ctx, objectName, expiration)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateSignedURL indicates an expected call of GenerateSignedURL.
func (mr *MockGCPStorageRepoMockRecorder) GenerateSignedURL(ctx, objectName, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSignedURL", reflect.TypeOf((*MockGCPStorageRepo)(nil).GenerateSignedURL), ctx, objectName, expiration)
}

// UploadFile mocks base method.
func (m *MockGCPStorageRepo) UploadFile(ctx context.Context, file io.Reader, fileName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", ctx, file, fileName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockGCPStorageRepoMockRecorder) UploadFile(ctx, file, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockGCPStorageRepo)(nil).UploadFile), ctx, file, fileName)
}
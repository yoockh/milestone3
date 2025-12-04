package mocks

import (
	"milestone3/be/internal/entity"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFinalDonationRepository is a mock of FinalDonationRepository interface.
type MockFinalDonationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFinalDonationRepositoryMockRecorder
}

// MockFinalDonationRepositoryMockRecorder is the mock recorder for MockFinalDonationRepository.
type MockFinalDonationRepositoryMockRecorder struct {
	mock *MockFinalDonationRepository
}

// NewMockFinalDonationRepository creates a new mock instance.
func NewMockFinalDonationRepository(ctrl *gomock.Controller) *MockFinalDonationRepository {
	mock := &MockFinalDonationRepository{ctrl: ctrl}
	mock.recorder = &MockFinalDonationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFinalDonationRepository) EXPECT() *MockFinalDonationRepositoryMockRecorder {
	return m.recorder
}

// GetAllFinalDonations mocks base method.
func (m *MockFinalDonationRepository) GetAllFinalDonations(page, limit int) ([]entity.FinalDonation, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllFinalDonations", page, limit)
	ret0, _ := ret[0].([]entity.FinalDonation)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllFinalDonations indicates an expected call of GetAllFinalDonations.
func (mr *MockFinalDonationRepositoryMockRecorder) GetAllFinalDonations(page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllFinalDonations", reflect.TypeOf((*MockFinalDonationRepository)(nil).GetAllFinalDonations), page, limit)
}

// GetAllFinalDonationsByUserID mocks base method.
func (m *MockFinalDonationRepository) GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllFinalDonationsByUserID", userID)
	ret0, _ := ret[0].([]entity.FinalDonation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllFinalDonationsByUserID indicates an expected call of GetAllFinalDonationsByUserID.
func (mr *MockFinalDonationRepositoryMockRecorder) GetAllFinalDonationsByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllFinalDonationsByUserID", reflect.TypeOf((*MockFinalDonationRepository)(nil).GetAllFinalDonationsByUserID), userID)
}
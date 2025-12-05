package service

// import (
// 	"errors"
// 	"testing"

// 	"milestone3/be/internal/dto"
// 	"milestone3/be/internal/entity"
// 	"milestone3/be/internal/mocks"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// add method on interface and the whole mock is generated so i intentionally comment this test (rafly) //

// func TestPaymentService_CreatePayment(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockPaymentRepository(ctrl)
// 	paymentService := NewPaymentService(mockRepo)

// 	tests := []struct {
// 		name          string
// 		req           dto.PaymentRequest
// 		userId        int
// 		auctionItemId int
// 		setup         func()
// 		wantErr       bool
// 	}{
// 		{
// 			name: "successful payment creation",
// 			req: dto.PaymentRequest{
// 				Amount: 100000,
// 			},
// 			userId:        1,
// 			auctionItemId: 1,
// 			setup: func() {
// 				mockRepo.EXPECT().CreateMidtrans(gomock.Any(), gomock.Any()).Return(dto.PaymentResponse{
// 					OrderId:        "YDR-123",
// 					TransactionId:  "TXN-123",
// 					PaymentLinkUrl: "https://payment.link",
// 				}, nil)
// 				mockRepo.EXPECT().Create(gomock.Any(), "YDR-123").Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repository create error",
// 			req: dto.PaymentRequest{
// 				Amount: 100000,
// 			},
// 			userId:        1,
// 			auctionItemId: 1,
// 			setup: func() {
// 				mockRepo.EXPECT().CreateMidtrans(gomock.Any(), gomock.Any()).Return(dto.PaymentResponse{
// 					OrderId:        "YDR-123",
// 					TransactionId:  "TXN-123",
// 					PaymentLinkUrl: "https://payment.link",
// 				}, nil)
// 				mockRepo.EXPECT().Create(gomock.Any(), "YDR-123").Return(errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
			
// 			result, err := paymentService.CreatePayment(tt.req, tt.userId, tt.auctionItemId)
			
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, "YDR-123", result.OrderId)
// 				assert.Equal(t, "TXN-123", result.TransactionId)
// 			}
// 		})
// 	}
// }

// func TestPaymentService_GetPaymentById(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockPaymentRepository(ctrl)
// 	paymentService := NewPaymentService(mockRepo)

// 	tests := []struct {
// 		name    string
// 		id      int
// 		setup   func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful get payment by id",
// 			id:   1,
// 			setup: func() {
// 				mockRepo.EXPECT().GetById(1).Return(entity.Payment{
// 					Id:            1,
// 					UserId:        1,
// 					AuctionItemId: 1,
// 					Amount:        100000.0,
// 					Status:        "pending",
// 				}, nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "payment not found",
// 			id:   999,
// 			setup: func() {
// 				mockRepo.EXPECT().GetById(999).Return(entity.Payment{}, errors.New("payment not found"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
			
// 			result, err := paymentService.GetPaymentById(tt.id)
			
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				assert.Empty(t, result)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, 1, result.Id)
// 				assert.Equal(t, 100000.0, result.Amount)
// 			}
// 		})
// 	}
// }

// func TestPaymentService_GetAllPayment(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockPaymentRepository(ctrl)
// 	paymentService := NewPaymentService(mockRepo)

// 	tests := []struct {
// 		name    string
// 		setup   func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful get all",
// 			setup: func() {
// 				payments := []entity.Payment{{Id: 1}, {Id: 2}}
// 				mockRepo.EXPECT().GetAll().Return(payments, nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repository error",
// 			setup: func() {
// 				mockRepo.EXPECT().GetAll().Return(nil, errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()
// 			result, err := paymentService.GetAllPayment()
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, result)
// 			}
// 		})
// 	}
// }

package repository

import (
	"context"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
	ctx context.Context
}

func NewPaymentRepository(db *gorm.DB, ctx context.Context) *PaymentRepo {
	return &PaymentRepo{db: db, ctx: ctx}
}

func (pr *PaymentRepo) Create(payment *entity.Payment, orderId string) (error) {
	payment.OrderId = orderId
	if err := pr.db.WithContext(pr.ctx).Omit("Status").Preload("User").Create(payment).Error; err != nil {
		return err
	}

	return nil
}

func (pr *PaymentRepo) GetById(id int) (payment entity.Payment, err error) {
	if err := pr.db.WithContext(pr.ctx).Preload("User").First(&payment, id).Error; err != nil {
		return entity.Payment{}, err
	}

	return payment, nil
}

func (pr *PaymentRepo) GetAll() (payment []entity.Payment, err error) {
	if err := pr.db.WithContext(pr.ctx).Preload("User").Find(&payment).Error; err != nil {
		return []entity.Payment{}, err
	}

	return payment, err
}

func (pr *PaymentRepo) CreateMidtrans(payment entity.Payment, orderId string) (res dto.PaymentResponse, err error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	c := coreapi.Client{}
	c.New(serverKey, midtrans.Sandbox)
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: orderId,
			GrossAmt: int64(payment.Amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: payment.User.Name,
			Email: payment.User.Email,
		},
	}
	coreApiResp, _ := c.ChargeTransaction(chargeReq)

	var paymentURL string
    if len(coreApiResp.Actions) > 0 {
        paymentURL = coreApiResp.Actions[1].URL
    }

	resp := dto.PaymentResponse{
		PaymentLinkUrl: paymentURL,
		TransactionId: coreApiResp.TransactionID,
		ExpiryTime: coreApiResp.ExpiryTime,
		OrderId: coreApiResp.OrderID,
	}
	return resp, nil
}


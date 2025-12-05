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

func (pr *PaymentRepo) GetBidByAuctionId(auctionItemId int) (bid entity.Bid, err error) {
	if err := pr.db.WithContext(context.Background()).First(&bid, "auction_item_id = ?", auctionItemId).Error; err != nil {
		return entity.Bid{}, err
	}

	return bid, nil
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

func (pr *PaymentRepo) CheckPaymentStatusMidtrans(orderId string) (res dto.CheckPaymentStatusResponse, err error) {
	var payment entity.Payment
	var auction entity.AuctionItem

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	c := coreapi.Client{}
	c.New(serverKey, midtrans.Sandbox)

	resp, _ := c.CheckTransaction(orderId)
	// if err != nil {
	// 	return res, err
	// }

	// assign into named variable res
	res = dto.CheckPaymentStatusResponse{
		OrderId:        resp.OrderID,
		TransactionId:  resp.TransactionID,
		PaymentStatus:  resp.TransactionStatus,
	}

	switch resp.TransactionStatus {
	case "settlement":
		pr.db.Model(&payment).
			WithContext(pr.ctx).
			Where("order_id = ?", orderId).
			Update("status", "paid")

	case "cancel", "expire":
		// update auction to scheduled
		pr.db.Model(&auction).
			Where("id = (?)",
				pr.db.Model(&payment).
					Select("auction_item_id").
					Where("order_id = ?", orderId),
			).Update("status", "scheduled")

		// update payment to failed
		pr.db.Model(&payment).
			WithContext(pr.ctx).
			Where("order_id = ?", orderId).
			Update("status", "failed")
	}

	return res, nil
}


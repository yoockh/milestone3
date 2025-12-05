package service

import (
	"fmt"
	"log"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(payment *entity.Payment, orderId string) (error)
	CreateMidtrans(payment entity.Payment, orderId string) (res dto.PaymentResponse, err error)
	CheckPaymentStatusMidtrans(orderId string) (res dto.CheckPaymentStatusResponse, err error)
	GetById(id int) (payment entity.Payment, err error)
	GetAll() (payment []entity.Payment, err error)
}

type PaymentServ struct {
	paymentRepo PaymentRepository
}

func NewPaymentService(pr PaymentRepository) *PaymentServ {
	return &PaymentServ{paymentRepo: pr}
}

func (ps *PaymentServ) CreatePayment(req dto.PaymentRequest, userId int, auctionItemId int) (res dto.PaymentResponse, err error) {
	//get auction for auctionItemId to check if auctionitemid exist or not
	//random id for order id
	uuid := uuid.New()
	orderId := fmt.Sprintf("YDR-%d", uuid.ID())
	
	payment := entity.Payment{
		Amount: req.Amount,
		UserId: userId,
		//hard code for now 
		AuctionItemId: auctionItemId,
	}

	
	log.Println("disini nih")
	resp, _ := ps.paymentRepo.CreateMidtrans(payment, orderId)
	
	if err := ps.paymentRepo.Create(&payment, resp.OrderId); err != nil {
		log.Printf("error create payment %s", err)
		return dto.PaymentResponse{}, err
	}
	return resp, nil
}

func (ps *PaymentServ) CheckPaymentStatusMidtrans(orderId string) (res dto.CheckPaymentStatusResponse, err error) {
	resp, _:= ps.paymentRepo.CheckPaymentStatusMidtrans(orderId)
	// if err != nil {
	// 	log.Printf("error check payment %s", err)
	// 	return dto.CheckPaymentStatusResponse{}, err
	// }

	return resp, nil
}

func (ps *PaymentServ) GetPaymentById(id int) (res dto.PaymentInfoResponse, err error) {
	resp, err := ps.paymentRepo.GetById(id)
	if err != nil {
		log.Printf("failed get payment by id %s", err)
		return dto.PaymentInfoResponse{}, err
	}

	res = dto.PaymentInfoResponse{
		Id: resp.Id,
		UserId: resp.UserId,
		User: resp.User,
		AuctionItemId: resp.AuctionItemId,
		Status: resp.Status,
		// PaymentStatus: resp.PaymentStatus,
		Amount: resp.Amount,
	}

	return res, nil
}

func (ps *PaymentServ) GetAllPayment() (res []dto.PaymentInfoResponse, err error) {
	resp, err := ps.paymentRepo.GetAll()
	if err != nil {
		log.Printf("failed get all payment info %s", err)
		return []dto.PaymentInfoResponse{}, err
	}

	for _, payment := range resp {
		res = append(res, dto.PaymentInfoResponse{
		Id: payment.Id,
		UserId: payment.UserId,
		User: payment.User,
		AuctionItemId: payment.AuctionItemId,
		Status: payment.Status,
		// PaymentStatus: payment.PaymentStatus,
		Amount: payment.Amount,
		})
	}

	return res, nil
}
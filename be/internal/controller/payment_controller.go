package controller

import (
	"milestone3/be/internal/dto"
	"milestone3/be/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PaymentService interface {
	CreatePayment(req dto.PaymentRequest, userId int, auctionItemId int) (res dto.PaymentResponse, err error)
	CheckPaymentStatusMidtrans(orderId string) (res dto.CheckPaymentStatusResponse, err error)
	GetPaymentById(id int) (res dto.PaymentInfoResponse, err error)
	GetAllPayment() (res []dto.PaymentInfoResponse, err error)
}

type PaymentController struct {
	paymentService PaymentService
	validate *validator.Validate
}

func NewPaymentController(validate *validator.Validate, ps PaymentService) *PaymentController {
	return &PaymentController{paymentService:ps, validate: validate}
}

// CreatePayment godoc
// @Summary Create payment for auction item
// @Description Create a payment transaction for a won auction item
// @Tags Your Donate Rise API - Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param auctionId path int true "Auction Item ID"
// @Param payment body dto.PaymentRequest true "Payment details"
// @Success 201 {object} utils.SuccessResponseData "create"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid payload or auction ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /payments/auction/{auctionId} [post]
func (pc *PaymentController) CreatePayment(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	userId := int(claim["id"].(float64))
	req := new(dto.PaymentRequest)

	if err := c.Bind(req); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := pc.validate.Struct(req); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	auctionIdStr := c.Param("auctionId")
	auctionId, err := strconv.Atoi(auctionIdStr)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	resp, err := pc.paymentService.CreatePayment(*req, userId, auctionId)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.CreatedResponse(c, "create", resp)
}

// CheckPaymentStatusMidtrans godoc
// @Summary Check payment status via Midtrans
// @Description Check the payment status of an order through Midtrans payment gateway
// @Tags Your Donate Rise API - Payments
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} utils.SuccessResponseData "ok"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /payments/status/{id} [get]
func (pc *PaymentController) CheckPaymentStatusMidtrans(c echo.Context) error {
	orderId := c.Param("id")
	resp, err := pc.paymentService.CheckPaymentStatusMidtrans(orderId)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "ok", resp)
}

// GetPaymentById godoc
// @Summary Get payment by ID
// @Description Retrieve payment information by payment ID
// @Tags Your Donate Rise API - Payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} utils.SuccessResponseData "ok"
// @Failure 400 {object} utils.ErrorResponse "Bad request - Invalid payment ID"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /payments/{id} [get]
func (pc *PaymentController) GetPaymentById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	resp, err := pc.paymentService.GetPaymentById(id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "ok", resp)
}

// GetAllPayment godoc
// @Summary Get all payments
// @Description Retrieve all payment transactions in the system
// @Tags Your Donate Rise API - Payments
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseData "ok"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /payments [get]
func (pc *PaymentController) GetAllPayment(c echo.Context) error {
	resp, err := pc.paymentService.GetAllPayment()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "ok", resp)
}
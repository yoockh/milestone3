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
	CreatePayment(req dto.PaymentRequest, userId int) (res dto.PaymentResponse, err error)
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

	resp, err := pc.paymentService.CreatePayment(*req, userId)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.CreatedResponse(c, "create", resp)
}

func (pc *PaymentController) CheckPaymentStatusMidtrans(c echo.Context) error {
	orderId := c.Param("id")
	resp, err := pc.paymentService.CheckPaymentStatusMidtrans(orderId)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "ok", resp)
}

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

func (pc *PaymentController) GetAllPayment(c echo.Context) error {
	resp, err := pc.paymentService.GetAllPayment()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "internal server error")
	}

	return utils.SuccessResponse(c, "ok", resp)
}
package routes

import (
	"milestone3/be/api/middleware"
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterPaymentRoutes(paymentCtrl *controller.PaymentController) {
	paymentRoutes := r.echo.Group("/payments")
	paymentRoutes.Use(middleware.JWTMiddleware)
	paymentRoutes.Use(middleware.LoggingMiddleware)

	//payment endpoint
	paymentRoutes.POST("/:auctionId", paymentCtrl.CreatePayment)
	paymentRoutes.GET("/status/:id", paymentCtrl.CheckPaymentStatusMidtrans)
	paymentRoutes.GET("/:id", paymentCtrl.GetPaymentById)
	paymentRoutes.GET("", paymentCtrl.GetAllPayment)
}
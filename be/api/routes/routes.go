package routes

import (
	"milestone3/be/internal/controller"

	"github.com/labstack/echo/v4"
)

type Router interface {
	RegisterArticleRoutes(articleCtrl *controller.ArticleController)
	RegisterDonationRoutes(donationCtrl *controller.DonationController)
	RegisterUserRoutes(userCtrl *controller.UserController)
	RegisterPaymentRoutes(paymentCtrl *controller.PaymentController)
	// RegisterBiddingRoutes(biddingCtrl *controller.BiddingController)
	RegisterFinalDonationRoutes(finalDonationCtrl *controller.FinalDonationController)
	// RegisterAuthRoutes(authCtrl *controller.AuthController)
	RegisterAuctionRoutes(auctionCtrl *controller.AuctionController)
	RegisterAuctionSessionRoutes(sessionCtrl *controller.AuctionSessionController)
}

type EchoRouter struct {
	echo *echo.Echo
}

func NewRouter(e *echo.Echo) Router {
	return &EchoRouter{echo: e}
}

// Example injection method in main:
//  router := routes.NewRouter(e)
//  router.RegisterArticleRoutes(articleCtrl)
//  router.RegisterDonationRoutes(donationCtrl)

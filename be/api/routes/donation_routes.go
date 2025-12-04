package routes

import (
	"milestone3/be/api/middleware"
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterDonationRoutes(donationCtrl *controller.DonationController) {
	donationRoutes := r.echo.Group("/donations")
	donationRoutes.Use(middleware.JWTMiddleware)

	donationRoutes.GET("", donationCtrl.GetAllDonations)
	donationRoutes.GET("/:id", donationCtrl.GetDonationByID)
	donationRoutes.POST("", donationCtrl.CreateDonation)
	donationRoutes.PUT("/:id", donationCtrl.UpdateDonation)
	donationRoutes.PATCH("/:id", donationCtrl.PatchDonation)
	donationRoutes.DELETE("/:id", donationCtrl.DeleteDonation)
}

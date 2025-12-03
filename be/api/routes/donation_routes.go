package routes

import (
	"milestone3/be/api/middleware"
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterDonationRoutes(donationCtrl *controller.DonationController) {
	donationRoutes := r.echo.Group("/donations")

	// public
	donationRoutes.GET("", donationCtrl.GetAllDonations)
	donationRoutes.GET("/:id", donationCtrl.GetDonationByID)

	// authenticated group
	auth := donationRoutes.Group("")
	auth.Use(middleware.JWTMiddleware)

	auth.POST("", donationCtrl.CreateDonation)
	auth.PUT("/:id", donationCtrl.UpdateDonation)
	auth.PATCH("/:id", donationCtrl.PatchDonation)
	auth.DELETE("/:id", donationCtrl.DeleteDonation)
}

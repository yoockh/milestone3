package routes

import (
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterDonationRoutes(donationCtrl *controller.DonationController) {
	donationRoutes := r.echo.Group("/donations")

	// public or auth-protected
	donationRoutes.GET("", donationCtrl.GetAllDonations)
	donationRoutes.GET("/:id", donationCtrl.GetDonationByID)

	// protected: create/update/delete should require auth (and owner/admin checks)
	donationRoutes.POST("", donationCtrl.CreateDonation)
	donationRoutes.PUT("/:id", donationCtrl.UpdateDonation)
	donationRoutes.DELETE("/:id", donationCtrl.DeleteDonation)
}

package routes

import (
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterFinalDonationRoutes(finalDonationCtrl *controller.FinalDonationController) {
	finalDonationRoutes := r.echo.Group("/final_donations")

	// public endpoints
	finalDonationRoutes.GET("", finalDonationCtrl.GetAllFinalDonations)
	finalDonationRoutes.GET("/user/:user_id", finalDonationCtrl.GetAllFinalDonationsByUserID)
}

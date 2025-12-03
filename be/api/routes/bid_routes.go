package routes

import (
	"milestone3/be/api/middleware"
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterBidRoutes(bidCtrl *controller.BidController) {
	g := r.echo.Group("/auction/sessions")

	g.Use(middleware.JWTMiddleware)
	g.Use(middleware.LoggingMiddleware)

	g.POST("/:sessionID/items/:itemID/bid", bidCtrl.PlaceBid)
	g.GET("/:sessionID/items/:itemID/highest-bid", bidCtrl.GetHighestBid)
	g.POST("/:sessionID/items/:itemID/sync", bidCtrl.SyncHighestBid)
}

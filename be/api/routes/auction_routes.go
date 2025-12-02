package routes

import (
	"milestone3/be/api/middleware"
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterAuctionRoutes(auctionCtrl *controller.AuctionController) {
	g := r.echo.Group("/auction/items")
	g.Use(middleware.JWTMiddleware)
	g.Use(middleware.LoggingMiddleware)

	g.GET("", auctionCtrl.GetAllAuctionItems)
	g.GET("/:id", auctionCtrl.GetAuctionItemByID)
	g.POST("", auctionCtrl.CreateAuctionItem)
	g.PUT("/:id", auctionCtrl.UpdateAuctionItem)
	g.DELETE("/:id", auctionCtrl.DeleteAuctionItem)
}

func (r *EchoRouter) RegisterAuctionSessionRoutes(sessionCtrl *controller.AuctionSessionController) {
	g := r.echo.Group("/auction/sessions")
	g.Use(middleware.JWTMiddleware)
	g.Use(middleware.LoggingMiddleware)

	g.GET("", sessionCtrl.GetAllAuctionSessions)
	g.GET("/:id", sessionCtrl.GetAuctionSessionByID)
	g.POST("", sessionCtrl.CreateAuctionSession)
	g.PUT("/:id", sessionCtrl.UpdateAuctionSession)
	g.DELETE("/:id", sessionCtrl.DeleteAuctionSession)
}

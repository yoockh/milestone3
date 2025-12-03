package routes

import (
	"milestone3/be/api/middleware" // import admin middleware
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterArticleRoutes(articleCtrl *controller.ArticleController) {
	articleRoutes := r.echo.Group("/articles")

	// public
	articleRoutes.GET("", articleCtrl.GetAllArticles)
	articleRoutes.GET("/:id", articleCtrl.GetArticleByID)

	// admin-only
	admin := articleRoutes.Group("")
	admin.Use(middleware.JWTMiddleware)
	admin.Use(middleware.RequireAdmin)

	admin.POST("", articleCtrl.CreateArticle)
	admin.PUT("/:id", articleCtrl.UpdateArticle)
	admin.DELETE("/:id", articleCtrl.DeleteArticle)
}

package routes

import (
	"milestone3/be/internal/controller"
)

func (r *EchoRouter) RegisterArticleRoutes(articleCtrl *controller.ArticleController) {
	articleRoutes := r.echo.Group("/articles")

	// public endpoints
	articleRoutes.GET("", articleCtrl.GetAllArticles)
	articleRoutes.GET("/:id", articleCtrl.GetArticleByID)

	// admin-protected (apply middleware at app init or check inside handler)
	articleRoutes.POST("", articleCtrl.CreateArticle)
	articleRoutes.PUT("/:id", articleCtrl.UpdateArticle)
	articleRoutes.DELETE("/:id", articleCtrl.DeleteArticle)
}

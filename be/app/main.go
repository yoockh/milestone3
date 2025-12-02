package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"cloud.google.com/go/storage"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"milestone3/be/api/routes"
	"milestone3/be/config"
	"milestone3/be/internal/controller"
	"milestone3/be/internal/repository"
	"milestone3/be/internal/service"
)

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

func main() {
	ctx := context.Background()
	db := config.ConnectionDb()
	validate := validator.New()

	// create GCS client if configured
	var gcsRepo repository.GCSStorageRepo
	if bucket := os.Getenv("GCS_BUCKET"); bucket != "" {
		gcsClient, err := storage.NewClient(ctx)
		if err != nil {
			// log.Fatalf("failed to create gcs client: %v", err)
		}
		gcsRepo = repository.NewGCSStorageRepo(gcsClient, bucket)
	} else {
		log.Println("GCS_BUCKET not set — file uploads to GCS will fail if used")
	}

	// repositories
	userRepo := repository.NewUserRepo(db, ctx)
	articleRepo := repository.NewArticleRepo(db)
	donationRepo := repository.NewDonationRepo(db)
	finalDonationRepo := repository.NewFinalDonationRepository(db)
	paymentRepo := repository.NewPaymentRepository(db, ctx)
	auctionItemRepo := repository.NewAuctionItemRepository(db)
	auctionSessionRepo := repository.NewAuctionSessionRepository(db)

	bidRepo := repository.NewBidRepository(db)

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	redisClient := redis.NewClient(opt)

	redisRepo := repository.NewBidRedisRepository(redisClient, ctx)
	auctionRedisRepo := repository.NewSessionRedisRepository(redisClient, ctx)

	aiRepo := repository.NewAIRepository(logger, os.Getenv("GEMINI_API_KEY"))

	// services
	userSvc := service.NewUserService(userRepo)
	articleSvc := service.NewArticleService(articleRepo)
	donationSvc := service.NewDonationService(donationRepo)
	finalDonationSvc := service.NewFinalDonationService(finalDonationRepo)
	paymentSvc := service.NewPaymentService(paymentRepo)
	auctionSvc := service.NewAuctionItemService(auctionItemRepo, aiRepo, logger)
	auctionSessionSvc := service.NewAuctionSessionService(auctionSessionRepo, auctionRedisRepo, logger)
	bidSvc := service.NewBidService(redisRepo, bidRepo, auctionItemRepo, logger)

	// controllers
	userCtrl := controller.NewUserController(validate, userSvc)
	articleCtrl := controller.NewArticleController(articleSvc)

	var donationCtrl *controller.DonationController
	if gcsRepo != nil {
		donationCtrl = controller.NewDonationController(donationSvc, gcsRepo)
	} else {
		donationCtrl = controller.NewDonationController(donationSvc, nil)
	}
	finalDonationCtrl := controller.NewFinalDonationController(finalDonationSvc)
	paymentCtrl := controller.NewPaymentController(validate, paymentSvc)
	auctionCtrl := controller.NewAuctionController(auctionSvc, validate)
	auctionSessionCtrl := controller.NewAuctionSessionController(auctionSessionSvc, validate)
	bidCtrl := controller.NewBidController(bidSvc, validate)

	// echo + router
	e := echo.New()
	router := routes.NewRouter(e)

	router.RegisterUserRoutes(userCtrl)
	router.RegisterArticleRoutes(articleCtrl)
	router.RegisterDonationRoutes(donationCtrl)
	router.RegisterFinalDonationRoutes(finalDonationCtrl)
	router.RegisterPaymentRoutes(paymentCtrl)
	router.RegisterAuctionRoutes(auctionCtrl)
	router.RegisterAuctionSessionRoutes(auctionSessionCtrl)
	router.RegisterBidRoutes(bidCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

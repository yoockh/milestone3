package main

import (
	"context"
	"log/slog"
	"milestone3/be/config"
	"milestone3/be/internal/repository"
	"milestone3/be/internal/service"
	"os"

	"github.com/robfig/cron/v3"
)

var (
	loggerOption = slog.HandlerOptions{AddSource: true}
	logger       = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))
)

func main() {
	db := config.ConnectionDb()
	ctx := context.Background()
	redisClient := config.ConnectRedis(ctx)

	redisRepo := repository.NewBidRedisRepository(redisClient, ctx)
	bidRepo := repository.NewBidRepository(db)
	auctionItemRepo := repository.NewAuctionItemRepository(db)

	bidSvc := service.NewBidService(redisRepo, bidRepo, auctionItemRepo, logger)

	c := cron.New()
	c.AddFunc("0 0 0 * * *", func() {
		bidSvc.SaveExpiredSessions(0)
	})
	c.Start()

	select {}
}

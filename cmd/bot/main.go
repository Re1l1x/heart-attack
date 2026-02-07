package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jus1d/kypidbot/internal/config"
	"github.com/jus1d/kypidbot/internal/delivery/telegram"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	"github.com/jus1d/kypidbot/internal/notifications"
	"github.com/jus1d/kypidbot/internal/repository/postgres"
	"github.com/jus1d/kypidbot/internal/usecase"
)

func main() {
	cfg := config.MustLoad()

	var level slog.Level
	switch cfg.Env {
	case config.EnvProduction:
		level = slog.LevelInfo
	default:
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		slog.Error("postgresql: failed to connect", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("postgresql: ok")

	userRepo := postgres.NewUserRepo(db)
	placeRepo := postgres.NewPlaceRepo(db)
	meetingRepo := postgres.NewMeetingRepo(db)

	registration := usecase.NewRegistration(userRepo)
	admin := usecase.NewAdmin(userRepo)
	matching := usecase.NewMatching(userRepo, meetingRepo, &cfg.Ollama)
	meeting := usecase.NewMeeting(userRepo, placeRepo, meetingRepo)

	bot, err := telegram.NewBot(
		cfg.Bot.Token,
		registration,
		admin,
		matching,
		meeting,
		userRepo,
	)
	if err != nil {
		slog.Error("failed to create the bot", sl.Err(err))
		os.Exit(1)
	}

	bot.Setup()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	notifCtx, notifCancel := context.WithCancel(context.Background())
	defer notifCancel()
	go notifications.Go(notifCtx, &cfg.Notifications, bot.TeleBot(), userRepo, placeRepo, meetingRepo)

	go bot.Start()

	<-stop
	slog.Info("bot: shutting down...")
	notifCancel()
	bot.Stop()
}

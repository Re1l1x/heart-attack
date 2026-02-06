package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jus1d/kypidbot/internal/config"
	"github.com/jus1d/kypidbot/internal/delivery/telegram"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	"github.com/jus1d/kypidbot/internal/repository/postgres"
	"github.com/jus1d/kypidbot/internal/usecase"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg := config.MustLoad()

	if err := telegram.LoadMessages("messages.yaml"); err != nil {
		log.Error("failed load message replics", sl.Err(err))
		os.Exit(1)
	}

	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		log.Error("postgresql: failed to connect", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	log.Info("postgresql: ok", sl.Err(err))

	userRepo := postgres.NewUserRepo(db)
	placeRepo := postgres.NewPlaceRepo(db)
	meetingRepo := postgres.NewMeetingRepo(db)

	registration := usecase.NewRegistration(userRepo)
	admin := usecase.NewAdmin(userRepo)
	matching := usecase.NewMatching(userRepo, meetingRepo, &cfg.Ollama)
	meeting := usecase.NewMeeting(userRepo, placeRepo, meetingRepo)

	bot, err := telegram.NewBot(
		cfg.Telegram.Token,
		registration,
		admin,
		matching,
		meeting,
		userRepo,
		log,
	)
	if err != nil {
		log.Error("failed to create the bot", sl.Err(err))
		os.Exit(1)
	}

	bot.Setup()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go bot.Start()

	<-stop
	log.Info("bot: shutting down...")
	bot.Stop()
}

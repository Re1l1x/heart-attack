package command

import (
	"github.com/jus1d/kypidbot/internal/infrastructure/s3"
	"github.com/jus1d/kypidbot/internal/usecase"
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	Registration *usecase.Registration
	Admin        *usecase.Admin
	Matching     *usecase.Matching
	Meeting      *usecase.Meeting
	Bot          *tele.Bot
	S3           *s3.Client
}

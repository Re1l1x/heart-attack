package callback

import (
	"sync"

	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/usecase"
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	Registration *usecase.Registration
	Meeting      *usecase.Meeting
	Users        domain.UserRepository
	Bot          *tele.Bot
	msgMutex     sync.Mutex
	msgIDs       map[string]int
}

func (h *Handler) DeleteAndSend(c tele.Context, what any, opts ...any) error {
	_ = c.Delete()
	return c.Send(what, opts...)
}

func (h *Handler) storeMessageID(key string, msgID int) {
	h.msgMutex.Lock()
	defer h.msgMutex.Unlock()
	if h.msgIDs == nil {
		h.msgIDs = make(map[string]int)
	}
	h.msgIDs[key] = msgID
}

func (h *Handler) getMessageID(key string) int {
	h.msgMutex.Lock()
	defer h.msgMutex.Unlock()
	return h.msgIDs[key]
}

package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type Process struct {
	BotUUID string
	UserID  int64
	Text    string
}

type ProcessHandler decorator.CommandHandler[Process]

type processHandler struct {
	bots         bots.Repository
	participants bots.ParticipantRepository
	msgSender    bots.BotMessageSender
}

func NewProcessHandler(
	bots bots.Repository,
	participants bots.ParticipantRepository,
	msgSender bots.BotMessageSender,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) ProcessHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if participants == nil {
		panic("participants repository is nil")
	}

	if msgSender == nil {
		panic("message sender is nil")
	}

	return decorator.ApplyCommandDecorators[Process](
		processHandler{bots: bots, participants: participants, msgSender: msgSender},
		logger,
		metricsClient,
	)
}

func (h processHandler) Handle(ctx context.Context, cmd Process) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	return h.participants.UpdateOrCreate(ctx, cmd.BotUUID, cmd.UserID, func(
		innerCtx context.Context, prt *bots.Participant,
	) error {
		messages, err := bot.Process(prt, cmd.Text)
		if err != nil {
			return err
		}

		for _, message := range messages {
			err = h.msgSender.Send(innerCtx, bot.Token, cmd.UserID, message)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

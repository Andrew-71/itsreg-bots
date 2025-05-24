package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type GetBot struct {
	UserUUID string
	BotUUID  string
}

type GetBotHandler decorator.QueryHandler[GetBot, Bot]

type getBotHandler struct {
	bots bots.Repository
}

func NewGetBotHandler(
	bots bots.Repository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetBotHandler {
	return decorator.ApplyQueryDecorators[GetBot, Bot](
		getBotHandler{bots: bots},
		logger,
		metricsClient,
	)
}

func (h getBotHandler) Handle(ctx context.Context, query GetBot) (Bot, error) {
	bot, err := h.bots.Bot(ctx, query.BotUUID)
	if err != nil {
		return Bot{}, err
	}

	if query.UserUUID != "" {
		if err = bot.CanSeeBot(query.UserUUID); err != nil {
			return Bot{}, err
		}
	}

	return MapBotFromDomain(bot), nil
}

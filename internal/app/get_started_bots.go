package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/pkg/decorator"
)

type GetStartedBots struct{}

type GetStartedBotsHandler decorator.QueryHandler[GetStartedBots, []Bot]

type getStartedBotsHandler struct {
	bots bots.Repository
}

func NewGetStartedBotsHandler(
	bots bots.Repository,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetStartedBotsHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	return decorator.ApplyQueryDecorators(
		getStartedBotsHandler{bots},
		logger, metricsClient,
	)
}

func (h getStartedBotsHandler) Handle(ctx context.Context, _ GetStartedBots) ([]Bot, error) {
	bs, err := h.bots.BotsWithStatus(ctx, bots.Started)
	if err != nil {
		return nil, err
	}

	return MapBotsFromDomain(bs), nil
}

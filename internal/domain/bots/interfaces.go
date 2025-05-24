package bots

import (
	"context"
	"errors"
)

var ErrRunningInstanceNotFound = errors.New("running instance not found")

type InstanceManager interface {
	Start(ctx context.Context, botId string, token string) error
	Stop(ctx context.Context, botId string) error
}

type BotMessageSender interface {
	Send(ctx context.Context, token string, userId int64, msg Message) error
}

type ProcessHandler interface {
	Process(ctx context.Context, botId string, userId int64, msg Message) error
}

type EntryHandler interface {
	Entry(ctx context.Context, botId string, userId int64, key string) error
}

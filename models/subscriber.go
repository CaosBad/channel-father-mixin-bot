package models

import (
	"context"
	"time"

	"github.com/go-pg/pg"

	"github.com/crossle/channel-father-mixin-bot/session"
)

const subscribers_DDL = `
CREATE TABLE IF NOT EXISTS subscribers (
  bot_id	        VARCHAR(36) NOT NULL,
  subscriber_id     VARCHAR(36) NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (bot_id, subscriber_id)
);
`

type Subscriber struct {
	BotId        string    `sql:"bot_id,pk"`
	SubscriberId string    `sql:"subscriber_id,pk"`
	CreatedAt    time.Time `sql:"created_at"`
}

func CreateSubscriber(ctx context.Context, botId, subscriberId string) (*Subscriber, error) {
	subscriber, err := FindSubscriber(ctx, botId, subscriberId)
	if err != pg.ErrNoRows {
		return subscriber, err
	}
	s := &Subscriber{
		BotId:        botId,
		SubscriberId: subscriberId,
		CreatedAt:    time.Now(),
	}
	if err := session.Database(ctx).Insert(s); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return s, nil
}

func RemoveSubscriber(ctx context.Context, botId, subscriberId string) error {
	_, err := FindSubscriber(ctx, botId, subscriberId)
	if err != nil {
		return err
	}
	s := &Subscriber{
		BotId:        botId,
		SubscriberId: subscriberId,
	}
	if err := session.Database(ctx).Delete(s); err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func FindSubscriber(ctx context.Context, botId, subscriberId string) (*Subscriber, error) {
	s := &Subscriber{
		BotId:        botId,
		SubscriberId: subscriberId,
	}
	if err := session.Database(ctx).Select(s); err != nil {
		return nil, err
	}

	return s, nil
}

func ListSubscribers(ctx context.Context, botId string) ([]*Subscriber, error) {
	var subscribers []*Subscriber
	err := session.Database(ctx).Model(&subscribers).Where("bot_id = ?", botId).Select()
	if err != nil {
		return subscribers, session.TransactionError(ctx, err)
	}
	return subscribers, nil
}

func CountSubscribers(ctx context.Context, botId string) int {
	count, err := session.Database(ctx).Model((*Subscriber)(nil)).Where("bot_id = ?", botId).Count()
	if err != nil {
		return 0
	}
	return count
}

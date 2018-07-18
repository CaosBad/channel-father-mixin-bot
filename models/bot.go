package models

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/go-pg/pg"

	"github.com/crossle/channel-father-mixin-bot/session"
	client "github.com/mixinmessenger/bot-api-go-client"
	uuid "github.com/satori/go.uuid"
)

const bots_DDL = `
CREATE TABLE IF NOT EXISTS bots (
  bot_id            VARCHAR(36) PRIMARY KEY,
  user_id	        VARCHAR(36) NOT NULL,
  trace_id          VARCHAR(36),
  client_id        	VARCHAR(36), 
  session_id	    VARCHAR(36),
  private_key	    VARCHAR(1024),
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  expire_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() 
);

CREATE INDEX ON bots (user_id);
`

type Bot struct {
	BotId      string    `sql:"bot_id,pk"`
	UserId     string    `sql:"user_id,notnull"`
	TraceId    string    `sql:"trace_id"`
	ExpireAt   time.Time `sql:"expire_at,notnull"`
	ClientId   string    `sql:"client_id"`
	SessionId  string    `sql:"session_id"`
	PrivateKey string    `sql:"private_key"`
	CreatedAt  time.Time `sql:"created_at,notnull"`
}

func CreateChannelBot(ctx context.Context, userId, traceId string) (*Bot, error) {
	bot, err := FindBotByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if bot == nil {
		bot = &Bot{
			BotId:     client.NewV4().String(),
			UserId:    userId,
			TraceId:   traceId,
			ExpireAt:  time.Now().AddDate(1, 0, 0),
			CreatedAt: time.Now(),
		}
		if err := session.Database(ctx).Insert(bot); err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		return bot, nil
	}
	bot.ExpireAt = bot.ExpireAt.AddDate(1, 0, 0)
	bot.TraceId = traceId
	if _, err := session.Database(ctx).Model(bot).Column("expire_at", "trace_id").WherePK().Update(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return bot, nil
}

func VerifyTrace(ctx context.Context, user *User, traceId string) (*Bot, error) {
	var bot Bot
	if _, err := session.Database(ctx).QueryOne(&bot, `SELECT bot_id, user_id, trace_id, expire_at, created_at FROM bots WHERE trace_id = ?`, traceId); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return &bot, nil
}

func PostKeys(ctx context.Context, userId, botId, clientId, sessionId, privateKey string) (*Bot, error) {
	bot, err := FindBotByBotId(ctx, botId)
	if err != nil {
		return nil, session.BadDataError(ctx)
	}
	if bot == nil {
		return nil, session.BadDataError(ctx)
	}
	if bot.UserId != userId {
		return nil, session.BadDataError(ctx)
	}
	if bot.ExpireAt.Before(time.Now()) {
		return nil, session.BadDataError(ctx)
	}
	if _, err := uuid.FromString(clientId); err != nil {
		return nil, session.BadDataError(ctx)
	}
	if _, err := uuid.FromString(sessionId); err != nil {
		return nil, session.BadDataError(ctx)
	}
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, session.BadDataError(ctx)
	}
	if _, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return nil, session.BadDataError(ctx)
	}

	bot.ClientId = clientId
	bot.SessionId = sessionId
	bot.PrivateKey = privateKey
	if _, err := session.Database(ctx).Model(bot).Column("client_id", "session_id", "private_key").WherePK().Update(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return bot, nil
}

func FindBotByUserId(ctx context.Context, userId string) (*Bot, error) {
	var bot Bot
	if _, err := session.Database(ctx).QueryOne(&bot, `SELECT bot_id, user_id, trace_id, expire_at, client_id, session_id, private_key, created_at FROM bots WHERE user_id = ?`, userId); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return &bot, nil
}

func FindBotByBotId(ctx context.Context, botId string) (*Bot, error) {
	bot := &Bot{
		BotId: botId,
	}
	if err := session.Database(ctx).Select(bot); err != nil {
		return nil, err
	}
	return bot, nil
}

func FindBotByClientId(ctx context.Context, clientId string) (*Bot, error) {
	var bot Bot
	if _, err := session.Database(ctx).QueryOne(&bot, `SELECT bot_id, user_id, trace_id, expire_at, client_id, session_id, private_key, created_at FROM bots WHERE client_id = ?`, clientId); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return &bot, nil
}

func ListBots(ctx context.Context) ([]*Bot, error) {
	var bots []*Bot
	err := session.Database(ctx).Model(&bots).Where("client_id IS NOT NULL AND session_id IS NOT NULL AND private_key IS NOT NULL AND expire_at > now()").Select()
	if err != nil {
		return bots, session.TransactionError(ctx, err)
	}
	return bots, nil
}

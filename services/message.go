package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/crossle/channel-father-mixin-bot/config"
	"github.com/crossle/channel-father-mixin-bot/models"
	"github.com/crossle/channel-father-mixin-bot/session"
	"github.com/mixinmessenger/bot-api-go-client"
	uuid "github.com/satori/go.uuid"
)

type BlazeService struct{}

func (service *BlazeService) Run(ctx context.Context) error {
	for {
		if err := bot.Loop(ctx, ResponseMessage{}, config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey); err != nil {
			session.Logger(ctx).Error(err)
		}
		session.Logger(ctx).Info("connection loop end")
		time.Sleep(time.Second)
	}
	return nil
}

type ResponseMessage struct{}

func (r ResponseMessage) OnMessage(ctx context.Context, mc *bot.MessageContext, msg bot.MessageView, uid string) error {
	if msg.Category == bot.MessageCategorySystemAccountSnapshot {
		data, err := base64.StdEncoding.DecodeString(msg.Data)
		if err != nil {
			return bot.BlazeServerError(ctx, err)
		}
		var transfer bot.TransferView
		err = json.Unmarshal(data, &transfer)
		if err != nil {
			return bot.BlazeServerError(ctx, err)
		}
		err = handleTransfer(ctx, transfer, msg.UserId)
		if err == nil {
			sendPaidMessage(ctx, mc, msg)
		}
	}
	return nil
}

func handleTransfer(ctx context.Context, transfer bot.TransferView, userId string) error {
	_, err := uuid.FromString(transfer.TraceId)
	if err != nil {
		return nil
	}
	if transfer.Amount != config.PayAmount || transfer.AssetId != config.PayAssetId {
		return session.BadDataError(ctx)
	}
	_, err = models.CreateChannelBot(ctx, userId, transfer.TraceId)
	return err
}
func sendPaidMessage(ctx context.Context, mc *bot.MessageContext, msg bot.MessageView) error {
	content := `您已付费，可以开始创建频道了, 登录 https://developers.mixin.one 创建一个机器人，复制 UserId, SessionId 和 PrivateKey 提交到网页 (如何获取，请看下图)`
	if err := bot.SendPlainText(ctx, mc, msg, content); err != nil {
		return bot.BlazeServerError(ctx, err)
	}

	imageMap := map[string]interface{}{
		"attachment_id": "2cd57e11-a58e-4705-a47e-77f0586c915e",
		"size":          316047,
		"width":         1532,
		"mime_type":     "image/jpeg",
		"height":        1098,
	}
	imageData, _ := json.Marshal(imageMap)
	if err := bot.SendMessage(ctx, mc, msg.ConversationId, msg.UserId, bot.MessageCategoryPlainImage, string(imageData), ""); err != nil {
		return bot.BlazeServerError(ctx, err)
	}
	return nil
}

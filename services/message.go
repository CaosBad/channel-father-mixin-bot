package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/crossle/channel-father-mixin-bot/config"
	"github.com/crossle/channel-father-mixin-bot/models"
	"github.com/crossle/channel-father-mixin-bot/session"
	uuid "github.com/satori/go.uuid"
)

type BlazeService struct {
}

type ResponseMessage struct {
	blazeClient *bot.BlazeClient
}

func (service *BlazeService) Run(ctx context.Context) error {
	for {
		blazeClient := bot.NewBlazeClient(config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
		response := ResponseMessage{
			blazeClient: blazeClient,
		}
		if err := blazeClient.Loop(ctx, response); err != nil {
			session.Logger(ctx).Error(err)
		}
		session.Logger(ctx).Info("connection loop end")
		time.Sleep(time.Second)
	}
	return nil
}

func (r ResponseMessage) OnMessage(ctx context.Context, msg bot.MessageView, uid string) error {
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
			sendPaidMessage(ctx, r.blazeClient, msg)
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

func sendPaidMessage(ctx context.Context, client *bot.BlazeClient, msg bot.MessageView) error {
	content := `您已付费，可以开始创建频道了, 登录 https://developers.mixin.one 创建一个机器人，复制 UserId, SessionId 和 PrivateKey 提交到网页 (如何获取，请看下图)`
	if err := client.SendPlainText(ctx, msg, content); err != nil {
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
	if err := client.SendMessage(ctx, msg.ConversationId, msg.UserId, bot.MessageCategoryPlainImage, string(imageData), ""); err != nil {
		return bot.BlazeServerError(ctx, err)
	}
	return nil
}

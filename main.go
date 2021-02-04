package main

import (
	"fmt"
	"os"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log := zap.NewProductionEncoderConfig()
	log.EncodeLevel = zapcore.CapitalLevelEncoder
	log.EncodeTime = zapcore.RFC3339TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(log),
		os.Stdout,
		zap.InfoLevel,
	))

	b, err := gotgbot.NewBot("YOUR_TELEGRAM_KEY")
	if err != nil {
		logger.Panic("failed to create new bot: " + err.Error())
		return
	}

	updater := ext.NewUpdater(b, nil)
	dispatcher := updater.Dispatcher

	// Add echo handler to reply to all messages.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("uid_callback"), uidCallback))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("help_callback"), helpCallback))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{Clean: true})
	if err != nil {
		logger.Panic("failed to start polling: " + err.Error())
		return
	}
	logger.Sugar().Info(fmt.Sprintf("%s has been started...\n", b.User.Username))

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

// start introduces the bot
func start(ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		ctx.Bot,
		fmt.Sprintf(
			"Hello, I'm @%s. Your <b>personal assistant</b>, how can i help you?",
			ctx.Bot.User.Username,
		),
		&gotgbot.SendMessageOpts{
			ParseMode: "html",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
					{Text: "My Telegram ID", CallbackData: "uid_callback"},
					{Text: "Please Help Me!", CallbackData: "help_callback"},
				}},
			},
		},
	)

	if err != nil {
		fmt.Println("failed to send: " + err.Error())
	}

	return nil
}

// callback
func uidCallback(ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery
	cb.Answer(ctx.Bot, nil)
	cb.Message.EditText(
		ctx.Bot,
		fmt.Sprintf(
			"Your Telegram ID : %d, \nYour Chat ID : %d",
			ctx.Bot.User.Id,
			ctx.EffectiveChat.Id,
		),
		nil)
	return nil
}

func helpCallback(ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery
	cb.Answer(ctx.Bot, nil)
	cb.Message.EditText(ctx.Bot, "Hello from help.", nil)
	return nil
}

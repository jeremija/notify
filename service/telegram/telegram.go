package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// ParseMode desecribes the message format.
type ParseMode string

const (
	HTML     ParseMode = tgbotapi.ModeHTML
	Markdown ParseMode = tgbotapi.ModeMarkdown
)

// Telegram struct holds necessary data to communicate with the Telegram API.
type Telegram struct {
	client    *tgbotapi.BotAPI
	chatIDs   []int64
	parseMode ParseMode
}

// New returns a new instance of a Telegram notification service.
// For more information about telegram api token:
//    -> https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api#NewBotAPI
func New(apiToken string, opts ...Opt) (*Telegram, error) {
	client, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}

	t := &Telegram{
		client:    client,
		chatIDs:   []int64{},
		parseMode: HTML,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t, nil
}

type Opt func(t *Telegram)

// WithParseMode sets the message parse mode. Default is HTML.
func WithParseMode(parseMode ParseMode) Opt {
	return func(t *Telegram) {
		t.parseMode = parseMode
	}
}

// AddReceivers takes Telegram chat IDs and adds them to the internal chat ID list. The Send method will send
// a given message to all those chats.
func (t *Telegram) AddReceivers(chatIDs ...int64) {
	t.chatIDs = append(t.chatIDs, chatIDs...)
}

// Send takes a message subject and a message body and sends them to all previously set chats. Message body supports
// html as markup language.
func (t Telegram) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	msg := tgbotapi.NewMessage(0, fullMessage)
	msg.ParseMode = string(t.parseMode)

	for _, chatID := range t.chatIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg.ChatID = chatID
			_, err := t.client.Send(msg)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Telegram chat '%d'", chatID)
			}
		}
	}

	return nil
}

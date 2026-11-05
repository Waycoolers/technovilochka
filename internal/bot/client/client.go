package client

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"technovilochka/internal/bot/transport"
	"technovilochka/internal/config"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func clearOldUpdates(api *tgbotapi.BotAPI) error {
	var lastID int

	for {
		updates, err := api.GetUpdates(tgbotapi.UpdateConfig{
			Offset:  0,
			Limit:   100,
			Timeout: 0,
		})
		if err != nil {
			return err
		}

		if len(updates) == 0 {
			break
		}

		lastID = updates[len(updates)-1].UpdateID
		_, _ = api.GetUpdates(tgbotapi.UpdateConfig{
			Offset: lastID + 1,
			Limit:  1,
		})
	}

	return nil
}

func convertUpdate(upd tgbotapi.Update) transport.Update {
	if upd.Message != nil {
		userName, firstName := userNames(upd.Message.From)
		return transport.Update{
			Message: &transport.Message{
				ChatID:    upd.Message.Chat.ID,
				UserName:  userName,
				FirstName: firstName,
				Text:      upd.Message.Text,
			},
			CallbackQuery: nil,
		}
	} else if upd.CallbackQuery != nil && upd.CallbackQuery.Message != nil {
		userName, _ := userNames(upd.CallbackQuery.From)
		return transport.Update{
			Message: nil,
			CallbackQuery: &transport.CallbackQuery{
				ChatID:    upd.CallbackQuery.Message.Chat.ID,
				MessageID: upd.CallbackQuery.Message.MessageID,
				Data:      upd.CallbackQuery.Data,
				UserName:  userName,
				Message:   upd.CallbackQuery.Message.Text,
			},
		}
	}
	return transport.Update{}
}

type Client struct {
	api            *tgbotapi.BotAPI
	updatesTimeout int
}

// NewClient теперь возвращает *Client, поэтому соответствие интерфейсу проверяем явно.
var _ transport.BotClient = (*Client)(nil)

func NewClient(cfg *config.BotConfig) (*Client, error) {
	// tgbotapi сам пишет ошибки поллинга в свой логгер, а в тексте ошибки лежит URL с токеном.
	_ = tgbotapi.SetLogger(tokenRedactingLogger{token: cfg.Token})

	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("telegram client init: %w", redactToken(err, cfg.Token))
	}
	return &Client{api: api, updatesTimeout: cfg.UpdatesTimeout}, nil
}

func (c *Client) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if _, err := c.api.Send(msg); err != nil {
		slog.Error("Telegram client send message error", "chatID", chatID, "error", c.redact(err))
		return
	}
	slog.Info("Bot send message", "chatID", chatID, "text", text)
}

func convertInlineKeyboard(keyboard transport.InlineKeyboard) tgbotapi.InlineKeyboardMarkup {
	var tgRaws [][]tgbotapi.InlineKeyboardButton
	for _, raw := range keyboard.Rows {
		var tgRaw []tgbotapi.InlineKeyboardButton
		for _, btn := range raw.Buttons {
			tgBtn := tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Data)
			tgRaw = append(tgRaw, tgBtn)
		}
		tgRaws = append(tgRaws, tgRaw)
	}
	tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgRaws...)
	return tgKeyboard
}

func (c *Client) SendWithInlineKeyboard(chatID int64, text string, keyboard transport.InlineKeyboard) error {
	markup := convertInlineKeyboard(keyboard)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	msg.ParseMode = "HTML"
	if _, err := c.api.Send(msg); err != nil {
		err = c.redact(err)
		slog.Error("Telegram client send message with buttons error", "chatID", chatID, "error", err)
		return err
	}
	slog.Info("Bot send message", "chatID", chatID, "text", text)
	return nil
}

func (c *Client) EditMessageReplyMarkup(chatID int64, messageID int, keyboard transport.InlineKeyboard) error {
	markup := convertInlineKeyboard(keyboard)

	edit := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, markup)
	_, err := c.api.Request(edit)
	return c.redact(err)
}

func (c *Client) DeleteMessageReplyMarkup(chatID int64, messageID int) error {
	empty := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}
	del := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, empty)
	_, err := c.api.Request(del)
	return c.redact(err)
}

func (c *Client) GetUpdatesChan() <-chan transport.Update {
	if err := clearOldUpdates(c.api); err != nil {
		slog.Error("Error while cleaning up old updates", "error", c.redact(err))
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = c.updatesTimeout
	tgChan := c.api.GetUpdatesChan(u)
	domainChan := make(chan transport.Update)

	go func() {
		defer close(domainChan)
		for update := range tgChan {
			domainChan <- convertUpdate(update)
		}
	}()

	return domainChan
}

func (c *Client) StopReceivingUpdates() {
	c.api.StopReceivingUpdates()
}

func (c *Client) SendKeyboard(chatID int64, text string, keyboard transport.Keyboard) (transport.Message, error) {
	markup := convertKeyboard(keyboard)

	markup.ResizeKeyboard = true
	markup.OneTimeKeyboard = false

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = markup

	botMessage, err := c.api.Send(msg)
	if err != nil {
		return transport.Message{}, c.redact(err)
	}
	userName, firstName := userNames(botMessage.From)
	message := transport.Message{
		ChatID:    botMessage.Chat.ID,
		UserName:  userName,
		FirstName: firstName,
		Text:      botMessage.Text,
	}
	return message, nil
}

func (c *Client) DeleteMessage(chatID int64, messageID int) error {
	req := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := c.api.Request(req)
	return c.redact(err)
}

func convertKeyboard(keyboard transport.Keyboard) tgbotapi.ReplyKeyboardMarkup {
	var tgRaws [][]tgbotapi.KeyboardButton
	for _, raw := range keyboard.Rows {
		var tgRaw []tgbotapi.KeyboardButton
		for _, btn := range raw.Buttons {
			tgBtn := tgbotapi.NewKeyboardButton(string(btn.Command))
			tgRaw = append(tgRaw, tgBtn)
		}
		tgRaws = append(tgRaws, tgRaw)
	}
	tgKeyboard := tgbotapi.NewReplyKeyboard(tgRaws...)
	return tgKeyboard
}

func userNames(user *tgbotapi.User) (userName, firstName string) {
	if user == nil {
		return "", ""
	}
	return user.UserName, user.FirstName
}

func (c *Client) redact(err error) error {
	return redactToken(err, c.api.Token)
}

// redactToken убирает токен из ошибки: tgbotapi возвращает *url.Error,
// у которого в URL лежит https://api.telegram.org/bot<token>/<method>.
func redactToken(err error, token string) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		urlErr.URL = redactString(urlErr.URL, token)
	}
	return err
}

func redactString(s, token string) string {
	if token == "" {
		return s
	}
	return strings.ReplaceAll(s, token, "<token>")
}

// tokenRedactingLogger перенаправляет внутренние логи tgbotapi в slog и вырезает из них токен.
type tokenRedactingLogger struct {
	token string
}

func (l tokenRedactingLogger) Println(v ...any) {
	l.log(fmt.Sprintln(v...))
}

func (l tokenRedactingLogger) Printf(format string, v ...any) {
	l.log(fmt.Sprintf(format, v...))
}

func (l tokenRedactingLogger) log(msg string) {
	slog.Warn("Telegram library message", "message", strings.TrimSpace(redactString(msg, l.token)))
}

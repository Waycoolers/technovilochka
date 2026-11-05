package transport

type Message struct {
	ChatID    int64
	UserName  string
	FirstName string
	Text      string
}

type CallbackQuery struct {
	ChatID    int64
	MessageID int
	Data      string
	UserName  string
	Message   string
}

type InlineKeyboard struct {
	Rows []InlineKeyboardRow
}

type InlineKeyboardRow struct {
	Buttons []InlineKeyboardButton
}

type InlineKeyboardButton struct {
	Text string
	Data string
}

type Keyboard struct {
	Rows []KeyboardRow
}

type KeyboardRow struct {
	Buttons []KeyboardButton
}

type KeyboardButton struct {
	Command Command
}

type Update struct {
	Message       *Message
	CallbackQuery *CallbackQuery
}

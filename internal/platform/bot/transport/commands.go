package transport

type Command string

type Callback string

const (
	Start Command = "/start"
	Help  Command = "/help"
)

var Commands = []Command{
	Start,
	Help,
}

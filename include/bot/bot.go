package bot

import (
	"github.com/mymmrac/telego"
)

func InitBot(token string) (*telego.Bot, error) {
	bot, err := telego.NewBot(token)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

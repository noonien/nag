package plugins

import "github.com/noonien/nag/bot"

type RegisterSnoonet struct {
	Password string
	bot      *bot.Bot
}

func (p *RegisterSnoonet) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b
	b.Handle("irc.001", p.welcome)

	return &bot.PluginInfo{
		Name:        "RegisterSnoonet",
		Author:      "noonien",
		Description: "Registers snoonet nickname",
		Version:     "1.0",
	}, nil
}

func (p *RegisterSnoonet) Unload() error {
	return nil
}

func (p *RegisterSnoonet) welcome(name string, params []interface{}) (bool, error) {
	if len(p.Password) > 0 {
		p.bot.Message(bot.PrivMsg("NickServ", "identify "+p.Password))
	}
	return false, nil
}

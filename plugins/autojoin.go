package plugins

import "github.com/noonien/nag/bot"

type AutoJoin struct {
	Channels []string

	bot *bot.Bot
}

func (p *AutoJoin) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b
	b.Handle("irc.900", p.welcome)

	return &bot.PluginInfo{
		Name:        "AutoJoin",
		Author:      "noonien",
		Description: "Auto joins channels upon connect",
		Version:     "1.0",
	}, nil
}

func (p *AutoJoin) welcome(name string, params []interface{}) (bool, error) {
	for _, channel := range p.Channels {
		err := p.bot.Message(bot.Join(channel))
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (p *AutoJoin) Unload() error {
	return nil
}

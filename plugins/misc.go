package plugins

import (
	"fmt"
	"strings"

	"github.com/noonien/nag/bot"
	"github.com/sorcix/irc"
)

type Misc struct {
	bot *bot.Bot
}

func (p *Misc) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b

	p.textCmd("cmd.macanache", "nu fi dilimache")

	p.textCmd("cmd.next", "Altă întrebare!")
	p.textReply("irc.privmsg", "Altă întrebare!", func(line string) bool {
		line = strings.ToLower(line)
		return strings.HasSuffix(line, ", next") || strings.HasSuffix(line, " next!")
	})

	p.bot.HandleCmdRateLimited("cmd.bullshit", p.bullshit)
	p.bot.HandleCmdRateLimited("cmd.bs", p.bullshit)

	return &bot.PluginInfo{
		Name:        "Misc",
		Author:      "noonien",
		Description: "Misc cmds",
		Version:     "1.0",
	}, nil
}

func (p *Misc) Unload() error {
	return nil
}

func (p *Misc) textCmd(cmd, text string) {
	handler := func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		p.bot.Message(bot.PrivMsg(target, text))

		return true, nil
	}

	p.bot.HandleCmdRateLimited(cmd, handler)
}

func (p *Misc) textReply(cmd, text string, check func(string) bool) {
	handler := func(msg *irc.Message) (bool, error) {
		if !check(msg.Trailing) {
			return false, nil
		}

		if p.bot.RateLimiter.Limited(msg.Params[0]) {
			return false, nil
		}

		p.bot.Message(bot.PrivMsg(msg.Params[0], text))

		return false, nil
	}

	p.bot.HandleIRC(cmd, handler)
}

func (p *Misc) bullshit(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
	if len(args) == 0 {
		return true, nil
	}

	msg := fmt.Sprintf("%s: Dragnea s-a bullshit pe tine!", args[0])
	p.bot.Message(bot.PrivMsg(target, msg))
	return true, nil
}

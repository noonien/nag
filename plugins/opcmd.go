package plugins

import (
	"fmt"

	"github.com/noonien/nag/bot"
	"github.com/sorcix/irc"
)

type OPCmd struct {
	bot *bot.Bot
}

func (p *OPCmd) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b
	b.HandleCmd("cmd.kb", p.kickban)
	b.HandleCmd("cmd.bafta", p.kickban)

	return &bot.PluginInfo{
		Name:        "OPCmd",
		Author:      "noonien",
		Description: "OP Commands",
		Version:     "1.0",
	}, nil
}

func (p *OPCmd) Unload() error {
	return nil
}

func (p *OPCmd) kickban(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
	if len(args) != 1 {
		return true, nil
	}

	perms, err := p.bot.Auth(source)
	if err != nil {
		return false, err
	}

	if perms == nil || !perms.Can("opcmds") {
		return true, nil
	}

	whom := args[0]
	p.bot.Message(bot.PrivMsg("ChanServ", fmt.Sprintf("ban %s %s", target, whom)))
	return true, nil
}

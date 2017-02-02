package bot

import (
	"strings"

	"github.com/sorcix/irc"
)

type Handler func(name string, params []interface{}) (bool, error)

type CmdHandler func(source *irc.Prefix, target string, cmd string, args []string) (bool, error)

type IRCHandler func(msg *irc.Message) (bool, error)

func (b *Bot) registerHandlers() {
	b.HandleIRC("irc.connect", b.connect)
	b.HandleIRC("irc.ping", b.ping)
	b.HandleIRC("irc.cap", b.caps)
	b.HandleIRC("irc.privmsg", b.cmd)
}

func (b *Bot) connect(msg *irc.Message) (bool, error) {
	messages := []*irc.Message{
		// list server capabilities
		&irc.Message{
			Command: irc.CAP,
			Params:  []string{irc.CAP_LS, "302"},
		},

		// register
		&irc.Message{
			Command: irc.NICK,
			Params:  []string{b.Config.Nickname},
		},
		&irc.Message{
			Command:  irc.USER,
			Params:   []string{b.Config.Username, "0", "*"},
			Trailing: b.Config.Username,
		},
	}

	for _, msg := range messages {
		err := b.Message(msg)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (b *Bot) ping(msg *irc.Message) (bool, error) {
	return true, b.Message(&irc.Message{
		Command:  irc.PONG,
		Params:   msg.Params,
		Trailing: msg.Trailing,
	})
}

func (b *Bot) caps(msg *irc.Message) (bool, error) {
	switch msg.Params[1] {
	case irc.CAP_LS:
		// FIXME: use theese?
		capabilities := strings.Fields(msg.Trailing)
		err := b.Set("server.capabilities", capabilities)
		if err != nil {
			return false, err
		}

		return false, b.Message(&irc.Message{
			Command: irc.CAP,
			Params:  []string{irc.CAP_END},
		})
	}
	return false, nil
}

func (b *Bot) cmd(msg *irc.Message) (bool, error) {
	source := msg.Prefix
	target := msg.Params[0]
	line := msg.Trailing

	// check if line has cmd prefix
	if !strings.HasPrefix(line, b.Config.CmdPrefix) {
		return false, nil
	}

	// remove prefix
	line = line[len(b.Config.CmdPrefix):]

	// split into parts
	parts := strings.Fields(line)
	if len(parts) == 0 || len(parts[0]) == 0 {
		return false, nil
	}

	return b.Event("cmd."+parts[0], source, target, parts[0], parts[1:])
}

func (b *Bot) HandleIRC(name string, handler IRCHandler) {
	b.Handle(name, func(name string, params []interface{}) (bool, error) {
		msg := params[0].(*irc.Message)
		return handler(msg)
	})
}

func (b *Bot) HandleCmd(name string, handler CmdHandler) {
	b.Handle(name, func(name string, params []interface{}) (bool, error) {
		source := params[0].(*irc.Prefix)
		target := params[1].(string)
		cmd := params[2].(string)
		args := params[3].([]string)
		return handler(source, target, cmd, args)
	})
}

func (b *Bot) HandleCmdRateLimited(name string, handler CmdHandler) {
	b.Handle(name, func(name string, params []interface{}) (bool, error) {
		source := params[0].(*irc.Prefix)
		target := params[1].(string)

		if b.RateLimiter.Limited(target) {
			return false, nil
		}

		cmd := params[2].(string)
		args := params[3].([]string)
		return handler(source, target, cmd, args)
	})
}

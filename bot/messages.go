package bot

import (
	"fmt"
	"strings"

	"github.com/sorcix/irc"
)

func PrivMsg(target, message string) *irc.Message {
	return &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   []string{target},
		Trailing: message,
	}
}

func Join(channel string) *irc.Message {
	return &irc.Message{
		Command: irc.JOIN,
		Params:  []string{channel},
	}
}

func Ban(channel string, masks ...string) *irc.Message {
	mode := fmt.Sprintf("%s +%s %s", channel, strings.Repeat("b", len(masks)), strings.Join(masks, ""))
	return &irc.Message{
		Command: irc.MODE,
		Params:  []string{mode},
	}
}

func Kick(channel string, nick string) *irc.Message {
	return &irc.Message{
		Command: irc.KICK,
		Params:  []string{channel, nick},
	}
}

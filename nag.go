package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/noonien/nag/bot"
	"github.com/noonien/nag/plugins"
	"github.com/sorcix/irc"
	"github.com/turnage/graw/reddit"
)

var noobs = map[string][]string{
	"user/nuunien":       []string{"all"},
	"fish.programmer":    []string{"annoy", "invite"},
	"user/victorrrrrr":   []string{"annoy", "invite"},
	"user/jupiter-crash": []string{"annoy", "invite"},
	"sch":                []string{"annoy"},
}

func main() {
	conf := bot.Config{
		Nickname:  "naagien",
		Username:  "naagien",
		CmdPrefix: ".",
	}

	auth := bot.AuthFunc(func(mask *irc.Prefix) (bot.Permissions, error) {
		perms, ok := noobs[mask.Host]
		if !ok {
			return nil, nil
		}

		return bot.PermissionsFunc(func(name string) bool {
			for _, perm := range perms {
				if perm == name || perm == "all" {
					return true
				}

			}
			return false
		}), nil
	})

	b := bot.New(conf, auth, make(MapStore))

	nsPass := os.Getenv("NAG_NSPASSWORD")
	b.LoadPlugin(&plugins.RegisterSnoonet{Password: nsPass})
	b.LoadPlugin(&plugins.AutoJoin{Channels: []string{"#Romania"}})
	b.LoadPlugin(&plugins.OPCmd{})
	b.LoadPlugin(&plugins.Misc{})

	rs, _ := reddit.NewScript("naggien, a bot by /u/nuunien", 3*time.Second)
	b.LoadPlugin(&plugins.RedditFilth{Lurker: rs})

	// print chat
	b.HandleIRC("irc.*", func(msg *irc.Message) (bool, error) {
		switch strings.ToLower(msg.Command) {
		case "privmsg":
			log.Printf("%s | <%s> %s\n", msg.Params[0], msg.Prefix.Name, msg.Trailing)
		}

		return false, nil
	})

	err := b.DialWithSSL("eu-irc.snoonet.org:6697", nil)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	_ = <-c
}

package bot

import (
	"crypto/tls"
	"io"
	"net"
	"strings"
	"time"

	"github.com/sorcix/irc"
)

type Config struct {
	Nickname  string
	Username  string
	CmdPrefix string

	RateLimitMessages int
	RateLimitDuration time.Duration
}

type Bot struct {
	Config Config

	Dispatcher
	Auther
	Store

	RateLimiter *RateLimiter

	plugins map[string]*botPlugin
	ic      *irc.Conn
}

func New(config Config, auth Auther, store Store) *Bot {
	if config.RateLimitMessages <= 0 || config.RateLimitDuration == 0 {
		config.RateLimitMessages = 3
		config.RateLimitDuration = 10 * time.Second

	}

	bot := &Bot{
		Config: config,

		Dispatcher: &trieDispatcher{},
		Auther:     auth,

		RateLimiter: NewRateLimiter(config.RateLimitMessages, config.RateLimitDuration),

		plugins: make(map[string]*botPlugin),
	}
	bot.Store = &botStore{bot, store}

	bot.registerHandlers()
	return bot
}

func (b *Bot) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	b.Connect(conn)
	return nil
}

func (b *Bot) DialWithSSL(addr string, config *tls.Config) error {
	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	b.Connect(conn)
	return nil
}

func (b *Bot) Connect(conn net.Conn) {
	b.ic = irc.NewConn(conn)

	go func() {
		b.Event("irc.connect", &irc.Message{})
		for {
			conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
			msg, err := b.ic.Decode()
			if err != nil {
				if err == io.EOF {
					b.Event("irc.disconnect")
				} else if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					b.Event("irc.timeout")
				} else {
					b.Event("irc.error", err)
				}

				return
			}

			_, err = b.Event("irc."+strings.ToLower(msg.Command), msg)
			if err != nil {
				b.Event("irc.error", err)
				return
			}
		}
	}()
}

func (b *Bot) Message(message *irc.Message) error {
	return b.ic.Encode(message)
}

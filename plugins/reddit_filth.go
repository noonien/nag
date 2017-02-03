package plugins

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/noonien/nag/bot"
	"github.com/sorcix/irc"
	"github.com/turnage/graw/reddit"
)

type command struct {
	Cmds  []string
	Subs  []string
	Ce    []string
	Check func(*reddit.Post) bool
}

var commands = []command{
	command{
		Cmds: []string{
			"buci", "booci", "cur",
		},
		Subs: []string{
			"AssOnTheGlass", "BoltedOnBooty", "ButtsAndBareFeet", "Cheeking",
			"HighResASS", "LoveToWatchYouLeave", "NoTorso", "SpreadEm",
			"TheUnderbun", "Top_Tier_Asses", "Tushy", "Underbun", "ass",
			"assgifs", "booty", "booty_gifs", "datass", "datbuttfromthefront",
			"twerking",
		},
		Ce: []string{
			"buci", "booci", "cur", "curuletz", "funduletz",
		},
		Check: checkIsImage,
	},
	command{
		Cmds: []string{
			"țț", "țâțe", "țațe", "tzatze", "boobs",
		},
		Subs: []string{
			"Bigtitssmalltits", "Boobies", "BreastEnvy", "EpicCleavage",
			"JustOneBoob", "OneInOneOut", "PM_ME_YOUR_TITS_GIRL",
			"PerfectTits", "Perky", "Rush_Boobs", "Saggy", "SloMoBoobs",
			"TheHangingBoobs", "TheUnderboob", "Titsgalore", "TittyDrop",
			"bananatits", "boobbounce", "boobgifs", "boobkarma", "boobland",
			"boobs", "breastplay", "cleavage", "feelthemup", "handbra",
			"hanging", "homegrowntits", "knockers", "naturaltitties",
			"sideboob", "tits", "titsagainstglass", "titties_n_kitties",
			"torpedotits", "underboob",
		},
		Ce: []string{
			"țț", "țâțe", "țațe", "tzatze",
		},
		Check: checkIsImage,
	},
	command{
		Cmds: []string{
			"lips", "buze",
		},
		Subs: []string{
			"lips", "lipsthatgrip",
		},
		Ce: []string{
			"buze", "lips",
		},
		Check: checkIsImage,
	},
}

type RedditFilth struct {
	PreloadCount int
	Lurker       reddit.Lurker

	bot   *bot.Bot
	close chan bool
}

func (p *RedditFilth) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b
	p.close = make(chan bool)

	if p.PreloadCount < 1 {
		p.PreloadCount = 10
	}

	for _, cmd := range commands {
		p.registerFilth(cmd.Cmds, cmd.Subs, cmd.Ce, cmd.Check)
	}

	return &bot.PluginInfo{
		Name:        "RedditFilth",
		Author:      "noonien",
		Description: "Reddit RedditFilth",
		Version:     "1.0",
	}, nil
}

func (p *RedditFilth) Unload() error {
	close(p.close)
	return nil
}

func (p *RedditFilth) newMizerie(subs []string, check func(*reddit.Post) bool) *mizerie {
	m := mizerie{
		Posts: make([]*reddit.Post, 0, p.PreloadCount),
		close: p.close,
	}

	go func() {
		for {
			select {
			case <-m.close:
				return
			case <-time.After(2 * time.Second):
				m.mu.Lock()
				full := len(m.Posts) == cap(m.Posts)
				m.mu.Unlock()
				if full {
					continue
				}

				sub := subs[rand.Intn(len(subs))]
				for {
					post, err := p.Lurker.Thread("/r/" + sub + "/random")
					if err != nil {
						log.Printf("error while getting random post from %s: %v\n", sub, err)
						sub = subs[rand.Intn(len(subs))]
						continue
					}

					if check != nil && !check(post) {
						continue
					}

					m.mu.Lock()
					m.Posts = append(m.Posts, post)
					m.mu.Unlock()
					break
				}
			}
		}
	}()

	return &m
}

type mizerie struct {
	Posts []*reddit.Post

	mu    sync.Mutex
	close chan bool
}

func (m *mizerie) Get() *reddit.Post {
	for i := 0; i < 20; i++ {
		m.mu.Lock()
		var post *reddit.Post
		if len(m.Posts) > 0 {
			post = m.Posts[len(m.Posts)-1]
			m.Posts = m.Posts[:len(m.Posts)-1]
		}
		m.mu.Unlock()

		if post != nil {
			return post
		}

		select {
		case <-m.close:
			return nil
		case <-time.After(time.Second):
		}
	}

	return nil
}

var imageHosts = []string{"imgur.com", "gfycat.com", "giphy.com", "i.redditmedia.com", "i.redd.it", "media.tumblr.com"}

func checkIsImage(post *reddit.Post) bool {
	linkURL, err := url.Parse(post.URL)
	if err != nil {
		return false
	}

	for _, host := range imageHosts {
		if strings.Contains(linkURL.Host, host) {
			return true
		}
	}

	return false
}

func chooseRandStr(opt []string) string {
	return opt[rand.Intn(len(opt))]
}

func (p *RedditFilth) registerFilth(cmds, subs, ce []string, check func(*reddit.Post) bool) {
	miz := p.newMizerie(subs, check)

	handler := func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		post := miz.Get()
		if post == nil {
			return true, nil
		}

		niste := chooseRandStr(ce)

		var msg string
		if len(args) > 0 {
			msg = fmt.Sprintf("%s, ia niste %s de la %s: %s NSFW (https://redd.it/%s)", args[0], niste, source.Name, post.URL, post.ID)
		} else {
			msg = fmt.Sprintf("%s, ia niste %s: %s NSFW (https://redd.it/%s)", source.Name, niste, post.URL, post.ID)
		}

		p.bot.Message(bot.PrivMsg(target, msg))
		return true, nil

	}

	for _, cmd := range cmds {
		p.bot.HandleCmdRateLimited("cmd."+cmd, handler)
	}
}

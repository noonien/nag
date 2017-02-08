package plugins

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/noonien/nag/bot"
	"github.com/sorcix/irc"
)

type Misc struct {
	bot *bot.Bot
}

func (p *Misc) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b

	p.bot.HandleCmdRateLimited("cmd.digi", p.digi())

	p.textCmd("cmd.macanache", []string{"nu fi dilimache"})
	p.textCmd("cmd.satraiesti", []string{"satz traiasca familia boss", "sa traiesti boss"})
	p.textCmd("cmd.noroc", []string{"hai noroc"})

	p.bot.HandleCmdRateLimited("cmd.ba", p.ba)
	p.bot.HandleCmdRateLimited("cmd.bă", p.ba)

	ai := "Altă întrebare! (https://soundcloud.com/armies/alta-intrebare-armies-edit)"
	p.textCmd("cmd.next", []string{ai})
	p.textReply("irc.privmsg", ai, func(line string) bool {
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

func (p *Misc) textCmd(cmd string, texts []string) {
	if len(texts) == 0 {
		return
	}

	handler := func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		text := texts[rand.Intn(len(texts))]

		if len(args) > 0 {
			text = args[0] + ": " + text
		}

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

func (p *Misc) digi() bot.CmdHandler {
	var mu sync.Mutex
	var shows []digiShow

	getShow := func() (digiShow, error) {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()

		for {
			for i, show := range shows {
				if show.when.Before(now) {
					continue
				}

				if i == 0 {
					break
				}

				return show, nil
			}

			var err error
			shows, err = getTodaysDigiShows()
			if err != nil {
				return digiShow{}, err
			}
		}
	}

	return func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		show, err := getShow()
		if err != nil {
			return false, err
		}

		p.bot.Message(bot.PrivMsg(target, "http://www.digi24.ro/live/digi24 | Emisiune: "+show.title))
		return true, nil
	}
}

type digiShow struct {
	when  time.Time
	title string
}

func getTodaysDigiShows() ([]digiShow, error) {
	today := time.Now().Truncate(24 * time.Hour)
	fmt.Println("getting shows")

	yesterdayShows, err := getDigiShows(today.Add(-24 * time.Hour))
	if err != nil {
		return nil, err
	}

	todayShows, err := getDigiShows(today)
	if err != nil {
		return nil, err
	}

	return append(yesterdayShows, todayShows...), nil
}

func getDigiShows(day time.Time) ([]digiShow, error) {
	day = day.Truncate(24 * time.Hour)
	today := day.Format("02/01/2006")

	url := fmt.Sprintf("http://www.rcs-rds.ro/asistenta/program-tv?cid=904&data=%d", day.Unix())
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var shows []digiShow
	var isTomorrow bool
	doc.Find("table.vtable tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")

		whenText := strings.TrimSpace(tds.Eq(0).Text())
		var titleParts []string
		tds.Eq(1).Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) != "#text" {
				return
			}

			titleParts = append(titleParts, strings.TrimSpace(s.Text()))
		})

		title := strings.Join(titleParts, " ")
		title = strings.Replace(title, "Cu", "cu", 1)
		title = strings.Replace(title, " si", ",", -1)

		when, _ := time.Parse("02/01/2006 15:04", today+" "+whenText)

		if isTomorrow || (len(shows) > 0 && shows[len(shows)-1].when.After(when)) {
			isTomorrow = true
			when = when.Add(24 * time.Hour)
		}

		shows = append(shows, digiShow{when: when, title: title})
	})

	if err != nil {
		return nil, err
	}

	return shows, nil
}

func (p *Misc) ba(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
	if len(args) == 0 {
		return true, nil
	}

	perms, err := p.bot.Auth(source)
	if err != nil {
		return false, err
	}

	if perms == nil || !perms.Can("annoy") {
		return true, nil
	}

	lines := []string{
		"%s",
		"%s!",
		"ba %s",
		"%s %[1]s %[1]s %[1]s %[1]s",
		"%s %[1]s %[1]s %[1]s",
		"%s anplm",
	}

	times := rand.Intn(3) + 3
	for i := 0; i < times; i++ {
		line := lines[rand.Intn(len(lines))]
		msg := fmt.Sprintf(line, args[0])
		p.bot.Message(bot.PrivMsg(target, msg))

		time.Sleep(time.Duration(rand.Intn(300)+300) * time.Millisecond)
	}
	return true, nil
}
